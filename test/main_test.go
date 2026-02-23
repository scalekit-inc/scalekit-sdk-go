package test

import (
	"context"
	"os"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
)

var (
	client scalekit.Scalekit
)

func TestMain(m *testing.M) {
	environmentUrl := os.Getenv("SCALEKIT_ENVIRONMENT_URL")
	clientId := os.Getenv("SCALEKIT_CLIENT_ID")
	apiSecret := os.Getenv("SCALEKIT_CLIENT_SECRET")
	client = scalekit.NewScalekitClient(environmentUrl, clientId, apiSecret)

	code := m.Run()
	os.Exit(code)
}

// SkipIfNoIntegrationEnv skips the test when SCALEKIT_* env vars are not set.
func SkipIfNoIntegrationEnv(t *testing.T) {
	t.Helper()
	if os.Getenv("SCALEKIT_ENVIRONMENT_URL") == "" || os.Getenv("SCALEKIT_CLIENT_ID") == "" || os.Getenv("SCALEKIT_CLIENT_SECRET") == "" {
		t.Skip("skipping: SCALEKIT_ENVIRONMENT_URL, SCALEKIT_CLIENT_ID, SCALEKIT_CLIENT_SECRET required")
	}
}

// createOrg creates an organization and returns its ID. Caller must defer DeleteTestOrganization(t, ctx, orgId).
// Only call after SkipIfNoIntegrationEnv(t).
func createOrg(t *testing.T, ctx context.Context, name, externalID string) string {
	t.Helper()
	resp, err := client.Organization().CreateOrganization(ctx, name, scalekit.CreateOrganizationOptions{ExternalId: externalID})
	if err != nil {
		t.Fatalf("createOrg: %v", err)
	}
	if resp == nil || resp.Organization == nil {
		t.Fatal("createOrg: nil response")
	}
	return resp.Organization.Id
}

func toPtr(s string) *string {
	return &s
}

func toInt32Ptr(i int32) *int32 {
	return &i
}
