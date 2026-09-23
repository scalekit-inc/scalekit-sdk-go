package test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"connectrpc.com/connect"
	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	clients "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// OtherResourceId is a syntactically valid but non-existent resource id, used
// to prove DeleteResourceClient's ownership check refuses a client/resource
// pairing that does not actually match.
const OtherResourceId = "res_999999999999999999"

// envResourceId returns the resource id configured for resource-client
// integration tests, skipping the calling test if it is not set. A resource
// must already exist in the test environment (resources are provisioned via
// the dashboard/API, not by this SDK), so this is env-driven rather than
// created inline like the shared test organization.
func envResourceId(t *testing.T) string {
	t.Helper()
	resourceId := os.Getenv("SCALEKIT_TEST_RESOURCE_ID")
	if resourceId == "" {
		t.Skip("set SCALEKIT_TEST_RESOURCE_ID to run resource client tests")
	}
	return resourceId
}

func TestGetResourceRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().GetResource(ctx, "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

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

func TestCreateResourceClientRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().CreateResourceClient(ctx, "", &clients.ResourceClient{Name: "Test"})
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestGetResourceClientRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().GetResourceClient(ctx, "", "m2m_dummy")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestGetResourceClientRequiresClientId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().GetResourceClient(ctx, "res_dummy", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestListResourceClientsRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().ListResourceClients(ctx, "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestUpdateResourceClientRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().UpdateResourceClient(ctx, "", "m2m_dummy", &clients.ResourceClient{Name: "Test"}, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestUpdateResourceClientRequiresClientId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().UpdateResourceClient(ctx, "res_dummy", "", &clients.ResourceClient{Name: "Test"}, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestDeleteResourceClientRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	err := client.Resources().DeleteResourceClient(ctx, "", "m2m_dummy")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestDeleteResourceClientRequiresClientId(t *testing.T) {
	ctx := context.Background()

	err := client.Resources().DeleteResourceClient(ctx, "res_dummy", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestCreateResourceClientSecretRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().CreateResourceClientSecret(ctx, "", "m2m_dummy")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestCreateResourceClientSecretRequiresClientId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().CreateResourceClientSecret(ctx, "res_dummy", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestDeleteResourceClientSecretRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	err := client.Resources().DeleteResourceClientSecret(ctx, "", "m2m_dummy", "secret_dummy")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestDeleteResourceClientSecretRequiresClientId(t *testing.T) {
	ctx := context.Background()

	err := client.Resources().DeleteResourceClientSecret(ctx, "res_dummy", "", "secret_dummy")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestDeleteResourceClientSecretRequiresSecretId(t *testing.T) {
	ctx := context.Background()

	err := client.Resources().DeleteResourceClientSecret(ctx, "res_dummy", "m2m_dummy", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrSecretIdRequired)
}

func TestListUserConsentsRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().ListUserConsents(ctx, "", scalekit.ListUserConsentsOptions{})
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestRevokeUserConsentRequiresClientId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().RevokeUserConsent(ctx, "", "usrcnst_dummy")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestRevokeUserConsentRequiresConsentId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resources().RevokeUserConsent(ctx, "m2m_dummy", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrConsentIdRequired)
}

func TestCreateGetUpdateDeleteResourceClient(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name:        "Go SDK Test Client",
		Description: "Integration test client",
		Scopes:      []string{"test:e2e_resource_scope"},
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.Client)
	assert.NotEmpty(t, created.PlainSecret)

	clientId := created.Client.ClientId
	deleted := false
	t.Cleanup(func() {
		if !deleted {
			_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
		}
	})

	assert.Equal(t, "Go SDK Test Client", created.Client.Name)
	assert.Equal(t, resourceId, created.Client.ResourceId)

	fetched, err := client.Resources().GetResourceClient(ctx, resourceId, clientId)
	require.NoError(t, err)
	require.NotNil(t, fetched.Client)
	assert.Equal(t, clientId, fetched.Client.ClientId)

	list, err := client.Resources().ListResourceClients(ctx, resourceId)
	require.NoError(t, err)
	found := false
	for _, c := range list.Clients {
		if c.ClientId == clientId {
			found = true
			break
		}
	}
	assert.True(t, found, "created client should appear in list")

	updated, err := client.Resources().UpdateResourceClient(ctx, resourceId, clientId, &clients.ResourceClient{
		Name: "Go SDK Test Client Updated",
	}, &fieldmaskpb.FieldMask{Paths: []string{"name"}})
	require.NoError(t, err)
	require.NotNil(t, updated.Client)
	assert.Equal(t, "Go SDK Test Client Updated", updated.Client.Name)

	err = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	require.NoError(t, err)
	deleted = true

	_, err = client.Resources().GetResourceClient(ctx, resourceId, clientId)
	assert.Error(t, err)
}

// TestDeleteResourceClientRefusesWrongResource proves that deleting a client
// under the wrong resource scope is refused rather than silently succeeding.
// OtherResourceId doesn't exist, so the client can't belong to it — the
// server's own resource-scoping on GetResourceClient refuses the delete
// (404) before it ever runs, and the SDK-side ownership check in
// DeleteResourceClient is the second line of defense for a backend that
// didn't enforce this.
func TestDeleteResourceClientRefusesWrongResource(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name: "Go SDK Ownership Test Client",
	})
	require.NoError(t, err)
	require.NotNil(t, created.Client)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})

	err = client.Resources().DeleteResourceClient(ctx, OtherResourceId, clientId)
	require.Error(t, err)

	// The client must still exist under its real resource.
	stillThere, err := client.Resources().GetResourceClient(ctx, resourceId, clientId)
	require.NoError(t, err)
	assert.Equal(t, clientId, stillThere.Client.ClientId)
}

