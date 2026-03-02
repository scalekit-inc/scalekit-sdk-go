package scalekit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/errdetails"
)

const maxTransientRetries = 3

type fn[TRequest interface{}, TResponse interface{}] func(
	context.Context,
	*connect.Request[TRequest],
) (*connect.Response[TResponse], error)

type connectExecuter[TRequest interface{}, TResponse interface{}] struct {
	coreClient       *coreClient
	data             *TRequest
	retries          int // retries for unauthenticated errors; compared against maxRetry.
	maxRetry         int
	transientRetries int // retries for transient Connect errors (ResourceExhausted, Unavailable); capped at maxTransientRetries with exponential backoff starting at 100ms.
	fn               fn[TRequest, TResponse]
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
				c.tokenMu.RLock()
				token := c.accessToken
				c.tokenMu.RUnlock()
				if token != nil {
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
		maxRetry:   1,
		fn:         fn,
	}
}

// validationErrorMessage builds a single string from a Connect InvalidArgument
// error message and its validation field violations.
func validationErrorMessage(ce *connect.Error) string {
	messages := []string{ce.Message()}
	for _, detail := range ce.Details() {
		msg, err := detail.Value()
		if err != nil {
			messages = append(messages, fmt.Sprintf("[unreadable validation detail: %v]", err))
			continue
		}
		info, ok := msg.(*errdetails.ErrorInfo)
		if !ok || info.ValidationErrorInfo == nil {
			continue
		}
		for _, field := range info.ValidationErrorInfo.FieldViolations {
			messages = append(messages, fmt.Sprintf("%s: %s", field.Field, field.Description))
		}
	}
	return strings.Join(messages, "\n")
}

// isTransientCode reports whether the Connect code indicates a transient error
// that may succeed on retry (ResourceExhausted, Unavailable).
func isTransientCode(code connect.Code) bool {
	switch code {
	case connect.CodeResourceExhausted, connect.CodeUnavailable:
		return true
	default:
		return false
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
		var connectErr *connect.Error
		_ = errors.As(err, &connectErr)

		if connectErr != nil && connectErr.Code() == connect.CodeInvalidArgument {
			return nil, fmt.Errorf("%s: %w", validationErrorMessage(connectErr), connectErr)
		}

		if connectErr != nil && isTransientCode(connectErr.Code()) && r.transientRetries < maxTransientRetries {
			backoff := time.Duration(1<<r.transientRetries) * 100 * time.Millisecond
			timer := time.NewTimer(backoff)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return nil, fmt.Errorf("context cancelled during retry backoff (attempt %d/%d): %w",
					r.transientRetries+1, maxTransientRetries, ctx.Err())
			}
			r.transientRetries++
			return r.exec(ctx)
		}

		if connectErr != nil && isTransientCode(connectErr.Code()) {
			return nil, fmt.Errorf("operation failed after %d transient retries: %w", r.transientRetries, err)
		}

		if r.maxRetry-r.retries > 0 && isUnauthenticated(err) {
			if authErr := r.coreClient.authenticateClient(ctx); authErr != nil {
				return nil, fmt.Errorf("reauthentication failed after %v: %w", err, authErr)
			}
			r.retries++
			r.transientRetries = 0 // reset for fresh attempt after re-auth
			return r.exec(ctx)
		}

		return nil, err
	}

	return data.Msg, nil
}

func (r *connectExecuter[TRequest, TResponse]) WithMaxRetry(retry int) *connectExecuter[TRequest, TResponse] {
	r.maxRetry = retry
	return r
}
