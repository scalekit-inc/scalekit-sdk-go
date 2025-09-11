package scalekit

import (
	"context"

	rolesv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/roles"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/roles/rolesconnect"
)

// Type aliases for permission-related responses
type CreatePermissionResponse = rolesv1.CreatePermissionResponse
type GetPermissionResponse = rolesv1.GetPermissionResponse
type ListPermissionsResponse = rolesv1.ListPermissionsResponse
type UpdatePermissionResponse = rolesv1.UpdatePermissionResponse
type ListRolePermissionsResponse = rolesv1.ListRolePermissionsResponse
type AddPermissionsToRoleResponse = rolesv1.AddPermissionsToRoleResponse
type ListEffectiveRolePermissionsResponse = rolesv1.ListEffectiveRolePermissionsResponse

// PermissionService defines the interface for permission management operations
type PermissionService interface {
	// Permission management
	CreatePermission(ctx context.Context, permission *rolesv1.CreatePermission) (*CreatePermissionResponse, error)
	GetPermission(ctx context.Context, permissionName string) (*GetPermissionResponse, error)
	ListPermissions(ctx context.Context, pageToken ...string) (*ListPermissionsResponse, error)
	UpdatePermission(ctx context.Context, permissionName string, permission *rolesv1.CreatePermission) (*UpdatePermissionResponse, error)
	DeletePermission(ctx context.Context, permissionName string) error

	// Role-Permission relationships
	ListRolePermissions(ctx context.Context, roleName string) (*ListRolePermissionsResponse, error)
	AddPermissionsToRole(ctx context.Context, roleName string, permissionNames []string) (*AddPermissionsToRoleResponse, error)
	RemovePermissionFromRole(ctx context.Context, roleName, permissionName string) error
	ListEffectiveRolePermissions(ctx context.Context, roleName string) (*ListEffectiveRolePermissionsResponse, error)
}

type permissionService struct {
	coreClient *coreClient
	client     rolesconnect.RolesServiceClient
}

func newPermissionService(coreClient *coreClient) PermissionService {
	return &permissionService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, rolesconnect.NewRolesServiceClient),
	}
}

// Permission management methods

// CreatePermission creates a new permission
func (p *permissionService) CreatePermission(ctx context.Context, permission *rolesv1.CreatePermission) (*CreatePermissionResponse, error) {
	return newConnectExecuter(
		p.coreClient,
		p.client.CreatePermission,
		&rolesv1.CreatePermissionRequest{
			Permission: permission,
		},
	).exec(ctx)
}

// GetPermission retrieves a permission by name
func (p *permissionService) GetPermission(ctx context.Context, permissionName string) (*GetPermissionResponse, error) {
	return newConnectExecuter(
		p.coreClient,
		p.client.GetPermission,
		&rolesv1.GetPermissionRequest{
			PermissionName: permissionName,
		},
	).exec(ctx)
}

// ListPermissions lists all permissions with optional pagination
func (p *permissionService) ListPermissions(ctx context.Context, pageToken ...string) (*ListPermissionsResponse, error) {
	req := &rolesv1.ListPermissionsRequest{}
	if len(pageToken) > 0 {
		req.PageToken = &pageToken[0]
	}
	return newConnectExecuter(
		p.coreClient,
		p.client.ListPermissions,
		req,
	).exec(ctx)
}

// UpdatePermission updates an existing permission by name
func (p *permissionService) UpdatePermission(ctx context.Context, permissionName string, permission *rolesv1.CreatePermission) (*UpdatePermissionResponse, error) {
	return newConnectExecuter(
		p.coreClient,
		p.client.UpdatePermission,
		&rolesv1.UpdatePermissionRequest{
			PermissionName: permissionName,
			Permission:     permission,
		},
	).exec(ctx)
}

// DeletePermission deletes a permission by name
func (p *permissionService) DeletePermission(ctx context.Context, permissionName string) error {
	_, err := newConnectExecuter(
		p.coreClient,
		p.client.DeletePermission,
		&rolesv1.DeletePermissionRequest{
			PermissionName: permissionName,
		},
	).exec(ctx)
	return err
}

// Role-Permission relationship methods

// ListRolePermissions lists all permissions associated with a role
func (p *permissionService) ListRolePermissions(ctx context.Context, roleName string) (*ListRolePermissionsResponse, error) {
	return newConnectExecuter(
		p.coreClient,
		p.client.ListRolePermissions,
		&rolesv1.ListRolePermissionsRequest{
			RoleName: roleName,
		},
	).exec(ctx)
}

// AddPermissionsToRole adds permissions to a role
func (p *permissionService) AddPermissionsToRole(ctx context.Context, roleName string, permissionNames []string) (*AddPermissionsToRoleResponse, error) {
	return newConnectExecuter(
		p.coreClient,
		p.client.AddPermissionsToRole,
		&rolesv1.AddPermissionsToRoleRequest{
			RoleName:        roleName,
			PermissionNames: permissionNames,
		},
	).exec(ctx)
}

// RemovePermissionFromRole removes a permission from a role
func (p *permissionService) RemovePermissionFromRole(ctx context.Context, roleName, permissionName string) error {
	_, err := newConnectExecuter(
		p.coreClient,
		p.client.RemovePermissionFromRole,
		&rolesv1.RemovePermissionFromRoleRequest{
			RoleName:       roleName,
			PermissionName: permissionName,
		},
	).exec(ctx)
	return err
}

// ListEffectiveRolePermissions lists all effective permissions for a role (including inherited permissions)
func (p *permissionService) ListEffectiveRolePermissions(ctx context.Context, roleName string) (*ListEffectiveRolePermissionsResponse, error) {
	return newConnectExecuter(
		p.coreClient,
		p.client.ListEffectiveRolePermissions,
		&rolesv1.ListEffectiveRolePermissionsRequest{
			RoleName: roleName,
		},
	).exec(ctx)
}
