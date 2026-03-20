package scalekit

import (
	"context"

	clientsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients/clientsconnect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type ListClientsResponse = clientsv1.ListClientsResponse
type GetClientResponse = clientsv1.GetClientResponse
type CreateClientResponse = clientsv1.CreateClientResponse
type UpdateClientResponse = clientsv1.UpdateClientResponse
type CreateClientSecretResponse = clientsv1.CreateClientSecretResponse
type ListClientsOptions = clientsv1.ListClientsRequest

type ClientService interface {
	CreateClient(ctx context.Context, client *clientsv1.CreateClient) (*CreateClientResponse, error)
	GetClient(ctx context.Context, clientId string) (*GetClientResponse, error)
	ListClients(ctx context.Context, options *ListClientsOptions) (*ListClientsResponse, error)
	UpdateClient(ctx context.Context, clientId string, client *clientsv1.UpdateClient, mask *fieldmaskpb.FieldMask) (*UpdateClientResponse, error)
	CreateClientSecret(ctx context.Context, clientId string) (*CreateClientSecretResponse, error)
	DeleteClientSecret(ctx context.Context, clientId string, secretId string) error
	DeleteClient(ctx context.Context, clientId string) error
}

type clientService struct {
	coreClient *coreClient
	client     clientsconnect.ClientServiceClient
}

func newClientService(coreClient *coreClient) ClientService {
	return &clientService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, clientsconnect.NewClientServiceClient),
	}
}

func (c *clientService) CreateClient(ctx context.Context, client *clientsv1.CreateClient) (*CreateClientResponse, error) {
	return newConnectExecuter(
		c.coreClient,
		c.client.CreateClient,
		&clientsv1.CreateClientRequest{
			Client: client,
		},
	).exec(ctx)
}

func (c *clientService) GetClient(ctx context.Context, clientId string) (*GetClientResponse, error) {
	return newConnectExecuter(
		c.coreClient,
		c.client.GetClient,
		&clientsv1.GetClientRequest{
			ClientId: clientId,
		},
	).exec(ctx)
}

func (c *clientService) ListClients(ctx context.Context, options *ListClientsOptions) (*ListClientsResponse, error) {
	request := &clientsv1.ListClientsRequest{}
	if options != nil {
		request = options
	}

	return newConnectExecuter(
		c.coreClient,
		c.client.ListClient,
		request,
	).exec(ctx)
}

func (c *clientService) UpdateClient(ctx context.Context, clientId string, client *clientsv1.UpdateClient, mask *fieldmaskpb.FieldMask) (*UpdateClientResponse, error) {
	return newConnectExecuter(
		c.coreClient,
		c.client.UpdateClient,
		&clientsv1.UpdateClientRequest{
			ClientId: clientId,
			Client:   client,
			Mask:     mask,
		},
	).exec(ctx)
}

func (c *clientService) DeleteClient(ctx context.Context, clientId string) error {
	_, err := newConnectExecuter(
		c.coreClient,
		c.client.DeleteClient,
		&clientsv1.DeleteClientRequest{
			ClientId: clientId,
		},
	).exec(ctx)
	return err
}

func (c *clientService) CreateClientSecret(ctx context.Context, clientId string) (*CreateClientSecretResponse, error) {
	return newConnectExecuter(
		c.coreClient,
		c.client.CreateClientSecret,
		&clientsv1.CreateClientSecretRequest{
			ClientId: clientId,
		},
	).exec(ctx)
}

func (c *clientService) DeleteClientSecret(ctx context.Context, clientId string, secretId string) error {
	_, err := newConnectExecuter(
		c.coreClient,
		c.client.DeleteClientSecret,
		&clientsv1.DeleteClientSecretRequest{
			ClientId: clientId,
			SecretId: secretId,
		},
	).exec(ctx)
	return err
}
