package test

import (
	"context"
	"strconv"
	"testing"

	"connectrpc.com/connect"
	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A real MCP server resource in the test environment. These tests assert the
// call shape and the pagination envelope, never the consent contents, so they
// hold whether or not the resource currently has consents.
const testResourceId = "res_142388234624696322"

func buildUserIds(n int) []string {
	userIds := make([]string, 0, n)
	for i := 0; i < n; i++ {
		userIds = append(userIds, "usr_"+strconv.Itoa(i))
	}
	return userIds
}

func TestResourceServiceListUserConsentsValidation(t *testing.T) {
	resourceService := client.Resources()
	ctx := context.Background()

	_, err := resourceService.ListUserConsents(ctx, "", scalekit.ListUserConsentsOptions{})
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestResourceServiceRevokeUserConsentValidation(t *testing.T) {
	resourceService := client.Resources()
	ctx := context.Background()

	tests := []struct {
		name      string
		clientId  string
		consentId string
		wantErr   error
	}{
		{name: "missing client id", clientId: "", consentId: "usrcnst_1234567890", wantErr: scalekit.ErrClientIdRequired},
		{name: "missing consent id", clientId: "m2m_1234567890", consentId: "", wantErr: scalekit.ErrConsentIdRequired},
		{name: "both missing", clientId: "", consentId: "", wantErr: scalekit.ErrClientIdRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resourceService.RevokeUserConsent(ctx, tt.clientId, tt.consentId)
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestResourceServiceListUserConsents(t *testing.T) {
	resourceService := client.Resources()
	ctx := context.Background()

	tests := []struct {
		name    string
		options scalekit.ListUserConsentsOptions
	}{
		{name: "no options", options: scalekit.ListUserConsentsOptions{}},
		{name: "with page size", options: scalekit.ListUserConsentsOptions{PageSize: 10}},
		{name: "with search", options: scalekit.ListUserConsentsOptions{Search: "usr_"}},
		{name: "with user ids", options: scalekit.ListUserConsentsOptions{UserIds: []string{"usr_does_not_exist"}}},
		// The filter wins over search server-side; both together must still be a
		// valid request rather than a 400.
		{name: "with user ids and search", options: scalekit.ListUserConsentsOptions{UserIds: []string{"usr_does_not_exist"}, Search: "usr_"}},
		{name: "with 25 user ids", options: scalekit.ListUserConsentsOptions{UserIds: buildUserIds(25)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := resourceService.ListUserConsents(ctx, testResourceId, tt.options)
			require.NoError(t, err)
			require.NotNil(t, response)
			// GetConsents() is nil, not empty, when the resource has no matching
			// consents, so assert what actually has to hold rather than non-nilness.
			consents := response.GetConsents()
			assert.LessOrEqual(t, len(consents), int(response.GetTotalSize()))
			if tt.options.PageSize != 0 {
				assert.LessOrEqual(t, uint32(len(consents)), tt.options.PageSize)
			}
			// Every returned consent must belong to one of the requested users —
			// this is the filter's contract, and it holds however many rows exist.
			if len(tt.options.UserIds) > 0 {
				for _, consent := range consents {
					assert.Contains(t, tt.options.UserIds, consent.GetExternalUserId())
				}
			}
		})
	}
}

// TestResourceServiceListUserConsentsRejectsMoreThan25UserIds asserts the
// server-side cap on the filter, so a caller batching user IDs learns about the
// limit instead of silently getting a truncated result.
func TestResourceServiceListUserConsentsRejectsMoreThan25UserIds(t *testing.T) {
	resourceService := client.Resources()
	ctx := context.Background()

	_, err := resourceService.ListUserConsents(ctx, testResourceId, scalekit.ListUserConsentsOptions{
		UserIds: buildUserIds(26),
	})
	require.Error(t, err)
	assert.Equal(t, connect.CodeInvalidArgument, connect.CodeOf(err))
}

// TestResourceServiceListUserConsentsUnknownResource asserts that an unknown
// resource id surfaces as an error rather than an empty success, so callers can
// tell "no consents" apart from "no such resource".
func TestResourceServiceListUserConsentsUnknownResource(t *testing.T) {
	resourceService := client.Resources()
	ctx := context.Background()

	_, err := resourceService.ListUserConsents(ctx, "res_000000000000000000", scalekit.ListUserConsentsOptions{})
	assert.Error(t, err)
}
