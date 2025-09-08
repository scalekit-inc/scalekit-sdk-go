package scalekit

import (
	"context"

	domainsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/domains"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/domains/domainsconnect"
)

type ListDomainResponse = domainsv1.ListDomainResponse
type GetDomainResponse = domainsv1.GetDomainResponse
type CreateDomainResponse = domainsv1.CreateDomainResponse

// DomainType is defined as a string type alias
type DomainType = string

// Domain type constants
const (
	DomainTypeUnspecified  DomainType = "DOMAIN_TYPE_UNSPECIFIED"
	DomainTypeAllowedEmail DomainType = "ALLOWED_EMAIL_DOMAIN"
	DomainTypeOrganization DomainType = "ORGANIZATION_DOMAIN"
)

// CreateDomainOptions represents optional parameters for creating a domain
type CreateDomainOptions struct {
	DomainType DomainType
}

type Domain interface {
	CreateDomain(ctx context.Context, organizationId, name string, options ...*CreateDomainOptions) (*CreateDomainResponse, error)
	GetDomain(ctx context.Context, id string, organizationId string) (*GetDomainResponse, error)
	ListDomains(ctx context.Context, organizationId string) (*ListDomainResponse, error)
	DeleteDomain(ctx context.Context, id string, organizationId string) error
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

func (d *domain) CreateDomain(ctx context.Context, organizationId, name string, options ...*CreateDomainOptions) (*CreateDomainResponse, error) {
	createDomain := &domainsv1.CreateDomain{
		Domain: name,
	}

	// Handle optional options - backward compatible
	if len(options) > 0 && options[0] != nil && options[0].DomainType != "" {
		// Simple map lookup for conversion (more efficient than switch)
		domainTypeMap := map[string]domainsv1.DomainType{
			"ALLOWED_EMAIL_DOMAIN":    domainsv1.DomainType_ALLOWED_EMAIL_DOMAIN,
			"ORGANIZATION_DOMAIN":     domainsv1.DomainType_ORGANIZATION_DOMAIN,
			"DOMAIN_TYPE_UNSPECIFIED": domainsv1.DomainType_DOMAIN_TYPE_UNSPECIFIED,
		}

		if domainType, exists := domainTypeMap[options[0].DomainType]; exists {
			createDomain.DomainType = domainType
		} else {
			createDomain.DomainType = domainsv1.DomainType_DOMAIN_TYPE_UNSPECIFIED
		}
	}

	return newConnectExecuter(
		d.coreClient,
		d.client.CreateDomain,
		&domainsv1.CreateDomainRequest{
			Identities: &domainsv1.CreateDomainRequest_OrganizationId{
				OrganizationId: organizationId,
			},
			Domain: createDomain,
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

func (d *domain) DeleteDomain(ctx context.Context, id string, organizationId string) error {
	_, err := newConnectExecuter(
		d.coreClient,
		d.client.DeleteDomain,
		&domainsv1.DeleteDomainRequest{
			Id: id,
			Identities: &domainsv1.DeleteDomainRequest_OrganizationId{
				OrganizationId: organizationId,
			},
		},
	).exec(ctx)
	return err
}
