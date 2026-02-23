package test

import (
	"context"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/directories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestDirectory_EndToEndIntegration(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	createResp, err := client.Directory().CreateDirectory(ctx, orgId, &directories.CreateDirectory{
		DirectoryType:     directories.DirectoryType_SCIM,
		DirectoryProvider: directories.DirectoryProvider_OKTA,
	})
	if err != nil {
		t.Skipf("CreateDirectory not supported or requires config: %v", err)
	}
	require.NotNil(t, createResp)
	require.NotNil(t, createResp.Directory)
	defer DeleteTestDirectory(t, ctx, orgId, createResp.Directory.Id)

	got, err := client.Directory().GetDirectory(ctx, orgId, createResp.Directory.Id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, createResp.Directory.Id, got.Directory.Id)
	assert.Equal(t, orgId, got.Directory.OrganizationId)
	assert.NotEmpty(t, got.Directory.Id)
	assert.NotNil(t, got.Directory.Stats)

	listResp, err := client.Directory().ListDirectories(ctx, orgId)
	require.NoError(t, err)
	require.True(t, len(listResp.Directories) > 0)
	var found bool
	for _, d := range listResp.Directories {
		if d.Id == createResp.Directory.Id {
			found = true
			break
		}
	}
	assert.True(t, found, "directory should be in list")

	usersResp, err := client.Directory().ListDirectoryUsers(ctx, orgId, createResp.Directory.Id, &scalekit.ListDirectoryUsersOptions{
		PageSize:      10,
		PageToken:     "",
		IncludeDetail: boolPtr(true),
	})
	assert.NoError(t, err)
	assert.NotNil(t, usersResp)

	groupsResp, err := client.Directory().ListDirectoryGroups(ctx, orgId, createResp.Directory.Id, &scalekit.ListDirectoryGroupsOptions{
		PageSize:      10,
		PageToken:     "",
		IncludeDetail: boolPtr(true),
	})
	assert.NoError(t, err)
	assert.NotNil(t, groupsResp)

	enableResp, err := client.Directory().EnableDirectory(ctx, orgId, createResp.Directory.Id)
	if err == nil {
		assert.True(t, enableResp.Enabled)
		disableResp, err := client.Directory().DisableDirectory(ctx, orgId, createResp.Directory.Id)
		if err == nil {
			assert.False(t, disableResp.Enabled)
			_, _ = client.Directory().EnableDirectory(ctx, orgId, createResp.Directory.Id)
		}
	}
}

func TestDirectory_GetPrimaryDirectoryByOrganizationId(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	createResp, err := client.Directory().CreateDirectory(ctx, orgId, &directories.CreateDirectory{
		DirectoryType:     directories.DirectoryType_SCIM,
		DirectoryProvider: directories.DirectoryProvider_OKTA,
	})
	if err != nil {
		t.Skipf("CreateDirectory not supported: %v", err)
	}
	require.NotNil(t, createResp)
	defer DeleteTestDirectory(t, ctx, orgId, createResp.Directory.Id)

	primary, err := client.Directory().GetPrimaryDirectoryByOrganizationId(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, primary)
	require.NotNil(t, primary.Directory)
	assert.Equal(t, orgId, primary.Directory.OrganizationId)

	byId, err := client.Directory().GetDirectory(ctx, orgId, primary.Directory.Id)
	require.NoError(t, err)
	assert.Equal(t, primary.Directory.Id, byId.Directory.Id)
}

func TestDirectory_ListDirectoryUsers_UpdatedAfter(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	createResp, err := client.Directory().CreateDirectory(ctx, orgId, &directories.CreateDirectory{
		DirectoryType:     directories.DirectoryType_SCIM,
		DirectoryProvider: directories.DirectoryProvider_OKTA,
	})
	if err != nil {
		t.Skipf("CreateDirectory not supported: %v", err)
	}
	require.NotNil(t, createResp)
	defer DeleteTestDirectory(t, ctx, orgId, createResp.Directory.Id)

	updatedAfter := time.Unix(1729851960, 0)
	usersResp, err := client.Directory().ListDirectoryUsers(ctx, orgId, createResp.Directory.Id, &scalekit.ListDirectoryUsersOptions{
		PageSize:      10,
		PageToken:     "",
		IncludeDetail: boolPtr(true),
		UpdatedAfter:  &updatedAfter,
	})
	require.NoError(t, err)
	require.NotNil(t, usersResp)
	for _, user := range usersResp.Users {
		assert.NotEmpty(t, user.Id)
		assert.NotEmpty(t, user.Email)
	}
}
