package test

import (
	"context"
	"testing"

	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A real MCP server resource in the test environment. It currently has no
// consents, which is fine — these tests assert the call shape and the
// pagination envelope, never the consent contents.
const testResourceId = "res_142145647087190278"

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := resourceService.ListUserConsents(ctx, testResourceId, tt.options)
			require.NoError(t, err)
			require.NotNil(t, response)
			assert.NotNil(t, response.GetConsents())
			if tt.options.PageSize != 0 {
				assert.LessOrEqual(t, uint32(len(response.GetConsents())), tt.options.PageSize)
			}
		})
	}
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
