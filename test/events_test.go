package test

import (
	"context"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvents_ListEvents_Basic(t *testing.T) {
	ctx := context.Background()
	resp, err := client.Events().ListEvents(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.GetEvents())
}

func TestEvents_ListEvents_WithFilters(t *testing.T) {
	ctx := context.Background()
	resp, err := client.Events().ListEvents(ctx, &scalekit.ListEventsOptions{
		OrganizationId: "org_does_not_exist",
		Source:         scalekit.EventSourceScalekit,
		PageSize:       5,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	// A nonexistent organization_id filters out all events, but the call still succeeds.
	assert.Equal(t, 0, len(resp.GetEvents()))
}

func TestEvents_ListEvents_AuthRequestIdFilter_NoMatch(t *testing.T) {
	ctx := context.Background()
	resp, err := client.Events().ListEvents(ctx, &scalekit.ListEventsOptions{
		AuthRequestId: "areq_does_not_exist",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.GetEvents()))
}

func TestEvents_ListEvents_NilOptions(t *testing.T) {
	ctx := context.Background()
	resp, err := client.Events().ListEvents(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
