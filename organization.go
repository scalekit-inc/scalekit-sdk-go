package scalekit

import (
	"context"
	"errors"

	commonsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/commons"
	organizationsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/organizations"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/organizations/organizationsconnect"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ListOrganizationsResponse = organizationsv1.ListOrganizationsResponse
type GetOrganizationResponse = organizationsv1.GetOrganizationResponse
type CreateOrganizationResponse = organizationsv1.CreateOrganizationResponse
type UpdateOrganizationResponse = organizationsv1.UpdateOrganizationResponse
type Link = organizationsv1.Link
type UpdateOrganization = organizationsv1.UpdateOrganization
type ListOrganizationOptions = organizationsv1.ListOrganizationsRequest
type SearchOrganizationsResponse = organizationsv1.SearchOrganizationsResponse
type GetOrganizationSessionPolicyResponse = organizationsv1.GetOrganizationSessionPolicyResponse
type UpdateOrganizationSessionPolicyResponse = organizationsv1.UpdateOrganizationSessionPolicyResponse
type GetApplicationSessionPolicyResponse = organizationsv1.GetApplicationSessionPolicyResponse
type GetOrganizationUserManagementSettingsResponse = organizationsv1.GetOrganizationUserManagementSettingsResponse

type OrganizationSettings struct {
	Features []Feature
}
type Feature struct {
	Name    string
	Enabled bool
}
type OrganizationUserManagementSettings struct {
	MaxAllowedUsers *int32
}

type CreateOrganizationOptions struct {
	ExternalId string
	Metadata   map[string]string
}

// OrganizationSessionPolicyOption configures optional fields on UpdateOrganizationSessionPolicyRequest.
type OrganizationSessionPolicyOption func(*organizationsv1.UpdateOrganizationSessionPolicyRequest)

// WithAbsoluteSessionTimeout sets the absolute session timeout value and unit.
func WithAbsoluteSessionTimeout(timeout int32, unit commonsv1.TimeUnit) OrganizationSessionPolicyOption {
	return func(req *organizationsv1.UpdateOrganizationSessionPolicyRequest) {
		req.AbsoluteSessionTimeout = wrapperspb.Int32(timeout)
		req.AbsoluteSessionTimeoutUnit = &unit
	}
}

// WithIdleSessionTimeout sets the idle session timeout value, unit, and whether it is enabled.
func WithIdleSessionTimeout(enabled bool, timeout int32, unit commonsv1.TimeUnit) OrganizationSessionPolicyOption {
	return func(req *organizationsv1.UpdateOrganizationSessionPolicyRequest) {
		req.IdleSessionTimeoutEnabled = wrapperspb.Bool(enabled)
		req.IdleSessionTimeout = wrapperspb.Int32(timeout)
		req.IdleSessionTimeoutUnit = &unit
	}
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
	UpsertUserManagementSettings(ctx context.Context, organizationId string, settings OrganizationUserManagementSettings) (*organizationsv1.OrganizationUserManagementSettings, error)
	SearchOrganization(ctx context.Context, query string, pageSize uint32, pageToken string) (*SearchOrganizationsResponse, error)
	GetOrganizationUserManagementSetting(ctx context.Context, organizationId string) (*GetOrganizationUserManagementSettingsResponse, error)
	UpdateOrganizationSessionPolicy(ctx context.Context, organizationId string, policySource organizationsv1.SessionPolicyType, opts ...OrganizationSessionPolicyOption) (*UpdateOrganizationSessionPolicyResponse, error)
	GetOrganizationSessionPolicy(ctx context.Context, organizationId string) (*GetOrganizationSessionPolicyResponse, error)
	GetApplicationSessionPolicy(ctx context.Context, organizationId string) (*GetApplicationSessionPolicyResponse, error)
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
	req := &organizationsv1.CreateOrganizationRequest{
		Organization: &organizationsv1.CreateOrganization{
			DisplayName: name,
			Metadata:    options.Metadata,
		},
	}
	if options.ExternalId != "" {
		req.Organization.ExternalId = &options.ExternalId
	}
	return newConnectExecuter(
		o.coreClient,
		o.client.CreateOrganization,
		req,
	).exec(ctx)
}

func (o *organization) ListOrganization(ctx context.Context, request *ListOrganizationOptions) (*ListOrganizationsResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.ListOrganization,
		request,
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
	if err != nil {
		return nil, err
	}
	if resp.Link == nil {
		return nil, errors.New("generate portal link: response missing link")
	}
	return resp.Link, nil
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

func (o *organization) UpsertUserManagementSettings(ctx context.Context, organizationId string, settings OrganizationUserManagementSettings) (*organizationsv1.OrganizationUserManagementSettings, error) {
	request := &organizationsv1.UpsertUserManagementSettingsRequest{
		OrganizationId: organizationId,
		Settings:       &organizationsv1.OrganizationUserManagementSettings{},
	}

	if settings.MaxAllowedUsers != nil {
		request.Settings.MaxAllowedUsers = wrapperspb.Int32(*settings.MaxAllowedUsers)
	}

	resp, err := newConnectExecuter(
		o.coreClient,
		o.client.UpsertUserManagementSettings,
		request,
	).exec(ctx)
	if err != nil {
		return nil, err
	}
	if resp.Settings == nil {
		return nil, errors.New("upsert user management settings: response missing settings")
	}
	return resp.Settings, nil
}

func (o *organization) SearchOrganization(ctx context.Context, query string, pageSize uint32, pageToken string) (*SearchOrganizationsResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.SearchOrganization,
		&organizationsv1.SearchOrganizationsRequest{
			Query:     query,
			PageSize:  pageSize,
			PageToken: pageToken,
		},
	).exec(ctx)
}

func (o *organization) GetOrganizationUserManagementSetting(ctx context.Context, organizationId string) (*GetOrganizationUserManagementSettingsResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.GetOrganizationUserManagementSetting,
		&organizationsv1.GetOrganizationUserManagementSettingsRequest{
			OrganizationId: organizationId,
		},
	).exec(ctx)
}

func (o *organization) UpdateOrganizationSessionPolicy(ctx context.Context, organizationId string, policySource organizationsv1.SessionPolicyType, opts ...OrganizationSessionPolicyOption) (*UpdateOrganizationSessionPolicyResponse, error) {
	request := &organizationsv1.UpdateOrganizationSessionPolicyRequest{
		OrganizationId: organizationId,
		PolicySource:   policySource,
	}
	for _, opt := range opts {
		opt(request)
	}
	return newConnectExecuter(
		o.coreClient,
		o.client.UpdateOrganizationSessionPolicy,
		request,
	).exec(ctx)
}

func (o *organization) GetOrganizationSessionPolicy(ctx context.Context, organizationId string) (*GetOrganizationSessionPolicyResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.GetOrganizationSessionPolicy,
		&organizationsv1.GetOrganizationSessionPolicyRequest{
			OrganizationId: organizationId,
		},
	).exec(ctx)
}

func (o *organization) GetApplicationSessionPolicy(ctx context.Context, organizationId string) (*GetApplicationSessionPolicyResponse, error) {
	return newConnectExecuter(
		o.coreClient,
		o.client.GetApplicationSessionPolicy,
		&organizationsv1.GetApplicationSessionPolicyRequest{
			OrganizationId: organizationId,
		},
	).exec(ctx)
}
