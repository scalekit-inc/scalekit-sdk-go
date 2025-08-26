package scalekit

import (
	"context"

	domainsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/domains"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/domains/domainsconnect"
)

type ListDomainResponse = domainsv1.ListDomainResponse
type GetDomainResponse = domainsv1.GetDomainResponse
type CreateDomainResponse = domainsv1.CreateDomainResponse

type Domain interface {
	CreateDomain(ctx context.Context, organizationId, name string) (*CreateDomainResponse, error)
	GetDomain(ctx context.Context, id string, organizationId string) (*GetDomainResponse, error)
	ListDomains(ctx context.Context, organizationId string) (*ListDomainResponse, error)
}

type domain struct {
	coreClient *coreClient
	client     domainsconnect.DomainServiceClient
}

func newDomainClient(coreClient *coreClient) Domain {
	return &domain{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, domainsconnect.NewDomainServiceClient),
	}
}

func (d *domain) CreateDomain(ctx context.Context, organizationId, name string) (*CreateDomainResponse, error) {
	return newConnectExecuter(
		d.coreClient,
		d.client.CreateDomain,
		&domainsv1.CreateDomainRequest{
			Identities: &domainsv1.CreateDomainRequest_OrganizationId{
				OrganizationId: organizationId,
			},
			Domain: &domainsv1.CreateDomain{
				Domain: name,
			},
		},
	).exec(ctx)
}

func (d *domain) GetDomain(ctx context.Context, id string, organizationId string) (*GetDomainResponse, error) {
	return newConnectExecuter(
		d.coreClient,
		d.client.GetDomain,
		&domainsv1.GetDomainRequest{
			Id: id,
			Identities: &domainsv1.GetDomainRequest_OrganizationId{
				OrganizationId: organizationId,
			},
		},
	).exec(ctx)
}

func (d *domain) ListDomains(ctx context.Context, organizationId string) (*ListDomainResponse, error) {
	return newConnectExecuter(
		d.coreClient,
		d.client.ListDomains,
		&domainsv1.ListDomainRequest{
			Identities: &domainsv1.ListDomainRequest_OrganizationId{
				OrganizationId: organizationId,
			},
		},
	).exec(ctx)
}
