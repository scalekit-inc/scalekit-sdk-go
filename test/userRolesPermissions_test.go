package test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// envOrSkip skips the test if the given environment variable is not set,
// and returns its value otherwise.
func envOrSkip(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("skipping: %s not set", key)
	}
	return v
}

func TestUser_ListUserRoles(t *testing.T) {
	ctx := context.Background()

	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	t.Cleanup(func() { DeleteTestOrganization(t, ctx, orgId) })

	uniqueEmail := fmt.Sprintf("list.roles.test.%d@example.com", time.Now().UnixNano()/1e6)
	createdUser, err := client.User().CreateUserAndMembership(ctx, orgId, &users.CreateUser{
		Email:    uniqueEmail,
		Metadata: map[string]string{"source": "list_roles_test"},
	}, false)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotEmpty(t, createdUser.GetUser().GetId())
	userId := createdUser.GetUser().GetId()
	t.Cleanup(func() { _ = client.User().DeleteUser(ctx, userId) })

	resp, err := client.User().ListUserRoles(ctx, orgId, userId)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Roles may be empty for a newly created user; just confirm the call succeeds
	_ = resp.GetRoles()
}

func TestUser_ListUserRoles_WithEnvUserId(t *testing.T) {
	ctx := context.Background()

	orgId := envOrSkip(t, "TEST_ORGANIZATION_ID")
	userId := envOrSkip(t, "TEST_USER_ID")

	resp, err := client.User().ListUserRoles(ctx, orgId, userId)
	require.NoError(t, err)
	require.NotNil(t, resp)
	t.Logf("ListUserRoles: %d role(s) returned", len(resp.GetRoles()))
	for _, role := range resp.GetRoles() {
		assert.NotEmpty(t, role.GetName(), "role name should not be empty")
	}
}

func TestUser_ListUserPermissions(t *testing.T) {
	ctx := context.Background()

	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	t.Cleanup(func() { DeleteTestOrganization(t, ctx, orgId) })

	uniqueEmail := fmt.Sprintf("list.perms.test.%d@example.com", time.Now().UnixNano()/1e6)
	createdUser, err := client.User().CreateUserAndMembership(ctx, orgId, &users.CreateUser{
		Email:    uniqueEmail,
		Metadata: map[string]string{"source": "list_perms_test"},
	}, false)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotEmpty(t, createdUser.GetUser().GetId())
	userId := createdUser.GetUser().GetId()
	t.Cleanup(func() { _ = client.User().DeleteUser(ctx, userId) })

	resp, err := client.User().ListUserPermissions(ctx, orgId, userId)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// Permissions may be empty for a newly created user; just confirm the call succeeds
	_ = resp.GetPermissions()
}

func TestUser_ListUserPermissions_WithEnvUserId(t *testing.T) {
	ctx := context.Background()

	orgId := envOrSkip(t, "TEST_ORGANIZATION_ID")
	userId := envOrSkip(t, "TEST_USER_ID")

	resp, err := client.User().ListUserPermissions(ctx, orgId, userId)
	require.NoError(t, err)
	require.NotNil(t, resp)
	t.Logf("ListUserPermissions: %d permission(s) returned", len(resp.GetPermissions()))
	for _, perm := range resp.GetPermissions() {
		assert.NotEmpty(t, perm.GetName(), "permission name should not be empty")
	}
}
