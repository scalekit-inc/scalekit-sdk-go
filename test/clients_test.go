package test

import (
	"context"
	"os"
	"strings"
	"testing"

	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	clients "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestClientLifecycleByClientType(t *testing.T) {
	ctx := context.Background()

	type testCase struct {
		name       string
		clientType string
	}

	tests := []testCase{
		{name: "spa client lifecycle", clientType: "SPA"},
		{name: "ntv client lifecycle", clientType: "NTV"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suffix := strings.ToLower(tt.clientType) + "-" + UniqueSuffix()
			initialName := "sdk-it-initial-" + suffix
			updatedName := "sdk-it-updated-" + suffix

			initialPostLogin := []string{
				"https://" + suffix + ".example.com/callback-1",
				"https://" + suffix + ".example.com/callback-2",
			}
			updatedPostLogin := []string{
				"https://" + suffix + ".example.com/new-callback-1",
				"https://" + suffix + ".example.com/new-callback-2",
			}
			initialPostLogout := []string{
				"https://" + suffix + ".example.com/logout-1",
				"https://" + suffix + ".example.com/logout-2",
			}
			updatedPostLogout := []string{
				"https://" + suffix + ".example.com/new-logout-1",
				"https://" + suffix + ".example.com/new-logout-2",
			}

			createReq := &clients.CreateClient{
				Name:                   initialName,
				ClientType:             tt.clientType,
				PostLoginUris:          initialPostLogin,
				PostLogoutRedirectUris: initialPostLogout,
				AccessTokenExpiry:      3600,
			}

			createResp, err := client.Client().CreateClient(ctx, createReq)
			require.NoError(t, err)
			require.NotNil(t, createResp)
			require.NotNil(t, createResp.GetClient())

			clientID := createResp.GetClient().GetId()
			require.NotEmpty(t, clientID)
			defer deleteClientIfExists(t, ctx, clientID)

			createdClient := mustGetClient(t, ctx, clientID)
			assert.Equal(t, initialName, createdClient.GetName())
			assert.Equal(t, tt.clientType, createdClient.GetClientType())
			assert.Equal(t, initialPostLogin, createdClient.GetPostLoginUris())
			assert.Equal(t, initialPostLogout, createdClient.GetPostLogoutRedirectUris())
			assert.Equal(t, int64(3600), createdClient.GetAccessTokenExpiry())
			assert.Len(t, createdClient.GetBackChannelLogoutUris(), 0)

			updateMask := &fieldmaskpb.FieldMask{
				Paths: []string{
					"name",
					"post_login_uris",
					"post_logout_redirect_uris",
					"access_token_expiry",
				},
			}
			updateReq := &clients.UpdateClient{
				Name:                   updatedName,
				PostLoginUris:          updatedPostLogin,
				PostLogoutRedirectUris: updatedPostLogout,
				AccessTokenExpiry:      1800,
			}
			_, err = client.Client().UpdateClient(ctx, clientID, updateReq, updateMask)
			require.NoError(t, err)

			updatedClient := mustGetClient(t, ctx, clientID)
			assert.Equal(t, updatedName, updatedClient.GetName())
			assert.Equal(t, updatedPostLogin, updatedClient.GetPostLoginUris())
			assert.Equal(t, updatedPostLogout, updatedClient.GetPostLogoutRedirectUris())
			assert.Equal(t, int64(1800), updatedClient.GetAccessTokenExpiry())
		})
	}
}

