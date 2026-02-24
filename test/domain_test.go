package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/domains"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// uniqueDomainName returns a unique domain name for tests (e.g. test-<suffix>.acmecorp.com).
func uniqueDomainName(prefix string) string {
	return fmt.Sprintf("%s-%s.%s", prefix, UniqueSuffix(), TestDomainRoot)
}

func TestDomain_CreateDomain(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		opts   *scalekit.CreateDomainOptions
		check  func(t *testing.T, d *domains.Domain)
	}{
		{
			name:   "default",
			domain: uniqueDomainName("test"),
			opts:   nil,
			check: func(t *testing.T, d *domains.Domain) {
				assert.Contains(t, []domains.DomainType{
					domains.DomainType_DOMAIN_TYPE_UNSPECIFIED,
					domains.DomainType_ALLOWED_EMAIL_DOMAIN,
					domains.DomainType_ORGANIZATION_DOMAIN,
				}, d.GetDomainType())
			},
		},
		{
			name:   "with_ALLOWED_EMAIL_DOMAIN",
			domain: uniqueDomainName("allowed-email"),
			opts:   &scalekit.CreateDomainOptions{DomainType: scalekit.DomainTypeAllowedEmail},
			check: func(t *testing.T, d *domains.Domain) {
				assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, d.GetDomainType())
			},
		},
		{
			name:   "with_ORGANIZATION_DOMAIN",
			domain: uniqueDomainName("org-domain"),
			opts:   &scalekit.CreateDomainOptions{DomainType: scalekit.DomainTypeOrganization},
			check: func(t *testing.T, d *domains.Domain) {
				assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, d.GetDomainType())
			},
		},
		{
			name:   "with_nil_options",
			domain: uniqueDomainName("nil-options"),
			opts:   nil,
			check: func(t *testing.T, d *domains.Domain) {
				assert.Contains(t, []domains.DomainType{
					domains.DomainType_DOMAIN_TYPE_UNSPECIFIED,
					domains.DomainType_ALLOWED_EMAIL_DOMAIN,
					domains.DomainType_ORGANIZATION_DOMAIN,
				}, d.GetDomainType())
			},
		},
		{
			name:   "with_string_type_ALLOWED_EMAIL_DOMAIN",
			domain: uniqueDomainName("string-domain"),
			opts:   &scalekit.CreateDomainOptions{DomainType: "ALLOWED_EMAIL_DOMAIN"},
			check: func(t *testing.T, d *domains.Domain) {
				assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, d.GetDomainType())
			},
		},
		{
			name:   "with_string_type_ORGANIZATION_DOMAIN",
			domain: uniqueDomainName("string-org"),
			opts:   &scalekit.CreateDomainOptions{DomainType: "ORGANIZATION_DOMAIN"},
			check: func(t *testing.T, d *domains.Domain) {
				assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, d.GetDomainType())
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
			defer DeleteTestOrganization(t, ctx, orgId)

			created, err := client.Domain().CreateDomain(ctx, orgId, tc.domain, tc.opts)
			require.NoError(t, err)
			require.NotNil(t, created)
			require.NotNil(t, created.GetDomain())
			defer DeleteTestDomain(t, ctx, orgId, created.GetDomain().GetId())

			assert.Equal(t, tc.domain, created.GetDomain().GetDomain())
			assert.Equal(t, orgId, created.GetDomain().GetOrganizationId())
			tc.check(t, created.GetDomain())
		})
	}
}

func TestDomain_EndToEndIntegration(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	domainName := uniqueDomainName("lifecycle")
	created, err := client.Domain().CreateDomain(ctx, orgId, domainName)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.GetDomain())
	domainId := created.GetDomain().GetId()
	require.NotEmpty(t, domainId)

	retrieved, err := client.Domain().GetDomain(ctx, domainId, orgId)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	require.NotNil(t, retrieved.GetDomain())
	assert.Equal(t, domainName, retrieved.GetDomain().GetDomain())

	err = client.Domain().DeleteDomain(ctx, domainId, orgId)
	require.NoError(t, err)

	_, err = client.Domain().GetDomain(ctx, domainId, orgId)
	assert.Error(t, err)
}

func TestDomain_ListDomainsWithFilters(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	orgDomainName := uniqueDomainName("list-org")
	orgDomain, err := client.Domain().CreateDomain(ctx, orgId, orgDomainName, &scalekit.CreateDomainOptions{
		DomainType: scalekit.DomainTypeOrganization,
	})
	require.NoError(t, err)
	require.NotNil(t, orgDomain)
	require.NotNil(t, orgDomain.GetDomain())
	defer DeleteTestDomain(t, ctx, orgId, orgDomain.GetDomain().GetId())

	allowedName := uniqueDomainName("list-allowed")
	allowedDomain, err := client.Domain().CreateDomain(ctx, orgId, allowedName, &scalekit.CreateDomainOptions{
		DomainType: scalekit.DomainTypeAllowedEmail,
	})
	require.NoError(t, err)
	require.NotNil(t, allowedDomain)
	require.NotNil(t, allowedDomain.GetDomain())
	defer DeleteTestDomain(t, ctx, orgId, allowedDomain.GetDomain().GetId())

	allList, err := client.Domain().ListDomains(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, allList)
	require.True(t, len(allList.GetDomains()) > 0)

	orgList, err := client.Domain().ListDomains(ctx, orgId, &scalekit.ListDomainOptions{
		DomainType: scalekit.DomainTypeOrganization,
	})
	require.NoError(t, err)
	require.NotNil(t, orgList)
	for _, d := range orgList.GetDomains() {
		assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, d.GetDomainType())
	}
	var foundOrg bool
	for _, d := range orgList.GetDomains() {
		if d.GetId() == orgDomain.GetDomain().GetId() {
			foundOrg = true
			break
		}
	}
	assert.True(t, foundOrg, "organization domain should be in filtered list")

	allowedList, err := client.Domain().ListDomains(ctx, orgId, &scalekit.ListDomainOptions{
		DomainType: scalekit.DomainTypeAllowedEmail,
	})
	require.NoError(t, err)
	require.NotNil(t, allowedList)
	for _, d := range allowedList.GetDomains() {
		assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, d.GetDomainType())
	}
	var foundAllowed bool
	for _, d := range allowedList.GetDomains() {
		if d.GetId() == allowedDomain.GetDomain().GetId() {
			foundAllowed = true
			break
		}
	}
	assert.True(t, foundAllowed, "allowed email domain should be in filtered list")
}

func TestDomain_CreateDomain_BackwardCompatibility(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	domainName := fmt.Sprintf("backward-compat-%d.com", time.Now().Unix())
	created, err := client.Domain().CreateDomain(ctx, orgId, domainName)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.GetDomain())
	defer DeleteTestDomain(t, ctx, orgId, created.GetDomain().GetId())
	assert.Equal(t, domainName, created.GetDomain().GetDomain())
	assert.Equal(t, orgId, created.GetDomain().GetOrganizationId())

	withOptsName := fmt.Sprintf("with-options-%d.com", time.Now().Unix())
	opts := &scalekit.CreateDomainOptions{DomainType: "ALLOWED_EMAIL_DOMAIN"}
	withOpts, err := client.Domain().CreateDomain(ctx, orgId, withOptsName, opts)
	require.NoError(t, err)
	require.NotNil(t, withOpts)
	require.NotNil(t, withOpts.GetDomain())
	defer DeleteTestDomain(t, ctx, orgId, withOpts.GetDomain().GetId())
	assert.Equal(t, withOptsName, withOpts.GetDomain().GetDomain())
	assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, withOpts.GetDomain().GetDomainType())
}
