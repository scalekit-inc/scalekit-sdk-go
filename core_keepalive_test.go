package scalekit

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
)

// backendKeepaliveMinTime is the backend's gRPC keepalive
// EnforcementPolicy.MinTime (scalekit's cmd/grpc.go, grpcKeepaliveMinTime).
// Any client-side ping interval at or below this, with insufficient margin,
// risks the server treating ordinary jitter as ping abuse and sending
// GOAWAY(ENHANCE_YOUR_CALM) mid-call or mid-idle. This mirrors the exact
// invariant scalekit-sdk-java and scalekit-sdk-python's keepalive fixes pin
// against, and the mistake caught in scalekit-sdk-python#195 (an initial
// default of exactly 30s, matching MinTime with zero margin).
const backendKeepaliveMinTime = 30 * time.Second

// gcpLoadBalancerBackendIdleTimeout is GCP's fixed backend idle timeout on the
// HTTPS load balancer fronting the Scalekit API. Not configurable from the
// client side.
const gcpLoadBalancerBackendIdleTimeout = 600 * time.Second

// Not run in parallel: one case uses t.Setenv (HTTPS_PROXY), which Go's
// testing package forbids once any test in the tree has called t.Parallel.
func TestGrpcKeepaliveInvariants(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			// grpcReadIdleTimeout governs how often a health-check PING fires
			// on the gRPC transport, whether or not a stream is open. Set it
			// too close to the backend's MinTime and ordinary timing jitter
			// can trip the server's ping-abuse detector.
			name: "ReadIdleTimeout clears the backend MinTime with real margin",
			run: func(t *testing.T) {
				_, transport := newGrpcHTTPClient()

				require.Greater(t, transport.ReadIdleTimeout, backendKeepaliveMinTime,
					"ReadIdleTimeout must stay above the backend's MinTime, or the server will treat this health-check ping as abuse")

				const minMargin = 10 * time.Second
				require.GreaterOrEqual(t, transport.ReadIdleTimeout, backendKeepaliveMinTime+minMargin,
					"ReadIdleTimeout should clear the backend's MinTime with real margin, not just barely")
			},
		},
		{
			name: "PingTimeout is a positive, bounded wait",
			run: func(t *testing.T) {
				_, transport := newGrpcHTTPClient()

				require.Greater(t, transport.PingTimeout, time.Duration(0))
				require.Less(t, transport.PingTimeout, transport.ReadIdleTimeout,
					"PingTimeout must be shorter than the interval between health checks, or checks could overlap")
			},
		},
		{
			// IdleConnTimeout must be set explicitly: http2.ConfigureTransports
			// leaves it at zero (no limit) unless told otherwise, so a fully
			// idle connection would never be closed client-side and could be
			// reused stale after GCP's LB drops it.
			name: "IdleConnTimeout is a finite bound below the GCP load balancer idle window",
			run: func(t *testing.T) {
				_, transport := newGrpcHTTPClient()

				require.Greater(t, transport.IdleConnTimeout, time.Duration(0),
					"IdleConnTimeout must be a finite bound; zero means no limit")
				require.Less(t, transport.IdleConnTimeout, gcpLoadBalancerBackendIdleTimeout,
					"IdleConnTimeout must stay below GCP's fixed LB backend idle timeout, or the LB drops the connection first")

				const minMargin = 2 * time.Minute
				require.LessOrEqual(t, transport.IdleConnTimeout, gcpLoadBalancerBackendIdleTimeout-minMargin,
					"IdleConnTimeout should clear the LB's window with real margin, not just barely")
			},
		},
		{
			// A bare *http2.Transport has no Proxy field at all and silently
			// ignores HTTPS_PROXY/NO_PROXY, breaking any customer behind a
			// corporate proxy. newGrpcHTTPClient must configure HTTP/2 onto a
			// real *http.Transport (via http2.ConfigureTransports) rather than
			// using a bare *http2.Transport directly, so proxy env vars still
			// apply to gRPC traffic.
			name: "gRPC client honors HTTPS_PROXY/NO_PROXY via the underlying http.Transport",
			run: func(t *testing.T) {
				client, _ := newGrpcHTTPClient()
				transport, ok := client.Transport.(*http.Transport)
				require.True(t, ok, "gRPC http.Client must be backed by *http.Transport (configured for HTTP/2 via http2.ConfigureTransports), not a bare *http2.Transport, or it loses proxy support")
				require.NotNil(t, transport.Proxy, "Transport.Proxy must be set (e.g. http.ProxyFromEnvironment) so HTTPS_PROXY/NO_PROXY are honored")

				req, err := http.NewRequest(http.MethodGet, "https://example.scalekit.dev", nil)
				require.NoError(t, err)
				t.Setenv("HTTPS_PROXY", "http://proxy.internal.example:8080")
				proxyURL, err := transport.Proxy(req)
				require.NoError(t, err)
				require.Equal(t, &url.URL{Scheme: "http", Host: "proxy.internal.example:8080"}, proxyURL,
					"Transport.Proxy must actually route through HTTPS_PROXY when it is set")
			},
		},
		{
			// The returned *http.Transport must still reach a plain "http://"
			// endpoint over HTTP/1.1 itself (its normal behavior for a
			// non-TLS request) — which is exactly what this SDK's own tests
			// use (httptest.NewServer, not NewTLSServer) to exercise the
			// connect-go client without a real backend. A bare *http2.Transport
			// rejects that scheme outright ("http2: unencrypted HTTP/2 not
			// enabled"), which broke CI once already (see git history).
			name: "gRPC client still reaches a plain HTTP test server",
			run: func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}))
				defer server.Close()

				client, _ := newGrpcHTTPClient()
				resp, err := client.Get(server.URL)
				require.NoError(t, err)
				defer resp.Body.Close()
				require.Equal(t, http.StatusNoContent, resp.StatusCode)
			},
		},
		{
			// All connect-go RPC service clients (see newConnectClient) must
			// share one coreClient's grpcHTTPClient, so they multiplex over
			// one pooled HTTP/2 connection instead of each opening their own.
			// A regression that has newConnectClient build a fresh client per
			// call (e.g. calling newGrpcHTTPClient() directly instead of
			// reading c.grpcHTTPClient) would silently multiply live
			// connections to the backend by the number of RPC services.
			name: "newConnectClient reuses coreClient.grpcHTTPClient, not a fresh one per call",
			run: func(t *testing.T) {
				c := newCoreClient("https://example.scalekit.dev", "client-id", "client-secret")
				require.NotNil(t, c.grpcHTTPClient)

				var gotFirst, gotSecond connect.HTTPClient
				capture := func(target *connect.HTTPClient) func(connect.HTTPClient, string, ...connect.ClientOption) *struct{} {
					return func(httpClient connect.HTTPClient, _ string, _ ...connect.ClientOption) *struct{} {
						*target = httpClient
						return &struct{}{}
					}
				}
				newConnectClient(c, capture(&gotFirst))
				newConnectClient(c, capture(&gotSecond))

				require.Same(t, c.grpcHTTPClient, gotFirst)
				require.Same(t, gotFirst, gotSecond,
					"two RPC service clients built from the same coreClient must receive the identical *http.Client")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
