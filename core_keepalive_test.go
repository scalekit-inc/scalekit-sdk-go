package scalekit

import (
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
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

func TestGrpcKeepaliveInvariants(t *testing.T) {
	t.Parallel()

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
				client := newGrpcHTTPClient()
				transport, ok := client.Transport.(*http2.Transport)
				require.True(t, ok, "gRPC http.Client must be backed by *http2.Transport to support ReadIdleTimeout/PingTimeout")

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
				client := newGrpcHTTPClient()
				transport := client.Transport.(*http2.Transport)

				require.Greater(t, transport.PingTimeout, time.Duration(0))
				require.Less(t, transport.PingTimeout, transport.ReadIdleTimeout,
					"PingTimeout must be shorter than the interval between health checks, or checks could overlap")
			},
		},
		{
			// IdleConnTimeout must be set explicitly: a bare *http2.Transport
			// (unlike http.DefaultTransport) defaults this to 0, meaning no
			// limit, so a fully idle connection would never be closed
			// client-side and could be reused stale after GCP's LB drops it.
			name: "IdleConnTimeout is a finite bound below the GCP load balancer idle window",
			run: func(t *testing.T) {
				client := newGrpcHTTPClient()
				transport := client.Transport.(*http2.Transport)

				require.Greater(t, transport.IdleConnTimeout, time.Duration(0),
					"IdleConnTimeout must be a finite bound; zero means no limit on a bare *http2.Transport")
				require.Less(t, transport.IdleConnTimeout, gcpLoadBalancerBackendIdleTimeout,
					"IdleConnTimeout must stay below GCP's fixed LB backend idle timeout, or the LB drops the connection first")

				const minMargin = 2 * time.Minute
				require.LessOrEqual(t, transport.IdleConnTimeout, gcpLoadBalancerBackendIdleTimeout-minMargin,
					"IdleConnTimeout should clear the LB's window with real margin, not just barely")
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
			t.Parallel()
			tt.run(t)
		})
	}
}
