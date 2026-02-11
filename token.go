package scalekit

import (
	"context"
	"errors"
	"fmt"

	tokensv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/tokens"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/tokens/tokensconnect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateTokenResponse = tokensv1.CreateTokenResponse
type ValidateTokenResponse = tokensv1.ValidateTokenResponse
type ListTokensResponse = tokensv1.ListTokensResponse
type TokenInfo = tokensv1.Token

type CreateTokenOptions struct {
	UserId       string
	CustomClaims map[string]string
	Expiry       *timestamppb.Timestamp
	Description  string
}

type ListTokensOptions struct {
	UserId    string
	PageSize  int32
	PageToken string
}

type TokenService interface {
	CreateToken(ctx context.Context, organizationId string, options CreateTokenOptions) (*CreateTokenResponse, error)
	ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error)
	InvalidateToken(ctx context.Context, token string) error
	ListTokens(ctx context.Context, organizationId string, options ListTokensOptions) (*ListTokensResponse, error)
}

type tokenService struct {
	coreClient *coreClient
	client     tokensconnect.ApiTokenServiceClient
}

func newTokenService(coreClient *coreClient) TokenService {
	return &tokenService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, tokensconnect.NewApiTokenServiceClient),
	}
}

func (t *tokenService) CreateToken(ctx context.Context, organizationId string, options CreateTokenOptions) (*CreateTokenResponse, error) {
	if organizationId == "" {
		return nil, errors.New("organizationId is required")
	}
	createToken := &tokensv1.CreateToken{
		OrganizationId: organizationId,
	}
	if options.UserId != "" {
		createToken.UserId = options.UserId
	}
	if options.CustomClaims != nil {
		createToken.CustomClaims = options.CustomClaims
	}
	if options.Expiry != nil {
		createToken.Expiry = options.Expiry
	}
	if options.Description != "" {
		createToken.Description = options.Description
	}

	return newConnectExecuter(
		t.coreClient,
		t.client.CreateToken,
		&tokensv1.CreateTokenRequest{
			Token: createToken,
		},
	).exec(ctx)
}

func (t *tokenService) ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	if token == "" {
		return nil, errors.New("token is required")
	}
	result, err := newConnectExecuter(
		t.coreClient,
		t.client.ValidateToken,
		&tokensv1.ValidateTokenRequest{
			Token: token,
		},
	).exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenValidationFailed, err)
	}
	return result, nil
}

func (t *tokenService) InvalidateToken(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("token is required")
	}
	_, err := newConnectExecuter(
		t.coreClient,
		t.client.InvalidateToken,
		&tokensv1.InvalidateTokenRequest{
			Token: token,
		},
	).exec(ctx)

	return err
}

func (t *tokenService) ListTokens(ctx context.Context, organizationId string, options ListTokensOptions) (*ListTokensResponse, error) {
	if organizationId == "" {
		return nil, errors.New("organizationId is required")
	}
	request := &tokensv1.ListTokensRequest{
		OrganizationId: organizationId,
	}
	if options.UserId != "" {
		request.UserId = options.UserId
	}
	if options.PageSize > 0 {
		request.PageSize = options.PageSize
	}
	if options.PageToken != "" {
		request.PageToken = options.PageToken
	}

	return newConnectExecuter(
		t.coreClient,
		t.client.ListTokens,
		request,
	).exec(ctx)
}
