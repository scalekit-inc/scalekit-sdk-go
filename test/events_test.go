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
	// Events may legitimately be empty for a fresh test environment; the page
	// tokens must be readable regardless of whether more pages exist.
	_ = resp.GetNextPageToken()
	_ = resp.GetPrevPageToken()
}

// TestListEventsPaginatedRejectsInvalidPageSize verifies that out-of-range page
// sizes are rejected with the ErrInvalidPageSize sentinel before any RPC is
// made: values beyond the uint32 range (which would otherwise wrap) and
// negative values (which must not silently fall through to a default page).
func TestListEventsPaginatedRejectsInvalidPageSize(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name     string
		pageSize int
	}{
		{name: "above uint32 max", pageSize: math.MaxUint32 + 1},
		{name: "negative", pageSize: -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Events().ListEventsPaginated(ctx, scalekit.ListEventsOptions{
				PageSize: tc.pageSize,
			})
			assert.Nil(t, resp)
			assert.ErrorIs(t, err, scalekit.ErrInvalidPageSize)
		})
	}
}
