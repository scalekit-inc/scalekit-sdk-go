package scalekit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// backendKeepaliveMaxConnectionIdle is the backend's own gRPC keepalive
// MaxConnectionIdle (scalekit's cmd/grpc.go, grpcKeepaliveMaxConnectionIdle).
// The derived idle-close timeout (idleConnTimeoutFor) must stay strictly
// below this, not merely below the looser GCP load balancer window
// core_keepalive_test.go already pins: landing exactly on (or above) the
// backend's own bound is a race, not a fix — whichever side's timer fires
// first wins, and the loser is a request written into a socket the other
// side just closed.
const backendKeepaliveMaxConnectionIdle = 5 * time.Minute

func TestGrpcIdleConnTimeoutBelowBackendMaxConnectionIdle(t *testing.T) {
	_, transport := newGrpcHTTPClient(grpcReadIdleTimeout, grpcPingTimeout)

	require.Less(t, transport.IdleConnTimeout, backendKeepaliveMaxConnectionIdle,
		"IdleConnTimeout must stay strictly below the backend's own MaxConnectionIdle, or the client and server race to close the connection first")

	const minMargin = 30 * time.Second
	require.LessOrEqual(t, transport.IdleConnTimeout, backendKeepaliveMaxConnectionIdle-minMargin,
		"IdleConnTimeout should clear the backend's MaxConnectionIdle with real margin, not just barely")
}

func TestRetryBackoff(t *testing.T) {
	for attempt := 0; attempt < 8; attempt++ {
		d := retryBackoff(attempt)
		require.Greater(t, d, time.Duration(0), "backoff must be positive")
		require.LessOrEqual(t, d, retryBackoffMax, "backoff must never exceed the ceiling")
	}

	// Attempt 0 must be jittered around retryBackoffBase (0.5x-1.0x), not the
	// full un-jittered value every time.
	sawBelowBase := false
	for i := 0; i < 50; i++ {
		if retryBackoff(0) < retryBackoffBase {
			sawBelowBase = true
			break
		}
	}
	require.True(t, sawBelowBase, "retryBackoff(0) should sometimes land below retryBackoffBase due to jitter")
}

// fakeUnavailable returns a *connect.Error with CodeUnavailable, matching
// what connect-go's own duplex_http_call.go produces for a raw transport
// failure (see isUnavailable's comment).
func fakeUnavailable() error {
	return connect.NewError(connect.CodeUnavailable, context.DeadlineExceeded)
}

// fakeCardinalityViolation returns the exact *connect.Error connect-go's
// receiveUnaryMessage synthesizes for a unary response that hits io.EOF
// before any message arrives — verified live against a real toxiproxy
// `timeout` toxic (force-closes the connection after N seconds of silence)
// in the fault-injection run this fix came from. See
// isZeroMessageCardinalityViolation.
func fakeCardinalityViolation() error {
	// Built from the shared constant, not re-derived — if connect-go's actual
	// wording (zeroMessageCardinalityViolation) ever changes, this fake
	// should change with it rather than silently keep passing against a
	// string that no longer matches reality. See also
	// TestIsZeroMessageCardinalityViolationAgainstRealConnectGo, which
	// doesn't use this fake at all and instead triggers the real
	// receiveUnaryMessage code path.
	return connect.NewError(connect.CodeUnimplemented, errors.New(zeroMessageCardinalityViolation))
}

// fakeGenuinelyUnimplemented returns a *connect.Error shaped like what a
// server would actually send back for an RPC it doesn't support — same code
// as fakeCardinalityViolation, deliberately different message, so
// isZeroMessageCardinalityViolation must tell the two apart.
func fakeGenuinelyUnimplemented() error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("method ListWidgets not implemented"))
}

func newTestExecuterWithFailure(failures int, failErr func() error) (*connectExecuter[string, string], *int) {
	calls := 0
	fn := func(_ context.Context, _ *connect.Request[string]) (*connect.Response[string], error) {
		calls++
		if calls <= failures {
			return nil, failErr()
		}
		return connect.NewResponse(new(string)), nil
	}
	c := newCoreClient("https://example.scalekit.dev", "client-id", "client-secret")
	token := "test-token"
	c.accessToken.Store(&token) // pre-seed so exec() skips the real authenticateClient network call
	req := ""
	return newConnectExecuter[string, string](c, fn, &req), &calls
}

