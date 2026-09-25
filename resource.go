package scalekit

import (
	"context"
	"fmt"

	clientsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients/clientsconnect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// Type aliases for public API surface
type ResourceClientInfo = clientsv1.ResourceClient
type CreateResourceClientResponse = clientsv1.CreateResourceClientResponse
type GetResourceClientResponse = clientsv1.GetResourceClientResponse
type UpdateResourceClientResponse = clientsv1.UpdateResourceClientResponse
type ListResourceClientsResponse = clientsv1.ListResourceClientsResponse
type ResourceUserConsent = clientsv1.ResourceUserConsent
type ListResourceUserConsentsResponse = clientsv1.ListResourceUserConsentsResponse
type RevokeUserConsentResponse = clientsv1.RevokeUserConsentResponse
type GetResourceResponse = clientsv1.GetResourceResponse
type ListResourcesResponse = clientsv1.ListResourcesResponse

// ResourceType identifies the kind of resource a client belongs to.
type ResourceType = clientsv1.ResourceType

// Enum constants for ResourceType.
const (
	ResourceTypeUnspecified = clientsv1.ResourceType_RESOURCE_TYPE_UNSPECIFIED
	ResourceTypeWeb         = clientsv1.ResourceType_WEB
	ResourceTypeMobile      = clientsv1.ResourceType_MOBILE
	ResourceTypeDesktop     = clientsv1.ResourceType_DESKTOP
	ResourceTypeServer      = clientsv1.ResourceType_SERVER
	ResourceTypeMcpServer   = clientsv1.ResourceType_MCP_SERVER
)

// ListUserConsentsOptions holds the filter and pagination parameters for
// listing the end-user consents granted against a resource.
type ListUserConsentsOptions struct {
	// Search is a case-insensitive substring match on external user IDs.
	// Ignored when UserIds is set.
	Search string
	// PageSize is the number of consents to return per page (max 30).
	PageSize uint32
	// PageToken is the pagination cursor.
	PageToken string
	// UserIds matches external user IDs exactly and case-sensitively, combining
	// the values with OR (max 25). Takes precedence over Search.
	UserIds []string
}

// ListResourcesOptions holds pagination parameters for listing resources.
type ListResourcesOptions struct {
	// PageSize is the page size for pagination (max 30).
	PageSize uint32
	// PageToken is the pagination cursor.
	PageToken string
}

// UpdateResourceClientOptions holds the fields to change on an existing
// resource client. Only the fields set here (non-nil) are sent to the
// server — there is no separate field mask to build; UpdateResourceClient
// derives it from whichever options are set.
type UpdateResourceClientOptions struct {
	// Name, if set, replaces the client's name. The server applies it only
	// when non-empty — an empty string is a no-op, not a clear.
	Name *string
	// Description, if set, replaces the client's description. The server
	// applies it only when non-empty — an empty string is a no-op, not a clear.
	Description *string
	// Scopes, if set, replaces the client's scopes. Set to an empty (non-nil)
	// slice to clear them.
	Scopes *[]string
	// CustomClaims, if set, replaces the client's custom claims. Set to an
	// empty (non-nil) slice to clear them.
	CustomClaims *[]*clientsv1.CustomClaim
	// Expiry, if set, replaces the access token lifetime in seconds.
	Expiry *int64
	// RedirectUris, if set, replaces the client's redirect URIs. Set to an
	// empty (non-nil) slice to clear them.
	RedirectUris *[]string
}

// ResourceService is a client for reading resources, managing the API
// clients scoped to a resource, and reading and revoking end-user consents
// granted against one.
//
// A resource (for example an MCP server) can have one or more API clients
// registered against it, each using the client_credentials OAuth flow scoped
// to that resource. A consent records that one of your end users allowed a
// specific client to act on their behalf against the resource, identified by
// ExternalUserId — the identifier your application supplied when the consent
// was granted.
type ResourceService interface {
	// GetResource retrieves a single resource by id.
	//
	// A resource client's scopes are only actually granted in an issued
	// token when they also appear in the resource's own scopes allowlist
	// (the server intersects requested scopes against the environment's
	// permissions, the resource's allowed scopes, and the client's own
	// scopes) — call this first to see what the resource actually allows
	// before creating or updating a resource client with scopes.
	//
	// The returned Resource.Scopes is every scope defined in the
	// environment, not just the ones this resource allows — each entry
	// carries an Enabled flag, and only the ones with Enabled: true are
	// actually usable on this resource. Filter on that flag to get the
	// actual allowlist.
	GetResource(ctx context.Context, resourceId string) (*GetResourceResponse, error)

	// ListResources lists resources of a given type in the environment,
	// with pagination.
	//
	// resourceType is required by the underlying API — there is no way to
	// list every type in one call; list each type separately if needed.
	ListResources(ctx context.Context, resourceType ResourceType, options ListResourcesOptions) (*ListResourcesResponse, error)

	// CreateResourceClient creates a new API client scoped to a resource.
	//
	// The response's PlainSecret is the plaintext client secret, only
	// available at creation time.
	//
	// Audience cannot be set through this SDK — it is always
	// server-determined, for any resource type. A non-empty client.Audience
	// returns ErrAudienceNotSettable rather than being silently forwarded.
	CreateResourceClient(ctx context.Context, resourceId string, client *clientsv1.ResourceClient) (*CreateResourceClientResponse, error)

	// GetResourceClient retrieves a single API client scoped to a resource,
	// along with the end-users who have granted it consent.
	GetResourceClient(ctx context.Context, resourceId string, clientId string) (*GetResourceClientResponse, error)

	// ListResourceClients lists every API client scoped to a resource.
	ListResourceClients(ctx context.Context, resourceId string) (*ListResourceClientsResponse, error)

	// UpdateResourceClient updates an existing API client scoped to a resource.
	//
	// Only the fields set on options are changed — set a field to update it,
	// leave it nil to leave it alone. There is no field mask to build
	// yourself; it's derived internally from whichever options are set.
	// Verified against a live environment: the server only actually honors
	// this for Scopes, CustomClaims and RedirectUris — set one to an empty
	// (non-nil) slice to clear it. Name/Description are applied whenever
	// non-empty regardless (an empty string is a no-op, not a clear).
	//
	// There is no Audience option — audience is always server-determined and
	// can never be set through this SDK, on create or update, for any
	// resource type.
	UpdateResourceClient(ctx context.Context, resourceId string, clientId string, options UpdateResourceClientOptions) (*UpdateResourceClientResponse, error)

	// DeleteResourceClient permanently deletes an API client scoped to a resource.
	//
	// DeleteResourceClient shares its underlying delete path with client
	// deletion in general, so nothing forces the given clientId to actually
	// belong to resourceId. Since this method lives on ResourceService,
	// callers reasonably expect it to only ever touch clients within that
	// resource — so this fetches the client first and verifies its own
	// resourceId matches before deleting, refusing instead of trusting the id
	// pair blindly.
	DeleteResourceClient(ctx context.Context, resourceId string, clientId string) error

	// CreateResourceClientSecret creates a new secret for an API client
	// scoped to a resource.
	//
	// The underlying secret-creation call is keyed by clientId alone — it
	// has no notion of a resource — so this fetches the client first and
	// verifies it belongs to resourceId before creating a secret for it,
	// the same ownership check DeleteResourceClient applies.
	//
	// The backend caps how many secrets a client can hold at once (a
	// configurable limit — 5 in Scalekit's own dev environment, verified
	// live; treat the exact number as environment-specific, not a fixed
	// constant). Exceeding it fails (the server rejects it as
	// INVALID_ARGUMENT, "only N secrets are allowed") — delete an existing
	// secret first via DeleteResourceClientSecret. The dashboard itself is
	// more conservative than the server limit: it only shows an "Add new
	// secret" action while a client has fewer than 2 secrets. Match
	// whichever threshold — the actual server limit or the dashboard's
	// stricter 2 — fits your own UX.
	CreateResourceClientSecret(ctx context.Context, resourceId string, clientId string) (*CreateClientSecretResponse, error)

	// DeleteResourceClientSecret permanently deletes a secret from an API
	// client scoped to a resource.
	//
	// Like CreateResourceClientSecret, the underlying delete call is keyed
	// by clientId alone, so this verifies the client belongs to resourceId
	// first rather than trusting the id pair blindly.
	//
	// A client must always keep at least 1 secret. Calling this on a
	// client's last remaining secret fails (the server rejects it as
	// INVALID_ARGUMENT, "at least one secret is required"). Mirror the
	// dashboard's own UX: only offer a "Revoke" action on a secret while the
	// client has more than 1.
	DeleteResourceClientSecret(ctx context.Context, resourceId string, clientId string, secretId string) error

	// ListUserConsents lists the end-user consents granted against a
	// resource, with pagination.
	//
	// Each returned consent carries Id, ExternalUserId, ClientId, ClientName,
	// Scopes and GrantedAt. The response also carries TotalSize plus
	// NextPageToken and PrevPageToken cursors.
	//
	// Set options.UserIds to match specific users exactly, or options.Search
	// for a case-insensitive substring match. When both are set, UserIds
	// wins and Search is ignored.
	ListUserConsents(ctx context.Context, resourceId string, options ListUserConsentsOptions) (*ListResourceUserConsentsResponse, error)

	// RevokeUserConsent revokes a single end-user consent held by an API client.
	//
	// It deletes the consent, so the client is prompted for consent again on
	// its next authorization attempt, and revokes every active refresh token
	// issued to that client for the same user. Access tokens already issued
	// stay valid until they expire.
	//
	// Note that clientId is the API client that holds the consent (m2m_
	// prefix), not the resource id.
	RevokeUserConsent(ctx context.Context, clientId string, consentId string) (*RevokeUserConsentResponse, error)
}

type resourceService struct {
	coreClient *coreClient
	client     clientsconnect.ClientServiceClient
}

func newResourceService(coreClient *coreClient) ResourceService {
	return &resourceService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, clientsconnect.NewClientServiceClient),
	}
}

