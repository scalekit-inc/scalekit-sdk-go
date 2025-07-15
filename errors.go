package scalekit

import "errors"

// Error types
var (
	ErrRefreshTokenRequired  = errors.New("refresh token is required")
	ErrTokenExpired          = errors.New("token has expired")
	ErrInvalidExpClaimFormat = errors.New("invalid exp claim format")
	ErrAuthRequestIdRequired = errors.New("authRequestId is required")
)
