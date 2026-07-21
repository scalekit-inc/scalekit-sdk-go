package scalekit

import (
	"context"

	eventsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/events"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/events/eventsconnect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListEventsResponse = eventsv1.ListEventsResponse
type EventFilter = eventsv1.EventFilter

// EventSource is defined as an int32 type alias for the Source enum.
type EventSource = eventsv1.Source

// Event source constants
const (
	EventSourceUnspecified EventSource = eventsv1.Source_SOURCE_UNSPECIFIED
	EventSourceScalekit    EventSource = eventsv1.Source_SCALEKIT
	EventSourceDirSync     EventSource = eventsv1.Source_DIR_SYNC
)

// ListEventsOptions represents optional filters and pagination for listing events.
type ListEventsOptions struct {
	EventTypes          []string
	StartTime           *timestamppb.Timestamp
	EndTime             *timestamppb.Timestamp
	OrganizationId      string
	Source              EventSource
	AuthRequestId       string
	InterceptorId       string
	InterceptorStatus   string
	InterceptorDecision string
	ConnectionId        string
	ConnectedAccountId  string
	PageSize            uint32
	PageToken           string
}

type Events interface {
	ListEvents(ctx context.Context, options ...*ListEventsOptions) (*ListEventsResponse, error)
}

type events struct {
	coreClient *coreClient
	client     eventsconnect.EventsServiceClient
}

func newEventsClient(coreClient *coreClient) Events {
	return &events{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, eventsconnect.NewEventsServiceClient),
	}
}

// ListEvents returns a page of environment events matching the given filter, ordered
// most-recent first, along with a total count of matching events. Pass AuthRequestId
// in ListEventsOptions to see every event a specific authentication request produced —
// correlate it with the AuthRequestId returned by AuditLogs().ListAuthRequests.
func (e *events) ListEvents(ctx context.Context, options ...*ListEventsOptions) (*ListEventsResponse, error) {
	request := &eventsv1.ListEventsRequest{}

	if len(options) > 0 && options[0] != nil {
		opts := options[0]
		filter := &EventFilter{
			OrganizationId: opts.OrganizationId,
			Source:         opts.Source,
			AuthRequestId:  opts.AuthRequestId,
			StartTime:      opts.StartTime,
			EndTime:        opts.EndTime,
		}
		if len(opts.EventTypes) > 0 {
			filter.EventTypes = opts.EventTypes
		}
		if opts.InterceptorId != "" {
			filter.InterceptorId = &opts.InterceptorId
		}
		if opts.InterceptorStatus != "" {
			filter.InterceptorStatus = &opts.InterceptorStatus
		}
		if opts.InterceptorDecision != "" {
			filter.InterceptorDecision = &opts.InterceptorDecision
		}
		if opts.ConnectionId != "" {
			filter.ConnectionId = &opts.ConnectionId
		}
		if opts.ConnectedAccountId != "" {
			filter.ConnectedAccountId = &opts.ConnectedAccountId
		}
		request.Filter = filter
		request.PageSize = opts.PageSize
		request.PageToken = opts.PageToken
	}

	return newConnectExecuter(
		e.coreClient,
		e.client.ListEvents,
		request,
	).exec(ctx)
}