func newTestExecuter(failures int) (*connectExecuter[string, string], *int) {
	return newTestExecuterWithFailure(failures, fakeUnavailable)
}

func TestExecRetriesOnUnavailableThenSucceeds(t *testing.T) {
	exec, calls := newTestExecuter(defaultMaxUnavailableRetries)
	exec.maxUnavailableRetries = defaultMaxUnavailableRetries

	// Speed the test up: don't actually wait out the real backoff ceiling.
	start := time.Now()
	_, err := exec.exec(context.Background())
	require.NoError(t, err)
	require.Equal(t, defaultMaxUnavailableRetries+1, *calls, "must retry exactly up to the budget, then succeed")
	require.Less(t, time.Since(start), retryBackoffMax*time.Duration(defaultMaxUnavailableRetries+1),
		"sanity: total wait should stay within a small multiple of the backoff ceiling")
}

func TestExecStopsRetryingAfterBudgetExhausted(t *testing.T) {
	exec, calls := newTestExecuter(defaultMaxUnavailableRetries + 5)

	_, err := exec.exec(context.Background())
	require.Error(t, err)
	require.True(t, isUnavailable(err))
	require.Equal(t, defaultMaxUnavailableRetries+1, *calls,
		"must stop after maxUnavailableRetries retries (maxUnavailableRetries+1 total attempts), surfacing the error")
}

func TestExecWithMaxUnavailableRetryZeroDisablesRetry(t *testing.T) {
	exec, calls := newTestExecuter(1)
	exec.WithMaxUnavailableRetry(0)

	_, err := exec.exec(context.Background())
	require.Error(t, err)
	require.True(t, isUnavailable(err))
	require.Equal(t, 1, *calls, "maxUnavailableRetries=0 must surface the first failure immediately, no retry")
}

func TestExecUnavailableRetryRespectsContextCancellation(t *testing.T) {
	exec, calls := newTestExecuter(defaultMaxUnavailableRetries + 5)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled before the first backoff sleep begins

	_, err := exec.exec(ctx)
	require.Error(t, err)
	// The first attempt still runs (exec always tries once before any retry
	// logic), but the backoff wait must observe the cancelled context instead
	// of sleeping the full delay.
	require.LessOrEqual(t, *calls, 2)
}

func TestIsZeroMessageCardinalityViolation(t *testing.T) {
	require.True(t, isZeroMessageCardinalityViolation(fakeCardinalityViolation()),
		"must recognize connect-go's client-synthesized zero-message error")

	require.False(t, isZeroMessageCardinalityViolation(fakeGenuinelyUnimplemented()),
		"must NOT treat a genuine server-sent Unimplemented as retryable — same code, different message")

	require.False(t, isZeroMessageCardinalityViolation(fakeUnavailable()),
		"must not match on an unrelated code")

	require.False(t, isZeroMessageCardinalityViolation(nil))
}

func TestIsRetryableTransientFailureCoversBothCases(t *testing.T) {
	require.True(t, isRetryableTransientFailure(fakeUnavailable()))
	require.True(t, isRetryableTransientFailure(fakeCardinalityViolation()))
	require.False(t, isRetryableTransientFailure(fakeGenuinelyUnimplemented()))
}

// TestExecRetriesOnCardinalityViolationThenSucceeds is the regression test
// for the toxiproxy `timeout` toxic finding: a clean mid-stream close (zero
// messages, no RST) surfaces as CodeUnimplemented, not CodeUnavailable, and
// must still be retried like any other dead-connection failure.
func TestExecRetriesOnCardinalityViolationThenSucceeds(t *testing.T) {
	exec, calls := newTestExecuterWithFailure(defaultMaxUnavailableRetries, fakeCardinalityViolation)

	_, err := exec.exec(context.Background())
	require.NoError(t, err)
	require.Equal(t, defaultMaxUnavailableRetries+1, *calls, "must retry exactly up to the budget, then succeed")
}

func TestExecDoesNotRetryGenuineUnimplemented(t *testing.T) {
	exec, calls := newTestExecuterWithFailure(1, fakeGenuinelyUnimplemented)

	_, err := exec.exec(context.Background())
	require.Error(t, err)
	require.Equal(t, 1, *calls, "a genuine Unimplemented must fail immediately, never retried")
}

