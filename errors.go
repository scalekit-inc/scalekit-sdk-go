package scalekit

import "errors"

// Error types
var (
	ErrRefreshTokenRequired       = errors.New("refresh token is required")
	ErrCodeAndRedirectUriRequired = errors.New("code and redirect uri is required")
	ErrMissingRequiredHeaders     = errors.New("missing required headers")
	ErrInvalidSecret              = errors.New("invalid secret")
	ErrMessageTimestampTooOld     = errors.New("message timestamp too old")
	ErrMessageTimestampTooNew     = errors.New("message timestamp too new")
	ErrInvalidSignature           = errors.New("invalid signature")
	ErrTokenExpired               = errors.New("token has expired")
	ErrInvalidExpClaimFormat      = errors.New("invalid exp claim format")
)