// TestCreateDeleteResourceClientSecret exercises the full lifecycle of a
// secret on a resource-scoped client: create a client, add a secret via
// CreateResourceClientSecret, verify the plaintext and the persisted secret
// record both come back, then remove that secret with
// DeleteResourceClientSecret.
func TestCreateDeleteResourceClientSecret(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name: "Go SDK Secret Lifecycle Client",
	})
	require.NoError(t, err)
	require.NotNil(t, created.Client)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})

	secretResp, err := client.Resources().CreateResourceClientSecret(ctx, resourceId, clientId)
	require.NoError(t, err)
	require.NotNil(t, secretResp)
	assert.NotEmpty(t, secretResp.GetPlainSecret())
	require.NotNil(t, secretResp.GetSecret())
	assert.NotEmpty(t, secretResp.GetSecret().GetId())

	err = client.Resources().DeleteResourceClientSecret(ctx, resourceId, clientId, secretResp.GetSecret().GetId())
	require.NoError(t, err)
}

// TestCreateResourceClientSecretRefusesWrongResource proves the ownership
// check applies to secret creation too: a client created under the real test
// resource can't have a secret created for it via a nonexistent OTHER
// resource id.
func TestCreateResourceClientSecretRefusesWrongResource(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name: "Go SDK Secret Ownership Test Client",
	})
	require.NoError(t, err)
	require.NotNil(t, created.Client)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})

	_, err = client.Resources().CreateResourceClientSecret(ctx, OtherResourceId, clientId)
	require.Error(t, err)
}

// TestDeleteResourceClientSecretRefusesWhenLastRemaining proves the server
// requires a resource client to always keep at least one secret: deleting
// the lone secret a client is created with is refused.
func TestDeleteResourceClientSecretRefusesWhenLastRemaining(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name: "Go SDK Min Secret Limit Client",
	})
	require.NoError(t, err)
	require.NotNil(t, created.Client)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})

	fetched, err := client.Resources().GetResourceClient(ctx, resourceId, clientId)
	require.NoError(t, err)
	require.NotEmpty(t, fetched.Client.GetSecrets())
	onlySecretId := fetched.Client.GetSecrets()[0].GetId()

	err = client.Resources().DeleteResourceClientSecret(ctx, resourceId, clientId, onlySecretId)
	require.Error(t, err)
}

// TestCreateResourceClientSecretRefusesPastLimit proves the server caps how
// many secrets a client can hold at once. The exact limit is
// environment-configurable (verified live: 5 in Scalekit's own dev
// environment, not the dashboard's stricter UI-only threshold of 2), so this
// probes until the server actually refuses rather than asserting a specific
// count.
func TestCreateResourceClientSecretRefusesPastLimit(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name: "Go SDK Max Secret Limit Client",
	})
	require.NoError(t, err)
	require.NotNil(t, created.Client)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})

	limitHit := false
	for i := 0; i < 20; i++ {
		if _, err := client.Resources().CreateResourceClientSecret(ctx, resourceId, clientId); err != nil {
			limitHit = true
			break
		}
	}
	assert.True(t, limitHit, "expected the server to eventually refuse creating another secret")
}

