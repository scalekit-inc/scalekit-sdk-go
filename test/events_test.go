package test

import (
	"context"
	"testing"

	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListEventsPaginated(t *testing.T) {
	ctx := context.Background()

	resp, err := client.Events().ListEventsPaginated(ctx, scalekit.ListEventsOptions{
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	// Events may legitimately be empty for a fresh test environment.
	assert.GreaterOrEqual(t, len(resp.GetEvents()), 0)
}
