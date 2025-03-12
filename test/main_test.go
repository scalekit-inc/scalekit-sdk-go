package test

import (
	"github.com/scalekit-inc/scalekit-sdk-go"
	"os"
	"testing"
)

var (
	client         scalekit.Scalekit
	domain         string
	testOrg        string
	testConnection string
	testDirectory  string
)

func TestMain(m *testing.M) {
	// Init client
	environmentUrl := os.Getenv("SCALEKIT_ENVIRONMENT_URL")
	clientId := os.Getenv("SCALEKIT_CLIENT_ID")
	apiSecret := os.Getenv("SCALEKIT_CLIENT_SECRET")
	domain = os.Getenv("TEST_DOMAIN")
	testOrg = os.Getenv("TEST_ORGANIZATION")
	testConnection = os.Getenv("TEST_CONNECTION")
	testDirectory = os.Getenv("TEST_DIRECTORY")

	client = scalekit.NewScalekitClient(environmentUrl, clientId, apiSecret)

	code := m.Run()
	os.Exit(code)
}

func toPtr(s string) *string {
	return &s
}
