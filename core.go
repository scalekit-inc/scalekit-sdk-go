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
	tokenEndpoint      = "oauth/token"
	jwksEndpoint       = "keys"
	sdkVersionNumber   = "2.8.0"
	sdkVersion         = "Scalekit-Go/" + sdkVersionNumber
	defaultHTTPTimeout = 10 * time.Second
	maxErrorBodyBytes  = 8 * 1024

	// grpcReadIdleTimeout/grpcPingTimeout mirror the keepalive settings the Java
	// and Python SDKs use on the same gRPC channel (ManagedChannelBuilder.keepAliveTime /
	// grpc.keepalive_time_ms): once no frame has been read for grpcReadIdleTimeout,
	// the transport sends an HTTP/2 PING to verify the connection, whether or not a
	// stream is open — this is how a connection silently killed by a network
	// intermediary (LB, NAT, proxy) gets detected and replaced instead of being
	// written to and failing with a raw transport error.
	//
	// 60s clears the backend's EnforcementPolicy.MinTime (30s, scalekit's
	// cmd/grpc.go) with the same margin Java and Python use. A value close to
	// 30s risks ordinary timing jitter tripping the server's ping-abuse
	// detector and getting GOAWAY'd — the exact bug this setting avoids.
	grpcReadIdleTimeout = 60 * time.Second
	grpcPingTimeout     = 10 * time.Second

	// grpcIdleConnTimeout must stay below GCP's fixed 600s HTTPS load balancer
	// backend idle timeout that fronts the Scalekit API, so a fully idle
	// connection is closed by this client before the LB silently drops it.
	// 5 minutes matches the margin the backend server uses for the same reason
	// (grpcKeepaliveMaxConnectionIdle in scalekit's cmd/grpc.go). Using a bare
	// *http2.Transport directly (see newGrpcHTTPClient) opts out of
	// http.Transport's own IdleConnTimeout default, so this must be set
	// explicitly or idle connections would never close on their own.
	grpcIdleConnTimeout = 5 * time.Minute
)

// withDefaultTimeout attaches a defaultHTTPTimeout deadline to ctx if it has
// no deadline yet, returning the wrapped context and its cancel function.
// If ctx already has a deadline it is returned unchanged alongside a no-op
// cancel, so callers can always safely defer cancel() in both cases.
func withDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, defaultHTTPTimeout)
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

	ctx, cancel := withDefaultTimeout(r.Context())
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
	apiVersion := "20260727"
	client := &coreClient{
		sdkVersion:   sdkVersion,
		apiVersion:   apiVersion,
		userAgent:    fmt.Sprintf("%s Go/%s (%s; %s)", sdkVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH),
		envUrl:       envUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
	client.httpClient = &http.Client{
		Transport: &headerInterceptor{
			t:      http.DefaultTransport,
			client: client,
		},
	}
	client.grpcHTTPClient = newGrpcHTTPClient(envUrl)

	return client
}

// newGrpcHTTPClient returns a dedicated *http.Client for the gRPC connect-go
// clients, backed by its own *http2.Transport rather than the shared
// http.DefaultClient/http.DefaultTransport singleton other packages in the same
// process may reconfigure.
//
// Every real Scalekit environment is HTTPS, where this tuned *http2.Transport
// applies. It requires TLS (AllowHTTP is false) and has no dial override, so it
// cannot reach a plain "http://" endpoint at all — which is exactly what our own
// tests use (httptest.NewServer, not NewTLSServer) to exercise the connect-go
// client without a real backend. Rather than teach the production transport to
// also speak cleartext HTTP/2 (which would mean hand-rolling the TLS dial
// override and carrying that risk into the path real traffic uses), a non-https
// envUrl falls back to http.DefaultTransport's own automatic negotiation
// (HTTP/2-over-TLS, HTTP/1.1 over plain), matching this SDK's behavior before
// this transport was introduced.
func newGrpcHTTPClient(envUrl string) *http.Client {
	if !strings.HasPrefix(envUrl, "https://") {
		return &http.Client{Transport: http.DefaultTransport}
	}
	return &http.Client{
		Transport: &http2.Transport{
			ReadIdleTimeout: grpcReadIdleTimeout,
			PingTimeout:     grpcPingTimeout,
			IdleConnTimeout: grpcIdleConnTimeout,
		},
	}
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
