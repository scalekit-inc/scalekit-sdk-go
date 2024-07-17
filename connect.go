package scalekit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/errdetails"
)

type fn[TRequest interface{}, TResponse interface{}] func(
	context.Context,
	*connect.Request[TRequest],
) (*connect.Response[TResponse], error)

type connectExecuter[TRequest interface{}, TResponse interface{}] struct {
	coreClient *coreClient
	data       *TRequest
	retries    int
	maxRetry   int
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
			if req.Spec().IsClient {
				req.Header().Set("user-agent", c.userAgent)
				req.Header().Set("x-sdk-version", c.sdkVersion)
				req.Header().Set("x-api-version", c.apiVersion)
				if c.accessToken != nil {
					req.Header().Set(
						"Authorization",
						fmt.Sprintf("Bearer %s", *c.accessToken),
					)
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
		maxRetry:   1,
		fn:         fn,
	}
}

func (r *connectExecuter[TRequest, TResponse]) exec(ctx context.Context) (*TResponse, error) {
	data, err := r.fn(ctx, connect.NewRequest(r.data))
	if err != nil {
		if r.maxRetry-r.retries > 0 {
			var isUnAuthenticatedError bool
			if httpErr, ok := err.(*httpError); ok {
				if httpErr.StatusCode != http.StatusUnauthorized {
					isUnAuthenticatedError = true
				}
			}
			if connectErr, ok := err.(*connect.Error); ok {
				if connectErr.Code() == connect.CodeUnauthenticated {
					isUnAuthenticatedError = true
				}
				if connectErr.Code() == connect.CodeInvalidArgument {
					messages := []string{connectErr.Message()}
					for _, detail := range connectErr.Details() {
						msg, _ := detail.Value()
						if info, ok := msg.(*errdetails.ErrorInfo); ok {
							if info.ValidationErrorInfo != nil {
								for _, field := range info.ValidationErrorInfo.FieldViolations {
									messages = append(messages, fmt.Sprintf("%s:  %s", field.Field, field.Description))
								}
							}
						}
					}

					return nil, errors.New(strings.Join(messages, "\n"))
				}
			}

			if isUnAuthenticatedError {
				_ = r.coreClient.authenticateClient()
				r.retries++
				return r.exec(ctx)
			}
		}

		return nil, err
	}

	return data.Msg, nil
}

func (r *connectExecuter[TRequest, TResponse]) WithMaxRetry(retry int) *connectExecuter[TRequest, TResponse] {
	r.maxRetry = retry
	return r
}