func (r *resourceService) GetResource(ctx context.Context, resourceId string) (*GetResourceResponse, error) {
	if resourceId == "" {
		return nil, ErrResourceIdRequired
	}
	return newConnectExecuter(
		r.coreClient,
		r.client.GetResource,
		&clientsv1.GetResourceRequest{
			ResourceId: resourceId,
		},
	).exec(ctx)
}

func (r *resourceService) ListResources(ctx context.Context, resourceType ResourceType, options ListResourcesOptions) (*ListResourcesResponse, error) {
	return newConnectExecuter(
		r.coreClient,
		r.client.ListResources,
		&clientsv1.ListResourcesRequest{
			ResourceType: resourceType,
			PageSize:     options.PageSize,
			PageToken:    options.PageToken,
		},
	).exec(ctx)
}

func (r *resourceService) CreateResourceClient(ctx context.Context, resourceId string, client *clientsv1.ResourceClient) (*CreateResourceClientResponse, error) {
	if resourceId == "" {
		return nil, ErrResourceIdRequired
	}
	if client != nil && len(client.GetAudience()) > 0 {
		return nil, ErrAudienceNotSettable
	}
	return newConnectExecuter(
		r.coreClient,
		r.client.CreateResourceClient,
		&clientsv1.CreateResourceClientRequest{
			ResourceId: resourceId,
			Client:     client,
		},
	).exec(ctx)
}

