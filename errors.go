package scalekit

import "errors"

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
	// being invalid, revoked, expired, or not found (wraps the original Connect error).
	ErrTokenValidationFailed = errors.New("token validation failed")

	// ErrCodeOrLinkTokenRequired is returned when VerifyPasswordlessEmail is called without
	// a Code or LinkToken in the options.
	ErrCodeOrLinkTokenRequired = errors.New("code or link token is required")

	// ErrOrganizationIdRequired is returned when an organizationId argument is required but was empty.
	ErrOrganizationIdRequired = errors.New("organizationId is required")

	// ErrDirectoryNotFound is returned when no directory exists for the given organization.
	ErrDirectoryNotFound = errors.New("directory does not exist for organization")

	// Deprecated: use ErrMissingExpClaim instead.
	ErrInvalidExpClaimFormat = ErrMissingExpClaim
)
