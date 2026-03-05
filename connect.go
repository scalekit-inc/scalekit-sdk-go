package scalekit

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
)

type fn[TRequest interface{}, TResponse interface{}] func(
	context.Context,
	*connect.Request[TRequest],
) (*connect.Response[TResponse], error)

type connectExecuter[TRequest interface{}, TResponse interface{}] struct {
	coreClient *coreClient
	data       *TRequest
	retries    int // retries for unauthenticated errors; compared against maxRetry.
	maxRetries int
	fn         fn[TRequest, TResponse]
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
		http.DefaultClient,
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
			ctx, cancel := withDefaultTimeout(ctx)
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
		coreClient: coreClient,
		data:       data,
		maxRetries:   1,
		fn:         fn,
	}
}

// isUnauthenticated reports whether err indicates an authentication failure
// (HTTP 401 or Connect CodeUnauthenticated).
func isUnauthenticated(err error) bool {
	var httpErr *httpError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnauthorized {
		return true
	}
	var connectErr *connect.Error
	return errors.As(err, &connectErr) && connectErr.Code() == connect.CodeUnauthenticated
}

func (r *connectExecuter[TRequest, TResponse]) exec(ctx context.Context) (*TResponse, error) {
	data, err := r.fn(ctx, connect.NewRequest(r.data))
	if err != nil {
		if r.maxRetries-r.retries > 0 && isUnauthenticated(err) {
			if authErr := r.coreClient.authenticateClient(ctx); authErr != nil {
				return nil, authErr
			}
			r.retries++
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
