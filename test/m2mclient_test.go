package test

import (
	"context"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrganizationClient(t *testing.T) {
	ctx := context.Background()

	resp, err := client.M2M().CreateOrganizationClient(ctx, testOrg, scalekit.CreateOrganizationClientOptions{
		Name:        "Test M2M Client",
		Description: "Integration test client",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Client)

	clientId := resp.Client.ClientId
	t.Cleanup(func() {
		_ = client.M2M().DeleteOrganizationClient(ctx, testOrg, clientId)
	})

	assert.NotEmpty(t, clientId)
	assert.Equal(t, "Test M2M Client", resp.Client.Name)
}

func TestCreateOrganizationClientRequiresOrgId(t *testing.T) {
	ctx := context.Background()

	_, err := client.M2M().CreateOrganizationClient(ctx, "", scalekit.CreateOrganizationClientOptions{
		Name: "Test",
	})
	require.Error(t, err)
}

func TestGetOrganizationClient(t *testing.T) {
	ctx := context.Background()

	created, err := client.M2M().CreateOrganizationClient(ctx, testOrg, scalekit.CreateOrganizationClientOptions{
		Name: "Get Test Client",
	})
	require.NoError(t, err)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.M2M().DeleteOrganizationClient(ctx, testOrg, clientId)
	})

	fetched, err := client.M2M().GetOrganizationClient(ctx, testOrg, clientId)
	require.NoError(t, err)
	require.NotNil(t, fetched)
	assert.Equal(t, clientId, fetched.Client.ClientId)
	assert.Equal(t, "Get Test Client", fetched.Client.Name)
}

func TestGetOrganizationClientRequiresClientId(t *testing.T) {
	ctx := context.Background()

	_, err := client.M2M().GetOrganizationClient(ctx, testOrg, "")
	require.Error(t, err)
}

func TestUpdateOrganizationClient(t *testing.T) {
	ctx := context.Background()

	created, err := client.M2M().CreateOrganizationClient(ctx, testOrg, scalekit.CreateOrganizationClientOptions{
		Name: "Original Name",
	})
	require.NoError(t, err)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.M2M().DeleteOrganizationClient(ctx, testOrg, clientId)
	})

	updated, err := client.M2M().UpdateOrganizationClient(ctx, testOrg, clientId, scalekit.UpdateOrganizationClientOptions{
		Name:        "Updated Name",
		Description: "Updated description",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Client.Name)
	assert.Equal(t, "Updated description", updated.Client.Description)
}

func TestCreateOrganizationClientSecret(t *testing.T) {
	ctx := context.Background()

	created, err := client.M2M().CreateOrganizationClient(ctx, testOrg, scalekit.CreateOrganizationClientOptions{
		Name: "Secret Test Client",
	})
	require.NoError(t, err)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.M2M().DeleteOrganizationClient(ctx, testOrg, clientId)
	})

	secretResp, err := client.M2M().CreateOrganizationClientSecret(ctx, testOrg, clientId)
	require.NoError(t, err)
	require.NotNil(t, secretResp)
	assert.NotEmpty(t, secretResp.SecretId)
	assert.NotEmpty(t, secretResp.Secret)

	t.Cleanup(func() {
		_ = client.M2M().DeleteOrganizationClientSecret(ctx, testOrg, clientId, secretResp.SecretId)
	})
}

func TestDeleteOrganizationClientSecret(t *testing.T) {
	ctx := context.Background()

	created, err := client.M2M().CreateOrganizationClient(ctx, testOrg, scalekit.CreateOrganizationClientOptions{
		Name: "Delete Secret Client",
	})
	require.NoError(t, err)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.M2M().DeleteOrganizationClient(ctx, testOrg, clientId)
	})

	secretResp, err := client.M2M().CreateOrganizationClientSecret(ctx, testOrg, clientId)
	require.NoError(t, err)

	err = client.M2M().DeleteOrganizationClientSecret(ctx, testOrg, clientId, secretResp.SecretId)
	require.NoError(t, err)
}

func TestDeleteOrganizationClientSecretRequiresSecretId(t *testing.T) {
	ctx := context.Background()

	err := client.M2M().DeleteOrganizationClientSecret(ctx, testOrg, "skc_dummy", "")
	require.Error(t, err)
}

func TestListOrganizationClients(t *testing.T) {
	ctx := context.Background()

	created, err := client.M2M().CreateOrganizationClient(ctx, testOrg, scalekit.CreateOrganizationClientOptions{
		Name: "List Test Client",
	})
	require.NoError(t, err)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.M2M().DeleteOrganizationClient(ctx, testOrg, clientId)
	})

	list, err := client.M2M().ListOrganizationClients(ctx, testOrg, scalekit.ListOrganizationClientsOptions{
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, list)
	assert.NotEmpty(t, list.Clients)

	found := false
	for _, c := range list.Clients {
		if c.ClientId == clientId {
			found = true
			break
		}
	}
	assert.True(t, found, "created client should appear in list")
}

func TestListOrganizationClientsRequiresOrgId(t *testing.T) {
	ctx := context.Background()

	_, err := client.M2M().ListOrganizationClients(ctx, "", scalekit.ListOrganizationClientsOptions{})
	require.Error(t, err)
}

func TestDeleteOrganizationClient(t *testing.T) {
	ctx := context.Background()

	created, err := client.M2M().CreateOrganizationClient(ctx, testOrg, scalekit.CreateOrganizationClientOptions{
		Name: "To Delete Client",
	})
	require.NoError(t, err)
	clientId := created.Client.ClientId

	err = client.M2M().DeleteOrganizationClient(ctx, testOrg, clientId)
	require.NoError(t, err)
}
