package scalekit

import (
	"context"

	organizationsv1 "github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/organizations"
	"github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/organizations/organizationsconnect"
)

type ListOrganizationsResponse = organizationsv1.ListOrganizationsResponse
type GetOrganizationResponse = organizationsv1.GetOrganizationResponse
type CreateOrganizationResponse = organizationsv1.CreateOrganizationResponse
type UpdateOrganizationResponse = organizationsv1.UpdateOrganizationResponse
type Link = organizationsv1.Link
type UpdateOrganization = organizationsv1.UpdateOrganization
type ListOrganizationOptions = organizationsv1.ListOrganizationsRequest
type OrganizationSettings struct {
	Features []Feature
}
type Feature struct {
	Name    string
	Enabled bool
}

type CreateOrganizationOptions struct {
	ExternalId string
}

type Organization interface {
	CreateOrganization(ctx context.Context, name string, options CreateOrganizationOptions) (*CreateOrganizationResponse, error)
	ListOrganization(ctx context.Context, options *ListOrganizationOptions) (*ListOrganizationsResponse, error)
	GetOrganization(ctx context.Context, id string) (*GetOrganizationResponse, error)
	GetOrganizationByExternalId(ctx context.Context, externalId string) (*GetOrganizationResponse, error)
	UpdateOrganization(ctx context.Context, id string, organization *UpdateOrganization) (*UpdateOrganizationResponse, error)
	UpdateOrganizationByExternalId(ctx context.Context, externalId string, organization *UpdateOrganization) (*UpdateOrganizationResponse, error)
	DeleteOrganization(ctx context.Context, id string) error
	GeneratePortalLink(ctx context.Context, organizationId string) (*Link, error)
	UpdateOrganizationSettings(ctx context.Context, id string, settings OrganizationSettings) (*GetOrganizationResponse, error)
}

type organization struct {
	coreClient *coreClient
	client     organizationsconnect.OrganizationServiceClient
}

func newOrganizationClient(coreClient *coreClient) Organization {
	return &organization{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, organizationsconnect.NewOrganizationServiceClient),
	}
}

func (o *organization) CreateOrganization(ctx context.Context, name string, options CreateOrganizationOptions) (*CreateOrganizationResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.CreateOrganization,
		&organizationsv1.CreateOrganizationRequest{
			Organization: &organizationsv1.CreateOrganization{
				DisplayName: name,
				ExternalId:  &options.ExternalId,
			},
		},
	).exec(ctx)
}

func (o *organization) ListOrganization(ctx context.Context, options *ListOrganizationOptions) (*ListOrganizationsResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.ListOrganization,
		&organizationsv1.ListOrganizationsRequest{
			PageSize:  options.PageSize,
			PageToken: options.PageToken,
		},
	).exec(ctx)
}

func (o *organization) GetOrganization(ctx context.Context, id string) (*GetOrganizationResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.GetOrganization,
		&organizationsv1.GetOrganizationRequest{
			Identities: &organizationsv1.GetOrganizationRequest_Id{
				Id: id,
			},
		},
	).exec(ctx)
}

func (o *organization) GetOrganizationByExternalId(ctx context.Context, externalId string) (*GetOrganizationResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.GetOrganization,
		&organizationsv1.GetOrganizationRequest{
			Identities: &organizationsv1.GetOrganizationRequest_ExternalId{
				ExternalId: externalId,
			},
		},
	).exec(ctx)
}

func (o *organization) UpdateOrganization(ctx context.Context, id string, organization *UpdateOrganization) (*UpdateOrganizationResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.UpdateOrganization,
		&organizationsv1.UpdateOrganizationRequest{
			Identities: &organizationsv1.UpdateOrganizationRequest_Id{
				Id: id,
			},
			Organization: organization,
		},
	).exec(ctx)
}

func (o *organization) UpdateOrganizationByExternalId(ctx context.Context, externalId string, organization *UpdateOrganization) (*UpdateOrganizationResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.UpdateOrganization,
		&organizationsv1.UpdateOrganizationRequest{
			Identities: &organizationsv1.UpdateOrganizationRequest_ExternalId{
				ExternalId: externalId,
			},
			Organization: organization,
		},
	).exec(ctx)
}

func (o *organization) DeleteOrganization(ctx context.Context, id string) error {
	_, err := newConnectExecuter(
		o.coreClient,
		o.client.DeleteOrganization,
		&organizationsv1.DeleteOrganizationRequest{
			Identities: &organizationsv1.DeleteOrganizationRequest_Id{
				Id: id,
			},
		},
	).exec(ctx)

	return err
}

func (o *organization) GeneratePortalLink(ctx context.Context, organizationId string) (*Link, error) {
	resp, err := newConnectExecuter(
		o.coreClient,
		o.client.GeneratePortalLink,
		&organizationsv1.GeneratePortalLinkRequest{
			Id: organizationId,
		},
	).exec(ctx)

	return resp.Link, err
}

func (o *organization) UpdateOrganizationSettings(ctx context.Context, id string, settings OrganizationSettings) (*GetOrganizationResponse, error) {
	request := &organizationsv1.UpdateOrganizationSettingsRequest{
		Id: id,
		Settings: &organizationsv1.OrganizationSettings{
			Features: []*organizationsv1.OrganizationSettingsFeature{},
		},
	}
	for _, feature := range settings.Features {
		request.Settings.Features = append(request.Settings.Features, &organizationsv1.OrganizationSettingsFeature{
			Name:    feature.Name,
			Enabled: feature.Enabled,
		})
	}

	return newConnectExecuter(
		o.coreClient,
		o.client.UpdateOrganizationSettings,
		request,
	).exec(ctx)
}
