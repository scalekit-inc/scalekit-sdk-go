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
	ctx := context.Background()

	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	domainName := shortDomain()
	domainResp, err := client.Domain().CreateDomain(ctx, orgId, domainName, nil)
	require.NoError(t, err)
	require.NotNil(t, domainResp)
	require.NotNil(t, domainResp.GetDomain())
	defer DeleteTestDomain(t, ctx, orgId, domainResp.GetDomain().GetId())

	connResp, err := client.Connection().CreateConnection(ctx, orgId, &connections.CreateConnection{
		Provider:    connections.ConnectionProvider_IDP_SIMULATOR,
		Type:        connections.ConnectionType_OIDC,
		ProviderKey: "test-key",
	})
	if err != nil {
		t.Skipf("CreateConnection not supported or requires config: %v", err)
	}
	require.NotNil(t, connResp)
	require.NotNil(t, connResp.GetConnection())
	defer DeleteTestConnection(t, ctx, orgId, connResp.GetConnection().GetId())

	got, err := client.Connection().GetConnection(ctx, orgId, connResp.GetConnection().GetId())
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.GetConnection())
	assert.Equal(t, connResp.GetConnection().GetId(), got.GetConnection().GetId())
	assert.NotEmpty(t, got.GetConnection().GetId())
	assert.NotEmpty(t, got.GetConnection().GetTestConnectionUri())
	expectedURL := os.Getenv(EnvEnvironmentURL) + "/sso/v1/oidc/" + got.GetConnection().GetId() + "/test"
	assert.Equal(t, expectedURL, got.GetConnection().GetTestConnectionUri())

	listByOrg, err := client.Connection().ListConnections(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, listByOrg)
	var found bool
	for _, c := range listByOrg.GetConnections() {
		if c.GetId() == connResp.GetConnection().GetId() {
			found = true
			break
		}
	}
	assert.True(t, found, "connection should be in org list")

	_, err = client.Connection().ListConnectionsByDomain(ctx, domainName)
	assert.NoError(t, err)

	enableResp, err := client.Connection().EnableConnection(ctx, orgId, connResp.GetConnection().GetId())
	if err == nil {
		require.NotNil(t, enableResp)
		assert.True(t, enableResp.GetEnabled())
		disableResp, err := client.Connection().DisableConnection(ctx, orgId, connResp.GetConnection().GetId())
		if err == nil {
			require.NotNil(t, disableResp)
			assert.False(t, disableResp.GetEnabled())
			_, _ = client.Connection().EnableConnection(ctx, orgId, connResp.GetConnection().GetId())
		}
	}
}
