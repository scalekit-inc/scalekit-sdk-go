package scalekit

import (
	"context"
	"fmt"
	"slices"

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
type ListResourceUserConsentsResponse = clientsv1.ListResourceUserConsentsResponse
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

// ListUserConsentsOptions holds pagination and search parameters for listing
// end-user consents granted against a resource.
type ListUserConsentsOptions struct {
	// Search is a case-insensitive substring match on external user IDs.
	Search string
	// PageSize is the page size for pagination (max 30).
	PageSize uint32
	// PageToken is the pagination cursor.
	PageToken string
}

// ListResourcesOptions holds pagination parameters for listing resources.
type ListResourcesOptions struct {
	// PageSize is the page size for pagination (max 30).
	PageSize uint32
	// PageToken is the pagination cursor.
	PageToken string
}

// ResourceService is a client for reading resources, managing the API
// clients scoped to a resource, and reading and revoking end-user consents
// granted against one.
//
// A resource (for example an MCP server) can have one or more API clients
// registered against it, each using the client_credentials OAuth flow scoped
// to that resource. A consent records that one of your end users allowed a
// specific client to act on their behalf against the resource.
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
	// available at creation time. Audience is ignored for MCP_SERVER/
	// MCP_GATEWAY resources, which get their audience from the resource itself.
	CreateResourceClient(ctx context.Context, resourceId string, client *clientsv1.ResourceClient) (*CreateResourceClientResponse, error)

	// GetResourceClient retrieves a single API client scoped to a resource,
	// along with the end-users who have granted it consent.
	GetResourceClient(ctx context.Context, resourceId string, clientId string) (*GetResourceClientResponse, error)

	// ListResourceClients lists every API client scoped to a resource.
	ListResourceClients(ctx context.Context, resourceId string) (*ListResourceClientsResponse, error)

	// UpdateResourceClient updates an existing API client scoped to a resource.
	//
	// mask lists which fields of client to change (e.g. &fieldmaskpb.FieldMask{
	// Paths: []string{"scopes", "custom_claims"}}). Verified against a live
	// environment: the server only actually honors the mask for scopes,
	// custom_claims and redirect_uris — include one of those paths with an
	// empty value (e.g. Scopes: []string{}) to clear it. Name/Description are
	// applied whenever non-empty regardless of mask (an empty string is a
	// no-op, not a clear).
	//
	// "audience" is not a supported mask path — a resource client's audience
	// is fixed at creation and can never be changed via update, for any
	// resource type, so this returns ErrAudienceNotUpdatable rather than
	// silently accepting a path that can never take effect.
	UpdateResourceClient(ctx context.Context, resourceId string, clientId string, client *clientsv1.ResourceClient, mask *fieldmaskpb.FieldMask) (*UpdateResourceClientResponse, error)

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
	// Each returned consent carries ConsentId, ExternalUserId, and Scopes.
	// The response also carries TotalSize plus NextPageToken/PrevPageToken
	// cursors.
	ListUserConsents(ctx context.Context, resourceId string, options ListUserConsentsOptions) (*ListResourceUserConsentsResponse, error)

	// RevokeUserConsent revokes a single end-user consent held by an API client.
	//
	// Deletes the consent, so the client is prompted for consent again on
	// its next authorization attempt, and revokes every active refresh token
	// issued to that client for the same user. Access tokens already issued
	// stay valid until they expire.
	//
	// Note that clientId is the API client that holds the consent, not the
	// resource id.
	RevokeUserConsent(ctx context.Context, clientId string, consentId string) error
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

func (r *resourceService) UpdateResourceClient(ctx context.Context, resourceId string, clientId string, client *clientsv1.ResourceClient, mask *fieldmaskpb.FieldMask) (*UpdateResourceClientResponse, error) {
	if resourceId == "" {
		return nil, ErrResourceIdRequired
	}
	if clientId == "" {
		return nil, ErrClientIdRequired
	}
	if mask != nil && slices.Contains(mask.GetPaths(), "audience") {
		return nil, ErrAudienceNotUpdatable
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

func (r *resourceService) ListUserConsents(ctx context.Context, resourceId string, options ListUserConsentsOptions) (*ListResourceUserConsentsResponse, error) {
	if resourceId == "" {
		return nil, ErrResourceIdRequired
	}
	return newConnectExecuter(
		r.coreClient,
		r.client.ListResourceUserConsents,
		&clientsv1.ListResourceUserConsentsRequest{
			ResourceId: resourceId,
			Search:     options.Search,
			PageSize:   options.PageSize,
			PageToken:  options.PageToken,
		},
	).exec(ctx)
}

func (r *resourceService) RevokeUserConsent(ctx context.Context, clientId string, consentId string) error {
	if clientId == "" {
		return ErrClientIdRequired
	}
	if consentId == "" {
		return ErrConsentIdRequired
	}
	_, err := newConnectExecuter(
		r.coreClient,
		r.client.RevokeUserConsent,
		&clientsv1.RevokeUserConsentRequest{
			ClientId:  clientId,
			ConsentId: consentId,
		},
	).exec(ctx)
	return err
}
