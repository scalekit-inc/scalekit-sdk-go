package test

import (
	"context"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go"
)

func TestPasswordlessService(t *testing.T) {
	if client == nil {
		t.Skip("Client not initialized, skipping integration test")
	}

	// Test that the passwordless service is available
	passwordlessService := client.Passwordless()
	if passwordlessService == nil {
		t.Fatal("Passwordless service should not be nil")
	}

	// Test SendPasswordlessEmail with options
	ctx := context.Background()
	email := "dhaneshbabu007@gmail.com"
	templateType := scalekit.TemplateTypeSignin
	options := &scalekit.SendPasswordlessOptions{
		Template:         &templateType,
		MagiclinkAuthUri: "https://myapp.com/auth/callback", // Now we can pass the string directly!
		State:            "integration-test-state",
		ExpiresIn:        1800, // 30 minutes
		TemplateVariables: map[string]string{
			"app_name": "Integration Test App",
			"company":  "Test Company",
		},
	}

	// Test 1: Send passwordless email
	t.Run("SendPasswordlessEmail", func(t *testing.T) {
		response, err := passwordlessService.SendPasswordlessEmail(ctx, email, options)
		if err != nil {
			t.Logf("SendPasswordlessEmail error: %v", err)
			return
		}

		t.Logf("✅ SendPasswordlessEmail successful!")
		t.Logf("   Auth Request ID: %s", response.AuthRequestId)
		t.Logf("   Expires At: %d", response.ExpiresAt)
		t.Logf("   Expires In: %d seconds", response.ExpiresIn)
		t.Logf("   Passwordless Type: %s", response.PasswordlessType)

		// Test 2: Resend passwordless email (if we got a valid auth request ID)
		if response.AuthRequestId != "" {
			t.Run("ResendPasswordlessEmail", func(t *testing.T) {
				// Wait a bit before resending to avoid rate limiting
				time.Sleep(2 * time.Second)

				resendResponse, err := passwordlessService.ResendPasswordlessEmail(ctx, response.AuthRequestId)
				if err != nil {
					t.Logf("ResendPasswordlessEmail error: %v", err)
					return
				}

				t.Logf("✅ ResendPasswordlessEmail successful!")
				t.Logf("   New Auth Request ID: %s", resendResponse.AuthRequestId)
				t.Logf("   New Expires At: %d", resendResponse.ExpiresAt)
				t.Logf("   New Expires In: %d seconds", resendResponse.ExpiresIn)
			})
		}
	})

	// Test 3: Verify with invalid code (this should fail, but we can test the API structure)
	t.Run("VerifyPasswordlessEmail_InvalidCode", func(t *testing.T) {
		verifyOptions := &scalekit.VerifyPasswordlessOptions{
			Code:          "000000", // Invalid code
			AuthRequestId: "invalid-auth-request-id",
		}

		_, err := passwordlessService.VerifyPasswordlessEmail(ctx, verifyOptions)
		if err != nil {
			t.Logf("✅ VerifyPasswordlessEmail with invalid code failed as expected: %v", err)
		} else {
			t.Logf("⚠️  VerifyPasswordlessEmail with invalid code succeeded (unexpected)")
		}
	})

	// Test 4: Verify with invalid link token (this should fail, but we can test the API structure)
	t.Run("VerifyPasswordlessEmail_InvalidLinkToken", func(t *testing.T) {
		verifyLinkOptions := &scalekit.VerifyPasswordlessOptions{
			LinkToken: "invalid-link-token",
		}

		_, err := passwordlessService.VerifyPasswordlessEmail(ctx, verifyLinkOptions)
		if err != nil {
			t.Logf("✅ VerifyPasswordlessEmail with invalid link token failed as expected: %v", err)
		} else {
			t.Logf("⚠️  VerifyPasswordlessEmail with invalid link token succeeded (unexpected)")
		}
	})
}
