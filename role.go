package scalekit

import (
	"context"

	rolesv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/roles"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/roles/rolesconnect"
)

// Type aliases for role-related responses
type CreateRoleResponse = rolesv1.CreateRoleResponse
type GetRoleResponse = rolesv1.GetRoleResponse
type ListRolesResponse = rolesv1.ListRolesResponse
type UpdateRoleResponse = rolesv1.UpdateRoleResponse
type CreateOrganizationRoleResponse = rolesv1.CreateOrganizationRoleResponse
type GetOrganizationRoleResponse = rolesv1.GetOrganizationRoleResponse
type ListOrganizationRolesResponse = rolesv1.ListOrganizationRolesResponse
type UpdateOrganizationRoleResponse = rolesv1.UpdateOrganizationRoleResponse
type GetRoleUsersCountResponse = rolesv1.GetRoleUsersCountResponse
type GetOrganizationRoleUsersCountResponse = rolesv1.GetOrganizationRoleUsersCountResponse
type UpdateDefaultOrganizationRolesResponse = rolesv1.UpdateDefaultOrganizationRolesResponse

// RoleService defines the interface for role management operations
type RoleService interface {
	// Environment-level role management
	CreateRole(ctx context.Context, role *rolesv1.CreateRole) (*CreateRoleResponse, error)
	GetRole(ctx context.Context, roleName string) (*GetRoleResponse, error)
	ListRoles(ctx context.Context) (*ListRolesResponse, error)
	UpdateRole(ctx context.Context, roleName string, role *rolesv1.UpdateRole) (*UpdateRoleResponse, error)
	DeleteRole(ctx context.Context, roleName string, reassignRoleName ...string) error
	GetRoleUsersCount(ctx context.Context, roleName string) (*GetRoleUsersCountResponse, error)

	// Organization-level role management
	CreateOrganizationRole(ctx context.Context, orgId string, role *rolesv1.CreateOrganizationRole) (*CreateOrganizationRoleResponse, error)
	GetOrganizationRole(ctx context.Context, orgId, roleName string) (*GetOrganizationRoleResponse, error)
	ListOrganizationRoles(ctx context.Context, orgId string) (*ListOrganizationRolesResponse, error)
	UpdateOrganizationRole(ctx context.Context, orgId, roleName string, role *rolesv1.UpdateRole) (*UpdateOrganizationRoleResponse, error)
	DeleteOrganizationRole(ctx context.Context, orgId, roleName string, reassignRoleName ...string) error
	GetOrganizationRoleUsersCount(ctx context.Context, orgId, roleName string) (*GetOrganizationRoleUsersCountResponse, error)
	UpdateDefaultOrganizationRoles(ctx context.Context, orgId, defaultMemberRole string) (*UpdateDefaultOrganizationRolesResponse, error)
	DeleteOrganizationRoleBase(ctx context.Context, orgId, roleName string) error
}

type roleService struct {
	coreClient *coreClient
	client     rolesconnect.RolesServiceClient
}

func newRoleService(coreClient *coreClient) RoleService {
	return &roleService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, rolesconnect.NewRolesServiceClient),
	}
}

// Environment-level role management methods

// CreateRole creates a new role in the environment
func (r *roleService) CreateRole(ctx context.Context, role *rolesv1.CreateRole) (*CreateRoleResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.CreateRole,
		&rolesv1.CreateRoleRequest{
			Role: role,
		},
	).exec(ctx)
}

// GetRole retrieves a role by name
func (r *roleService) GetRole(ctx context.Context, roleName string) (*GetRoleResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.GetRole,
		&rolesv1.GetRoleRequest{
			RoleName: roleName,
		},
	).exec(ctx)
}

// ListRoles lists all roles in the environment
func (r *roleService) ListRoles(ctx context.Context) (*ListRolesResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.ListRoles,
		&rolesv1.ListRolesRequest{},
	).exec(ctx)
}

// UpdateRole updates an existing role by name
func (r *roleService) UpdateRole(ctx context.Context, roleName string, role *rolesv1.UpdateRole) (*UpdateRoleResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.UpdateRole,
		&rolesv1.UpdateRoleRequest{
			RoleName: roleName,
			Role:     role,
		},
	).exec(ctx)
}

