package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
)

// Required integration test environment variables (single source of truth).
const (
	EnvEnvironmentURL = "SCALEKIT_ENVIRONMENT_URL"
	EnvClientID       = "SCALEKIT_CLIENT_ID"
	EnvClientSecret   = "SCALEKIT_CLIENT_SECRET"
)

var (
	client  scalekit.Scalekit
	testOrg string
)

func TestMain(m *testing.M) {
	environmentUrl := os.Getenv(EnvEnvironmentURL)
	clientId := os.Getenv(EnvClientID)
	apiSecret := os.Getenv(EnvClientSecret)
	if environmentUrl == "" || clientId == "" || apiSecret == "" {
		fmt.Fprintf(os.Stderr, "integration tests require %s, %s, %s\n", EnvEnvironmentURL, EnvClientID, EnvClientSecret)
		os.Exit(1)
	}
	client = scalekit.NewScalekitClient(environmentUrl, clientId, apiSecret)

	ctx := context.Background()
	orgResp, err := client.Organization().CreateOrganization(ctx, TestOrgName, scalekit.CreateOrganizationOptions{ExternalId: UniqueSuffix()})
	if err != nil || orgResp == nil || orgResp.GetOrganization() == nil {
		fmt.Fprintf(os.Stderr, "failed to create shared test org: %v\n", err)
		os.Exit(1)
	}
	testOrg = orgResp.GetOrganization().GetId()

	code := m.Run()
	client.Organization().DeleteOrganization(ctx, testOrg)
	os.Exit(code)
}

// createOrg creates an organization and returns its ID. Caller must defer DeleteTestOrganization(t, ctx, orgId).
func createOrg(t *testing.T, ctx context.Context, name, externalID string) string {
	t.Helper()
	resp, err := client.Organization().CreateOrganization(ctx, name, scalekit.CreateOrganizationOptions{ExternalId: externalID})
	if err != nil {
		t.Fatalf("createOrg: %v", err)
	}
	if resp == nil || resp.GetOrganization() == nil {
		t.Fatal("createOrg: nil response")
	}
	return resp.GetOrganization().GetId()
}

func toPtr(s string) *string {
	return &s
}

func toInt32Ptr(i int32) *int32 {
	return &i
}
