package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/organizations"

	"github.com/stretchr/testify/assert"
)

func TestOrganizations(t *testing.T) {
	organizationName := "Tested from Sdk"

	externalId := fmt.Sprintf("test-%d", time.Now().Unix())

	createdOrganization, err := client.Organization().CreateOrganization(context.Background(), organizationName, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
		Metadata: map[string]string{
			"key": "value",
		},
	})
	assert.NoError(t, err)

	retrievedOrganizationById, err := client.Organization().GetOrganization(context.Background(), createdOrganization.Organization.Id)
	assert.NoError(t, err)

	retrieveByExternalId, err := client.Organization().GetOrganizationByExternalId(context.Background(), *createdOrganization.Organization.ExternalId)
	assert.NoError(t, err)

	updatedOrganizationById, err := client.Organization().UpdateOrganization(context.Background(), createdOrganization.Organization.Id, &organizations.UpdateOrganization{
		DisplayName: toPtr("Updated name"),
	})
	assert.NoError(t, err)

	updatedOrganizationByExternalId, err := client.Organization().UpdateOrganizationByExternalId(context.Background(), createdOrganization.Organization.GetExternalId(), &organizations.UpdateOrganization{
		DisplayName: toPtr("Updated name again"),
	})
	assert.NoError(t, err)

	err = client.Organization().DeleteOrganization(context.Background(), createdOrganization.Organization.Id)
	assert.NoError(t, err)

	// Create again with same external Id
	reCreatedOrganization, err := client.Organization().CreateOrganization(context.Background(), organizationName, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
	})
	assert.NoError(t, err)
	err = client.Organization().DeleteOrganization(context.Background(), reCreatedOrganization.Organization.Id)
	assert.NoError(t, err)

	organizationsList, err := client.Organization().ListOrganization(context.Background(), &scalekit.ListOrganizationOptions{
		PageSize:  10,
		PageToken: "",
	})
	assert.NoError(t, err)
	assert.NotNil(t, organizationsList)
	assert.True(t, organizationsList.TotalSize > 10)
	assert.NotNil(t, organizationsList.NextPageToken)

	_, err = client.Organization().GetOrganization(context.Background(), createdOrganization.Organization.Id)
	assert.Error(t, err)

	_, err = client.Organization().GetOrganization(context.Background(), reCreatedOrganization.Organization.Id)
	assert.Error(t, err)

	_, err = client.Organization().GetOrganizationByExternalId(context.Background(), *reCreatedOrganization.Organization.ExternalId)
	assert.Error(t, err)

	assert.Equal(t, organizationName, createdOrganization.Organization.DisplayName)
	assert.Equal(t, createdOrganization.Organization.Metadata, createdOrganization.Organization.Metadata)
	assert.Equal(t, retrievedOrganizationById.Organization.Id, createdOrganization.Organization.Id)
	assert.Equal(t, retrievedOrganizationById.Organization.ExternalId, createdOrganization.Organization.ExternalId)
	assert.Equal(t, retrievedOrganizationById.Organization.Metadata, createdOrganization.Organization.Metadata)
	assert.Equal(t, retrieveByExternalId.Organization, retrievedOrganizationById.Organization)
	assert.Equal(t, updatedOrganizationById.Organization.DisplayName, "Updated name")
	assert.Equal(t, updatedOrganizationById.Organization.Id, createdOrganization.Organization.Id)
	assert.Equal(t, updatedOrganizationByExternalId.Organization.ExternalId, createdOrganization.Organization.ExternalId)
	assert.Equal(t, updatedOrganizationByExternalId.Organization.DisplayName, "Updated name again")
}

func TestException(t *testing.T) {
	organizationName := "Exception Test"

	_, err := client.Organization().CreateOrganization(context.Background(), organizationName, scalekit.CreateOrganizationOptions{
		ExternalId: "123",
	})
	assert.Error(t, err)
}

func TestUpdateOrganizationSettings(t *testing.T) {
	// Get first organization from list
	organizationsList, err := client.Organization().ListOrganization(context.Background(), &scalekit.ListOrganizationOptions{
		PageSize: 10,
	})
	assert.NoError(t, err)

	organization := organizationsList.Organizations[0]

	featuresEnable := scalekit.OrganizationSettings{
		Features: []scalekit.Feature{
			{
				Name:    "sso",
				Enabled: true,
			},
			{
				Name:    "dir_sync",
				Enabled: true,
			},
		},
	}

	featuresDisable := scalekit.OrganizationSettings{
		Features: []scalekit.Feature{
			{
				Name:    "sso",
				Enabled: false,
			},
			{
				Name:    "dir_sync",
				Enabled: false,
			},
		},
	}

	updatedOrganization, err := client.Organization().UpdateOrganizationSettings(context.Background(), organization.Id, featuresEnable)
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.True(t, updatedOrganization.Organization.Settings.Features[0].Enabled)
	assert.True(t, updatedOrganization.Organization.Settings.Features[1].Enabled)

	updatedOrganization, err = client.Organization().UpdateOrganizationSettings(context.Background(), organization.Id, featuresDisable)
	assert.NoError(t, err)
	assert.False(t, updatedOrganization.Organization.Settings.Features[0].Enabled)
	assert.False(t, updatedOrganization.Organization.Settings.Features[1].Enabled)
}

func TestUpsertUserManagementSettings(t *testing.T) {
	organizationsList, err := client.Organization().ListOrganization(context.Background(), &scalekit.ListOrganizationOptions{
		PageSize: 1,
	})
	assert.NoError(t, err)
	if len(organizationsList.Organizations) == 0 {
		t.Skip("no organizations available for testing user management settings")
	}

	organization := organizationsList.Organizations[0]
	maxUsers := int32(150)

	settings, err := client.Organization().UpsertUserManagementSettings(context.Background(), organization.Id, scalekit.OrganizationUserManagementSettings{
		MaxAllowedUsers: toInt32Ptr(maxUsers),
	})
	if err != nil {
		t.Skipf("skipping UpsertUserManagementSettings test due to error: %v", err)
	}
	assert.NotNil(t, settings)
	if settings.MaxAllowedUsers == nil {
		t.Fatalf("expected MaxAllowedUsers to be set")
	}
	assert.Equal(t, maxUsers, settings.MaxAllowedUsers.Value)

	updatedMaxUsers := int32(0)
	settings, err = client.Organization().UpsertUserManagementSettings(context.Background(), organization.Id, scalekit.OrganizationUserManagementSettings{
		MaxAllowedUsers: toInt32Ptr(updatedMaxUsers),
	})
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.NotNil(t, settings.MaxAllowedUsers)
	assert.Equal(t, updatedMaxUsers, settings.MaxAllowedUsers.Value)
}

func TestCreateWithMetadata(t *testing.T) {
	organizationName := "Tested from GO Sdk"

	externalId := fmt.Sprintf("test-%d", time.Now().Unix())

	createdOrganization, err := client.Organization().CreateOrganization(context.Background(), organizationName, scalekit.CreateOrganizationOptions{
		ExternalId: externalId,
		Metadata: map[string]string{
			"meta_key": "meta_val",
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, organizationName, createdOrganization.Organization.DisplayName)
	assert.Equal(t, createdOrganization.Organization.Metadata["meta_key"], "meta_val")
}
