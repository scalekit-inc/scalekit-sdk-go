package scalekit

import (
	"context"

	clientsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients/clientsconnect"
)

// Type aliases for public API surface
type ResourceUserConsent = clientsv1.ResourceUserConsent
type ListResourceUserConsentsResponse = clientsv1.ListResourceUserConsentsResponse
type RevokeUserConsentResponse = clientsv1.RevokeUserConsentResponse

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

// ResourceService defines the interface for reading and revoking the end-user
// consents granted against a resource, such as an MCP server.
//
// A consent records that one end user allowed a specific API client to act on
// their behalf. Each consent identifies the user by ExternalUserId — the
// identifier your application supplied when the consent was granted.
type ResourceService interface {
	ListUserConsents(ctx context.Context, resourceId string, options ListUserConsentsOptions) (*ListResourceUserConsentsResponse, error)
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
