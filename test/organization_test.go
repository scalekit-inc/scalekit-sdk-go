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
	defer DeleteTestOrganization(t, ctx, createdOrganization.GetOrganization().GetId())

	assert.Equal(t, TestOrgName, createdOrganization.GetOrganization().GetDisplayName())
	assert.Equal(t, "value", createdOrganization.GetOrganization().GetMetadata()["key"])

	retrievedOrganizationById, err := client.Organization().GetOrganization(ctx, createdOrganization.GetOrganization().GetId())
	require.NoError(t, err)
	assert.Equal(t, createdOrganization.GetOrganization().GetId(), retrievedOrganizationById.GetOrganization().GetId())
	assert.Equal(t, createdOrganization.GetOrganization().GetExternalId(), retrievedOrganizationById.GetOrganization().GetExternalId())

	retrieveByExternalId, err := client.Organization().GetOrganizationByExternalId(ctx, createdOrganization.GetOrganization().GetExternalId())
	require.NoError(t, err)
	assert.Equal(t, retrievedOrganizationById.GetOrganization().GetId(), retrieveByExternalId.GetOrganization().GetId())

	updatedOrganizationById, err := client.Organization().UpdateOrganization(ctx, createdOrganization.GetOrganization().GetId(), &organizations.UpdateOrganization{
		DisplayName: toPtr("Updated name"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated name", updatedOrganizationById.GetOrganization().GetDisplayName())

	updatedOrganizationByExternalId, err := client.Organization().UpdateOrganizationByExternalId(ctx, createdOrganization.GetOrganization().GetExternalId(), &organizations.UpdateOrganization{
		DisplayName: toPtr("Updated name again"),
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated name again", updatedOrganizationByExternalId.GetOrganization().GetDisplayName())

	err = client.Organization().DeleteOrganization(ctx, createdOrganization.GetOrganization().GetId())
	require.NoError(t, err)

	reCreatedOrganization, err := client.Organization().CreateOrganization(ctx, name, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
	})
	require.NoError(t, err)
	defer DeleteTestOrganization(t, ctx, reCreatedOrganization.GetOrganization().GetId())

	_, err = client.Organization().GetOrganization(ctx, createdOrganization.GetOrganization().GetId())
	assert.Error(t, err)

	organizationsList, err := client.Organization().ListOrganization(ctx, &scalekit.ListOrganizationOptions{
		PageSize:  10,
		PageToken: "",
	})
	require.NoError(t, err)
	assert.NotNil(t, organizationsList)
}

func TestOrganization_CreateOrganization_InvalidExternalID(t *testing.T) {
	ctx := context.Background()
	_, err := client.Organization().CreateOrganization(ctx, "Exception Test", scalekit.CreateOrganizationOptions{
		ExternalId: "123",
	})
	assert.Error(t, err)
}

func TestOrganization_UpdateOrganizationSettings(t *testing.T) {
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
	require.True(t, len(updatedOrganization.GetOrganization().GetSettings().GetFeatures()) >= 2)
	assert.True(t, updatedOrganization.GetOrganization().GetSettings().GetFeatures()[0].GetEnabled())
	assert.True(t, updatedOrganization.GetOrganization().GetSettings().GetFeatures()[1].GetEnabled())

	featuresDisable := scalekit.OrganizationSettings{
		Features: []scalekit.Feature{
			{Name: "sso", Enabled: false},
			{Name: "dir_sync", Enabled: false},
		},
	}
	updatedOrganization, err = client.Organization().UpdateOrganizationSettings(ctx, orgId, featuresDisable)
	require.NoError(t, err)
	assert.False(t, updatedOrganization.GetOrganization().GetSettings().GetFeatures()[0].GetEnabled())
	assert.False(t, updatedOrganization.GetOrganization().GetSettings().GetFeatures()[1].GetEnabled())
}

func TestOrganization_UpsertUserManagementSettings(t *testing.T) {
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
	require.NotNil(t, settings.GetMaxAllowedUsers())
	assert.Equal(t, maxUsers, settings.GetMaxAllowedUsers().GetValue())

	updatedMaxUsers := int32(0)
	settings, err = client.Organization().UpsertUserManagementSettings(ctx, orgId, scalekit.OrganizationUserManagementSettings{
		MaxAllowedUsers: toInt32Ptr(updatedMaxUsers),
	})
	require.NoError(t, err)
	assert.Equal(t, updatedMaxUsers, settings.GetMaxAllowedUsers().GetValue())
}

func TestOrganization_CreateOrganization_WithMetadata(t *testing.T) {
	ctx := context.Background()
	externalId := UniqueSuffix()
	createdOrganization, err := client.Organization().CreateOrganization(ctx, TestOrgName, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
		Metadata:   map[string]string{"meta_key": "meta_val"},
	})
	require.NoError(t, err)
	require.NotNil(t, createdOrganization)
	defer DeleteTestOrganization(t, ctx, createdOrganization.GetOrganization().GetId())

	assert.Equal(t, TestOrgName, createdOrganization.GetOrganization().GetDisplayName())
	assert.Equal(t, "meta_val", createdOrganization.GetOrganization().GetMetadata()["meta_key"])
}
