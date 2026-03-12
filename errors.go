package scalekit

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Sentinel errors returned by SDK methods. Use errors.Is to check for specific conditions.
var (
	// ErrRefreshTokenRequired is returned when RefreshAccessToken is called with an empty refresh token.
	ErrRefreshTokenRequired = errors.New("refresh token is required")

	// ErrTokenRequired is returned when a token argument is required but was empty.
	ErrTokenRequired = errors.New("token is required")

	// ErrTokenExpired is returned when a JWT's exp claim is in the past.
	ErrTokenExpired = errors.New("token has expired")

	// ErrMissingExpClaim is returned when a JWT has no exp claim.
	ErrMissingExpClaim = errors.New("token missing required exp claim")

	// ErrAuthRequestIdRequired is returned when ResendPasswordlessEmail is called without an authRequestId.
	ErrAuthRequestIdRequired = errors.New("authRequestId is required")

	// ErrTokenValidationFailed is returned when opaque-token validation fails due to the token
	// being invalid, revoked, expired, or not found (joined with the original Connect error
	// via errors.Join; use errors.Is or errors.As to inspect the underlying cause).
	ErrTokenValidationFailed = errors.New("token validation failed")

	// ErrCodeOrLinkTokenRequired is returned when VerifyPasswordlessEmail is called without
	// a Code or LinkToken in the options.
	ErrCodeOrLinkTokenRequired = errors.New("code or link token is required")

	// ErrOrganizationIdRequired is returned when an organizationId argument is required but was empty.
	ErrOrganizationIdRequired = errors.New("organizationId is required")

	// ErrClientIdRequired is returned when a clientId argument is required but was empty.
	ErrClientIdRequired = errors.New("clientId is required")

	// ErrSecretIdRequired is returned when a secretId argument is required but was empty.
	ErrSecretIdRequired = errors.New("secretId is required")

	// ErrDirectoryNotFound is returned when no directory exists for the given organization.
	ErrDirectoryNotFound = errors.New("directory does not exist for organization")

	// Deprecated: ErrInvalidExpClaimFormat is now an alias for ErrMissingExpClaim; they are
	// the same error value so existing errors.Is checks continue to work unchanged.
	// Note: the previous error fired when exp was present but had an unexpected type; that
	// path is eliminated — exp parse errors now surface as json.Unmarshal errors.
	// Use ErrMissingExpClaim to check for an absent exp claim.
	ErrInvalidExpClaimFormat = ErrMissingExpClaim

	// ErrCodeRequired is returned when AuthenticateWithCode is called with an empty code.
	ErrCodeRequired = errors.New("code is required")

	// ErrRedirectUriRequired is returned when AuthenticateWithCode is called with an empty redirectUri.
	ErrRedirectUriRequired = errors.New("redirectUri is required")

	// ErrAuthenticationResponseMissingIdToken is returned when the auth response has no id_token.
	ErrAuthenticationResponseMissingIdToken = errors.New("authentication response missing id_token")

	// ErrMissingRequiredHeaders is returned when webhook verification is missing required headers.
	ErrMissingRequiredHeaders = errors.New("missing required headers")

	// ErrInvalidSecret is returned when the webhook secret format is invalid.
	ErrInvalidSecret = errors.New("invalid secret")

	// ErrInvalidSignature is returned when the webhook signature does not match.
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrMessageTimestampTooOld is returned when the webhook message timestamp is too far in the past.
	ErrMessageTimestampTooOld = errors.New("message timestamp too old")

	// ErrMessageTimestampTooNew is returned when the webhook message timestamp is too far in the future.
	ErrMessageTimestampTooNew = errors.New("message timestamp too new")

	// ErrJwksFunctionRequired is returned when ValidateToken is called with a nil jwks function.
	ErrJwksFunctionRequired = errors.New("jwks function is required")

	// ErrAuthenticationResponseMissingAccessToken is returned when the authentication response has no access_token.
	ErrAuthenticationResponseMissingAccessToken = errors.New("authentication response missing access_token")

	// ErrJwksEmptyKeySet is returned when the JWKS endpoint returns a key set with no keys.
	ErrJwksEmptyKeySet = errors.New("JWKS endpoint returned empty key set")
)

// errorCore holds the common fields for SDK errors. Unexported so HTTPError can embed it without
// the embedded field name shadowing the Error() method (embedding "Error" would shadow).
type errorCore struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *errorCore) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *errorCore) Unwrap() error {
	return e.Err
}

// Error represents an error returned by the SDK. Currently it only covers non-2xx HTTP responses.
// Use errors.As(err, &e) to extract: var e *scalekit.Error; if errors.As(err, &e) { code := e.StatusCode }.
type Error struct {
	errorCore
}

// httpErrorFromResponse reads the response body (capped at maxErrorBodyBytes to avoid
// unbounded memory use on server-controlled error payloads) and returns an HTTPError for
// non-success responses. The prefix is used in the error message (e.g. "authentication failed").
// The caller is responsible for closing resp.Body; this function may consume part or all of it.
func httpErrorFromResponse(resp *http.Response, prefix string) *Error {
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes+1))
	if readErr != nil {
		return &Error{errorCore: errorCore{StatusCode: resp.StatusCode, Message: readErr.Error(), Err: readErr}}
	}
	msg := strings.TrimSpace(string(body))
	if len(body) > maxErrorBodyBytes {
		msg = strings.TrimSpace(string(body[:maxErrorBodyBytes])) + " …(truncated)"
	}
	err := fmt.Errorf("%s: HTTP %d: %s", prefix, resp.StatusCode, msg)
	return &Error{errorCore: errorCore{StatusCode: resp.StatusCode, Message: err.Error(), Err: err}}
}
