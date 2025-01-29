package test

import (
	"github.com/scalekit-inc/scalekit-sdk-go"
	"os"
	"testing"
)

var client scalekit.Scalekit

func TestMain(m *testing.M) {
	environmentUrl := os.Getenv("SCALEKIT_ENVIRONMENT_URL")
	clientId := os.Getenv("SCALEKIT_CLIENT_ID")
	apiSecret := os.Getenv("SCALEKIT_CLIENT_SECRET")

	client = scalekit.NewScalekitClient(environmentUrl, clientId, apiSecret)

	code := m.Run()
	os.Exit(code)
}

func toPtr(s string) *string {
	return &s
}
