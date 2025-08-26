package scalekit

import (
	"context"

	authv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/auth"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/auth/authconnect"
)

// Type aliases for response types
type SendPasswordlessResponse = authv1.SendPasswordlessResponse
type VerifyPasswordLessResponse = authv1.VerifyPasswordLessResponse

// Type aliases for enum types
type TemplateType = authv1.TemplateType
type PasswordlessType = authv1.PasswordlessType

// Enum constants for TemplateType
const (
	TemplateTypeUnspecified = authv1.TemplateType_UNSPECIFIED
	TemplateTypeSignin      = authv1.TemplateType_SIGNIN
	TemplateTypeSignup      = authv1.TemplateType_SIGNUP
)

// Enum constants for PasswordlessType
const (
	PasswordlessTypeUnspecified = authv1.PasswordlessType_PASSWORDLESS_TYPE_UNSPECIFIED
	PasswordlessTypeOtp         = authv1.PasswordlessType_OTP
	PasswordlessTypeLink        = authv1.PasswordlessType_LINK
	PasswordlessTypeLinkOtp     = authv1.PasswordlessType_LINK_OTP
)

// SendPasswordlessOptions represents optional parameters for sending passwordless authentication
type SendPasswordlessOptions struct {
	Template          *TemplateType
	MagiclinkAuthUri  string // Use empty string for no magic link URI, or specify the authentication URI
	State             string // Use empty string for no state, or specify a custom state value
	ExpiresIn         uint32 // Use 0 for server default, or specify seconds (e.g., 3600 for 1 hour)
	TemplateVariables map[string]string
}

// VerifyPasswordlessOptions represents options for verifying passwordless authentication
type VerifyPasswordlessOptions struct {
	Code          string // Use empty string for no code, or specify the OTP code
	LinkToken     string // Use empty string for no link token, or specify the link token
	AuthRequestId string // Use empty string for no auth request id, or specify the id
}

// PasswordlessService interface defines the methods for passwordless authentication
type PasswordlessService interface {
	SendPasswordlessEmail(ctx context.Context, email string, options *SendPasswordlessOptions) (*SendPasswordlessResponse, error)
	VerifyPasswordlessEmail(ctx context.Context, options *VerifyPasswordlessOptions) (*VerifyPasswordLessResponse, error)
	ResendPasswordlessEmail(ctx context.Context, authRequestId string) (*SendPasswordlessResponse, error)
}

type passwordlessService struct {
	coreClient *coreClient
	client     authconnect.PasswordlessServiceClient
}

// newPasswordlessClient creates a new passwordless client
func newPasswordlessClient(coreClient *coreClient) PasswordlessService {
	return &passwordlessService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, authconnect.NewPasswordlessServiceClient),
	}
}

// SendPasswordlessEmail sends a passwordless authentication email
func (p *passwordlessService) SendPasswordlessEmail(ctx context.Context, email string, options *SendPasswordlessOptions) (*SendPasswordlessResponse, error) {
	request := &authv1.SendPasswordlessRequest{
		Email: email,
	}

	if options != nil {
		request.Template = options.Template
		// Convert string to *string for the protobuf request
		if options.MagiclinkAuthUri != "" {
			request.MagiclinkAuthUri = &options.MagiclinkAuthUri
		}
		// Convert string to *string for the protobuf request
		if options.State != "" {
			request.State = &options.State
		}
		// Convert uint32 to *uint32 for the protobuf request
		if options.ExpiresIn > 0 {
			request.ExpiresIn = &options.ExpiresIn
		}
		request.TemplateVariables = options.TemplateVariables
	}

	return newConnectExecuter(
		p.coreClient,
		p.client.SendPasswordlessEmail,
		request,
	).exec(ctx)
}

// VerifyPasswordlessEmail verifies a passwordless authentication
func (p *passwordlessService) VerifyPasswordlessEmail(ctx context.Context, options *VerifyPasswordlessOptions) (*VerifyPasswordLessResponse, error) {
	request := &authv1.VerifyPasswordLessRequest{}

	if options != nil {
		if options.Code != "" {
			request.AuthCredential = &authv1.VerifyPasswordLessRequest_Code{
				Code: options.Code,
			}
		} else if options.LinkToken != "" {
			request.AuthCredential = &authv1.VerifyPasswordLessRequest_LinkToken{
				LinkToken: options.LinkToken,
			}
		}
		if options.AuthRequestId != "" {
			request.AuthRequestId = &options.AuthRequestId
		}
	}

	return newConnectExecuter(
		p.coreClient,
		p.client.VerifyPasswordlessEmail,
		request,
	).exec(ctx)
}

// ResendPasswordlessEmail resends a passwordless authentication email
func (p *passwordlessService) ResendPasswordlessEmail(ctx context.Context, authRequestId string) (*SendPasswordlessResponse, error) {
	if authRequestId == "" {
		return nil, ErrAuthRequestIdRequired
	}

	request := &authv1.ResendPasswordlessRequest{
		AuthRequestId: authRequestId,
	}

	return newConnectExecuter(
		p.coreClient,
		p.client.ResendPasswordlessEmail,
		request,
	).exec(ctx)
}