func TestWebClientLifecycleWithPKCEAndSecrets(t *testing.T) {
	ctx := context.Background()
	environmentURL := os.Getenv(EnvEnvironmentURL)

	type testCase struct {
		name                string
		enforcePkceOnCreate *wrapperspb.BoolValue
		expectEnforcePkce   bool
	}

	tests := []testCase{
		{
			name:                "web client with pkce enforced",
			enforcePkceOnCreate: wrapperspb.Bool(true),
			expectEnforcePkce:   true,
		},
		{
			name:                "web client without enforce_pkce flag",
			enforcePkceOnCreate: nil,
			expectEnforcePkce:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suffix := "web-" + UniqueSuffix()
			initialName := "sdk-it-web-initial-" + suffix
			updatedName := "sdk-it-web-updated-" + suffix

			initialPostLogin := []string{
				"https://" + suffix + ".example.com/callback-1",
				"https://" + suffix + ".example.com/callback-2",
			}
			updatedPostLogin := []string{
				"https://" + suffix + ".example.com/new-callback-1",
				"https://" + suffix + ".example.com/new-callback-2",
			}
			initialPostLogout := []string{
				"https://" + suffix + ".example.com/logout-1",
				"https://" + suffix + ".example.com/logout-2",
			}
			updatedPostLogout := []string{
				"https://" + suffix + ".example.com/new-logout-1",
				"https://" + suffix + ".example.com/new-logout-2",
			}
			initialBackchannel := []string{
				"https://example.com/backchannel-1",
				"https://example.com/backchannel-2",
			}
			updatedBackchannel := []string{
				"https://example.com/new-backchannel-1",
				"https://example.com/new-backchannel-2",
			}

			createResp, err := client.Client().CreateClient(ctx, &clients.CreateClient{
				Name:                   initialName,
				ClientType:             "WEB_APP",
				PostLoginUris:          initialPostLogin,
				PostLogoutRedirectUris: initialPostLogout,
				BackChannelLogoutUris:  initialBackchannel,
				AccessTokenExpiry:      3600,
				EnforcePkce:            tt.enforcePkceOnCreate,
			})
			require.NoError(t, err)
			require.NotNil(t, createResp)
			require.NotNil(t, createResp.GetClient())

			clientID := createResp.GetClient().GetId()
			require.NotEmpty(t, clientID)
			defer deleteClientIfExists(t, ctx, clientID)

			createdClient := mustGetClient(t, ctx, clientID)
			assert.Equal(t, initialName, createdClient.GetName())
			assert.Equal(t, "WEB_APP", createdClient.GetClientType())
			assert.Equal(t, initialPostLogin, createdClient.GetPostLoginUris())
			assert.Equal(t, initialPostLogout, createdClient.GetPostLogoutRedirectUris())
			assert.Equal(t, initialBackchannel, createdClient.GetBackChannelLogoutUris())
			assert.Equal(t, tt.expectEnforcePkce, createdClient.GetEnforcePkce())

			updateMask := &fieldmaskpb.FieldMask{
				Paths: []string{
					"name",
					"post_login_uris",
					"post_logout_redirect_uris",
					"back_channel_logout_uris",
					"access_token_expiry",
				},
			}
			_, err = client.Client().UpdateClient(ctx, clientID, &clients.UpdateClient{
				Name:                   updatedName,
				PostLoginUris:          updatedPostLogin,
				PostLogoutRedirectUris: updatedPostLogout,
				BackChannelLogoutUris:  updatedBackchannel,
				AccessTokenExpiry:      1800,
			}, updateMask)
			require.NoError(t, err)

			updatedClient := mustGetClient(t, ctx, clientID)
			assert.Equal(t, updatedName, updatedClient.GetName())
			assert.Equal(t, updatedPostLogin, updatedClient.GetPostLoginUris())
			assert.Equal(t, updatedPostLogout, updatedClient.GetPostLogoutRedirectUris())
			assert.Equal(t, updatedBackchannel, updatedClient.GetBackChannelLogoutUris())

			secret1, err := client.Client().CreateClientSecret(ctx, clientID)
			require.NoError(t, err)
			require.NotNil(t, secret1)
			require.NotEmpty(t, secret1.GetPlainSecret())
			require.NotNil(t, secret1.GetSecret())
			require.NotEmpty(t, secret1.GetSecret().GetId())

			secret2, err := client.Client().CreateClientSecret(ctx, clientID)
			require.NoError(t, err)
			require.NotNil(t, secret2)
			require.NotEmpty(t, secret2.GetPlainSecret())
			require.NotNil(t, secret2.GetSecret())
			require.NotEmpty(t, secret2.GetSecret().GetId())
			require.NotEqual(t, secret1.GetSecret().GetId(), secret2.GetSecret().GetId())

			withTwoSecrets := mustGetClient(t, ctx, clientID)
			require.GreaterOrEqual(t, len(withTwoSecrets.GetSecrets()), 2)

			err = client.Client().DeleteClientSecret(ctx, clientID, secret1.GetSecret().GetId())
			require.NoError(t, err)

			withOneSecret := mustGetClient(t, ctx, clientID)
			require.Len(t, withOneSecret.GetSecrets(), 1)
			assert.Equal(t, secret2.GetSecret().GetId(), withOneSecret.GetSecrets()[0].GetId())

			webClient := scalekit.NewScalekitClient(environmentURL, clientID, secret2.GetPlainSecret())

			authURL, err := webClient.GetAuthorizationUrl("https://example.com/callback", scalekit.AuthorizationUrlOptions{
				State: "state-123",
			})
			require.NoError(t, err)
			require.NotNil(t, authURL)
			assert.Contains(t, authURL.String(), "client_id="+clientID)
			assert.Contains(t, authURL.String(), "response_type=code")

			isValid, err := webClient.ValidateAccessToken(ctx, "invalid.token.value")
			assert.Error(t, err)
			assert.False(t, isValid)

			logoutURL, err := webClient.GetLogoutUrl(scalekit.LogoutUrlOptions{
				PostLogoutRedirectUri: "https://example.com/post-logout",
				State:                 "logout-state",
			})
			require.NoError(t, err)
			require.NotNil(t, logoutURL)
			assert.Contains(t, logoutURL.String(), "post_logout_redirect_uri=https%3A%2F%2Fexample.com%2Fpost-logout")
			assert.Contains(t, logoutURL.String(), "state=logout-state")
		})
	}
}

