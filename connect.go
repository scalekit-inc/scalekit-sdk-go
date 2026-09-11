package scalekit

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
)

// retryBackoffBase/retryBackoffMax and the jitter formula in retryBackoff
// below intentionally match the Node and Python SDKs' UNAVAILABLE-retry
// backoff exactly (scalekit-sdk-node's core.ts, scalekit-sdk-python#195/#198):
// base doubles each attempt up to a 30s ceiling, then half-jittered (0.5-1.0x)
// so a fleet of clients retrying the same overloaded backend doesn't retry in
// lockstep.
const (
	retryBackoffBase = 1 * time.Second
	retryBackoffMax  = 30 * time.Second

	// defaultMaxUnavailableRetries matches the Python SDK's default (grpc_exec's
	// retry=2). Independent of maxRetries (the unauthenticated-retry budget)
	// since the two failure modes aren't related — a caller retrying auth
	// shouldn't be charged against its transient-connection retry budget or
	// vice versa.
	defaultMaxUnavailableRetries = 2
)

type fn[TRequest interface{}, TResponse interface{}] func(
	context.Context,
	*connect.Request[TRequest],
) (*connect.Response[TResponse], error)

type connectExecuter[TRequest interface{}, TResponse interface{}] struct {
	coreClient *coreClient
	data       *TRequest
	retries    int // retries for unauthenticated errors; compared against maxRetries.
	maxRetries int

	// unavailableRetries/maxUnavailableRetries retry a transient CodeUnavailable
	// (e.g. a dead/reset connection — see isUnavailable) with backoff. Set
	// maxUnavailableRetries to 0 via WithMaxUnavailableRetry at a call site
	// where even one retry risks double-executing a non-idempotent operation,
	// mirroring the Python SDK's retry_on_unavailable=False escape hatch.
	unavailableRetries    int
	maxUnavailableRetries int

	fn fn[TRequest, TResponse]
}

func newConnectClient[T interface{}](
	c *coreClient,
	fn func(
		httpClient connect.HTTPClient,
		baseURL string,
		opts ...connect.ClientOption,
	) T,
) T {
	return fn(
		c.grpcHTTPClient,
		c.envUrl,
		connect.WithGRPC(),
		connect.WithInterceptors(newHeaderInterceptor(c)),
	)
}

func newHeaderInterceptor(c *coreClient) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			ctx, cancel := c.withDefaultTimeout(ctx)
			defer cancel()
			if req.Spec().IsClient {
				req.Header().Set("user-agent", c.userAgent)
				req.Header().Set("x-sdk-version", c.sdkVersion)
				req.Header().Set("x-api-version", c.apiVersion)
				if token := c.accessToken.Load(); token != nil {
					req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", *token))
				}
			}
			return next(ctx, req)
		})
	})
}

func newConnectExecuter[TRequest interface{}, TResponse interface{}](
	coreClient *coreClient,
	fn fn[TRequest, TResponse],
	data *TRequest,
) *connectExecuter[TRequest, TResponse] {
	return &connectExecuter[TRequest, TResponse]{
		coreClient:            coreClient,
		data:                  data,
		maxRetries:            1,
		maxUnavailableRetries: defaultMaxUnavailableRetries,
		fn:                    fn,
	}
}

// isUnauthenticated reports whether err indicates an authentication failure
// (HTTP 401 or Connect CodeUnauthenticated).
func isUnauthenticated(err error) bool {
	var httpErr *Error
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnauthorized {
		return true
	}
	var connectErr *connect.Error
	return errors.As(err, &connectErr) && connectErr.Code() == connect.CodeUnauthenticated
}

