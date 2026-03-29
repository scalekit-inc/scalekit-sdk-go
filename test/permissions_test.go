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

func uniquePermissionName() string {
	return fmt.Sprintf("test_permission_%d", time.Now().UnixNano()/1e6)
}

func createPermission(t *testing.T, ctx context.Context, name string) {
	t.Helper()
	_, err := client.Permission().CreatePermission(ctx, &rolesv1.CreatePermission{
		Name:        name,
		Description: "Integration test permission",
	})
	require.NoError(t, err)
}

func deletePermission(t *testing.T, ctx context.Context, name string) {
	t.Helper()
	if err := client.Permission().DeletePermission(ctx, name); err != nil {
		t.Errorf("failed to delete permission %s: %v", name, err)
	}
}

func TestPermission_ListPermissions_NoParams(t *testing.T) {
	ctx := context.Background()

	resp, err := client.Permission().ListPermissions(ctx, "", 0)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPermission_ListPermissions_WithPageSize(t *testing.T) {
	ctx := context.Background()

	// Create a permission to ensure there is at least one
	name := uniquePermissionName()
	createPermission(t, ctx, name)
	defer deletePermission(t, ctx, name)

	resp, err := client.Permission().ListPermissions(ctx, "", 10)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPermission_ListPermissions_WithPageToken(t *testing.T) {
	ctx := context.Background()

	// Get first page
	first, err := client.Permission().ListPermissions(ctx, "", 0)
	require.NoError(t, err)

	// If a next page token exists, fetch the next page using it
	if first.GetNextPageToken() != "" {
		second, err := client.Permission().ListPermissions(ctx, first.GetNextPageToken(), 0)
		assert.NoError(t, err)
		assert.NotNil(t, second)
	}
}

func TestPermission_ListPermissions_WithBothParams(t *testing.T) {
	ctx := context.Background()

	resp, err := client.Permission().ListPermissions(ctx, "", 5)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	if resp.GetNextPageToken() != "" {
		next, err := client.Permission().ListPermissions(ctx, resp.GetNextPageToken(), 5)
		assert.NoError(t, err)
		assert.NotNil(t, next)
	}
}
