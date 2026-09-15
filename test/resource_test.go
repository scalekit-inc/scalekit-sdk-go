package test

import (
	"context"
	"os"
	"testing"

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

// testResourceId returns the resource id configured for resource-client
// integration tests, skipping the calling test if it is not set. A resource
// must already exist in the test environment (resources are provisioned via
// the dashboard/API, not by this SDK), so this is env-driven rather than
// created inline like the shared test organization.
func testResourceId(t *testing.T) string {
	t.Helper()
	resourceId := os.Getenv("SCALEKIT_TEST_RESOURCE_ID")
	if resourceId == "" {
		t.Skip("set SCALEKIT_TEST_RESOURCE_ID to run resource client tests")
	}
	return resourceId
}

func TestCreateResourceClientRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resource().CreateResourceClient(ctx, "", &clients.ResourceClient{Name: "Test"})
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestGetResourceClientRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resource().GetResourceClient(ctx, "", "m2m_dummy")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestGetResourceClientRequiresClientId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resource().GetResourceClient(ctx, "res_dummy", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestListResourceClientsRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resource().ListResourceClients(ctx, "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestUpdateResourceClientRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resource().UpdateResourceClient(ctx, "", "m2m_dummy", &clients.ResourceClient{Name: "Test"}, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestUpdateResourceClientRequiresClientId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resource().UpdateResourceClient(ctx, "res_dummy", "", &clients.ResourceClient{Name: "Test"}, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestDeleteResourceClientRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	err := client.Resource().DeleteResourceClient(ctx, "", "m2m_dummy")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestDeleteResourceClientRequiresClientId(t *testing.T) {
	ctx := context.Background()

	err := client.Resource().DeleteResourceClient(ctx, "res_dummy", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestListUserConsentsRequiresResourceId(t *testing.T) {
	ctx := context.Background()

	_, err := client.Resource().ListUserConsents(ctx, "", scalekit.ListUserConsentsOptions{})
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrResourceIdRequired)
}

func TestRevokeUserConsentRequiresClientId(t *testing.T) {
	ctx := context.Background()

	err := client.Resource().RevokeUserConsent(ctx, "", "usrcnst_dummy")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
}

func TestRevokeUserConsentRequiresConsentId(t *testing.T) {
	ctx := context.Background()

	err := client.Resource().RevokeUserConsent(ctx, "m2m_dummy", "")
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrConsentIdRequired)
}

func TestCreateGetUpdateDeleteResourceClient(t *testing.T) {
	resourceId := testResourceId(t)
	ctx := context.Background()

	created, err := client.Resource().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
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
			_ = client.Resource().DeleteResourceClient(ctx, resourceId, clientId)
		}
	})

	assert.Equal(t, "Go SDK Test Client", created.Client.Name)
	assert.Equal(t, resourceId, created.Client.ResourceId)

	fetched, err := client.Resource().GetResourceClient(ctx, resourceId, clientId)
	require.NoError(t, err)
	require.NotNil(t, fetched.Client)
	assert.Equal(t, clientId, fetched.Client.ClientId)

	list, err := client.Resource().ListResourceClients(ctx, resourceId)
	require.NoError(t, err)
	found := false
	for _, c := range list.Clients {
		if c.ClientId == clientId {
			found = true
			break
		}
	}
	assert.True(t, found, "created client should appear in list")

	updated, err := client.Resource().UpdateResourceClient(ctx, resourceId, clientId, &clients.ResourceClient{
		Name: "Go SDK Test Client Updated",
	}, &fieldmaskpb.FieldMask{Paths: []string{"name"}})
	require.NoError(t, err)
	require.NotNil(t, updated.Client)
	assert.Equal(t, "Go SDK Test Client Updated", updated.Client.Name)

	err = client.Resource().DeleteResourceClient(ctx, resourceId, clientId)
	require.NoError(t, err)
	deleted = true

	_, err = client.Resource().GetResourceClient(ctx, resourceId, clientId)
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
	resourceId := testResourceId(t)
	ctx := context.Background()

	created, err := client.Resource().CreateResourceClient(ctx, resourceId, &clients.ResourceClient{
		Name: "Go SDK Ownership Test Client",
	})
	require.NoError(t, err)
	require.NotNil(t, created.Client)
	clientId := created.Client.ClientId
	t.Cleanup(func() {
		_ = client.Resource().DeleteResourceClient(ctx, resourceId, clientId)
	})

	err = client.Resource().DeleteResourceClient(ctx, OtherResourceId, clientId)
	require.Error(t, err)

	// The client must still exist under its real resource.
	stillThere, err := client.Resource().GetResourceClient(ctx, resourceId, clientId)
	require.NoError(t, err)
	assert.Equal(t, clientId, stillThere.Client.ClientId)
}

func TestListUserConsents(t *testing.T) {
	resourceId := testResourceId(t)
	ctx := context.Background()

	list, err := client.Resource().ListUserConsents(ctx, resourceId, scalekit.ListUserConsentsOptions{
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, list)
}