// isUnavailable reports whether err is a transient CodeUnavailable failure —
// notably, connect-go's own duplex_http_call.go falls back to CodeUnavailable
// whenever the underlying http.Client.Do round-trip itself fails (a dead/reset
// connection, e.g. "use of closed network connection") and no more specific
// wrapper (context error, RST_STREAM reason, etc.) already classified it.
// That is exactly the failure mode this retry exists for.
//
// This is NOT retry-scope parity with the Node SDK, despite both retrying a
// code spelled Unavailable/CodeUnavailable: connect-node's error mapping
// reclassifies a transport reset as Code.Aborted (see scalekit-sdk-node's
// core.ts _connectExec / scalekit-sdk-python#195's grpc_exec comment for the
// documented cross-SDK divergence), so Node's Unavailable branch never
// actually fires for that failure. connect-go does not have that
// reclassification — a dead connection genuinely surfaces as CodeUnavailable
// here — so this branch is verified-correct for Go specifically, not merely
// copied from Node's naming.
func isUnavailable(err error) bool {
	var connectErr *connect.Error
	return errors.As(err, &connectErr) && connectErr.Code() == connect.CodeUnavailable
}

// zeroMessageCardinalityViolation is the exact message connect-go's
// receiveUnaryMessage (connect.go) synthesizes for a unary client-side
// receive that hits io.EOF before any message arrives. Matched verbatim
// against this SDK's pinned connectrpc.com/connect version — see
// isZeroMessageCardinalityViolation.
const zeroMessageCardinalityViolation = "unary response has zero messages"

// isZeroMessageCardinalityViolation reports whether err is connect-go's
// CodeUnimplemented "cardinality violation" for a unary response — NOT a
// genuine "the server doesn't implement this RPC" response, but a
// client-side error connect-go synthesizes itself whenever the underlying
// stream closes (io.EOF) before any message arrives: a clean-looking
// mid-stream disconnect (verified live: a toxiproxy `timeout` toxic that
// force-closes the connection after N seconds of silence produces exactly
// this) rather than an RST or a raw transport read error (both of which
// isUnavailable already catches). Per the gRPC status-code spec, both
// clients and servers are required to use CodeUnimplemented for a
// cardinality violation — see receiveUnaryMessage's own comment — so this
// exact code is legitimately reused for two unrelated situations, and only
// the message text tells them apart.
//
// String-matching an error message is fragile in principle, but connect-go
// exposes no distinct sentinel for this specific condition (it's
// deliberately reusing the standard code per spec), so this is the only way
// to distinguish "dead connection" from "genuinely unimplemented" from
// outside the connect package. The asymmetry matters: a false negative here
// (connect-go rewords the message in a future version) just means this
// specific fault mode stops being retried — no worse than before this fix
// existed. A false positive would be far worse — silently retrying a call
// against an RPC the server genuinely doesn't support — and is not possible
// here, since a real server-sent Unimplemented carries the server's own
// grpc-message, never this literal client-synthesized string.
func isZeroMessageCardinalityViolation(err error) bool {
	var connectErr *connect.Error
	if !errors.As(err, &connectErr) || connectErr.Code() != connect.CodeUnimplemented {
		return false
	}
	return strings.Contains(connectErr.Message(), zeroMessageCardinalityViolation)
}

// isRetryableTransientFailure reports whether err represents a transient
// connection-level failure worth retrying with backoff — see isUnavailable
// and isZeroMessageCardinalityViolation for the two distinct ways connect-go
// surfaces "the connection died," and exec's comment for the
// non-idempotency caveat that applies to both.
func isRetryableTransientFailure(err error) bool {
	return isUnavailable(err) || isZeroMessageCardinalityViolation(err)
}

// retryBackoff returns a jittered exponential backoff delay for unavailable
// retry attempt N (0-indexed): base doubles each attempt up to
// retryBackoffMax, then half-jittered (0.5-1.0x) — see the package-level
// comment on retryBackoffBase for why this matches Node/Python exactly.
func retryBackoff(attempt int) time.Duration {
	backoff := retryBackoffBase * time.Duration(1<<attempt)
	if backoff > retryBackoffMax || backoff <= 0 {
		backoff = retryBackoffMax
	}
	return time.Duration(float64(backoff) * (0.5 + rand.Float64()*0.5))
}

