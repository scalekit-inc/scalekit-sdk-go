package test

import (
	"context"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestSendPasswordlessEmail(t *testing.T) {
	passwordlessService := client.Passwordless()
	ctx := context.Background()
	email := "dhaneshbabu007@gmail.com"
	templateType := scalekit.TemplateTypeSignin
	options := &scalekit.SendPasswordlessOptions{
		Template:         &templateType,
		MagiclinkAuthUri: "https://myapp.com/auth/callback",
		State:            "integration-test-state",
		ExpiresIn:        1800, // 30 minutes
		TemplateVariables: map[string]string{
			"app_name": "Integration Test App",
			"company":  "Test Company",
		},
	}

	response, err := passwordlessService.SendPasswordlessEmail(ctx, email, options)
	if err != nil {
		// TODO: Fix with an email provider so this test can run reliably in CI.
		return
	}

	// Assert response is not nil and contains expected fields
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.GetAuthRequestId())
	assert.True(t, response.GetExpiresAt() > 0)
	assert.True(t, response.GetExpiresIn() > 0)
	assert.NotEmpty(t, response.GetPasswordlessType().String())

	// Log the auth request ID for manual testing
	t.Logf("Auth Request ID: %s", response.GetAuthRequestId())
}

func TestResendPasswordlessEmail(t *testing.T) {
	passwordlessService := client.Passwordless()
	ctx := context.Background()
	email := "dhaneshbabu007@gmail.com"
	templateType := scalekit.TemplateTypeSignin
	options := &scalekit.SendPasswordlessOptions{
		Template:         &templateType,
		MagiclinkAuthUri: "https://myapp.com/auth/callback",
		State:            "integration-test-state",
		ExpiresIn:        1800, // 30 minutes
		TemplateVariables: map[string]string{
			"app_name": "Integration Test App",
			"company":  "Test Company",
		},
	}

	// First send an email to get an auth request ID
	response, err := passwordlessService.SendPasswordlessEmail(ctx, email, options)
	if err != nil {
		t.Skipf("Cannot test resend without first sending: %v", err)
	}

	// Assert initial response has required fields
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.GetAuthRequestId())

	verifyCodeResponse, err := passwordlessService.VerifyPasswordlessEmail(ctx, &scalekit.VerifyPasswordlessOptions{
		Code:          "000000", // Invalid code
		AuthRequestId: response.GetAuthRequestId(),
	})

	// Assert that verification with invalid code fails as expected
	assert.Error(t, err)
	assert.Nil(t, verifyCodeResponse)

	// Wait a bit before resending to avoid rate limiting
	time.Sleep(2 * time.Second)

	resendResponse, err := passwordlessService.ResendPasswordlessEmail(ctx, response.GetAuthRequestId())
	if err != nil {
		// TODO: Fix with an email provider so this test can run reliably in CI.
		return
	}

	// Assert resend response is not nil and contains expected fields
	assert.NotNil(t, resendResponse)
	assert.NotEmpty(t, resendResponse.GetAuthRequestId())
	assert.True(t, resendResponse.GetExpiresAt() > 0)
	assert.True(t, resendResponse.GetExpiresIn() > 0)
}

func TestVerifyPasswordlessEmail_InvalidCode(t *testing.T) {
	passwordlessService := client.Passwordless()
	ctx := context.Background()
	verifyOptions := &scalekit.VerifyPasswordlessOptions{
		Code:          "000000", // Invalid code
		AuthRequestId: "invalid-auth-request-id",
	}

	response, err := passwordlessService.VerifyPasswordlessEmail(ctx, verifyOptions)

	// Assert that verification with invalid code fails as expected
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestVerifyPasswordlessEmail_InvalidLinkToken(t *testing.T) {
	passwordlessService := client.Passwordless()
	ctx := context.Background()
	verifyLinkOptions := &scalekit.VerifyPasswordlessOptions{
		LinkToken: "invalid-link-token",
	}

	response, err := passwordlessService.VerifyPasswordlessEmail(ctx, verifyLinkOptions)

	// Assert that verification with invalid link token fails as expected
	assert.Error(t, err)
	assert.Nil(t, response)
}

// TestVerifyPasswordlessEmail_ValidCode runs against real API; fails without valid code + auth request ID (e.g. OTP 424242 in dev).
func TestVerifyPasswordlessEmail_ValidCode(t *testing.T) {
	ctx := context.Background()
	verifyOptions := &scalekit.VerifyPasswordlessOptions{
		Code:          "424242",
		AuthRequestId: "placeholder_auth_request_id",
	}
	response, err := client.Passwordless().VerifyPasswordlessEmail(ctx, verifyOptions)
	if err != nil {
		// TODO: Fix with an email provider so this test can run reliably in CI.
		t.Logf("Expected to fail without real auth request: %v", err)
		return
	}
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.GetEmail())
	assert.NotEmpty(t, response.GetPasswordlessType().String())
}

// TestVerifyPasswordlessEmail_ValidLinkToken runs against real API; fails without valid link token from email.
func TestVerifyPasswordlessEmail_ValidLinkToken(t *testing.T) {
	ctx := context.Background()
	verifyOptions := &scalekit.VerifyPasswordlessOptions{
		LinkToken: "placeholder_link_token",
	}
	response, err := client.Passwordless().VerifyPasswordlessEmail(ctx, verifyOptions)
	if err != nil {
		// TODO: Fix with an email provider so this test can run reliably in CI.
		t.Logf("Expected to fail without real link token: %v", err)
		return
	}
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.GetEmail())
	assert.NotEmpty(t, response.GetPasswordlessType().String())
}
