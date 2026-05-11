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

// SessionPolicySource indicates whether an organization uses its own policy or inherits the application default.
type SessionPolicySource = organizationsv1.SessionPolicyType

const (
	SessionPolicySourceApplication = organizationsv1.SessionPolicyType_APPLICATION
	SessionPolicySourceCustom      = organizationsv1.SessionPolicyType_CUSTOM
)

// TimeUnit for session timeout fields accepted in UpdateOrganizationSessionPolicy.
type TimeUnit = commonsv1.TimeUnit

const (
	TimeUnitMinutes = commonsv1.TimeUnit_MINUTES
	TimeUnitHours   = commonsv1.TimeUnit_HOURS
	TimeUnitDays    = commonsv1.TimeUnit_DAYS
)

// OrganizationSessionPolicy is the input type for UpdateOrganizationSessionPolicy.
// Set PolicySource to SessionPolicySourceApplication to revert the organization to application defaults.
// Set PolicySource to SessionPolicySourceCustom and supply timeout values to activate a custom policy.
type OrganizationSessionPolicy struct {
	PolicySource               SessionPolicySource
	AbsoluteSessionTimeout     *int32
	AbsoluteSessionTimeoutUnit TimeUnit
	IdleSessionTimeoutEnabled  *bool
	IdleSessionTimeout         *int32
	IdleSessionTimeoutUnit     TimeUnit
}

// OrganizationSessionPolicySettings is the response type for session policy operations.
type OrganizationSessionPolicySettings = organizationsv1.OrganizationSessionPolicySettings

type CreateOrganizationOptions struct {
	ExternalId string
	Metadata   map[string]string
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
	GetOrganizationSessionPolicy(ctx context.Context, organizationId string) (*OrganizationSessionPolicySettings, error)
	UpdateOrganizationSessionPolicy(ctx context.Context, organizationId string, policy OrganizationSessionPolicy) (*OrganizationSessionPolicySettings, error)
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

func (o *organization) GetOrganizationSessionPolicy(ctx context.Context, organizationId string) (*OrganizationSessionPolicySettings, error) {
	resp, err := newConnectExecuter(
		o.coreClient,
		o.client.GetOrganizationSessionPolicy,
		&organizationsv1.GetOrganizationSessionPolicyRequest{
			OrganizationId: organizationId,
		},
	).exec(ctx)
	if err != nil {
		return nil, err
	}

	return resp.Policy, nil
}

func (o *organization) UpdateOrganizationSessionPolicy(ctx context.Context, organizationId string, policy OrganizationSessionPolicy) (*OrganizationSessionPolicySettings, error) {
	req := &organizationsv1.UpdateOrganizationSessionPolicyRequest{
		OrganizationId: organizationId,
		PolicySource:   policy.PolicySource,
	}
	if policy.AbsoluteSessionTimeout != nil {
		req.AbsoluteSessionTimeout = wrapperspb.Int32(*policy.AbsoluteSessionTimeout)
	}
	if policy.AbsoluteSessionTimeoutUnit != commonsv1.TimeUnit_SESSION_TIME_UNIT_UNSPECIFIED {
		u := policy.AbsoluteSessionTimeoutUnit
		req.AbsoluteSessionTimeoutUnit = &u
	}
	if policy.IdleSessionTimeoutEnabled != nil {
		req.IdleSessionTimeoutEnabled = wrapperspb.Bool(*policy.IdleSessionTimeoutEnabled)
	}
	if policy.IdleSessionTimeout != nil {
		req.IdleSessionTimeout = wrapperspb.Int32(*policy.IdleSessionTimeout)
	}
	if policy.IdleSessionTimeoutUnit != commonsv1.TimeUnit_SESSION_TIME_UNIT_UNSPECIFIED {
		u := policy.IdleSessionTimeoutUnit
		req.IdleSessionTimeoutUnit = &u
	}

	resp, err := newConnectExecuter(
		o.coreClient,
		o.client.UpdateOrganizationSessionPolicy,
		req,
	).exec(ctx)
	if err != nil {
		return nil, err
	}

	return resp.Policy, nil
}
