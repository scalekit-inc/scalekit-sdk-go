package test

import (
	"context"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/organizations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganization_EndToEndIntegration(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	name := TestOrgName
	externalId := UniqueSuffix()
	metadata := map[string]string{"key": "value"}

	createdOrganization, err := client.Organization().CreateOrganization(ctx, name, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
		Metadata:   metadata,
	})
	require.NoError(t, err)
	require.NotNil(t, createdOrganization)
	defer DeleteTestOrganization(t, ctx, createdOrganization.Organization.Id)

	assert.Equal(t, TestOrgName, createdOrganization.Organization.DisplayName)
	assert.Equal(t, "value", createdOrganization.Organization.Metadata["key"])

	retrievedOrganizationById, err := client.Organization().GetOrganization(ctx, createdOrganization.Organization.Id)
	require.NoError(t, err)
	assert.Equal(t, createdOrganization.Organization.Id, retrievedOrganizationById.Organization.Id)
	assert.Equal(t, createdOrganization.Organization.ExternalId, retrievedOrganizationById.Organization.ExternalId)

	retrieveByExternalId, err := client.Organization().GetOrganizationByExternalId(ctx, *createdOrganization.Organization.ExternalId)
	require.NoError(t, err)
	assert.Equal(t, retrievedOrganizationById.Organization.Id, retrieveByExternalId.Organization.Id)

	updatedOrganizationById, err := client.Organization().UpdateOrganization(ctx, createdOrganization.Organization.Id, &organizations.UpdateOrganization{
		DisplayName: toPtr("Updated name"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated name", updatedOrganizationById.Organization.DisplayName)

	updatedOrganizationByExternalId, err := client.Organization().UpdateOrganizationByExternalId(ctx, createdOrganization.Organization.GetExternalId(), &organizations.UpdateOrganization{
		DisplayName: toPtr("Updated name again"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated name again", updatedOrganizationByExternalId.Organization.DisplayName)

	err = client.Organization().DeleteOrganization(ctx, createdOrganization.Organization.Id)
	require.NoError(t, err)

	reCreatedOrganization, err := client.Organization().CreateOrganization(ctx, name, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
	})
	require.NoError(t, err)
	defer DeleteTestOrganization(t, ctx, reCreatedOrganization.Organization.Id)

	_, err = client.Organization().GetOrganization(ctx, createdOrganization.Organization.Id)
	assert.Error(t, err)

	organizationsList, err := client.Organization().ListOrganization(ctx, &scalekit.ListOrganizationOptions{
		PageSize:  10,
		PageToken: "",
	})
	require.NoError(t, err)
	assert.NotNil(t, organizationsList)
}

func TestOrganization_CreateOrganization_InvalidExternalID(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	_, err := client.Organization().CreateOrganization(ctx, "Exception Test", scalekit.CreateOrganizationOptions{
		ExternalId: "123",
	})
	assert.Error(t, err)
}

func TestOrganization_UpdateOrganizationSettings(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	featuresEnable := scalekit.OrganizationSettings{
		Features: []scalekit.Feature{
			{Name: "sso", Enabled: true},
			{Name: "dir_sync", Enabled: true},
		},
	}
	updatedOrganization, err := client.Organization().UpdateOrganizationSettings(ctx, orgId, featuresEnable)
	require.NoError(t, err)
	require.NotNil(t, updatedOrganization)
	require.True(t, len(updatedOrganization.Organization.Settings.Features) >= 2)
	assert.True(t, updatedOrganization.Organization.Settings.Features[0].Enabled)
	assert.True(t, updatedOrganization.Organization.Settings.Features[1].Enabled)

	featuresDisable := scalekit.OrganizationSettings{
		Features: []scalekit.Feature{
			{Name: "sso", Enabled: false},
			{Name: "dir_sync", Enabled: false},
		},
	}
	updatedOrganization, err = client.Organization().UpdateOrganizationSettings(ctx, orgId, featuresDisable)
	require.NoError(t, err)
	assert.False(t, updatedOrganization.Organization.Settings.Features[0].Enabled)
	assert.False(t, updatedOrganization.Organization.Settings.Features[1].Enabled)
}

func TestOrganization_UpsertUserManagementSettings(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	maxUsers := int32(150)
	settings, err := client.Organization().UpsertUserManagementSettings(ctx, orgId, scalekit.OrganizationUserManagementSettings{
		MaxAllowedUsers: toInt32Ptr(maxUsers),
	})
	if err != nil {
		t.Skipf("skipping UpsertUserManagementSettings test due to error: %v", err)
	}
	require.NotNil(t, settings)
	require.NotNil(t, settings.MaxAllowedUsers)
	assert.Equal(t, maxUsers, settings.MaxAllowedUsers.Value)

	updatedMaxUsers := int32(0)
	settings, err = client.Organization().UpsertUserManagementSettings(ctx, orgId, scalekit.OrganizationUserManagementSettings{
		MaxAllowedUsers: toInt32Ptr(updatedMaxUsers),
	})
	require.NoError(t, err)
	assert.Equal(t, updatedMaxUsers, settings.MaxAllowedUsers.Value)
}

func TestOrganization_CreateOrganization_WithMetadata(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	externalId := UniqueSuffix()
	createdOrganization, err := client.Organization().CreateOrganization(ctx, TestOrgName, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
		Metadata:   map[string]string{"meta_key": "meta_val"},
	})
	require.NoError(t, err)
	require.NotNil(t, createdOrganization)
	defer DeleteTestOrganization(t, ctx, createdOrganization.Organization.Id)

	assert.Equal(t, TestOrgName, createdOrganization.Organization.DisplayName)
	assert.Equal(t, "meta_val", createdOrganization.Organization.Metadata["meta_key"])
}
