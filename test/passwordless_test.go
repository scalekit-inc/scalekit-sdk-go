package test

import (
	"context"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestSendPasswordlessEmail(t *testing.T) {
	if client == nil {
		t.Skip("Client not initialized, skipping integration test")
	}

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
		return
	}

	// Assert response is not nil and contains expected fields
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AuthRequestId)
	assert.True(t, response.ExpiresAt > 0)
	assert.True(t, response.ExpiresIn > 0)
	assert.NotEmpty(t, response.PasswordlessType)
}

func TestResendPasswordlessEmail(t *testing.T) {
	if client == nil {
		t.Skip("Client not initialized, skipping integration test")
	}

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
	assert.NotEmpty(t, response.AuthRequestId)

	// Wait a bit before resending to avoid rate limiting
	time.Sleep(2 * time.Second)

	resendResponse, err := passwordlessService.ResendPasswordlessEmail(ctx, response.AuthRequestId)
	if err != nil {
		return
	}

	// Assert resend response is not nil and contains expected fields
	assert.NotNil(t, resendResponse)
	assert.NotEmpty(t, resendResponse.AuthRequestId)
	assert.True(t, resendResponse.ExpiresAt > 0)
	assert.True(t, resendResponse.ExpiresIn > 0)
}

func TestVerifyPasswordlessEmail_InvalidCode(t *testing.T) {
	if client == nil {
		t.Skip("Client not initialized, skipping integration test")
	}

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
	if client == nil {
		t.Skip("Client not initialized, skipping integration test")
	}

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

func TestVerifyPasswordlessEmail_ValidCode(t *testing.T) {
	if client == nil {
		t.Skip("Client not initialized, skipping integration test")
	}

	passwordlessService := client.Passwordless()
	ctx := context.Background()

	// TODO: Replace with actual valid code and auth request ID from email
	validCode := "123456"                              // Change this to the actual code from email
	validAuthRequestId := "auth_request_id_from_email" // Change this to the actual auth request ID from email

	verifyOptions := &scalekit.VerifyPasswordlessOptions{
		Code:          validCode,
		AuthRequestId: validAuthRequestId,
	}

	response, err := passwordlessService.VerifyPasswordlessEmail(ctx, verifyOptions)

	// Assert that verification with valid code succeeds
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Assert response contains expected fields
	if response != nil {
		assert.NotEmpty(t, response.Email)
		assert.NotEmpty(t, response.PasswordlessType)
	}
}

func TestVerifyPasswordlessEmail_ValidLinkToken(t *testing.T) {
	if client == nil {
		t.Skip("Client not initialized, skipping integration test")
	}

	passwordlessService := client.Passwordless()
	ctx := context.Background()

	// TODO: Replace with actual valid link token from email
	validLinkToken := "link_token_from_email" // Change this to the actual link token from email

	verifyOptions := &scalekit.VerifyPasswordlessOptions{
		LinkToken: validLinkToken,
	}

	response, err := passwordlessService.VerifyPasswordlessEmail(ctx, verifyOptions)

	// Assert that verification with valid link token succeeds
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Assert response contains expected fields
	if response != nil {
		assert.NotEmpty(t, response.Email)
		assert.NotEmpty(t, response.PasswordlessType)
	}
}
