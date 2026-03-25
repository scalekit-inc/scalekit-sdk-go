package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	rolesv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/roles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createChildRole creates a role that extends the given base role.
func createChildRole(t *testing.T, ctx context.Context, name, extendsName string) {
	t.Helper()
	desc := fmt.Sprintf("Integration test child role for %s", name)
	_, err := client.Role().CreateRole(ctx, &rolesv1.CreateRole{
		Name:        name,
		DisplayName: "Test Child Role " + name,
		Description: &desc,
		Extends:     &extendsName,
	})
	require.NoError(t, err)
}

func TestRole_UpdateDefaultRoles_WithEnvRoles(t *testing.T) {
	ctx := context.Background()

	creatorRole := envOrSkip(t, "TEST_DEFAULT_CREATOR_ROLE")
	memberRole := envOrSkip(t, "TEST_DEFAULT_MEMBER_ROLE")

	resp, err := client.Role().UpdateDefaultRoles(ctx, creatorRole, memberRole)
	require.NoError(t, err)
	require.NotNil(t, resp)
	t.Logf("UpdateDefaultRoles: defaultCreator=%q defaultMember=%q",
		resp.GetDefaultCreator().GetName(), resp.GetDefaultMember().GetName())
}

func TestRole_UpdateDefaultRoles_UsesListedRoles(t *testing.T) {
	ctx := context.Background()

	// List existing roles to find valid role names to use as defaults.
	rolesResp, err := client.Role().ListRoles(ctx)
	require.NoError(t, err)
	require.NotNil(t, rolesResp)

	roles := rolesResp.GetRoles()
	if len(roles) < 2 {
		t.Skip("skipping: need at least 2 roles in the environment to test UpdateDefaultRoles")
	}

	creatorRole := roles[0].GetName()
	memberRole := roles[1].GetName()
	if creatorRole == "" || memberRole == "" {
		t.Skip("skipping: role names are empty, cannot set defaults")
	}

	resp, err := client.Role().UpdateDefaultRoles(ctx, creatorRole, memberRole)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestRole_ListDependentRoles_Empty(t *testing.T) {
	ctx := context.Background()

	// Create a standalone role with no children.
	baseName := uniqueRoleName()
	createEnvRole(t, ctx, baseName)
	t.Cleanup(func() { deleteEnvRole(t, ctx, baseName) })

	resp, err := client.Role().ListDependentRoles(ctx, baseName)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// A brand-new role has no dependents.
	assert.Empty(t, resp.GetRoles(), "expected no dependent roles for a new standalone role")
}

func TestRole_ListDependentRoles_WithChild(t *testing.T) {
	ctx := context.Background()

	// Create a base role.
	baseName := uniqueRoleName()
	createEnvRole(t, ctx, baseName)
	t.Cleanup(func() { deleteEnvRole(t, ctx, baseName) })

	// Create a child role that extends the base.
	childName := uniqueRoleName()
	createChildRole(t, ctx, childName, baseName)
	t.Cleanup(func() { deleteEnvRole(t, ctx, childName) })

	resp, err := client.Role().ListDependentRoles(ctx, baseName)
	require.NoError(t, err)
	require.NotNil(t, resp)
	t.Logf("ListDependentRoles(%q): %d dependent role(s)", baseName, len(resp.GetRoles()))

	found := false
	for _, r := range resp.GetRoles() {
		if r.GetName() == childName {
			found = true
			break
		}
	}
	assert.True(t, found, "expected child role %q in dependent roles list", childName)
}

func TestRole_ListDependentRoles_WithEnvRole(t *testing.T) {
	ctx := context.Background()

	roleName := envOrSkip(t, "TEST_ROLE_NAME")

	resp, err := client.Role().ListDependentRoles(ctx, roleName)
	require.NoError(t, err)
	require.NotNil(t, resp)
	t.Logf("ListDependentRoles(%q): %d dependent role(s)", roleName, len(resp.GetRoles()))
	for _, r := range resp.GetRoles() {
		assert.NotEmpty(t, r.GetName(), "role name should not be empty")
	}
}

func TestRole_UpdateDefaultRoles_RequiresDefaultCreatorRole(t *testing.T) {
	ctx := context.Background()

	_, err := client.Role().UpdateDefaultRoles(ctx, "", "some-role")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrDefaultCreatorRoleRequired)
}

func TestRole_UpdateDefaultRoles_RequiresDefaultMemberRole(t *testing.T) {
	ctx := context.Background()

	_, err := client.Role().UpdateDefaultRoles(ctx, "some-role", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrDefaultMemberRoleRequired)
}

func TestRole_ListDependentRoles_RequiresRoleName(t *testing.T) {
	ctx := context.Background()

	_, err := client.Role().ListDependentRoles(ctx, "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrRoleNameRequired)
}
