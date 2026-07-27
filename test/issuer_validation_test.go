package test

import (
	"context"
	"os"
	"testing"

	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestValidateTokenOptionsIssuerField is a pure-unit assertion that the Issuer
// field exists on ValidateTokenOptions and is settable.
func TestValidateTokenOptionsIssuerField(t *testing.T) {
	options := scalekit.ValidateTokenOptions{Issuer: "https://example.scalekit.dev"}
	assert.Equal(t, "https://example.scalekit.dev", options.Issuer)
}

// TestValidateTokenWithOptionsWrongIssuer verifies that supplying an issuer that
// does not match the token's iss claim fails validation. Guarded with t.Skip
// when the required environment values are absent.
func TestValidateTokenWithOptionsWrongIssuer(t *testing.T) {
	accessToken := os.Getenv("SCALEKIT_TEST_ACCESS_TOKEN")
	issuer := os.Getenv("SCALEKIT_TEST_ISSUER")
	if accessToken == "" || issuer == "" {
		t.Skip("set SCALEKIT_TEST_ACCESS_TOKEN and SCALEKIT_TEST_ISSUER to run this test")
	}

	ctx := context.Background()

	// Correct issuer should pass.
	valid, err := client.ValidateTokenWithOptions(ctx, accessToken, &scalekit.ValidateTokenOptions{
		Issuer: issuer,
	})
	require.NoError(t, err)
	assert.True(t, valid)

	// Wrong issuer must fail validation with the issuer-mismatch sentinel, so
	// this branch cannot pass on an unrelated token or JWKS failure.
	valid, err = client.ValidateTokenWithOptions(ctx, accessToken, &scalekit.ValidateTokenOptions{
		Issuer: issuer + "-wrong",
	})
	assert.ErrorIs(t, err, scalekit.ErrIssuerMismatch)
	assert.False(t, valid)
}