// DeleteRole deletes a role by name
func (r *roleService) DeleteRole(ctx context.Context, roleName string, reassignRoleName ...string) error {
	req := &rolesv1.DeleteRoleRequest{
		RoleName: roleName,
	}
	if len(reassignRoleName) > 0 {
		req.ReassignRoleName = &reassignRoleName[0]
	}
	_, err := newConnectExecuter(
		r.coreClient,
		r.client.DeleteRole,
		req,
	).exec(ctx)
	return err
}

// GetRoleUsersCount gets the count of users associated with a role
func (r *roleService) GetRoleUsersCount(ctx context.Context, roleName string) (*GetRoleUsersCountResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.GetRoleUsersCount,
		&rolesv1.GetRoleUsersCountRequest{
			RoleName: roleName,
		},
	).exec(ctx)
}

// Organization-level role management methods

// CreateOrganizationRole creates a new role in an organization
func (r *roleService) CreateOrganizationRole(ctx context.Context, orgId string, role *rolesv1.CreateOrganizationRole) (*CreateOrganizationRoleResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.CreateOrganizationRole,
		&rolesv1.CreateOrganizationRoleRequest{
			OrgId: orgId,
			Role:  role,
		},
	).exec(ctx)
}

// GetOrganizationRole retrieves an organization role by name
func (r *roleService) GetOrganizationRole(ctx context.Context, orgId, roleName string) (*GetOrganizationRoleResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.GetOrganizationRole,
		&rolesv1.GetOrganizationRoleRequest{
			OrgId:    orgId,
			RoleName: roleName,
		},
	).exec(ctx)
}

// ListOrganizationRoles lists all roles in an organization
func (r *roleService) ListOrganizationRoles(ctx context.Context, orgId string) (*ListOrganizationRolesResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.ListOrganizationRoles,
		&rolesv1.ListOrganizationRolesRequest{
			OrgId: orgId,
		},
	).exec(ctx)
}

// UpdateOrganizationRole updates an existing organization role by name
func (r *roleService) UpdateOrganizationRole(ctx context.Context, orgId, roleName string, role *rolesv1.UpdateRole) (*UpdateOrganizationRoleResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.UpdateOrganizationRole,
		&rolesv1.UpdateOrganizationRoleRequest{
			OrgId:    orgId,
			RoleName: roleName,
			Role:     role,
		},
	).exec(ctx)
}

// DeleteOrganizationRole deletes an organization role by name
func (r *roleService) DeleteOrganizationRole(ctx context.Context, orgId, roleName string, reassignRoleName ...string) error {
	req := &rolesv1.DeleteOrganizationRoleRequest{
		OrgId:    orgId,
		RoleName: roleName,
	}
	if len(reassignRoleName) > 0 {
		req.ReassignRoleName = &reassignRoleName[0]
	}
	_, err := newConnectExecuter(
		r.coreClient,
		r.client.DeleteOrganizationRole,
		req,
	).exec(ctx)
	return err
}

// GetOrganizationRoleUsersCount gets the count of users associated with an organization role
func (r *roleService) GetOrganizationRoleUsersCount(ctx context.Context, orgId, roleName string) (*GetOrganizationRoleUsersCountResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.GetOrganizationRoleUsersCount,
		&rolesv1.GetOrganizationRoleUsersCountRequest{
			OrgId:    orgId,
			RoleName: roleName,
		},
	).exec(ctx)
}

// UpdateDefaultOrganizationRoles updates the default roles for an organization
func (r *roleService) UpdateDefaultOrganizationRoles(ctx context.Context, orgId, defaultMemberRole string) (*UpdateDefaultOrganizationRolesResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.UpdateDefaultOrganizationRoles,
		&rolesv1.UpdateDefaultOrganizationRolesRequest{
			OrgId:             orgId,
			DefaultMemberRole: defaultMemberRole,
		},
	).exec(ctx)
}

// DeleteOrganizationRoleBase deletes the base relationship for an organization role
func (r *roleService) DeleteOrganizationRoleBase(ctx context.Context, orgId, roleName string) error {
	_, err := newConnectExecuter(
		r.coreClient,
		r.client.DeleteOrganizationRoleBase,
		&rolesv1.DeleteOrganizationRoleBaseRequest{
			OrgId:    orgId,
			RoleName: roleName,
		},
	).exec(ctx)
	return err
}
