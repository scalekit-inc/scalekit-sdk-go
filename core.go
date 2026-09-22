package scalekit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-jose/go-jose/v4"
	"golang.org/x/net/http2"
	"golang.org/x/sync/singleflight"
)

const (
	tokenEndpoint     = "oauth/token"
	jwksEndpoint      = "keys"
	sdkVersionNumber  = "2.8.1"
	sdkVersion        = "Scalekit-Go/" + sdkVersionNumber
	maxErrorBodyBytes = 8 * 1024

	// defaultCallTimeout bounds every call (REST token/JWKS calls and every
	// gRPC RPC) that doesn't already carry its own context deadline — without
	// it, a call can block forever on a connection that looks fine to the
	// client but is silently dead. Matches the Python SDK's call_timeout_s and
	// the Node SDK's timeoutMs default (both 20s) exactly. Override per
	// instance with WithCallTimeout, or per call by passing a context that
	// already has a deadline (checked first, see withDefaultTimeout) — Go
	// callers have that idiom natively; Python/Node don't, which is why their
	// only lever is the constructor-level default.
	defaultCallTimeout = 20 * time.Second

	// grpcReadIdleTimeout/grpcPingTimeout are this client's active
	// connection-liveness check: once no frame has been read for
	// grpcReadIdleTimeout, the transport sends an HTTP/2 PING to verify the
	// connection, whether or not a stream is open — this is how a connection
	// silently killed by a network intermediary (LB, NAT, proxy) gets detected
	// and replaced instead of being written to and failing with a raw
	// transport error. Mirrors the keepalive settings the Java and Python
	// SDKs use on the same gRPC channel (ManagedChannelBuilder.keepAliveTime /
	// grpc.keepalive_time_ms).
	//
	// 60s was originally chosen to clear the backend's EnforcementPolicy.MinTime
	// (30s, scalekit's cmd/grpc.go) with margin. That policy turned out to be
	// inert in production: the backend serves all traffic through
	// grpc.Server.ServeHTTP, and grpc-go's handler-transport for ServeHTTP
	// does not accept or apply KeepaliveEnforcementPolicy/KeepaliveParams at
	// all (verified against grpc-go's vendored source — NewServerHandlerTransport's
	// signature takes no keepalive arguments, and internal/transport/handler_server.go
	// has no reference to enforcement or idle handling). See scalekit's
	// cmd/server.go IdleTimeout comment for the investigation that surfaced
	// this. The VALUE is kept unchanged regardless: 60s remains a reasonable,
	// conservative liveness-ping cadence on its own terms, independent of
	// whether the backend penalizes faster ones, and stays consistent across
	// all four SDKs.
	// These are the DEFAULTS; override per instance with WithKeepAlive.
	grpcReadIdleTimeout = 60 * time.Second
	grpcPingTimeout     = 10 * time.Second

	// grpcMinPingInterval is the floor WithKeepAlive enforces on a
	// caller-supplied ping interval (0 is separately allowed as the "disabled"
	// escape hatch — see validateKeepAlive). Matches the Python SDK's
	// MIN_KEEPALIVE_TIME_MS and the Node SDK's MIN_PING_INTERVAL_MS exactly.
	// Originally floored to clear the backend's 30s EnforcementPolicy.MinTime
	// with margin — see grpcReadIdleTimeout's comment above for why that
	// specific enforcement turned out to be inert in production. The floor is
	// kept regardless, as a sane minimum independent of backend enforcement.
	grpcMinPingInterval = 60 * time.Second

	// grpcIdleConnCeiling bounds how long a fully idle gRPC connection is kept
	// before this client proactively closes it (see idleConnTimeoutFor).
	//
	// This client never connects to the Scalekit backend directly — GCLB
	// terminates and re-originates the connection, so this setting governs
	// only the SDK<->GCLB leg, never the GCLB<->pod leg. The real constraint
	// on THIS leg is GCLB's client-facing (GFE) keepalive timeout: 610s by
	// default (confirmed as the live, unmodified value on the Scalekit API's
	// production target-https-proxy — httpKeepAliveTimeoutSec is unset),
	// configurable up to 1200s. Staying below it means this client always
	// closes an idle connection first, so it is never blindsided by GCLB
	// silently retiring a connection this client still considers valid.
	//
	// (An earlier revision of this comment instead cited the backend's own
	// grpcKeepaliveMaxConnectionIdle, cmd/grpc.go. That setting governs a
	// DIFFERENT connection — GCLB<->pod, which this client never touches —
	// and is additionally inert in production regardless of leg; see
	// grpcReadIdleTimeout's comment above.)
	//
	// 4 minutes clears GCLB's 610s default with a comfortable margin. Mirrors
	// the Node SDK's IDLE_CONNECTION_TIMEOUT_CEILING_MS (connect.ts) exactly.
	grpcIdleConnCeiling = 4 * time.Minute

	// grpcIdleConnPingCycles is the multiplier idleConnTimeoutFor applies to a
	// caller-supplied ping interval before clamping to grpcIdleConnCeiling —
	// see that function's comment. Matches the Node SDK's
	// IDLE_PING_CYCLES_BEFORE_CLOSE exactly.
	grpcIdleConnPingCycles = 5
)

