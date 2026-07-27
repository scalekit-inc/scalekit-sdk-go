package scalekit

import (
	"context"
	"math"

	eventsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/events"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/events/eventsconnect"
)

// Type aliases for event-related types.
type EventFilter = eventsv1.EventFilter
type ListEventsPaginatedResponse = eventsv1.ListEventsPaginatedResponse

// ListEventsOptions defines the optional inputs for listing events with pagination.
type ListEventsOptions struct {
	// Filter is the optional set of filters (event types, time range,
	// organization, source, etc.) applied to the event query.
	Filter *EventFilter
	// PageSize is the optional maximum number of events to return per page.
	PageSize int
	// PageToken is the optional cursor returned by a previous call used to
	// fetch the next page of results.
	PageToken string
}

// EventsService provides helper methods for querying the Events gRPC surface.
type EventsService interface {
	ListEventsPaginated(ctx context.Context, options ListEventsOptions) (*ListEventsPaginatedResponse, error)
}

type eventsService struct {
	coreClient *coreClient
	client     eventsconnect.EventsServiceClient
}

func newEventsService(coreClient *coreClient) EventsService {
	return &eventsService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, eventsconnect.NewEventsServiceClient),
	}
}

func (e *eventsService) ListEventsPaginated(ctx context.Context, options ListEventsOptions) (*ListEventsPaginatedResponse, error) {
	request := &eventsv1.ListEventsPaginatedRequest{
		Filter:    options.Filter,
		PageToken: options.PageToken,
	}
	if options.PageSize < 0 {
		return nil, ErrInvalidPageSize
	}
	if options.PageSize > 0 {
		if uint64(options.PageSize) > math.MaxUint32 {
			return nil, ErrInvalidPageSize
		}
		request.PageSize = uint32(options.PageSize)
	}
	return newConnectExecuter(
		e.coreClient,
		e.client.ListEventsPaginated,
		request,
	).exec(ctx)
}
