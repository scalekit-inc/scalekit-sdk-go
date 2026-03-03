package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	rolesv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/roles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uniqueRoleName() string {
	return fmt.Sprintf("test_role_%d", time.Now().UnixNano()/1e6)
}

func createEnvRole(t *testing.T, ctx context.Context, name string) {
	t.Helper()
	desc := "Integration test role"
	_, err := client.Role().CreateRole(ctx, &rolesv1.CreateRole{
		Name:        name,
		DisplayName: fmt.Sprintf("Test Role %s", name),
		Description: &desc,
	})
	require.NoError(t, err)
}

func deleteEnvRole(t *testing.T, ctx context.Context, name string) {
	t.Helper()
	_ = client.Role().DeleteRole(ctx, name)
}

func TestRole_DeleteRoleBase(t *testing.T) {
	ctx := context.Background()

	// Create a base role
	baseName := uniqueRoleName()
	createEnvRole(t, ctx, baseName)
	defer deleteEnvRole(t, ctx, baseName)

	// Create a child role that extends the base
	childName := uniqueRoleName()
	childDesc := "Child role for DeleteRoleBase test"
	_, err := client.Role().CreateRole(ctx, &rolesv1.CreateRole{
		Name:        childName,
		DisplayName: fmt.Sprintf("Child Role %s", childName),
		Description: &childDesc,
		Extends:     &baseName,
	})
	require.NoError(t, err)
	defer deleteEnvRole(t, ctx, childName)

	// Delete the base relationship from the child role
	err = client.Role().DeleteRoleBase(ctx, childName)
	assert.NoError(t, err)
}
