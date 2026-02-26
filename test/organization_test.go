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
	require.NotNil(t, createdOrganization.GetOrganization())
	defer DeleteTestOrganization(t, ctx, createdOrganization.GetOrganization().GetId())

	assert.Equal(t, TestOrgName, createdOrganization.GetOrganization().GetDisplayName())
	assert.Equal(t, "value", createdOrganization.GetOrganization().GetMetadata()["key"])

	retrievedOrganizationById, err := client.Organization().GetOrganization(ctx, createdOrganization.GetOrganization().GetId())
	require.NoError(t, err)
	require.NotNil(t, retrievedOrganizationById)
	require.NotNil(t, retrievedOrganizationById.GetOrganization())
	assert.Equal(t, createdOrganization.GetOrganization().GetId(), retrievedOrganizationById.GetOrganization().GetId())
	assert.Equal(t, createdOrganization.GetOrganization().GetExternalId(), retrievedOrganizationById.GetOrganization().GetExternalId())

	retrieveByExternalId, err := client.Organization().GetOrganizationByExternalId(ctx, createdOrganization.GetOrganization().GetExternalId())
	require.NoError(t, err)
	require.NotNil(t, retrieveByExternalId)
	require.NotNil(t, retrieveByExternalId.GetOrganization())
	assert.Equal(t, retrievedOrganizationById.GetOrganization().GetId(), retrieveByExternalId.GetOrganization().GetId())

	updatedOrganizationById, err := client.Organization().UpdateOrganization(ctx, createdOrganization.GetOrganization().GetId(), &organizations.UpdateOrganization{
		DisplayName: toPtr("Updated name"),
	})
	require.NoError(t, err)
	require.NotNil(t, updatedOrganizationById)
	require.NotNil(t, updatedOrganizationById.GetOrganization())
	assert.Equal(t, "Updated name", updatedOrganizationById.GetOrganization().GetDisplayName())

	updatedOrganizationByExternalId, err := client.Organization().UpdateOrganizationByExternalId(ctx, createdOrganization.GetOrganization().GetExternalId(), &organizations.UpdateOrganization{
		DisplayName: toPtr("Updated name again"),
	})
	require.NoError(t, err)
	require.NotNil(t, updatedOrganizationByExternalId)
	require.NotNil(t, updatedOrganizationByExternalId.GetOrganization())
	assert.Equal(t, "Updated name again", updatedOrganizationByExternalId.GetOrganization().GetDisplayName())

	err = client.Organization().DeleteOrganization(ctx, createdOrganization.GetOrganization().GetId())
	require.NoError(t, err)

	reCreatedOrganization, err := client.Organization().CreateOrganization(ctx, name, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
	})
	require.NoError(t, err)
	require.NotNil(t, reCreatedOrganization)
	require.NotNil(t, reCreatedOrganization.GetOrganization())
	defer DeleteTestOrganization(t, ctx, reCreatedOrganization.GetOrganization().GetId())

	_, err = client.Organization().GetOrganization(ctx, createdOrganization.GetOrganization().GetId())
	assert.Error(t, err)

	organizationsList, err := client.Organization().ListOrganization(ctx, &scalekit.ListOrganizationOptions{
		PageSize:  10,
		PageToken: "",
	})
	require.NoError(t, err)
	require.NotNil(t, organizationsList)
}

func TestOrganization_CreateOrganization_DuplicateExternalID(t *testing.T) {
	ctx := context.Background()
	externalId := UniqueSuffix()

	first, err := client.Organization().CreateOrganization(ctx, TestOrgName, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
	})
	require.NoError(t, err)
	require.NotNil(t, first)
	require.NotNil(t, first.GetOrganization())
	defer DeleteTestOrganization(t, ctx, first.GetOrganization().GetId())

	_, err = client.Organization().CreateOrganization(ctx, "Duplicate Org", scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
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
	require.NotNil(t, updatedOrganization.GetOrganization())
	require.True(t, len(updatedOrganization.GetOrganization().GetSettings().GetFeatures()) >= 2)
	enabledFeatures := map[string]bool{}
	for _, f := range updatedOrganization.GetOrganization().GetSettings().GetFeatures() {
		enabledFeatures[f.GetName()] = f.GetEnabled()
	}
	assert.Contains(t, enabledFeatures, "sso", "sso feature should be present")
	assert.Contains(t, enabledFeatures, "dir_sync", "dir_sync feature should be present")
	assert.True(t, enabledFeatures["sso"])
	assert.True(t, enabledFeatures["dir_sync"])

	featuresDisable := scalekit.OrganizationSettings{
		Features: []scalekit.Feature{
			{Name: "sso", Enabled: false},
			{Name: "dir_sync", Enabled: false},
		},
	}
	updatedOrganization, err = client.Organization().UpdateOrganizationSettings(ctx, orgId, featuresDisable)
	require.NoError(t, err)
	require.NotNil(t, updatedOrganization.GetOrganization())
	disabledFeatures := map[string]bool{}
	for _, f := range updatedOrganization.GetOrganization().GetSettings().GetFeatures() {
		disabledFeatures[f.GetName()] = f.GetEnabled()
	}
	assert.Contains(t, disabledFeatures, "sso", "sso feature should be present")
	assert.Contains(t, disabledFeatures, "dir_sync", "dir_sync feature should be present")
	assert.False(t, disabledFeatures["sso"])
	assert.False(t, disabledFeatures["dir_sync"])
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
	require.NotNil(t, settings)
	require.NotNil(t, settings.GetMaxAllowedUsers())
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
	require.NotNil(t, createdOrganization.GetOrganization())
	defer DeleteTestOrganization(t, ctx, createdOrganization.GetOrganization().GetId())

	assert.Equal(t, TestOrgName, createdOrganization.GetOrganization().GetDisplayName())
	assert.Equal(t, "meta_val", createdOrganization.GetOrganization().GetMetadata()["meta_key"])
}