// TestGetResourceAndListResources exercises the two read-only resource
// methods against a real MCP_SERVER resource in the test environment:
// GetResource returns the resource itself (including its scopes allowlist),
// and ListResources, filtered by that same resource type, includes it.
func TestGetResourceAndListResources(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	got, err := client.Resources().GetResource(ctx, resourceId)
	require.NoError(t, err)
	require.NotNil(t, got.Resource)
	assert.Equal(t, resourceId, got.Resource.Id)
	assert.Equal(t, scalekit.ResourceTypeMcpServer, got.Resource.ResourceType)

	list, err := client.Resources().ListResources(ctx, scalekit.ResourceTypeMcpServer, scalekit.ListResourcesOptions{
		PageSize: 30,
	})
	require.NoError(t, err)
	require.NotNil(t, list)
	found := false
	for _, res := range list.Resources {
		if res.Id == resourceId {
			found = true
			break
		}
	}
	assert.True(t, found, "test resource should appear in ListResources")
}

// TestListResourcesRejectsUnspecifiedType confirms the server rejects
// RESOURCE_TYPE_UNSPECIFIED rather than treating it as "list every type" —
// resource_type is marked required at the proto level.
func TestListResourcesRejectsUnspecifiedType(t *testing.T) {
	envResourceId(t) // ensure the live test environment is configured
	ctx := context.Background()

	_, err := client.Resources().ListResources(ctx, scalekit.ResourceTypeUnspecified, scalekit.ListResourcesOptions{})
	assert.Error(t, err)
}

func TestListUserConsents(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	list, err := client.Resources().ListUserConsents(ctx, resourceId, scalekit.ListUserConsentsOptions{
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, list)
}

// TestCreateResourceClientAllFields exercises every settable field on create.
// Audience is deliberately excluded — it's not settable at all (see
// TestCreateResourceClientRejectsAudience).
func TestCreateResourceClientAllFields(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name:         "Go SDK All Fields Client",
		Description:  "exercises every field",
		Scopes:       []string{"test:e2e_resource_scope"},
		CustomClaims: []*clients.CustomClaim{{Key: "team", Value: "sdk"}},
		Expiry:       3600,
		RedirectUris: []string{"https://example.com/callback"},
	})
	require.NoError(t, err)
	require.NotNil(t, created.Client)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})

	assert.Equal(t, []string{"test:e2e_resource_scope"}, created.Client.Scopes)
	require.Len(t, created.Client.CustomClaims, 1)
	assert.Equal(t, "team", created.Client.CustomClaims[0].Key)
	assert.Equal(t, "sdk", created.Client.CustomClaims[0].Value)
	assert.Equal(t, int64(3600), created.Client.Expiry)
	assert.Equal(t, []string{"https://example.com/callback"}, created.Client.RedirectUris)
}

// TestCreateResourceClientRejectsAudience confirms audience can't be set via
// create, by design — the SDK rejects a non-empty Audience outright, rather
// than sending a request that used to be honored server-side for non-MCP
// resource types.
func TestCreateResourceClientRejectsAudience(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	_, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name:     "Audience Reject Test",
		Audience: []string{"https://example.com/should-not-apply"},
	})
	require.ErrorIs(t, err, scalekit.ErrAudienceNotSettable)
}

// TestUpdateResourceClientNameDescriptionAppliedRegardlessOfMask confirms
// name/description are truthy-gated, not mask-gated: a non-empty value
// applies even when its path isn't in the mask.
func TestUpdateResourceClientNameDescriptionAppliedRegardlessOfMask(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{Name: "Original Name"})
	require.NoError(t, err)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})

	updated, err := client.Resources().UpdateResourceClient(ctx, resourceId, clientId, &clients.ResourceClient{
		Name:        "Applied Despite Missing From Mask",
		Description: "also applied",
	}, &fieldmaskpb.FieldMask{Paths: []string{"description"}}) // "name" deliberately left out of the mask
	require.NoError(t, err)
	assert.Equal(t, "Applied Despite Missing From Mask", updated.Client.Name)
	assert.Equal(t, "also applied", updated.Client.Description)
}

