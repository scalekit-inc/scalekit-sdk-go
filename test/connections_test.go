package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/connections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// shortDomain returns a domain name of at most 32 characters for connection tests (API limit).
func shortDomain() string {
	return fmt.Sprintf("c%06d.acme.com", time.Now().UnixNano()%1000000)
}

func TestConnection_EndToEndIntegration(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()

	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	domainName := shortDomain()
	domainResp, err := client.Domain().CreateDomain(ctx, orgId, domainName, nil)
	require.NoError(t, err)
	require.NotNil(t, domainResp)
	defer DeleteTestDomain(t, ctx, orgId, domainResp.Domain.Id)

	connResp, err := client.Connection().CreateConnection(ctx, orgId, &connections.CreateConnection{
		Provider:    connections.ConnectionProvider_IDP_SIMULATOR,
		Type:        connections.ConnectionType_OIDC,
		ProviderKey: "test-key",
	})
	if err != nil {
		t.Skipf("CreateConnection not supported or requires config: %v", err)
	}
	require.NotNil(t, connResp)
	require.NotNil(t, connResp.Connection)
	defer DeleteTestConnection(t, ctx, orgId, connResp.Connection.Id)

	got, err := client.Connection().GetConnection(ctx, orgId, connResp.Connection.Id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, connResp.Connection.Id, got.Connection.Id)
	assert.NotEmpty(t, got.Connection.Id)
	assert.NotEmpty(t, got.Connection.TestConnectionUri)
	expectedURL := os.Getenv("SCALEKIT_ENVIRONMENT_URL") + "/sso/v1/oidc/" + got.Connection.Id + "/test"
	assert.Equal(t, expectedURL, got.Connection.TestConnectionUri)

	listByOrg, err := client.Connection().ListConnections(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, listByOrg)
	var found bool
	for _, c := range listByOrg.Connections {
		if c.Id == connResp.Connection.Id {
			found = true
			break
		}
	}
	assert.True(t, found, "connection should be in org list")

	_, err = client.Connection().ListConnectionsByDomain(ctx, domainName)
	assert.NoError(t, err)

	enableResp, err := client.Connection().EnableConnection(ctx, orgId, connResp.Connection.Id)
	if err == nil {
		assert.True(t, enableResp.Enabled)
		disableResp, err := client.Connection().DisableConnection(ctx, orgId, connResp.Connection.Id)
		if err == nil {
			assert.False(t, disableResp.Enabled)
			_, _ = client.Connection().EnableConnection(ctx, orgId, connResp.Connection.Id)
		}
	}
}