// withDefaultTimeout attaches c.callTimeout as a deadline to ctx if it has no
// deadline yet, returning the wrapped context and its cancel function. If ctx
// already has a deadline it is returned unchanged alongside a no-op cancel,
// so callers can always safely defer cancel() in both cases.
func (c *coreClient) withDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, c.callTimeout)
}

// idleConnTimeoutFor derives the proactive idle-connection-close timeout from
// a ping interval, exactly mirroring the Node SDK's
// idleConnectionTimeoutMsFor (connect.ts): pingInterval * grpcIdleConnPingCycles,
// clamped to grpcIdleConnCeiling. The multiplier is intentionally superseded
// by the ceiling for every currently-valid non-zero ping interval —
// grpcMinPingInterval (60s) * 5 = 300s already exceeds the 4-min ceiling — but
// stays in the formula (rather than being dropped for a bare constant) so
// this still scales down correctly should the floor on ping interval itself
// ever be lowered.
func idleConnTimeoutFor(pingInterval time.Duration) time.Duration {
	scaled := pingInterval * grpcIdleConnPingCycles
	if scaled <= 0 || scaled > grpcIdleConnCeiling {
		return grpcIdleConnCeiling
	}
	return scaled
}

// validateKeepAlive panics if pingInterval/pingTimeout is an invalid
// combination — see WithKeepAlive. Mirrors the Python SDK's keepalive
// validation (grpc silently clamps sub-10s values instead of rejecting them,
// which is exactly the silent-misconfiguration failure mode this guards
// against) and the Node SDK's assertValidPingInterval/assertValidTimeout.
func validateKeepAlive(pingInterval, pingTimeout time.Duration) {
	if pingInterval == 0 {
		return // the deliberate "disabled" escape hatch; pingTimeout is unused in this state.
	}
	if pingInterval < grpcMinPingInterval {
		panic(fmt.Sprintf(
			"scalekit: WithKeepAlive: pingInterval must be 0 (disabled) or >= %s; got %s. "+
				"A value below the default leaves too little margin over the Scalekit server's "+
				"30s keepalive MinTime — early pings are struck as abusive, and enough strikes "+
				"GOAWAYs the connection mid-call.", grpcMinPingInterval, pingInterval))
	}
	if pingTimeout <= 0 {
		panic(fmt.Sprintf("scalekit: WithKeepAlive: pingTimeout must be positive, got %s", pingTimeout))
	}
	if pingTimeout >= pingInterval {
		panic(fmt.Sprintf(
			"scalekit: WithKeepAlive: pingTimeout (%s) must be less than pingInterval (%s), or the "+
				"interval timer can fire again before a hung ping would ever be detected as hung.",
			pingTimeout, pingInterval))
	}
}

type coreClient struct {
	envUrl       string
	clientId     string
	clientSecret string
	sdkVersion   string
	apiVersion   string
	userAgent    string

	accessToken atomic.Pointer[string]
	authGroup   singleflight.Group

	jwksGroup     singleflight.Group
	jsonWebKeySet atomic.Pointer[jose.JSONWebKeySet]

	httpClient *http.Client
	// grpcHTTPClient backs every connect-go RPC service client built via
	// newConnectClient. Built once here and shared across all of them (see
	// newConnectClient), the same way http.DefaultClient used to be shared —
	// except scoped to this coreClient instead of the whole process, and with
	// keepalive settings tuned for the backend's enforcement policy.
	grpcHTTPClient *http.Client

	// pingInterval/pingTimeout are the values grpcHTTPClient's transport was
	// last built with — kept here (rather than only inside the *http2.Transport)
	// so WithSecret can carry them forward onto a new coreClient instead of
	// silently resetting a caller's WithKeepAlive override back to the
	// defaults. Set via WithKeepAlive; default to grpcReadIdleTimeout/
	// grpcPingTimeout.
	pingInterval time.Duration
	pingTimeout  time.Duration

	// callTimeout is the deadline withDefaultTimeout applies to a call whose
	// context has none. Set via WithCallTimeout; defaults to defaultCallTimeout.
	callTimeout time.Duration
}

