package scalekit

import (
	"connectrpc.com/connect"
	"context"
	authv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/auth"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/auth/authconnect"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// Type aliases for auth requests.
type UpdateLoginUserDetailsRequest = authv1.UpdateLoginUserDetailsRequest
type LoggedInUserDetails = authv1.User

// AuthService provides helper methods for interacting with the Auth gRPC surface.
type AuthService interface {
	UpdateLoginUserDetails(ctx context.Context, req *UpdateLoginUserDetailsRequest) error
}

type authServiceClient interface {
	UpdateLoginUserDetails(context.Context, *connect.Request[authv1.UpdateLoginUserDetailsRequest]) (*connect.Response[emptypb.Empty], error)
}

type authService struct {
	coreClient *coreClient
	client     authServiceClient
}

func newAuthService(coreClient *coreClient) AuthService {
	return &authService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, authconnect.NewAuthServiceClient),
	}
}

func (a *authService) UpdateLoginUserDetails(ctx context.Context, req *UpdateLoginUserDetailsRequest) error {
	_, err := newConnectExecuter(
		a.coreClient,
		a.client.UpdateLoginUserDetails,
		req,
	).exec(ctx)
	return err
}