// TestExecBoundsEntireCallIncludingBackoff is the regression test for the PR
// #89 review finding: exec() previously applied callTimeout fresh to every
// individual attempt (via newHeaderInterceptor), and the backoff waits
// between retries were unbounded entirely — so with a persistently failing
// call, an exec() invocation could run to roughly (retries+1)*callTimeout
// instead of callTimeout, contradicting WithCallTimeout's documented "bounds
// every call" contract. With a short callTimeout and a fn that always fails
// transiently, the very first backoff wait (averaging ~0.75s at the default
// jitter) must blow through the deadline and exec() must fail with
// ctx.Err() well before exhausting the retry budget.
func TestExecBoundsEntireCallIncludingBackoff(t *testing.T) {
	exec, calls := newTestExecuterWithFailure(defaultMaxUnavailableRetries+5, fakeUnavailable)
	exec.coreClient.callTimeout = 50 * time.Millisecond

	start := time.Now()
	_, err := exec.exec(context.Background())
	elapsed := time.Since(start)

	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, elapsed, 500*time.Millisecond,
		"must fail within roughly callTimeout, not run through the full backoff/retry sequence")
	require.LessOrEqual(t, *calls, 2,
		"should fail during the first backoff wait, not exhaust all retry attempts")
}

// TestExecStillHonorsCallerSuppliedDeadline confirms exec()'s new deadline
// wrap is a no-op when the caller's context already has one — withDefaultTimeout
// only attaches callTimeout when ctx has no deadline of its own, so a caller
// that wants a different (shorter) bound than the client-wide default isn't
// silently overridden by it.
func TestExecStillHonorsCallerSuppliedDeadline(t *testing.T) {
	calls := 0
	fn := func(ctx context.Context, _ *connect.Request[string]) (*connect.Response[string], error) {
		calls++
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		return connect.NewResponse(new(string)), nil
	}
	c := newCoreClient("https://example.scalekit.dev", "client-id", "client-secret")
	c.callTimeout = 20 * time.Second // must NOT be what actually fires below
	token := "test-token"
	c.accessToken.Store(&token)
	req := ""
	exec := newConnectExecuter[string, string](c, fn, &req)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	time.Sleep(15 * time.Millisecond) // let the caller-supplied (short) deadline actually pass

	_, err := exec.exec(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Equal(t, 1, calls,
		"the already-expired caller deadline must surface from the call itself, not from a 20s callTimeout override")
}

// TestIsZeroMessageCardinalityViolationAgainstRealConnectGo drives an actual
// connect-go gRPC unary call against a real (test) server that responds with
// a clean, successful (grpc-status: 0) stream end but zero messages — the
// same well-formed-but-cardinality-violating shape a toxiproxy `timeout`
// toxic produces (verified live against production in the fault-injection
// pass this fix came from). Unlike fakeCardinalityViolation (which
// constructs the *connect.Error by hand), this exercises connect-go's real
// receiveUnaryMessage code path directly, so a future connectrpc.com/connect
// upgrade that reworks this error shape gets caught by this test failing,
// not by isZeroMessageCardinalityViolation silently going stale.
func TestIsZeroMessageCardinalityViolationAgainstRealConnectGo(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/grpc")
		w.Header().Set("Trailer", "Grpc-Status")
		w.WriteHeader(http.StatusOK)
		// Deliberately zero bytes of body — no gRPC data frame at all — then
		// a trailer signaling a clean, successful end of stream. A unary
		// call requires exactly one message; ending cleanly with none is
		// exactly the cardinality violation isZeroMessageCardinalityViolation
		// exists to catch.
		w.Header().Set("Grpc-Status", "0")
	}))
	server.EnableHTTP2 = true
	server.StartTLS()
	defer server.Close()

	client := connect.NewClient[wrapperspb.StringValue, wrapperspb.StringValue](
		server.Client(),
		server.URL+"/scalekit.v1.test.TestService/TestMethod",
		connect.WithGRPC(),
	)
	_, err := client.CallUnary(context.Background(), connect.NewRequest(&wrapperspb.StringValue{}))

	require.Error(t, err)
	require.True(t, isZeroMessageCardinalityViolation(err),
		"a real connect-go unary call ending with zero messages must be recognized as a cardinality violation")
	require.False(t, isUnavailable(err),
		"this specific shape must not also be classified as CodeUnavailable")
}
