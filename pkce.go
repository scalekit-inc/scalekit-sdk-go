package scalekit

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const (
	// https://datatracker.ietf.org/doc/html/rfc7636
	pkceCodeChallengeMethodS256 = "S256"
	pkceDefaultVerifierLength   = 64
	pkceMinVerifierLength       = 43
	pkceMaxVerifierLength       = 128
	pkceVerifierCharset         = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
)

// PKCEOptions configures PKCE code verifier and challenge generation.
type PKCEOptions struct {
	// Optional: Sets the PKCE challenge method.
	// Defaults to "S256" when omitted.
	// Required when provided: must be "S256".
	CodeChallengeMethod string
	// Optional: Sets generated code verifier length when CodeVerifier is empty.
	// Must be between 43 and 128. Defaults to 64.
	VerifierLength int
	// Optional: Provides a precomputed code verifier.
	// Required when provided: must be 43-128 chars and contain only
	// unreserved URI characters [A-Z a-z 0-9 - . _ ~].
	// VerifierLength is ignored when this field is set.
	CodeVerifier string
}

// PKCEConfiguration contains generated PKCE parameters for auth code flow.
type PKCEConfiguration struct {
	CodeChallenge       string
	CodeVerifier        string
	CodeChallengeMethod string
}

// GeneratePKCEConfiguration generates PKCE values for OAuth authorization code flow.
//
// Optional: options.CodeChallengeMethod defaults to "S256" when omitted.
// Required when provided: options.CodeChallengeMethod must be "S256".
//
// Optional: options.VerifierLength controls generated code verifier length when
// options.CodeVerifier is empty. It must be between 43 and 128 and defaults to 64.
//
// Optional: options.CodeVerifier can be provided directly. When provided,
// options.VerifierLength is ignored. The value must be 43-128 characters and
// use only unreserved URI characters [A-Z a-z 0-9 - . _ ~].
func (s *scalekitClient) GeneratePKCEConfiguration(options PKCEOptions) (*PKCEConfiguration, error) {
	method, err := normalizeCodeChallengeMethod(options.CodeChallengeMethod)
	if err != nil {
		return nil, err
	}

	verifier := options.CodeVerifier
	if verifier != "" {
		if err := validateCodeVerifier(verifier); err != nil {
			return nil, err
		}
	} else {
		verifierLength := options.VerifierLength
		if verifierLength == 0 {
			verifierLength = pkceDefaultVerifierLength
		}
		verifier, err = generateCodeVerifier(verifierLength)
		if err != nil {
			return nil, err
		}
	}

	challenge := verifier
	if method == pkceCodeChallengeMethodS256 {
		hash := sha256.Sum256([]byte(verifier))
		challenge = base64.RawURLEncoding.EncodeToString(hash[:])
	}

	return &PKCEConfiguration{
		CodeChallenge:       challenge,
		CodeVerifier:        verifier,
		CodeChallengeMethod: method,
	}, nil
}

func normalizeCodeChallengeMethod(method string) (string, error) {
	if method == "" {
		return pkceCodeChallengeMethodS256, nil
	}
	switch strings.ToUpper(method) {
	case pkceCodeChallengeMethodS256:
		return pkceCodeChallengeMethodS256, nil
	default:
		return "", fmt.Errorf("unsupported code challenge method: %s (only S256 is supported)", method)
	}
}

func generateCodeVerifier(length int) (string, error) {
	if length < pkceMinVerifierLength || length > pkceMaxVerifierLength {
		return "", fmt.Errorf("code verifier length must be between %d and %d", pkceMinVerifierLength, pkceMaxVerifierLength)
	}

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	verifier := make([]byte, length)
	for i, b := range randomBytes {
		verifier[i] = pkceVerifierCharset[int(b)%len(pkceVerifierCharset)]
	}

	return string(verifier), nil
}

func validateCodeVerifier(verifier string) error {
	if len(verifier) < pkceMinVerifierLength || len(verifier) > pkceMaxVerifierLength {
		return fmt.Errorf("code verifier length must be between %d and %d", pkceMinVerifierLength, pkceMaxVerifierLength)
	}
	for _, ch := range verifier {
		if !strings.ContainsRune(pkceVerifierCharset, ch) {
			return errors.New("code verifier must contain only unreserved URI characters [A-Z a-z 0-9 - . _ ~]")
		}
	}
	return nil
}