// TestUpdateResourceClientEmptyStringIsNoOp confirms an empty name/description
// does not clear the field, even when its path is in the mask — there is
// currently no way to clear either field via update.
func TestUpdateResourceClientEmptyStringIsNoOp(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name:        "Keep This Name",
		Description: "keep this description",
	})
	require.NoError(t, err)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})

	updated, err := client.Resources().UpdateResourceClient(ctx, resourceId, clientId, &clients.ResourceClient{
		Name:        "",
		Description: "",
	}, &fieldmaskpb.FieldMask{Paths: []string{"name", "description"}})
	require.NoError(t, err)
	assert.Equal(t, "Keep This Name", updated.Client.Name)
	assert.Equal(t, "keep this description", updated.Client.Description)
}

// TestUpdateResourceClientRejectsAudienceInMask confirms audience can't be
// changed via update, by design — the SDK rejects an "audience" mask path
// outright, rather than sending a request that the server would just ignore.
func TestUpdateResourceClientRejectsAudienceInMask(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{Name: "Audience Immutable Test"})
	require.NoError(t, err)
	clientId := created.Client.ClientId
	originalAudience := created.Client.Audience
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})

	_, err = client.Resources().UpdateResourceClient(ctx, resourceId, clientId, &clients.ResourceClient{
		Audience: []string{"https://example.com/should-not-apply"},
	}, &fieldmaskpb.FieldMask{Paths: []string{"audience"}})
	require.ErrorIs(t, err, scalekit.ErrAudienceNotSettable)

	fetched, err := client.Resources().GetResourceClient(ctx, resourceId, clientId)
	require.NoError(t, err)
	assert.Equal(t, originalAudience, fetched.Client.Audience)
}

// TestUpdateResourceClientClearsListFields confirms scopes/customClaims/
// redirectUris — the three fields the mask actually governs — can be cleared
// by passing an empty value with the path included in the mask.
func TestUpdateResourceClientClearsListFields(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name:         "Clear Fields Test",
		Scopes:       []string{"test:e2e_resource_scope"},
		CustomClaims: []*clients.CustomClaim{{Key: "k", Value: "v"}},
		RedirectUris: []string{"https://example.com/callback"},
	})
	require.NoError(t, err)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	})
	require.NotEmpty(t, created.Client.Scopes)
	require.NotEmpty(t, created.Client.CustomClaims)
	require.NotEmpty(t, created.Client.RedirectUris)

	updated, err := client.Resources().UpdateResourceClient(ctx, resourceId, clientId, &clients.ResourceClient{
		Scopes:       []string{},
		CustomClaims: []*clients.CustomClaim{},
		RedirectUris: []string{},
	}, &fieldmaskpb.FieldMask{Paths: []string{"scopes", "custom_claims", "redirect_uris"}})
	require.NoError(t, err)
	assert.Empty(t, updated.Client.Scopes)
	assert.Empty(t, updated.Client.CustomClaims)
	assert.Empty(t, updated.Client.RedirectUris)
}

// TestCreateResourceClientRejectsJavascriptRedirectUri and the schemeless
// case below document the server's actual redirect URI validation, verified
// live: a javascript: scheme and a URI with no scheme at all are both
// rejected; a plain (non-TLS) http:// URI, perhaps surprisingly, is not.
func TestCreateResourceClientRejectsJavascriptRedirectUri(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	_, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name:         "Bad Redirect",
		RedirectUris: []string{"javascript:alert(1)"},
	})
	require.Error(t, err)
}

func TestCreateResourceClientRejectsSchemelessRedirectUri(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	_, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name:         "Bad Redirect",
		RedirectUris: []string{"not-a-uri"},
	})
	require.Error(t, err)
}

// TestDoubleDeleteResourceClient confirms a second delete on an
// already-deleted client fails rather than silently no-op-ing.
func TestDoubleDeleteResourceClient(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	created, err := client.Resources().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{Name: "Double Delete Test"})
	require.NoError(t, err)
	clientId := created.Client.ClientId

	require.NoError(t, client.Resources().DeleteResourceClient(ctx, resourceId, clientId))
	err = client.Resources().DeleteResourceClient(ctx, resourceId, clientId)
	assert.Error(t, err)
}

// TestGetResourceClientRejectsMalformedClientId confirms a syntactically
// invalid client id is rejected by the server rather than treated as a
// not-found.
func TestGetResourceClientRejectsMalformedClientId(t *testing.T) {
	resourceId := envResourceId(t)
	ctx := context.Background()

	_, err := client.Resources().GetResourceClient(ctx, resourceId, "not-a-real-client-id")
	assert.Error(t, err)
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
