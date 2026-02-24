package test

import (
	"context"
	"os"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/connections"

	"github.com/stretchr/testify/assert"
)

func TestConnections(t *testing.T) {
	if domain == "" {
		t.Skip("TEST_DOMAIN is not set")
	}
	if testConnection == "" {
		t.Skip("TEST_CONNECTION is not set")
	}

	// Test listing connections by domain
	connectionsByDomain, err := client.Connection().ListConnectionsByDomain(context.Background(), domain)
	assert.NoError(t, err)
	assert.True(t, len(connectionsByDomain.Connections) > 0)

	// Test listing connections by organization
	connectionsByOrg, err := client.Connection().ListConnections(context.Background(), testOrg)
	assert.NoError(t, err)
	assert.True(t, len(connectionsByOrg.Connections) > 0)

	// Test getting connection by ID
	connection, err := client.Connection().GetConnection(context.Background(), testOrg, testConnection)
	assert.NoError(t, err)
	assert.Equal(t, testConnection, connection.Connection.Id)
	assert.Equal(t, connections.ConnectionProvider_OKTA, connection.Connection.Provider)

	expectedConnectionURL := os.Getenv("SCALEKIT_ENVIRONMENT_URL") + "/sso/v1/oidc/" + connection.Connection.Id + "/test"
	assert.Equal(t, expectedConnectionURL, connection.Connection.TestConnectionUri)
}
