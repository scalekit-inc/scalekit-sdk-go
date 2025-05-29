package scalekit

import (
	"context"

	usersv1 "github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/users"
	"github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/users/usersconnect"
)

// Type aliases for response types
type CreateUserResponse = usersv1.CreateUserResponse
type UpdateUserResponse = usersv1.UpdateUserResponse
type GetUserResponse = usersv1.GetUserResponse
type ListUserResponse = usersv1.ListUserResponse
type AddUserResponse = usersv1.AddUserResponse

// ListUsersOptions represents optional parameters for listing users
type ListUsersOptions struct {
	PageSize  uint32
	PageToken string
}

type UserService interface {
	CreateUser(ctx context.Context, organizationId string, user *usersv1.User) (*CreateUserResponse, error)
	UpdateUser(ctx context.Context, organizationId string, userId string, updateUser *usersv1.UpdateUser) (*UpdateUserResponse, error)
	GetUser(ctx context.Context, organizationId string, userId string) (*GetUserResponse, error)
	ListUsers(ctx context.Context, organizationId string, options *ListUsersOptions) (*ListUserResponse, error)
	DeleteUser(ctx context.Context, organizationId string, userId string) error
	AddUserToOrganization(ctx context.Context, organizationId string, userId string) (*AddUserResponse, error)
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

// CreateUser creates a new user in the organization
func (u *userService) CreateUser(ctx context.Context, organizationId string, user *usersv1.User) (*CreateUserResponse, error) {
	return newConnectExecuter(
		u.coreClient,
		u.client.CreateUser,
		&usersv1.CreateUserRequest{
			OrganizationId: organizationId,
			User:           user,
		},
	).exec(ctx)
}

// UpdateUser updates an existing user
func (u *userService) UpdateUser(ctx context.Context, organizationId string, userId string, updateUser *usersv1.UpdateUser) (*UpdateUserResponse, error) {
	request := &usersv1.UpdateUserRequest{
		OrganizationId: organizationId,
		User:           updateUser,
	}
	request.Identities = &usersv1.UpdateUserRequest_Id{Id: userId}

	return newConnectExecuter(
		u.coreClient,
		u.client.UpdateUser,
		request,
	).exec(ctx)
}

// GetUser retrieves a user by ID
func (u *userService) GetUser(ctx context.Context, organizationId string, userId string) (*GetUserResponse, error) {
	request := &usersv1.GetUserRequest{
		OrganizationId: organizationId,
	}
	request.Identities = &usersv1.GetUserRequest_Id{Id: userId}

	return newConnectExecuter(
		u.coreClient,
		u.client.GetUser,
		request,
	).exec(ctx)
}

// ListUsers retrieves a list of users in the organization
func (u *userService) ListUsers(ctx context.Context, organizationId string, options *ListUsersOptions) (*ListUserResponse, error) {
	request := &usersv1.ListUserRequest{
		OrganizationId: organizationId,
	}
	if options != nil {
		request.PageSize = options.PageSize
		request.PageToken = options.PageToken
	}

	return newConnectExecuter(
		u.coreClient,
		u.client.ListUsers,
		request,
	).exec(ctx)
}

// DeleteUser deletes a user from the organization
func (u *userService) DeleteUser(ctx context.Context, organizationId string, userId string) error {
	request := &usersv1.DeleteUserRequest{
		OrganizationId: organizationId,
	}
	request.Identities = &usersv1.DeleteUserRequest_Id{Id: userId}

	_, err := newConnectExecuter(
		u.coreClient,
		u.client.DeleteUser,
		request,
	).exec(ctx)
	return err
}

// AddUserToOrganization adds an existing user to an organization
func (u *userService) AddUserToOrganization(ctx context.Context, organizationId string, userId string) (*AddUserResponse, error) {
	request := &usersv1.AddUserRequest{
		OrganizationId: organizationId,
	}
	request.Identities = &usersv1.AddUserRequest_Id{Id: userId}

	return newConnectExecuter(
		u.coreClient,
		u.client.AddUserToOrganization,
		request,
	).exec(ctx)
}
