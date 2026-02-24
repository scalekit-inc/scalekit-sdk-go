package test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
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
	client scalekit.Scalekit
)

func TestMain(m *testing.M) {
	loadEnvIfPresent("test/.env", ".env")

	environmentUrl := os.Getenv(EnvEnvironmentURL)
	clientId := os.Getenv(EnvClientID)
	apiSecret := os.Getenv(EnvClientSecret)
	if environmentUrl == "" || clientId == "" || apiSecret == "" {
		fmt.Fprintf(os.Stderr, "integration tests require %s, %s, %s\n", EnvEnvironmentURL, EnvClientID, EnvClientSecret)
		os.Exit(1)
	}
	client = scalekit.NewScalekitClient(environmentUrl, clientId, apiSecret)
	code := m.Run()
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

func loadEnvIfPresent(paths ...string) {
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}

			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			value = strings.Trim(value, `"'`)
			if key == "" {
				continue
			}

			// Keep explicitly provided environment values.
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
			_ = os.Setenv(key, value)
		}

		_ = file.Close()
	}
}
