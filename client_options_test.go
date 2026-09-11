package scalekit

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestScalekitClient(opts ...any) *scalekitClient {
	c := NewScalekitClient("https://example.scalekit.dev", "client-id", opts...)
	sc, ok := c.(*scalekitClient)
	if !ok {
		panic("NewScalekitClient did not return *scalekitClient")
	}
	return sc
}

func TestNewScalekitClientDefaults(t *testing.T) {
	c := newTestScalekitClient()

	require.Equal(t, grpcReadIdleTimeout, c.coreClient.pingInterval)
	require.Equal(t, grpcPingTimeout, c.coreClient.pingTimeout)
	require.Equal(t, defaultCallTimeout, c.coreClient.callTimeout)

	_, transport := newGrpcHTTPClient(c.coreClient.pingInterval, c.coreClient.pingTimeout)
	require.Equal(t, grpcReadIdleTimeout, transport.ReadIdleTimeout)
	require.Equal(t, grpcPingTimeout, transport.PingTimeout)
}

func TestNewScalekitClientWithKeepAlive(t *testing.T) {
	pingInterval := 90 * time.Second
	pingTimeout := 15 * time.Second
	c := newTestScalekitClient(WithKeepAlive(pingInterval, pingTimeout))

	require.Equal(t, pingInterval, c.coreClient.pingInterval)
	require.Equal(t, pingTimeout, c.coreClient.pingTimeout)

	// c.coreClient.grpcHTTPClient.Transport is the *http.Transport
	// (http2.ConfigureTransports upgrades it in place); the *http2.Transport
	// handle carrying the actual ReadIdleTimeout/PingTimeout fields isn't
	// retained on the client (see newGrpcHTTPClient), so re-derive it with the
	// same inputs WithKeepAlive used, mirroring core_keepalive_test.go's
	// pattern.
	_, transport := newGrpcHTTPClient(c.coreClient.pingInterval, c.coreClient.pingTimeout)
	require.Equal(t, pingInterval, transport.ReadIdleTimeout)
	require.Equal(t, pingTimeout, transport.PingTimeout)
	require.Equal(t, idleConnTimeoutFor(pingInterval), transport.IdleConnTimeout)
}

func TestNewScalekitClientWithKeepAliveDisabled(t *testing.T) {
	c := newTestScalekitClient(WithKeepAlive(0, 0))

	require.Equal(t, time.Duration(0), c.coreClient.pingInterval)
	require.Equal(t, time.Duration(0), c.coreClient.pingTimeout)
}

func TestNewScalekitClientWithCallTimeout(t *testing.T) {
	c := newTestScalekitClient(WithCallTimeout(5 * time.Second))
	require.Equal(t, 5*time.Second, c.coreClient.callTimeout)
}

func TestWithKeepAlivePanicsBelowFloor(t *testing.T) {
	require.Panics(t, func() {
		WithKeepAlive(30*time.Second, 5*time.Second).apply(newCoreClient("https://example.scalekit.dev", "id", "secret"))
	})
}

func TestWithKeepAlivePanicsWhenTimeoutNotBelowInterval(t *testing.T) {
	require.Panics(t, func() {
		WithKeepAlive(60*time.Second, 60*time.Second).apply(newCoreClient("https://example.scalekit.dev", "id", "secret"))
	})
}

func TestWithCallTimeoutPanicsOnNonPositive(t *testing.T) {
	require.Panics(t, func() {
		WithCallTimeout(0).apply(newCoreClient("https://example.scalekit.dev", "id", "secret"))
	})
	require.Panics(t, func() {
		WithCallTimeout(-1 * time.Second).apply(newCoreClient("https://example.scalekit.dev", "id", "secret"))
	})
}

func TestWithSecretCarriesForwardOptions(t *testing.T) {
	pingInterval := 120 * time.Second
	pingTimeout := 20 * time.Second
	callTimeout := 7 * time.Second
	original := newTestScalekitClient(WithKeepAlive(pingInterval, pingTimeout), WithCallTimeout(callTimeout))

	updated := original.WithSecret("new-secret").(*scalekitClient)

	require.Equal(t, pingInterval, updated.coreClient.pingInterval)
	require.Equal(t, pingTimeout, updated.coreClient.pingTimeout)
	require.Equal(t, callTimeout, updated.coreClient.callTimeout)
	require.Equal(t, "new-secret", updated.coreClient.clientSecret)
}

func TestIdleConnTimeoutForFormula(t *testing.T) {
	require.Equal(t, grpcIdleConnCeiling, idleConnTimeoutFor(grpcMinPingInterval))
	require.Equal(t, grpcIdleConnCeiling, idleConnTimeoutFor(0))

	small := 30 * time.Second // below the floor, hypothetically, to exercise the scaling branch
	require.Equal(t, small*grpcIdleConnPingCycles, idleConnTimeoutFor(small))
}

// TestHTTPClientDoesNotShareProcessGlobalTransport is the regression test for
// the PR #89 review finding: httpClient (used for /oauth/token and JWKS) must
// not wrap the shared http.DefaultTransport singleton, or unrelated code
// elsewhere in the same process reconfiguring/replacing http.DefaultTransport
// silently affects every Scalekit REST call too — the same isolation concern
// newGrpcHTTPClient already solved for the gRPC client.
func TestHTTPClientDoesNotShareProcessGlobalTransport(t *testing.T) {
	c := newCoreClient("https://example.scalekit.dev", "client-id", "client-secret")

	interceptor, ok := c.httpClient.Transport.(*headerInterceptor)
	require.True(t, ok, "httpClient.Transport must be a *headerInterceptor wrapping the real transport")
	require.NotSame(t, http.DefaultTransport, interceptor.t,
		"must not share the process-global http.DefaultTransport singleton")

	transport, ok := interceptor.t.(*http.Transport)
	require.True(t, ok, "the wrapped transport must be a real *http.Transport, not some other RoundTripper")
	require.NotNil(t, transport.Proxy, "must still honor HTTPS_PROXY/NO_PROXY via Proxy: http.ProxyFromEnvironment")

	// Two clients must not accidentally share the same transport instance
	// either — each coreClient gets its own.
	c2 := newCoreClient("https://example.scalekit.dev", "client-id", "client-secret")
	interceptor2 := c2.httpClient.Transport.(*headerInterceptor)
	require.NotSame(t, interceptor.t, interceptor2.t)
}

func TestHTTPClientHonorsHTTPSProxy(t *testing.T) {
	c := newCoreClient("https://example.scalekit.dev", "client-id", "client-secret")
	interceptor := c.httpClient.Transport.(*headerInterceptor)
	transport := interceptor.t.(*http.Transport)

	req, err := http.NewRequest(http.MethodGet, "https://example.scalekit.dev", nil)
	require.NoError(t, err)
	t.Setenv("HTTPS_PROXY", "http://proxy.internal.example:8080")
	proxyURL, err := transport.Proxy(req)
	require.NoError(t, err)
	require.Equal(t, &url.URL{Scheme: "http", Host: "proxy.internal.example:8080"}, proxyURL)
}
