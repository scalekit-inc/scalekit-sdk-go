package scalekit

import (
	"context"

	auditlogsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/auditlogs"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/auditlogs/auditlogsconnect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListAuthRequestsResponse = auditlogsv1.ListAuthLogResponse
type AuthRequestLog = auditlogsv1.AuthLogRequest

// ListAuthRequestsOptions represents optional filters and pagination for listing
// authentication request logs.
type ListAuthRequestsOptions struct {
	PageSize                   uint32
	PageToken                  string
	Email                      string
	Status                     []string
	StartTime                  *timestamppb.Timestamp
	EndTime                    *timestamppb.Timestamp
	ResourceId                 string
	ConnectedAccountIdentifier string
	ClientId                   string
}

type AuditLogs interface {
	ListAuthRequests(ctx context.Context, options ...*ListAuthRequestsOptions) (*ListAuthRequestsResponse, error)
}

type auditLogs struct {
	coreClient *coreClient
	client     auditlogsconnect.AuditLogsServiceClient
}

func newAuditLogsClient(coreClient *coreClient) AuditLogs {
	return &auditLogs{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, auditlogsconnect.NewAuditLogsServiceClient),
	}
}

// ListAuthRequests returns a page of authentication request logs for the environment,
// ordered most-recent first. Each log's AuthRequestId can be passed to
// Events().ListEvents via ListEventsOptions.AuthRequestId to see every event a specific
// login produced.
func (a *auditLogs) ListAuthRequests(ctx context.Context, options ...*ListAuthRequestsOptions) (*ListAuthRequestsResponse, error) {
	request := &auditlogsv1.ListAuthLogRequest{}

	if len(options) > 0 && options[0] != nil {
		opts := options[0]
		request.PageSize = opts.PageSize
		request.PageToken = opts.PageToken
		request.Email = opts.Email
		request.Status = opts.Status
		request.StartTime = opts.StartTime
		request.EndTime = opts.EndTime
		request.ResourceId = opts.ResourceId
		request.ConnectedAccountIdentifier = opts.ConnectedAccountIdentifier
		request.ClientId = opts.ClientId
	}

	return newConnectExecuter(
		a.coreClient,
		a.client.ListAuthRequests,
		request,
	).exec(ctx)
}
