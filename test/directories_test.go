package test

import (
	"context"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"

	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/directories"
	"github.com/stretchr/testify/assert"
)

func TestGetDirectory(t *testing.T) {
	// List directories
	directoriesResp, err := client.Directory().ListDirectories(context.Background(), testOrg)
	assert.NoError(t, err)
	assert.True(t, len(directoriesResp.Directories) > 0)

	// Get directory by ID
	directory, err := client.Directory().GetDirectory(context.Background(), testOrg, testDirectory)
	assert.NoError(t, err)

	firstDirectory := directoriesResp.Directories[0]
	assert.NotNil(t, firstDirectory)
	assert.Equal(t, testDirectory, firstDirectory.Id)
	assert.Equal(t, testOrg, firstDirectory.OrganizationId)
	assert.Equal(t, directories.DirectoryProvider_OKTA, firstDirectory.DirectoryProvider)
	assert.True(t, firstDirectory.Stats.TotalGroups > 0)
	assert.True(t, firstDirectory.Stats.TotalUsers > 0)
	assert.Equal(t, directory.Directory.Id, firstDirectory.Id)
	assert.Equal(t, directory.Directory.OrganizationId, firstDirectory.OrganizationId)
}

func TestListDirectoryUsers(t *testing.T) {
	// List users with options
	options := &scalekit.ListDirectoryUsersOptions{
		PageSize:      10,
		PageToken:     "",
		IncludeDetail: boolPtr(true),
	}

	usersResp, err := client.Directory().ListDirectoryUsers(context.Background(), testOrg, testDirectory, options)
	assert.NoError(t, err)
	assert.True(t, len(usersResp.Users) > 1)

	for _, user := range usersResp.Users {
		assert.NotNil(t, user)
		assert.NotEmpty(t, user.Id)
		assert.NotEmpty(t, user.Email)
		assert.NotEmpty(t, user.UserDetail)
	}
}

func TestListDirectoryUsersUpdatedAfter(t *testing.T) {
	updatedAfter := time.Unix(1729851960, 0)
	options := &scalekit.ListDirectoryUsersOptions{
		PageSize:      10,
		PageToken:     "",
		IncludeDetail: boolPtr(true),
		UpdatedAfter:  &updatedAfter,
	}

	usersResp, err := client.Directory().ListDirectoryUsers(context.Background(), testOrg, testDirectory, options)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(usersResp.Users))

	user := usersResp.Users[0]
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.Id)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.UserDetail)
}

func TestEnableDisableDirectory(t *testing.T) {
	// Enable directory
	enableResp, err := client.Directory().EnableDirectory(context.Background(), testOrg, testDirectory)
	assert.NoError(t, err)
	assert.True(t, enableResp.Enabled)

	// Disable directory
	disableResp, err := client.Directory().DisableDirectory(context.Background(), testOrg, testDirectory)
	assert.NoError(t, err)
	assert.False(t, disableResp.Enabled)

	// Cleanup: re-enable directory
	_, err = client.Directory().EnableDirectory(context.Background(), testOrg, testDirectory)
	assert.NoError(t, err)
}

func TestListDirectoryGroups(t *testing.T) {
	options := &scalekit.ListDirectoryGroupsOptions{
		PageSize:      10,
		PageToken:     "",
		IncludeDetail: boolPtr(true),
	}

	groupsResp, err := client.Directory().ListDirectoryGroups(context.Background(), testOrg, testDirectory, options)
	assert.NoError(t, err)
	assert.True(t, len(groupsResp.Groups) > 0)

	for _, group := range groupsResp.Groups {
		assert.NotNil(t, group)
		assert.NotEmpty(t, group.Id)
		assert.NotEmpty(t, group.DisplayName)
		assert.NotEmpty(t, group.GroupDetail)
	}
}

func TestGetPrimaryDirectoryByOrganizationId(t *testing.T) {
	directory, err := client.Directory().GetPrimaryDirectoryByOrganizationId(context.Background(), testOrg)
	assert.NoError(t, err)

	directoryById, err := client.Directory().GetDirectory(context.Background(), testOrg, directory.Directory.Id)
	assert.NoError(t, err)

	assert.NotNil(t, directory)
	assert.NotNil(t, directoryById)
	assert.Equal(t, directory.Directory.Id, directoryById.Directory.Id)
	assert.Equal(t, testOrg, directory.Directory.OrganizationId)
	assert.Equal(t, directories.DirectoryProvider_OKTA, directory.Directory.DirectoryProvider)
	assert.True(t, directory.Directory.Stats.TotalGroups > 0)
	assert.True(t, directory.Directory.Stats.TotalUsers > 0)
}

func boolPtr(b bool) *bool {
	return &b
}