// exec runs the Connect RPC. Errors (including validation/CodeInvalidArgument) are returned
// as-is; use errors.As(err, &connectErr) with *connect.Error to inspect Code() and Details().
//
// Bounds the ENTIRE call — every retry attempt and every backoff wait — with
// one deadline, matching WithCallTimeout's documented "bounds every call"
// contract; without this, each retry got its own fresh callTimeout and the
// backoff waits were unbounded, so one exec() call could run to roughly
// (retries+1)*callTimeout instead of callTimeout. This only takes effect on
// the outermost invocation: once withDefaultTimeout attaches a deadline here,
// every per-attempt wrap further down (this function's own recursive
// re-entry on retry, newHeaderInterceptor's wrap around the gRPC attempt
// itself) sees ctx already has one and is a no-op.
//
// One deliberate exception: authenticateClient uses context.WithoutCancel
// internally (so one caller's cancellation can't fail an in-flight auth
// request other callers are sharing via singleflight — see core.go). That
// strips this deadline, so a slow initial authentication (the
// !hasAccessToken() call below, or the re-auth on a 401 retry) is NOT bounded
// by it and can consume part of the budget before the actual RPC attempt
// even starts. Accepted as consistent with "bound the whole call": if
// authentication itself is slow, the overall call should still fail within
// its documented budget rather than silently running long, not receive a
// separate unbounded allowance.
func (r *connectExecuter[TRequest, TResponse]) exec(ctx context.Context) (*TResponse, error) {
	ctx, cancel := r.coreClient.withDefaultTimeout(ctx)
	defer cancel()

	if r.coreClient.clientSecret == "" {
		return nil, ErrClientSecretRequired
	}
	if !r.coreClient.hasAccessToken() {
		// explicitly making a call before the actual call to ensure that the access token is available. Not consuming the error, so that the call can go into retries
		_ = r.coreClient.authenticateClient(ctx)
	}
	data, err := r.fn(ctx, connect.NewRequest(r.data))
	if err != nil {
		if r.maxRetries-r.retries > 0 && isUnauthenticated(err) {
			if authErr := r.coreClient.authenticateClient(ctx); authErr != nil {
				return nil, authErr
			}
			r.retries++
			return r.exec(ctx)
		}
		// A transient connection failure (UNAVAILABLE, or the zero-message
		// cardinality violation isZeroMessageCardinalityViolation catches) can
		// mean the request already reached and was processed by the server — a
		// dead/refused connection, a stream torn down mid-flight, or a keepalive
		// ping timeout on a still-in-progress call are all indistinguishable
		// from "the server did the work but the response never made it back."
		// Retrying risks double-executing a non-idempotent call; set
		// maxUnavailableRetries to 0 via WithMaxUnavailableRetry at a call site
		// where that's a real concern (see the Python SDK's
		// ToolsClient.execute_tool for the precedent). Backed off (jittered
		// exponential) rather than retried immediately: on a backend returning
		// UNAVAILABLE because it's overloaded, every client retrying instantly
		// just triples the load it's already struggling with.
		if r.maxUnavailableRetries-r.unavailableRetries > 0 && isRetryableTransientFailure(err) {
			delay := retryBackoff(r.unavailableRetries)
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
			r.unavailableRetries++
			return r.exec(ctx)
		}
		return nil, err
	}
	return data.Msg, nil
}

func (r *connectExecuter[TRequest, TResponse]) WithMaxRetry(retry int) *connectExecuter[TRequest, TResponse] {
	r.maxRetries = retry
	return r
}

// WithMaxUnavailableRetry overrides the number of times a transient
// CodeUnavailable failure is retried with backoff (default
// defaultMaxUnavailableRetries). Pass 0 to disable — see isUnavailable's
// comment on exec for why a call site might need to opt out.
func (r *connectExecuter[TRequest, TResponse]) WithMaxUnavailableRetry(retry int) *connectExecuter[TRequest, TResponse] {
	r.maxUnavailableRetries = retry
	return r
}
