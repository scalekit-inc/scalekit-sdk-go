package test

import (
	"context"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditLogs_ListAuthRequests_Basic(t *testing.T) {
	ctx := context.Background()
	resp, err := client.AuditLogs().ListAuthRequests(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.GetAuthRequests())
}

func TestAuditLogs_ListAuthRequests_WithFilters(t *testing.T) {
	ctx := context.Background()
	resp, err := client.AuditLogs().ListAuthRequests(ctx, &scalekit.ListAuthRequestsOptions{
		Email:    "nobody-matching@example.com",
		Status:   []string{"SUCCESS"},
		PageSize: 5,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	// A non-matching email filters out all logs, but the call still succeeds.
	assert.Equal(t, 0, len(resp.GetAuthRequests()))
}

func TestAuditLogs_ListAuthRequests_NilOptions(t *testing.T) {
	ctx := context.Background()
	resp, err := client.AuditLogs().ListAuthRequests(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestAuditLogs_EventsCorrelateWithAuthRequestId fetches real authentication request logs,
// takes a real auth_request_id from the results, and confirms the Events API returns at
// least one event for that same auth_request_id. No IDs are hardcoded — everything is
// fetched live from the environment. If the environment has no authentication request
// history, there is nothing to correlate, so the test is skipped rather than failed.
func TestAuditLogs_EventsCorrelateWithAuthRequestId(t *testing.T) {
	ctx := context.Background()

	authResp, err := client.AuditLogs().ListAuthRequests(ctx, &scalekit.ListAuthRequestsOptions{PageSize: 50})
	require.NoError(t, err)
	require.NotNil(t, authResp)

	var authRequestId string
	for _, entry := range authResp.GetAuthRequests() {
		if entry.GetAuthRequestId() != "" {
			authRequestId = entry.GetAuthRequestId()
			break
		}
	}
	if authRequestId == "" {
		t.Skip("no authentication request logs with an auth_request_id were found in this environment; nothing to correlate against the Events API")
	}

	eventsResp, err := client.Events().ListEvents(ctx, &scalekit.ListEventsOptions{AuthRequestId: authRequestId})
	require.NoError(t, err)
	require.NotNil(t, eventsResp)
	assert.Greater(t, len(eventsResp.GetEvents()), 0,
		"expected at least one event for auth_request_id=%q, which was returned by ListAuthRequests, but the Events API returned none", authRequestId)
}
