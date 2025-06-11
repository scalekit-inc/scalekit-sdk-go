package scalekit

import (
	"context"

	usersv1 "github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/users"
	"github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/users/usersconnect"
)

// Type aliases for response types
type CreateUserAndMembershipResponse = usersv1.CreateUserAndMembershipResponse
type UpdateUserResponse = usersv1.UpdateUserResponse
type GetUserResponse = usersv1.GetUserResponse
type ListOrganizationUsersResponse = usersv1.ListOrganizationUsersResponse
type CreateMembershipResponse = usersv1.CreateMembershipResponse
type UpdateMembershipResponse = usersv1.UpdateMembershipResponse

// ListUsersOptions represents optional parameters for listing users
type ListUsersOptions struct {
	PageSize  uint32
	PageToken string
}

type UserService interface {
	CreateUserAndMembership(ctx context.Context, organizationId string, user *usersv1.CreateUser, sendActivationEmail bool) (*CreateUserAndMembershipResponse, error)
	UpdateUser(ctx context.Context, userId string, updateUser *usersv1.UpdateUser) (*UpdateUserResponse, error)
	GetUser(ctx context.Context, userId string) (*GetUserResponse, error)
	ListOrganizationUsers(ctx context.Context, organizationId string, options *ListUsersOptions) (*ListOrganizationUsersResponse, error)
	DeleteUser(ctx context.Context, userId string) error
	CreateMembership(ctx context.Context, organizationId string, userId string, membership *usersv1.CreateMembership, sendActivationEmail bool) (*CreateMembershipResponse, error)
	UpdateMembership(ctx context.Context, organizationId string, userId string, membership *usersv1.UpdateMembership) (*UpdateMembershipResponse, error)
	DeleteMembership(ctx context.Context, organizationId string, userId string, cascade bool) error
}

type userService struct {
	coreClient *coreClient
	client     usersconnect.UserServiceClient
}

// newUserClient creates a new user client
func newUserClient(coreClient *coreClient) UserService {
	return &userService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, usersconnect.NewUserServiceClient),
	}
}

// CreateUserAndMembership creates a new user with membership in the organization
func (u *userService) CreateUserAndMembership(ctx context.Context, organizationId string, user *usersv1.CreateUser, sendActivationEmail bool) (*CreateUserAndMembershipResponse, error) {
	return newConnectExecuter(
		u.coreClient,
		u.client.CreateUserAndMembership,
		&usersv1.CreateUserAndMembershipRequest{
			OrganizationId:      organizationId,
			User:                user,
			SendActivationEmail: sendActivationEmail,
		},
	).exec(ctx)
}

// UpdateUser updates an existing user
func (u *userService) UpdateUser(ctx context.Context, userId string, updateUser *usersv1.UpdateUser) (*UpdateUserResponse, error) {
	request := &usersv1.UpdateUserRequest{
		User: updateUser,
	}
	request.Identities = &usersv1.UpdateUserRequest_Id{Id: userId}

	return newConnectExecuter(
		u.coreClient,
		u.client.UpdateUser,
		request,
	).exec(ctx)
}

// GetUser retrieves a user by ID
func (u *userService) GetUser(ctx context.Context, userId string) (*GetUserResponse, error) {
	request := &usersv1.GetUserRequest{}
	request.Identities = &usersv1.GetUserRequest_Id{Id: userId}

	return newConnectExecuter(
		u.coreClient,
		u.client.GetUser,
		request,
	).exec(ctx)
}

// ListOrganizationUsers retrieves a list of users in the organization
func (u *userService) ListOrganizationUsers(ctx context.Context, organizationId string, options *ListUsersOptions) (*ListOrganizationUsersResponse, error) {
	request := &usersv1.ListOrganizationUsersRequest{
		OrganizationId: organizationId,
	}
	if options != nil {
		request.PageSize = options.PageSize
		request.PageToken = options.PageToken
	}

	return newConnectExecuter(
		u.coreClient,
		u.client.ListOrganizationUsers,
		request,
	).exec(ctx)
}

// DeleteUser deletes a user
func (u *userService) DeleteUser(ctx context.Context, userId string) error {
	request := &usersv1.DeleteUserRequest{}
	request.Identities = &usersv1.DeleteUserRequest_Id{Id: userId}

	_, err := newConnectExecuter(
		u.coreClient,
		u.client.DeleteUser,
		request,
	).exec(ctx)
	return err
}

// CreateMembership creates a membership for a user in an organization
func (u *userService) CreateMembership(ctx context.Context, organizationId string, userId string, membership *usersv1.CreateMembership, sendActivationEmail bool) (*CreateMembershipResponse, error) {
	request := &usersv1.CreateMembershipRequest{
		OrganizationId:      organizationId,
		Membership:          membership,
		SendActivationEmail: sendActivationEmail,
	}
	request.Identities = &usersv1.CreateMembershipRequest_Id{Id: userId}

	return newConnectExecuter(
		u.coreClient,
		u.client.CreateMembership,
		request,
	).exec(ctx)
}

// UpdateMembership updates a user's membership in an organization
func (u *userService) UpdateMembership(ctx context.Context, organizationId string, userId string, membership *usersv1.UpdateMembership) (*UpdateMembershipResponse, error) {
	request := &usersv1.UpdateMembershipRequest{
		OrganizationId: organizationId,
		Membership:     membership,
	}
	request.Identities = &usersv1.UpdateMembershipRequest_Id{Id: userId}

	return newConnectExecuter(
		u.coreClient,
		u.client.UpdateMembership,
		request,
	).exec(ctx)
}

// DeleteMembership deletes a user's membership from an organization
func (u *userService) DeleteMembership(ctx context.Context, organizationId string, userId string, cascade bool) error {
	request := &usersv1.DeleteMembershipRequest{
		OrganizationId: organizationId,
		Cascade:        &cascade,
	}
	request.Identities = &usersv1.DeleteMembershipRequest_Id{Id: userId}

	_, err := newConnectExecuter(
		u.coreClient,
		u.client.DeleteMembership,
		request,
	).exec(ctx)
	return err
}