func (r *resourceService) GetResourceClient(ctx context.Context, resourceId string, clientId string) (*GetResourceClientResponse, error) {
	if resourceId == "" {
		return nil, ErrResourceIdRequired
	}
	if clientId == "" {
		return nil, ErrClientIdRequired
	}
	return newConnectExecuter(
		r.coreClient,
		r.client.GetResourceClient,
		&clientsv1.GetResourceClientRequest{
			ResourceId: resourceId,
			ClientId:   clientId,
		},
	).exec(ctx)
}

func (r *resourceService) ListResourceClients(ctx context.Context, resourceId string) (*ListResourceClientsResponse, error) {
	if resourceId == "" {
		return nil, ErrResourceIdRequired
	}
	return newConnectExecuter(
		r.coreClient,
		r.client.ListResourceClients,
		&clientsv1.ListResourceClientsRequest{
			ResourceId: resourceId,
		},
	).exec(ctx)
}

func (r *resourceService) UpdateResourceClient(ctx context.Context, resourceId string, clientId string, options UpdateResourceClientOptions) (*UpdateResourceClientResponse, error) {
	if resourceId == "" {
		return nil, ErrResourceIdRequired
	}
	if clientId == "" {
		return nil, ErrClientIdRequired
	}

	var paths []string
	client := &clientsv1.ResourceClient{}
	if options.Name != nil {
		client.Name = *options.Name
		paths = append(paths, "name")
	}
	if options.Description != nil {
		client.Description = *options.Description
		paths = append(paths, "description")
	}
	if options.Scopes != nil {
		client.Scopes = *options.Scopes
		paths = append(paths, "scopes")
	}
	if options.CustomClaims != nil {
		client.CustomClaims = *options.CustomClaims
		paths = append(paths, "custom_claims")
	}
	if options.Expiry != nil {
		client.Expiry = *options.Expiry
		paths = append(paths, "expiry")
	}
	if options.RedirectUris != nil {
		client.RedirectUris = *options.RedirectUris
		paths = append(paths, "redirect_uris")
	}

	var mask *fieldmaskpb.FieldMask
	if len(paths) > 0 {
		mask = &fieldmaskpb.FieldMask{Paths: paths}
	}

	return newConnectExecuter(
		r.coreClient,
		r.client.UpdateResourceClient,
		&clientsv1.UpdateResourceClientRequest{
			ResourceId: resourceId,
			ClientId:   clientId,
			Client:     client,
			UpdateMask: mask,
		},
	).exec(ctx)
}