type authenticationResponse struct {
	IdToken      string `json:"id_token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type headerInterceptor struct {
	t      http.RoundTripper
	client *coreClient
}

// cancelOnClose wraps an io.ReadCloser and invokes cancel when Close is
// called.  This defers context cancellation until after the caller has
// finished reading the response body, rather than cancelling at the point
// RoundTrip returns (before the body has been consumed).
type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c *cancelOnClose) Close() error {
	defer c.cancel()
	return c.ReadCloser.Close()
}

func (h *headerInterceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("user-agent", h.client.userAgent)
	r.Header.Add("x-sdk-version", h.client.sdkVersion)
	r.Header.Add("x-api-version", h.client.apiVersion)
	if token := h.client.accessToken.Load(); token != nil {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	ctx, cancel := h.client.withDefaultTimeout(r.Context())
	resp, err := h.t.RoundTrip(r.WithContext(ctx))
	if err != nil {
		cancel()
		return nil, err
	}
	resp.Body = &cancelOnClose{ReadCloser: resp.Body, cancel: cancel}
	return resp, nil
}

func newCoreClient(envUrl, clientId, clientSecret string) *coreClient {
	sdkVersion := sdkVersion
	apiVersion := "20260922"
	client := &coreClient{
		sdkVersion:   sdkVersion,
		apiVersion:   apiVersion,
		userAgent:    fmt.Sprintf("%s Go/%s (%s; %s)", sdkVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH),
		envUrl:       envUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
		pingInterval: grpcReadIdleTimeout,
		pingTimeout:  grpcPingTimeout,
		callTimeout:  defaultCallTimeout,
	}
	client.httpClient = &http.Client{
		Transport: &headerInterceptor{
			// A dedicated *http.Transport, not the shared http.DefaultTransport
			// singleton — the same process-global isolation concern
			// newGrpcHTTPClient already addresses for the gRPC client. Unrelated
			// code elsewhere in the same process reconfiguring/replacing
			// http.DefaultTransport (a real thing customers do, e.g. to inject
			// instrumentation) would otherwise silently affect every Scalekit
			// REST call (token exchange, JWKS) too. Proxy: http.ProxyFromEnvironment
			// preserves the HTTPS_PROXY/NO_PROXY support http.DefaultTransport
			// provided.
			t:      &http.Transport{Proxy: http.ProxyFromEnvironment},
			client: client,
		},
	}
	client.grpcHTTPClient, _ = newGrpcHTTPClient(client.pingInterval, client.pingTimeout)

	return client
}

// newGrpcHTTPClient returns a dedicated *http.Client for the gRPC connect-go
// clients, instead of the shared http.DefaultClient/http.DefaultTransport
// singleton other packages in the same process may reconfigure. It also
// returns the *http2.Transport handle backing that client, purely so tests can
// assert on the timeout values below without duplicating them.
//
// pingInterval/pingTimeout come from WithKeepAlive (or its defaults,
// grpcReadIdleTimeout/grpcPingTimeout) and must already be validated — see
// validateKeepAlive. pingInterval == 0 is the deliberate "disabled" escape
// hatch: ReadIdleTimeout/PingTimeout/IdleConnTimeout are all left at their Go
// zero values (no health-check ping, no proactive idle close), matching
// connect-node's own untuned defaults — for a network path that rejects our
// probing pattern entirely, mirroring the Python SDK's keepalive_time_ms=0
// and the Node SDK's pingIntervalMs=0.
//
// It configures HTTP/2 onto a private *http.Transport (http2.ConfigureTransports)
// rather than using a bare *http2.Transport directly, for two reasons:
//   - A bare *http2.Transport has no Proxy field at all, silently ignoring
//     HTTPS_PROXY/NO_PROXY for customers behind a corporate proxy — proxy
//     support lives on *http.Transport, which is what ConfigureTransports
//     upgrades in place while returning the *http2.Transport handle used below
//     only to set the timeouts.
//   - The returned *http.Transport still handles a plain "http://" envUrl over
//     HTTP/1.1 itself (its normal behavior for a non-TLS request), which is
//     exactly what our own tests use (httptest.NewServer, not NewTLSServer) to
//     exercise the connect-go client without a real backend — a bare
//     *http2.Transport rejects that scheme outright ("http2: unencrypted
//     HTTP/2 not enabled") since AllowHTTP defaults false and it has no
//     cleartext dial override.
func newGrpcHTTPClient(pingInterval, pingTimeout time.Duration) (*http.Client, *http2.Transport) {
	t1 := &http.Transport{Proxy: http.ProxyFromEnvironment}
	t2, err := http2.ConfigureTransports(t1)
	if err != nil {
		// Can only fail if t1 were already HTTP/2-enabled, which a freshly
		// constructed *http.Transport never is.
		panic(fmt.Sprintf("scalekit: unreachable: configuring HTTP/2 on a fresh transport failed: %v", err))
	}
	if pingInterval != 0 {
		t2.ReadIdleTimeout = pingInterval
		t2.PingTimeout = pingTimeout
	}
	// Independent of pingInterval: unlike ReadIdleTimeout/PingTimeout, this
	// doesn't send anything over the wire — it's a purely local pool-hygiene
	// timer — so disabling the active health-check ping (pingInterval == 0,
	// for a network path that rejects that probing pattern) is no reason to
	// also leave idle connections completely unbounded. idleConnTimeoutFor(0)
	// already resolves to grpcIdleConnCeiling.
	t2.IdleConnTimeout = idleConnTimeoutFor(pingInterval)
	return &http.Client{Transport: t1}, t2
}

func (c *coreClient) authenticateClient(ctx context.Context) error {
	if c.clientSecret == "" {
		return ErrClientSecretRequired
	}
	_, err, _ := c.authGroup.Do("auth", func() (any, error) {
		requestData := url.Values{}
		requestData.Set("grant_type", "client_credentials")
		requestData.Set("client_id", c.clientId)
		requestData.Set("client_secret", c.clientSecret)
		// Use WithoutCancel so one caller's context cancellation does not fail all waiters.
		res, err := c.authenticate(context.WithoutCancel(ctx), requestData)
		if err != nil {
			return nil, err
		}
		c.accessToken.Store(&res.AccessToken)
		return nil, nil
	})
	return err
}

func (c *coreClient) authenticate(ctx context.Context, requestData url.Values) (*authenticationResponse, error) {
	request, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s", c.envUrl, tokenEndpoint),
		strings.NewReader(requestData.Encode()),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Add(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	// Close errors are intentionally ignored; the response body is fully consumed or discarded below.
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, httpErrorFromResponse(response, "authentication failed")
	}
	var responseData authenticationResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}
	if responseData.AccessToken == "" {
		return nil, ErrAuthenticationResponseMissingAccessToken
	}

	return &responseData, nil
}

func (c *coreClient) GetJwks(ctx context.Context) (*jose.JSONWebKeySet, error) {
	if cached := c.jsonWebKeySet.Load(); cached != nil {
		return copyJSONWebKeySet(cached), nil
	}
	v, err, _ := c.jwksGroup.Do("jwks", func() (any, error) {
		if cached := c.jsonWebKeySet.Load(); cached != nil {
			return copyJSONWebKeySet(cached), nil
		}
		// Use WithoutCancel so one caller's context cancellation does not fail all waiters.
		request, err := http.NewRequestWithContext(context.WithoutCancel(ctx),
			http.MethodGet,
			fmt.Sprintf("%s/%s", c.envUrl, jwksEndpoint),
			nil,
		)
		if err != nil {
			return nil, err
		}
		response, err := c.httpClient.Do(request)
		if err != nil {
			return nil, err
		}
		// Close errors are intentionally ignored; the response body is fully consumed or discarded below.
		defer func() { _ = response.Body.Close() }()
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return nil, httpErrorFromResponse(response, "failed to fetch JWKS")
		}
		var responseData jose.JSONWebKeySet
		err = json.NewDecoder(response.Body).Decode(&responseData)
		if err != nil {
			return nil, err
		}
		if len(responseData.Keys) == 0 {
			return nil, ErrJwksEmptyKeySet
		}
		c.jsonWebKeySet.Store(&responseData)
		return copyJSONWebKeySet(&responseData), nil
	})
	if err != nil {
		return nil, err
	}
	jwks, ok := v.(*jose.JSONWebKeySet)
	if !ok {
		return nil, errors.New("internal: unexpected JWKS result type")
	}
	return jwks, nil
}

// copyJSONWebKeySet returns a shallow copy of the key set so callers cannot mutate the internal cache (e.g. the Keys slice).
func copyJSONWebKeySet(src *jose.JSONWebKeySet) *jose.JSONWebKeySet {
	if src == nil {
		return nil
	}
	keys := make([]jose.JSONWebKey, len(src.Keys))
	copy(keys, src.Keys)
	return &jose.JSONWebKeySet{Keys: keys}
}

func (c *coreClient) hasAccessToken() bool {
	token := c.accessToken.Load()
	return token != nil && *token != ""
}
