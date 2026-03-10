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
	"golang.org/x/sync/singleflight"
)

const (
	tokenEndpoint      = "oauth/token"
	jwksEndpoint       = "keys"
	sdkVersion         = "Scalekit-Go/2.2.0"
	defaultHTTPTimeout = 10 * time.Second
	maxErrorBodyBytes  = 8 * 1024
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

	jwksGroup    singleflight.Group
	jsonWebKeySet atomic.Pointer[jose.JSONWebKeySet]

	httpClient *http.Client
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
	apiVersion := "20260310"
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

	return client
}

func (c *coreClient) authenticateClient(ctx context.Context) error {
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