func (r *resourceService) DeleteResourceClient(ctx context.Context, resourceId string, clientId string) error {
	if resourceId == "" {
		return ErrResourceIdRequired
	}
	if clientId == "" {
		return ErrClientIdRequired
	}

	fetched, err := r.GetResourceClient(ctx, resourceId, clientId)
	if err != nil {
		return err
	}
	if fetched.Client == nil || fetched.Client.ResourceId != resourceId {
		return fmt.Errorf("%w: client %q does not belong to resource %q", ErrClientNotInResource, clientId, resourceId)
	}

	_, err = newConnectExecuter(
		r.coreClient,
		r.client.DeleteResourceClient,
		&clientsv1.DeleteResourceClientRequest{
			ResourceId: resourceId,
			ClientId:   clientId,
		},
	).exec(ctx)
	return err
}

func (r *resourceService) CreateResourceClientSecret(ctx context.Context, resourceId string, clientId string) (*CreateClientSecretResponse, error) {
	if resourceId == "" {
		return nil, ErrResourceIdRequired
	}
	if clientId == "" {
		return nil, ErrClientIdRequired
	}

	fetched, err := r.GetResourceClient(ctx, resourceId, clientId)
	if err != nil {
		return nil, err
	}
	if fetched.Client == nil || fetched.Client.ResourceId != resourceId {
		return nil, fmt.Errorf("%w: client %q does not belong to resource %q", ErrClientNotInResource, clientId, resourceId)
	}

	return newConnectExecuter(
		r.coreClient,
		r.client.CreateClientSecret,
		&clientsv1.CreateClientSecretRequest{
			ClientId: clientId,
		},
	).exec(ctx)
}

func (r *resourceService) DeleteResourceClientSecret(ctx context.Context, resourceId string, clientId string, secretId string) error {
	if resourceId == "" {
		return ErrResourceIdRequired
	}
	if clientId == "" {
		return ErrClientIdRequired
	}
	if secretId == "" {
		return ErrSecretIdRequired
	}

	fetched, err := r.GetResourceClient(ctx, resourceId, clientId)
	if err != nil {
		return err
	}
	if fetched.Client == nil || fetched.Client.ResourceId != resourceId {
		return fmt.Errorf("%w: client %q does not belong to resource %q", ErrClientNotInResource, clientId, resourceId)
	}

	_, err = newConnectExecuter(
		r.coreClient,
		r.client.DeleteClientSecret,
		&clientsv1.DeleteClientSecretRequest{
			ClientId: clientId,
			SecretId: secretId,
		},
	).exec(ctx)
	return err
}

// ListUserConsents lists the end-user consents granted against a resource,
// with pagination.
//
// Each returned consent carries Id, ExternalUserId, ClientId, ClientName,
// Scopes and GrantedAt. The response also carries TotalSize plus NextPageToken
// and PrevPageToken cursors.
//
// Set options.UserIds to match specific users exactly, or options.Search for a
// case-insensitive substring match. When both are set, UserIds wins and Search
// is ignored.
func (r *resourceService) ListUserConsents(ctx context.Context, resourceId string, options ListUserConsentsOptions) (*ListResourceUserConsentsResponse, error) {
	if resourceId == "" {
		return nil, ErrResourceIdRequired
	}
	request := &clientsv1.ListResourceUserConsentsRequest{
		ResourceId: resourceId,
	}
	if options.Search != "" {
		request.Search = options.Search
	}
	if options.PageSize != 0 {
		request.PageSize = options.PageSize
	}
	if options.PageToken != "" {
		request.PageToken = options.PageToken
	}
	// The filter takes precedence over search server-side, so only attach it when
	// the caller actually supplied user IDs — an empty filter would otherwise
	// suppress a search the caller did supply.
	if len(options.UserIds) > 0 {
		request.Filter = &clientsv1.ResourceUserConsentFilter{
			ExternalUserId: options.UserIds,
		}
	}
	return newConnectExecuter(
		r.coreClient,
		r.client.ListResourceUserConsents,
		request,
	).exec(ctx)
}

// RevokeUserConsent revokes a single end-user consent held by an API client.
//
// It deletes the consent, so the client is prompted for consent again on its
// next authorization attempt, and revokes every active refresh token issued to
// that client for the same user. Access tokens already issued stay valid until
// they expire.
//
// Note that clientId is the API client that holds the consent (m2m_ prefix),
// not the resource id.
func (r *resourceService) RevokeUserConsent(ctx context.Context, clientId string, consentId string) (*RevokeUserConsentResponse, error) {
	if clientId == "" {
		return nil, ErrClientIdRequired
	}
	if consentId == "" {
		return nil, ErrConsentIdRequired
	}
	return newConnectExecuter(
		r.coreClient,
		r.client.RevokeUserConsent,
		&clientsv1.RevokeUserConsentRequest{
			ClientId:  clientId,
			ConsentId: consentId,
		},
	).exec(ctx)
}
