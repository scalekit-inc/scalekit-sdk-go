package scalekit

import (
	"context"

	clientsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients/clientsconnect"
)

// Type aliases for public API surface
type OrganizationClientInfo = clientsv1.OrganizationClient
type CreateOrganizationClientResponse = clientsv1.CreateOrganizationClientResponse
type GetOrganizationClientResponse = clientsv1.GetOrganizationClientResponse
type UpdateOrganizationClientResponse = clientsv1.UpdateOrganizationClientResponse
type CreateOrganizationClientSecretResponse = clientsv1.CreateOrganizationClientSecretResponse
type ListOrganizationClientsResponse = clientsv1.ListOrganizationClientsResponse

// CreateOrganizationClientOptions holds optional fields for creating an org client.
type CreateOrganizationClientOptions struct {
	Name         string
	Description  string
	CustomClaims []*clientsv1.CustomClaim
	Audience     []string
	Scopes       []string
}

// UpdateOrganizationClientOptions holds optional fields for updating an org client.
type UpdateOrganizationClientOptions struct {
	Name         string
	Description  string
	CustomClaims []*clientsv1.CustomClaim
	Audience     []string
	Scopes       []string
}

// ListOrganizationClientsOptions holds pagination parameters for listing org clients.
type ListOrganizationClientsOptions struct {
	PageSize  uint32
	PageToken string
}

// M2MService defines the interface for managing M2M organization clients.
type M2MService interface {
	CreateOrganizationClient(ctx context.Context, organizationId string, options CreateOrganizationClientOptions) (*CreateOrganizationClientResponse, error)
	GetOrganizationClient(ctx context.Context, organizationId string, clientId string) (*GetOrganizationClientResponse, error)
	UpdateOrganizationClient(ctx context.Context, organizationId string, clientId string, options UpdateOrganizationClientOptions) (*UpdateOrganizationClientResponse, error)
	DeleteOrganizationClient(ctx context.Context, organizationId string, clientId string) error
	CreateOrganizationClientSecret(ctx context.Context, organizationId string, clientId string) (*CreateOrganizationClientSecretResponse, error)
	DeleteOrganizationClientSecret(ctx context.Context, organizationId string, clientId string, secretId string) error
	ListOrganizationClients(ctx context.Context, organizationId string, options ListOrganizationClientsOptions) (*ListOrganizationClientsResponse, error)
}

type m2mService struct {
	coreClient *coreClient
	client     clientsconnect.ClientServiceClient
}

func newM2MService(coreClient *coreClient) M2MService {
	return &m2mService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, clientsconnect.NewClientServiceClient),
	}
}

// buildOrganizationClient constructs an OrganizationClient protobuf from options.
func buildOrganizationClient(name, description string, customClaims []*clientsv1.CustomClaim, audience, scopes []string) *clientsv1.OrganizationClient {
	client := &clientsv1.OrganizationClient{}
	if name != "" {
		client.Name = name
	}
	if description != "" {
		client.Description = description
	}
	if customClaims != nil {
		client.CustomClaims = customClaims
	}
	if audience != nil {
		client.Audience = audience
	}
	if scopes != nil {
		client.Scopes = scopes
	}
	return client
}

func (m *m2mService) CreateOrganizationClient(ctx context.Context, organizationId string, options CreateOrganizationClientOptions) (*CreateOrganizationClientResponse, error) {
	if organizationId == "" {
		return nil, ErrOrganizationIdRequired
	}
	client := buildOrganizationClient(options.Name, options.Description, options.CustomClaims, options.Audience, options.Scopes)
	return newConnectExecuter(
		m.coreClient,
		m.client.CreateOrganizationClient,
		&clientsv1.CreateOrganizationClientRequest{
			OrganizationId: organizationId,
			Client:         client,
		},
	).exec(ctx)
}

func (m *m2mService) GetOrganizationClient(ctx context.Context, organizationId string, clientId string) (*GetOrganizationClientResponse, error) {
	if organizationId == "" {
		return nil, ErrOrganizationIdRequired
	}
	if clientId == "" {
		return nil, ErrClientIdRequired
	}
	return newConnectExecuter(
		m.coreClient,
		m.client.GetOrganizationClient,
		&clientsv1.GetOrganizationClientRequest{
			OrganizationId: organizationId,
			ClientId:       clientId,
		},
	).exec(ctx)
}

func (m *m2mService) UpdateOrganizationClient(ctx context.Context, organizationId string, clientId string, options UpdateOrganizationClientOptions) (*UpdateOrganizationClientResponse, error) {
	if organizationId == "" {
		return nil, ErrOrganizationIdRequired
	}
	if clientId == "" {
		return nil, ErrClientIdRequired
	}
	client := buildOrganizationClient(options.Name, options.Description, options.CustomClaims, options.Audience, options.Scopes)
	return newConnectExecuter(
		m.coreClient,
		m.client.UpdateOrganizationClient,
		&clientsv1.UpdateOrganizationClientRequest{
			OrganizationId: organizationId,
			ClientId:       clientId,
			Client:         client,
		},
	).exec(ctx)
}

func (m *m2mService) DeleteOrganizationClient(ctx context.Context, organizationId string, clientId string) error {
	if organizationId == "" {
		return ErrOrganizationIdRequired
	}
	if clientId == "" {
		return ErrClientIdRequired
	}
	_, err := newConnectExecuter(
		m.coreClient,
		m.client.DeleteOrganizationClient,
		&clientsv1.DeleteOrganizationClientRequest{
			OrganizationId: organizationId,
			ClientId:       clientId,
		},
	).exec(ctx)
	return err
}

func (m *m2mService) CreateOrganizationClientSecret(ctx context.Context, organizationId string, clientId string) (*CreateOrganizationClientSecretResponse, error) {
	if organizationId == "" {
		return nil, ErrOrganizationIdRequired
	}
	if clientId == "" {
		return nil, ErrClientIdRequired
	}
	return newConnectExecuter(
		m.coreClient,
		m.client.CreateOrganizationClientSecret,
		&clientsv1.CreateOrganizationClientSecretRequest{
			OrganizationId: organizationId,
			ClientId:       clientId,
		},
	).exec(ctx)
}

func (m *m2mService) DeleteOrganizationClientSecret(ctx context.Context, organizationId string, clientId string, secretId string) error {
	if organizationId == "" {
		return ErrOrganizationIdRequired
	}
	if clientId == "" {
		return ErrClientIdRequired
	}
	if secretId == "" {
		return ErrSecretIdRequired
	}
	_, err := newConnectExecuter(
		m.coreClient,
		m.client.DeleteOrganizationClientSecret,
		&clientsv1.DeleteOrganizationClientSecretRequest{
			OrganizationId: organizationId,
			ClientId:       clientId,
			SecretId:       secretId,
		},
	).exec(ctx)
	return err
}

func (m *m2mService) ListOrganizationClients(ctx context.Context, organizationId string, options ListOrganizationClientsOptions) (*ListOrganizationClientsResponse, error) {
	if organizationId == "" {
		return nil, ErrOrganizationIdRequired
	}
	request := &clientsv1.ListOrganizationClientsRequest{
		OrganizationId: organizationId,
	}
	if options.PageSize > 0 {
		request.PageSize = options.PageSize
	}
	if options.PageToken != "" {
		request.PageToken = options.PageToken
	}
	return newConnectExecuter(
		m.coreClient,
		m.client.ListOrganizationClients,
		request,
	).exec(ctx)
}
