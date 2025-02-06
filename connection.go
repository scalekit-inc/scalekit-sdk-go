package scalekit

import (
	"context"
	connectionsv1 "github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/connections"
	"github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/connections/connectionsconnect"
)

type ListConnectionsResponse = connectionsv1.ListConnectionsResponse
type GetConnectionResponse = connectionsv1.GetConnectionResponse
type ToggleConnectionResponse = connectionsv1.ToggleConnectionResponse

type Connection interface {
	GetConnection(ctx context.Context, organizationId string, id string) (*GetConnectionResponse, error)
	ListConnectionsByDomain(ctx context.Context, domain string) (*ListConnectionsResponse, error)
	ListConnections(ctx context.Context, organizationId string) (*ListConnectionsResponse, error)
	EnableConnection(ctx context.Context, organizationId string, id string) (*ToggleConnectionResponse, error)
	DisableConnection(ctx context.Context, organizationId string, id string) (*ToggleConnectionResponse, error)
}

type connection struct {
	coreClient *coreClient
	client     connectionsconnect.ConnectionServiceClient
}

func newConnectionClient(coreClient *coreClient) Connection {
	return &connection{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, connectionsconnect.NewConnectionServiceClient),
	}
}

func (c *connection) GetConnection(ctx context.Context, organizationId string, id string) (*GetConnectionResponse, error) {
	return newConnectExecuter(
		c.coreClient,
		c.client.GetConnection,
		&connectionsv1.GetConnectionRequest{
			Id:             id,
			OrganizationId: organizationId,
		},
	).exec(ctx)
}

func (c *connection) ListConnectionsByDomain(ctx context.Context, domain string) (*ListConnectionsResponse, error) {
	include := "all"
	return newConnectExecuter(
		c.coreClient,
		c.client.ListConnections,
		&connectionsv1.ListConnectionsRequest{
			Domain:  &domain,
			Include: &include,
		},
	).exec(ctx)
}

func (c *connection) ListConnections(ctx context.Context, organizationId string) (*ListConnectionsResponse, error) {
	include := "all"
	return newConnectExecuter(
		c.coreClient,
		c.client.ListConnections,
		&connectionsv1.ListConnectionsRequest{
			OrganizationId: &organizationId,
			Include:        &include,
		},
	).exec(ctx)
}

func (c *connection) EnableConnection(ctx context.Context, organizationId string, id string) (*ToggleConnectionResponse, error) {
	return newConnectExecuter(
		c.coreClient,
		c.client.EnableConnection,
		&connectionsv1.ToggleConnectionRequest{
			Id:             id,
			OrganizationId: organizationId,
		},
	).exec(ctx)
}

func (c *connection) DisableConnection(ctx context.Context, organizationId string, id string) (*ToggleConnectionResponse, error) {
	return newConnectExecuter(
		c.coreClient,
		c.client.DisableConnection,
		&connectionsv1.ToggleConnectionRequest{
			Id:             id,
			OrganizationId: organizationId,
		},
	).exec(ctx)
}