func TestWebClientDisallowScalekitApiAccess(t *testing.T) {
	ctx := context.Background()
	environmentURL := os.Getenv(EnvEnvironmentURL)

	type testCase struct {
		name              string
		disallowApiAccess bool
		expectAuthError   bool
	}

	tests := []testCase{
		{
			name:              "disallow scalekit api access true returns unauthenticated",
			disallowApiAccess: true,
			expectAuthError:   true,
		},
		{
			name:              "disallow scalekit api access false allows api access",
			disallowApiAccess: false,
			expectAuthError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suffix := "web-api-" + UniqueSuffix()
			createResp, err := client.Client().CreateClient(ctx, &clients.CreateClient{
				Name:                      "sdk-it-web-api-" + suffix,
				ClientType:                "WEB_APP",
				PostLoginUris:             []string{"https://" + suffix + ".example.com/callback-1", "https://" + suffix + ".example.com/callback-2"},
				PostLogoutRedirectUris:    []string{"https://" + suffix + ".example.com/logout-1", "https://" + suffix + ".example.com/logout-2"},
				BackChannelLogoutUris:     []string{"https://example.com/backchannel-1", "https://example.com/backchannel-2"},
				AccessTokenExpiry:         3600,
				DisallowScalekitApiAccess: wrapperspb.Bool(tt.disallowApiAccess),
			})
			require.NoError(t, err)
			require.NotNil(t, createResp)
			require.NotNil(t, createResp.GetClient())

			clientID := createResp.GetClient().GetId()
			require.NotEmpty(t, clientID)
			defer deleteClientIfExists(t, ctx, clientID)

			secretResp, err := client.Client().CreateClientSecret(ctx, clientID)
			require.NoError(t, err)
			require.NotNil(t, secretResp)
			require.NotEmpty(t, secretResp.GetPlainSecret())

			webClient := scalekit.NewScalekitClient(environmentURL, clientID, secretResp.GetPlainSecret())

			assertApiAccessExpectation(t, ctx, webClient, tt.expectAuthError)
		})
	}
}

func assertApiAccessExpectation(t *testing.T, ctx context.Context, sdkClient scalekit.Scalekit, expectAuthError bool) {
	t.Helper()

	resp, err := sdkClient.Organization().ListOrganization(ctx, &scalekit.ListOrganizationOptions{PageSize: 1})

	if expectAuthError {
		require.Error(t, err)
		assert.True(t, isUnauthenticatedOr401(err), "expected unauthenticated/401 error, got: %v", err)
		return
	}

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func isUnauthenticatedOr401(err error) bool {
	if err == nil {
		return false
	}
	lowerErr := strings.ToLower(err.Error())
	return strings.Contains(lowerErr, "unauthenticated") ||
		strings.Contains(lowerErr, "401") ||
		strings.Contains(lowerErr, "unauthorized")
}

func mustGetClient(t *testing.T, ctx context.Context, clientID string) *clients.Client {
	t.Helper()

	resp, err := client.Client().GetClient(ctx, clientID)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetClient())
	return resp.GetClient()
}

func deleteClientIfExists(t *testing.T, ctx context.Context, clientID string) {
	t.Helper()
	if clientID == "" || client == nil {
		return
	}
	err := client.Client().DeleteClient(ctx, clientID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return
		}
		t.Logf("deleteClientIfExists %s: %v", clientID, err)
	}
}
