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
	require.NotNil(t, createResp.GetDirectory())
	defer DeleteTestDirectory(t, ctx, orgId, createResp.GetDirectory().GetId())

	got, err := client.Directory().GetDirectory(ctx, orgId, createResp.GetDirectory().GetId())
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.GetDirectory())
	assert.Equal(t, createResp.GetDirectory().GetId(), got.GetDirectory().GetId())
	assert.Equal(t, orgId, got.GetDirectory().GetOrganizationId())
	assert.NotEmpty(t, got.GetDirectory().GetId())
	require.NotNil(t, got.GetDirectory().GetStats())

	listResp, err := client.Directory().ListDirectories(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.True(t, len(listResp.GetDirectories()) > 0)
	var found bool
	for _, d := range listResp.GetDirectories() {
		if d.GetId() == createResp.GetDirectory().GetId() {
			found = true
			break
		}
	}
	assert.True(t, found, "directory should be in list")

	usersResp, err := client.Directory().ListDirectoryUsers(ctx, orgId, createResp.GetDirectory().GetId(), &scalekit.ListDirectoryUsersOptions{
		PageSize:      10,
		PageToken:     "",
		IncludeDetail: boolPtr(true),
	})
	require.NoError(t, err)
	require.NotNil(t, usersResp)

	groupsResp, err := client.Directory().ListDirectoryGroups(ctx, orgId, createResp.GetDirectory().GetId(), &scalekit.ListDirectoryGroupsOptions{
		PageSize:      10,
		PageToken:     "",
		IncludeDetail: boolPtr(true),
	})
	require.NoError(t, err)
	require.NotNil(t, groupsResp)

	enableResp, err := client.Directory().EnableDirectory(ctx, orgId, createResp.GetDirectory().GetId())
	if err == nil {
		require.NotNil(t, enableResp)
		assert.True(t, enableResp.GetEnabled())
		disableResp, err := client.Directory().DisableDirectory(ctx, orgId, createResp.GetDirectory().GetId())
		if err == nil {
			require.NotNil(t, disableResp)
			assert.False(t, disableResp.GetEnabled())
			_, _ = client.Directory().EnableDirectory(ctx, orgId, createResp.GetDirectory().GetId())
		}
	}
}

func TestDirectory_GetPrimaryDirectoryByOrganizationId(t *testing.T) {
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
	require.NotNil(t, createResp.GetDirectory())
	defer DeleteTestDirectory(t, ctx, orgId, createResp.GetDirectory().GetId())

	primary, err := client.Directory().GetPrimaryDirectoryByOrganizationId(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, primary)
	require.NotNil(t, primary.GetDirectory())
	assert.Equal(t, orgId, primary.GetDirectory().GetOrganizationId())

	byId, err := client.Directory().GetDirectory(ctx, orgId, primary.GetDirectory().GetId())
	require.NoError(t, err)
	require.NotNil(t, byId)
	require.NotNil(t, byId.GetDirectory())
	assert.Equal(t, primary.GetDirectory().GetId(), byId.GetDirectory().GetId())
}

func TestDirectory_ListDirectoryUsers_UpdatedAfter(t *testing.T) {
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
	require.NotNil(t, createResp.GetDirectory())
	defer DeleteTestDirectory(t, ctx, orgId, createResp.GetDirectory().GetId())

	updatedAfter := time.Unix(1729851960, 0)
	usersResp, err := client.Directory().ListDirectoryUsers(ctx, orgId, createResp.GetDirectory().GetId(), &scalekit.ListDirectoryUsersOptions{
		PageSize:      10,
		PageToken:     "",
		IncludeDetail: boolPtr(true),
		UpdatedAfter:  &updatedAfter,
	})
	require.NoError(t, err)
	require.NotNil(t, usersResp)
	for _, user := range usersResp.GetUsers() {
		assert.NotEmpty(t, user.GetId())
		assert.NotEmpty(t, user.GetEmail())
	}
}
