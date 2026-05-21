package test

import (
	"context"
	"strings"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	organizations "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/organizations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testLogoURL = "https://cdn.acmecorp.com/logo.png"
	testSlug    = "app.acmecorp.com"
)

func TestOrganization_CreateWithLogoUrl(t *testing.T) {
	ctx := context.Background()

	resp, err := client.Organization().CreateOrganization(ctx, TestOrgName, scalekit.CreateOrganizationOptions{
		ExternalId: UniqueSuffix(),
		LogoUrl:    toPtr(testLogoURL),
	})
	if err != nil {
		if strings.Contains(err.Error(), "logo_url is not allowed") {
			t.Skipf("skipping logo_url test: environment does not support logo_url (%v)", err)
		}
		require.NoError(t, err)
	}
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetOrganization())
	defer DeleteTestOrganization(t, ctx, resp.GetOrganization().GetId())

	assert.Equal(t, testLogoURL, resp.GetOrganization().GetLogoUrl())
}

func TestOrganization_CreateWithSlug(t *testing.T) {
	ctx := context.Background()

	resp, err := client.Organization().CreateOrganization(ctx, TestOrgName, scalekit.CreateOrganizationOptions{
		ExternalId: UniqueSuffix(),
		Slug:       toPtr(testSlug),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetOrganization())
	defer DeleteTestOrganization(t, ctx, resp.GetOrganization().GetId())

	assert.Equal(t, testSlug, resp.GetOrganization().GetSlug())
}

func TestOrganization_UpdateLogoUrl(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	updated, err := client.Organization().UpdateOrganization(ctx, orgId, &organizations.UpdateOrganization{
		LogoUrl: toPtr(testLogoURL),
	})
	if err != nil {
		if strings.Contains(err.Error(), "logo_url is not allowed") {
			t.Skipf("skipping logo_url test: environment does not support logo_url (%v)", err)
		}
		require.NoError(t, err)
	}
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
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	updated, err := client.Organization().UpdateOrganization(ctx, orgId, &organizations.UpdateOrganization{
		Slug:     toPtr(testSlug),
		Metadata: map[string]string{"custom_domain": "app.acmecorp.com"},
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.NotNil(t, updated.GetOrganization())
	assert.Equal(t, testSlug, updated.GetOrganization().GetSlug())
	assert.Equal(t, "app.acmecorp.com", updated.GetOrganization().GetMetadata()["custom_domain"])

	fetched, err := client.Organization().GetOrganization(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, fetched)
	require.NotNil(t, fetched.GetOrganization())
	assert.Equal(t, testSlug, fetched.GetOrganization().GetSlug())
	assert.Equal(t, "app.acmecorp.com", fetched.GetOrganization().GetMetadata()["custom_domain"])
}
