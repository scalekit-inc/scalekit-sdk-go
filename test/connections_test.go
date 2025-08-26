package test

import (
	"context"
	"github.com/scalekit-inc/scalekit-sdk-go"
	"github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/connections"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnections(t *testing.T) {
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

func TestCreateDeleteConnection(t *testing.T) {

	// Test creating connection
	createdConnection, err := client.Connection().CreateConnection(context.Background(), testOrg, scalekit.CreateConnectionOptions{
		Provider: connections.ConnectionProvider_OKTA,
		Type:     connections.ConnectionType_OIDC,
	})
	assert.NoError(t, err)
	assert.NotNil(t, createdConnection)
	assert.Equal(t, createdConnection.Connection.Provider, connections.ConnectionProvider_OKTA)
	assert.Equal(t, createdConnection.Connection.Type, connections.ConnectionType_OIDC)

	getConnection, err := client.Connection().GetConnection(context.Background(), testOrg, createdConnection.Connection.Id)
	assert.NoError(t, err)
	assert.NotNil(t, getConnection)
	assert.Equal(t, getConnection.Connection.Provider, connections.ConnectionProvider_OKTA)
	assert.Equal(t, getConnection.Connection.Type, connections.ConnectionType_OIDC)

	// Test deleting connection
	err = client.Connection().DeleteConnection(context.Background(), testOrg, createdConnection.Connection.Id)
	assert.NoError(t, err)

	_, err = client.Connection().GetConnection(context.Background(), testOrg, createdConnection.Connection.Id)
	assert.Error(t, err)

}
