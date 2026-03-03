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
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
)

const (
	tokenEndpoint      = "oauth/token"
	jwksEndpoint       = "keys"
	sdkVersion         = "Scalekit-Go/2.2.0"
	defaultHTTPTimeout = 10 * time.Second
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

	tokenMu     sync.RWMutex
	accessToken *string

	jwksMu        sync.RWMutex
	jsonWebKeySet *jose.JSONWebKeySet

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

type httpError struct {
	err        error
	StatusCode int
}

func (h *httpError) Error() string {
	return h.err.Error()
}

// Unwrap exposes the underlying error so errors.As and errors.Is can traverse the chain.
func (h *httpError) Unwrap() error {
	return h.err
}

// httpErrorFromResponse reads the response body and returns an httpError for
// non-success responses. The prefix is used in the error message (e.g. "authentication failed").
// The caller is responsible for closing resp.Body; this function reads but does not close it.
func httpErrorFromResponse(resp *http.Response, prefix string) *httpError {
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return &httpError{
			err:        fmt.Errorf("%s: HTTP %d: body read error: %w", prefix, resp.StatusCode, readErr),
			StatusCode: resp.StatusCode,
		}
	}
	return &httpError{
		err:        fmt.Errorf("%s: HTTP %d: %s", prefix, resp.StatusCode, strings.TrimSpace(string(body))),
		StatusCode: resp.StatusCode,
	}
}

func (h *headerInterceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("user-agent", h.client.userAgent)
	r.Header.Add("x-sdk-version", h.client.sdkVersion)
	r.Header.Add("x-api-version", h.client.apiVersion)
	// Read the token pointer under a read lock. Defer is deliberately not used:
	// using defer would hold the lock across the subsequent network call.
	h.client.tokenMu.RLock()
	token := h.client.accessToken
	h.client.tokenMu.RUnlock()
	if token != nil {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	ctx, cancel := withDefaultTimeout(r.Context())
	defer cancel()

	return h.t.RoundTrip(r.WithContext(ctx))
}

func newCoreClient(envUrl, clientId, clientSecret string) *coreClient {
	sdkVersion := sdkVersion
	apiVersion := "20260226"
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
	requestData := url.Values{}
	requestData.Set("grant_type", "client_credentials")
	requestData.Set("client_id", c.clientId)
	requestData.Set("client_secret", c.clientSecret)
	res, err := c.authenticate(ctx, requestData)
	if err != nil {
		return err
	}
	// Lock scope is a single pointer assignment; explicit unlock is used for clarity rather than defer.
	c.tokenMu.Lock()
	c.accessToken = &res.AccessToken
	c.tokenMu.Unlock()

	return nil
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
		return nil, errors.New("authentication response missing access_token")
	}

	return &responseData, nil
}

func (c *coreClient) GetJwks(ctx context.Context) (*jose.JSONWebKeySet, error) {
	// Read lock is released explicitly rather than deferred because this goroutine
	// acquires a write lock below — promoting from RLock to Lock on the same mutex deadlocks.
	// The pointer is copied under the read lock so the returned value remains valid after the lock is released.
	c.jwksMu.RLock()
	if c.jsonWebKeySet != nil {
		jwks := c.jsonWebKeySet
		c.jwksMu.RUnlock()
		return jwks, nil
	}
	c.jwksMu.RUnlock()

	c.jwksMu.Lock()
	defer c.jwksMu.Unlock()
	// Double-checked locking: another goroutine may have populated jsonWebKeySet while this goroutine waited for the write lock.
	if c.jsonWebKeySet != nil {
		return c.jsonWebKeySet, nil
	}

	request, err := http.NewRequestWithContext(ctx,
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
		return nil, errors.New("JWKS endpoint returned empty key set")
	}
	c.jsonWebKeySet = &responseData

	return &responseData, nil
}
