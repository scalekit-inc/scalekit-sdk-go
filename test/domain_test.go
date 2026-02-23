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
				}, d.DomainType)
			},
		},
		{
			name:   "with_ALLOWED_EMAIL_DOMAIN",
			domain: uniqueDomainName("allowed-email"),
			opts:   &scalekit.CreateDomainOptions{DomainType: scalekit.DomainTypeAllowedEmail},
			check: func(t *testing.T, d *domains.Domain) {
				assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, d.DomainType)
			},
		},
		{
			name:   "with_ORGANIZATION_DOMAIN",
			domain: uniqueDomainName("org-domain"),
			opts:   &scalekit.CreateDomainOptions{DomainType: scalekit.DomainTypeOrganization},
			check: func(t *testing.T, d *domains.Domain) {
				assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, d.DomainType)
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
				}, d.DomainType)
			},
		},
		{
			name:   "with_string_type_ALLOWED_EMAIL_DOMAIN",
			domain: uniqueDomainName("string-domain"),
			opts:   &scalekit.CreateDomainOptions{DomainType: "ALLOWED_EMAIL_DOMAIN"},
			check: func(t *testing.T, d *domains.Domain) {
				assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, d.DomainType)
			},
		},
		{
			name:   "with_string_type_ORGANIZATION_DOMAIN",
			domain: uniqueDomainName("string-org"),
			opts:   &scalekit.CreateDomainOptions{DomainType: "ORGANIZATION_DOMAIN"},
			check: func(t *testing.T, d *domains.Domain) {
				assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, d.DomainType)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			SkipIfNoIntegrationEnv(t)
			ctx := context.Background()
			orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
			defer DeleteTestOrganization(t, ctx, orgId)

			created, err := client.Domain().CreateDomain(ctx, orgId, tc.domain, tc.opts)
			require.NoError(t, err)
			require.NotNil(t, created)
			require.NotNil(t, created.Domain)
			defer DeleteTestDomain(t, ctx, orgId, created.Domain.Id)

			assert.Equal(t, tc.domain, created.Domain.Domain)
			assert.Equal(t, orgId, created.Domain.OrganizationId)
			tc.check(t, created.Domain)
		})
	}
}

func TestDomain_EndToEndIntegration(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	domainName := uniqueDomainName("lifecycle")
	created, err := client.Domain().CreateDomain(ctx, orgId, domainName)
	require.NoError(t, err)
	require.NotNil(t, created)
	domainId := created.Domain.Id
	require.NotEmpty(t, domainId)

	retrieved, err := client.Domain().GetDomain(ctx, domainId, orgId)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	assert.Equal(t, domainName, retrieved.Domain.Domain)

	err = client.Domain().DeleteDomain(ctx, domainId, orgId)
	require.NoError(t, err)

	_, err = client.Domain().GetDomain(ctx, domainId, orgId)
	assert.Error(t, err)
}

func TestDomain_ListDomainsWithFilters(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	orgDomainName := uniqueDomainName("list-org")
	orgDomain, err := client.Domain().CreateDomain(ctx, orgId, orgDomainName, &scalekit.CreateDomainOptions{
		DomainType: scalekit.DomainTypeOrganization,
	})
	require.NoError(t, err)
	require.NotNil(t, orgDomain)
	defer DeleteTestDomain(t, ctx, orgId, orgDomain.Domain.Id)

	allowedName := uniqueDomainName("list-allowed")
	allowedDomain, err := client.Domain().CreateDomain(ctx, orgId, allowedName, &scalekit.CreateDomainOptions{
		DomainType: scalekit.DomainTypeAllowedEmail,
	})
	require.NoError(t, err)
	require.NotNil(t, allowedDomain)
	defer DeleteTestDomain(t, ctx, orgId, allowedDomain.Domain.Id)

	allList, err := client.Domain().ListDomains(ctx, orgId)
	require.NoError(t, err)
	require.True(t, len(allList.Domains) > 0)

	orgList, err := client.Domain().ListDomains(ctx, orgId, &scalekit.ListDomainOptions{
		DomainType: scalekit.DomainTypeOrganization,
	})
	require.NoError(t, err)
	for _, d := range orgList.Domains {
		assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, d.DomainType)
	}
	var foundOrg bool
	for _, d := range orgList.Domains {
		if d.Id == orgDomain.Domain.Id {
			foundOrg = true
			break
		}
	}
	assert.True(t, foundOrg, "organization domain should be in filtered list")

	allowedList, err := client.Domain().ListDomains(ctx, orgId, &scalekit.ListDomainOptions{
		DomainType: scalekit.DomainTypeAllowedEmail,
	})
	require.NoError(t, err)
	for _, d := range allowedList.Domains {
		assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, d.DomainType)
	}
	var foundAllowed bool
	for _, d := range allowedList.Domains {
		if d.Id == allowedDomain.Domain.Id {
			foundAllowed = true
			break
		}
	}
	assert.True(t, foundAllowed, "allowed email domain should be in filtered list")
}

func TestDomain_CreateDomain_BackwardCompatibility(t *testing.T) {
	SkipIfNoIntegrationEnv(t)
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	domainName := fmt.Sprintf("backward-compat-%d.com", time.Now().Unix())
	created, err := client.Domain().CreateDomain(ctx, orgId, domainName)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.Domain)
	defer DeleteTestDomain(t, ctx, orgId, created.Domain.Id)
	assert.Equal(t, domainName, created.Domain.Domain)
	assert.Equal(t, orgId, created.Domain.OrganizationId)

	withOptsName := fmt.Sprintf("with-options-%d.com", time.Now().Unix())
	opts := &scalekit.CreateDomainOptions{DomainType: "ALLOWED_EMAIL_DOMAIN"}
	withOpts, err := client.Domain().CreateDomain(ctx, orgId, withOptsName, opts)
	require.NoError(t, err)
	require.NotNil(t, withOpts)
	defer DeleteTestDomain(t, ctx, orgId, withOpts.Domain.Id)
	assert.Equal(t, withOptsName, withOpts.Domain.Domain)
	assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, withOpts.Domain.DomainType)
}
