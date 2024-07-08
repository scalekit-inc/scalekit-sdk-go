package scalekit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
)

const tokenEndpoint = "oauth/token"
const jwksEndpoint = "keys"

type coreClient struct {
	envUrl        string
	clientId      string
	clientSecret  string
	sdkVersion    string
	apiVersion    string
	userAgent     string
	accessToken   *string
	httpClient    *http.Client
	jsonWebKeySet *jose.JSONWebKeySet
}

type authenticationResponse struct {
	IdToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type headerInterceptor struct {
	t      http.RoundTripper
	client *coreClient
}

type httpError struct {
	err        error
	StatusCode int
}

// Error implements error.
func (h *httpError) Error() string {
	return h.err.Error()
}

func (h *headerInterceptor) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("user-agent", h.client.userAgent)
	r.Header.Add("x-sdk-version", h.client.sdkVersion)
	r.Header.Add("x-api-version", h.client.apiVersion)
	if h.client.accessToken != nil {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *h.client.accessToken))
	}

	resp, err := h.t.RoundTrip(r)
	if err != nil {
		return nil, &httpError{
			err: err,
		}
	}

	return resp, nil
}

func newCoreClient(envUrl, clientId, clientSecret string) *coreClient {
	sdkVersion := "Scalekit-Go/v1.0.0"
	apiVersion := "20240430"
	client := &coreClient{
		sdkVersion:   sdkVersion,
		apiVersion:   apiVersion,
		userAgent:    fmt.Sprintf("%s Go/%s (%s; %s)", sdkVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH),
		envUrl:       envUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
	client.httpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &headerInterceptor{
			t:      http.DefaultTransport,
			client: client,
		},
	}

	return client
}

func (c *coreClient) authenticateClient() error {
	requestData := url.Values{}
	requestData.Set("grant_type", "client_credentials")
	requestData.Set("client_id", c.clientId)
	requestData.Set("client_secret", c.clientSecret)
	res, err := c.authenticate(requestData)
	if err != nil {
		return err
	}
	c.accessToken = &res.AccessToken

	return nil
}

func (c *coreClient) authenticate(requestData url.Values) (*authenticationResponse, error) {
	request, err := http.NewRequest(
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
	defer response.Body.Close()
	var responseData authenticationResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

func (c *coreClient) getJwks() error {
	if c.jsonWebKeySet != nil {
		return nil
	}
	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%s", c.envUrl, jwksEndpoint),
		nil,
	)
	if err != nil {
		return err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var responseData jose.JSONWebKeySet
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return err
	}
	c.jsonWebKeySet = &responseData

	return nil
}
