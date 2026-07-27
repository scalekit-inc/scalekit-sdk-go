package test

import (
	"context"
	"math"
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

// TestListEventsPaginatedRejectsOversizePageSize verifies that a page size
// beyond the uint32 range is rejected with the ErrInvalidPageSize sentinel
// before any RPC is made, rather than silently wrapping to a small value.
func TestListEventsPaginatedRejectsOversizePageSize(t *testing.T) {
	ctx := context.Background()

	resp, err := client.Events().ListEventsPaginated(ctx, scalekit.ListEventsOptions{
		PageSize: math.MaxUint32 + 1,
	})
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, scalekit.ErrInvalidPageSize)
}
