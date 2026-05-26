package test

import (
	"context"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	organizations "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/organizations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testLogoURL = "https://example.com/logo.png"
)

func TestOrganization_CreateWithLogoUrl(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	resp, err := client.Organization().CreateOrganization(ctx, TestOrgName, scalekit.CreateOrganizationOptions{
		ExternalId: UniqueSuffix(),
		LogoUrl:    testLogoURL,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetOrganization())
	defer DeleteTestOrganization(t, ctx, resp.GetOrganization().GetId())

	assert.Equal(t, testLogoURL, resp.GetOrganization().GetLogoUrl())
}

func TestOrganization_CreateWithSlug(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	slug := generateTestSlug()

	resp, err := client.Organization().CreateOrganization(ctx, TestOrgName, scalekit.CreateOrganizationOptions{
		ExternalId: UniqueSuffix(),
		Slug:       slug,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetOrganization())
	defer DeleteTestOrganization(t, ctx, resp.GetOrganization().GetId())

	assert.Equal(t, slug, resp.GetOrganization().GetSlug())
}

func TestOrganization_UpdateLogoUrl(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	updated, err := client.Organization().UpdateOrganization(ctx, orgId, &organizations.UpdateOrganization{
		LogoUrl: toPtr(testLogoURL),
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.NotNil(t, updated.GetOrganization())
	assert.Equal(t, testLogoURL, updated.GetOrganization().GetLogoUrl())

	fetched, err := client.Organization().GetOrganization(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, fetched)
	require.NotNil(t, fetched.GetOrganization())
	assert.Equal(t, testLogoURL, fetched.GetOrganization().GetLogoUrl())
}

func TestOrganization_UpdateSlug(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)
	slug := generateTestSlug()

	updated, err := client.Organization().UpdateOrganization(ctx, orgId, &organizations.UpdateOrganization{
		Slug:     toPtr(slug),
		Metadata: map[string]string{"custom_domain": slug},
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.NotNil(t, updated.GetOrganization())
	assert.Equal(t, slug, updated.GetOrganization().GetSlug())
	assert.Equal(t, slug, updated.GetOrganization().GetMetadata()["custom_domain"])

	fetched, err := client.Organization().GetOrganization(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, fetched)
	require.NotNil(t, fetched.GetOrganization())
	assert.Equal(t, slug, fetched.GetOrganization().GetSlug())
	assert.Equal(t, slug, fetched.GetOrganization().GetMetadata()["custom_domain"])
}
