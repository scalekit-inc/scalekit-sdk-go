package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/domains"
	"github.com/stretchr/testify/assert"
)

func TestDomains(t *testing.T) {
	// Test creating a domain without options (backward compatibility)
	domainName := fmt.Sprintf("test-domain-%d.com", time.Now().Unix())
	domain, err := client.Domain().CreateDomain(context.Background(), testOrg, domainName, nil)
	assert.NoError(t, err)
	assert.NotNil(t, domain)
	assert.NotNil(t, domain.Domain)
	assert.Equal(t, domainName, domain.Domain.Domain)
	assert.Equal(t, testOrg, domain.Domain.OrganizationId)
	// Note: The API might set a default domain type, so we just verify it's a valid enum value
	assert.Contains(t, []domains.DomainType{
		domains.DomainType_DOMAIN_TYPE_UNSPECIFIED,
		domains.DomainType_ALLOWED_EMAIL_DOMAIN,
		domains.DomainType_ORGANIZATION_DOMAIN,
	}, domain.Domain.DomainType)

	// Test creating a domain with ALLOWED_EMAIL_DOMAIN type
	allowedEmailDomainName := fmt.Sprintf("allowed-email-%d.com", time.Now().Unix())
	allowedEmailDomain, err := client.Domain().CreateDomain(context.Background(), testOrg, allowedEmailDomainName, &scalekit.CreateDomainOptions{
		DomainType: scalekit.DomainTypeAllowedEmail,
	})
	assert.NoError(t, err)
	assert.NotNil(t, allowedEmailDomain)
	assert.NotNil(t, allowedEmailDomain.Domain)
	assert.Equal(t, allowedEmailDomainName, allowedEmailDomain.Domain.Domain)
	assert.Equal(t, testOrg, allowedEmailDomain.Domain.OrganizationId)
	// The API should respect the requested domain type
	assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, allowedEmailDomain.Domain.DomainType)

	// Test creating a domain with ORGANIZATION_DOMAIN type
	orgDomainName := fmt.Sprintf("org-domain-%d.com", time.Now().Unix())
	orgDomain, err := client.Domain().CreateDomain(context.Background(), testOrg, orgDomainName, &scalekit.CreateDomainOptions{
		DomainType: scalekit.DomainTypeOrganization,
	})
	assert.NoError(t, err)
	assert.NotNil(t, orgDomain)
	assert.NotNil(t, orgDomain.Domain)
	assert.Equal(t, orgDomainName, orgDomain.Domain.Domain)
	assert.Equal(t, testOrg, orgDomain.Domain.OrganizationId)
	// The API should respect the requested domain type
	assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, orgDomain.Domain.DomainType)

	// Test getting domain by ID
	retrievedDomain, err := client.Domain().GetDomain(context.Background(), domain.Domain.Id, testOrg)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedDomain)
	assert.NotNil(t, retrievedDomain.Domain)
	assert.Equal(t, domain.Domain.Id, retrievedDomain.Domain.Id)
	assert.Equal(t, domainName, retrievedDomain.Domain.Domain)
	assert.Equal(t, testOrg, retrievedDomain.Domain.OrganizationId)

	// Test listing domains
	domainsList, err := client.Domain().ListDomains(context.Background(), testOrg)
	assert.NoError(t, err)
	assert.NotNil(t, domainsList)
	assert.True(t, len(domainsList.Domains) > 0)

	// Verify that our created domains are in the list
	foundCreatedDomain := false
	foundAllowedEmailDomain := false
	foundOrgDomain := false

	for _, d := range domainsList.Domains {
		if d.Id == domain.Domain.Id {
			foundCreatedDomain = true
			assert.Equal(t, domainName, d.Domain)
		}
		if d.Id == allowedEmailDomain.Domain.Id {
			foundAllowedEmailDomain = true
			assert.Equal(t, allowedEmailDomainName, d.Domain)
			assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, d.DomainType)
		}
		if d.Id == orgDomain.Domain.Id {
			foundOrgDomain = true
			assert.Equal(t, orgDomainName, d.Domain)
			assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, d.DomainType)
		}
	}

	assert.True(t, foundCreatedDomain, "Created domain should be in the list")
	assert.True(t, foundAllowedEmailDomain, "Allowed email domain should be in the list")
	assert.True(t, foundOrgDomain, "Organization domain should be in the list")
}

func TestCreateDomainWithNilOptions(t *testing.T) {
	// Test creating a domain with nil options (should work like CreateDomain)
	domainName := fmt.Sprintf("nil-options-%d.com", time.Now().Unix())
	domain, err := client.Domain().CreateDomain(context.Background(), testOrg, domainName, nil)
	assert.NoError(t, err)
	assert.NotNil(t, domain)
	assert.NotNil(t, domain.Domain)
	assert.Equal(t, domainName, domain.Domain.Domain)
	assert.Equal(t, testOrg, domain.Domain.OrganizationId)
	// Note: The API might set a default domain type, so we just verify it's a valid enum value
	assert.Contains(t, []domains.DomainType{
		domains.DomainType_DOMAIN_TYPE_UNSPECIFIED,
		domains.DomainType_ALLOWED_EMAIL_DOMAIN,
		domains.DomainType_ORGANIZATION_DOMAIN,
	}, domain.Domain.DomainType)
}

func TestCreateDomainBackwardCompatibility(t *testing.T) {
	// Test backward compatibility - 3 parameters (no options)
	domainName := fmt.Sprintf("backward-compat-%d.com", time.Now().Unix())
	domain, err := client.Domain().CreateDomain(context.Background(), testOrg, domainName)
	assert.NoError(t, err)
	assert.NotNil(t, domain)
	assert.NotNil(t, domain.Domain)
	assert.Equal(t, domainName, domain.Domain.Domain)
	assert.Equal(t, testOrg, domain.Domain.OrganizationId)

	// Test new functionality - 4 parameters (with options)
	domainWithOptionsName := fmt.Sprintf("with-options-%d.com", time.Now().Unix())
	options := &scalekit.CreateDomainOptions{
		DomainType: "ALLOWED_EMAIL_DOMAIN",
	}
	domainWithOptions, err := client.Domain().CreateDomain(context.Background(), testOrg, domainWithOptionsName, options)
	assert.NoError(t, err)
	assert.NotNil(t, domainWithOptions)
	assert.NotNil(t, domainWithOptions.Domain)
	assert.Equal(t, domainWithOptionsName, domainWithOptions.Domain.Domain)
	assert.Equal(t, testOrg, domainWithOptions.Domain.OrganizationId)
	assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, domainWithOptions.Domain.DomainType)
}

func TestCreateDomainWithStringTypes(t *testing.T) {
	// Test creating a domain with string domain type
	stringDomainName := fmt.Sprintf("string-domain-%d.com", time.Now().Unix())
	options := &scalekit.CreateDomainOptions{
		DomainType: "ALLOWED_EMAIL_DOMAIN",
	}
	stringDomain, err := client.Domain().CreateDomain(context.Background(), testOrg, stringDomainName, options)
	assert.NoError(t, err)
	assert.NotNil(t, stringDomain)
	assert.NotNil(t, stringDomain.Domain)
	assert.Equal(t, stringDomainName, stringDomain.Domain.Domain)
	assert.Equal(t, testOrg, stringDomain.Domain.OrganizationId)
	assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, stringDomain.Domain.DomainType)

	// Test creating a domain with another string domain type
	stringOrgDomainName := fmt.Sprintf("string-org-%d.com", time.Now().Unix())
	orgOptions := &scalekit.CreateDomainOptions{
		DomainType: "ORGANIZATION_DOMAIN",
	}
	stringOrgDomain, err := client.Domain().CreateDomain(context.Background(), testOrg, stringOrgDomainName, orgOptions)
	assert.NoError(t, err)
	assert.NotNil(t, stringOrgDomain)
	assert.NotNil(t, stringOrgDomain.Domain)
	assert.Equal(t, stringOrgDomainName, stringOrgDomain.Domain.Domain)
	assert.Equal(t, testOrg, stringOrgDomain.Domain.OrganizationId)
	assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, stringOrgDomain.Domain.DomainType)
}

func TestDeleteDomain(t *testing.T) {
	// First create a domain to delete
	domainName := fmt.Sprintf("delete-test-%d.com", time.Now().Unix())
	domain, err := client.Domain().CreateDomain(context.Background(), testOrg, domainName)
	assert.NoError(t, err)
	assert.NotNil(t, domain)
	assert.NotNil(t, domain.Domain)

	domainId := domain.Domain.Id
	assert.NotEmpty(t, domainId)

	// Verify the domain exists
	retrievedDomain, err := client.Domain().GetDomain(context.Background(), domainId, testOrg)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedDomain)
	assert.Equal(t, domainName, retrievedDomain.Domain.Domain)

	// Delete the domain
	err = client.Domain().DeleteDomain(context.Background(), domainId, testOrg)
	assert.NoError(t, err)

	// Verify the domain is deleted by trying to get it
	_, err = client.Domain().GetDomain(context.Background(), domainId, testOrg)
	assert.Error(t, err) // Should return an error since domain is deleted
}

func TestListDomains(t *testing.T) {
	// Create a domain with ORGANIZATION_DOMAIN type
	orgDomainName := fmt.Sprintf("list-org-%d.com", time.Now().Unix())
	orgDomain, err := client.Domain().CreateDomain(context.Background(), testOrg, orgDomainName, &scalekit.CreateDomainOptions{
		DomainType: scalekit.DomainTypeOrganization,
	})
	assert.NoError(t, err)
	assert.NotNil(t, orgDomain)
	assert.NotNil(t, orgDomain.Domain)
	assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, orgDomain.Domain.DomainType)

	orgDomainId := orgDomain.Domain.Id
	assert.NotEmpty(t, orgDomainId)

	// Create a domain with ALLOWED_EMAIL_DOMAIN type
	allowedEmailDomainName := fmt.Sprintf("list-allowed-email-%d.com", time.Now().Unix())
	allowedEmailDomain, err := client.Domain().CreateDomain(context.Background(), testOrg, allowedEmailDomainName, &scalekit.CreateDomainOptions{
		DomainType: scalekit.DomainTypeAllowedEmail,
	})
	assert.NoError(t, err)
	assert.NotNil(t, allowedEmailDomain)
	assert.NotNil(t, allowedEmailDomain.Domain)
	assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, allowedEmailDomain.Domain.DomainType)

	allowedEmailDomainId := allowedEmailDomain.Domain.Id
	assert.NotEmpty(t, allowedEmailDomainId)

	// Verify the domains exist
	retrievedOrgDomain, err := client.Domain().GetDomain(context.Background(), orgDomainId, testOrg)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedOrgDomain)
	assert.Equal(t, orgDomainName, retrievedOrgDomain.Domain.Domain)

	retrievedAllowedEmailDomain, err := client.Domain().GetDomain(context.Background(), allowedEmailDomainId, testOrg)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedAllowedEmailDomain)
	assert.Equal(t, allowedEmailDomainName, retrievedAllowedEmailDomain.Domain.Domain)

	// Test listing all domains (no filter)
	allDomainsList, err := client.Domain().ListDomains(context.Background(), testOrg)
	assert.NoError(t, err)
	assert.True(t, len(allDomainsList.Domains) > 0)

	foundOrgDomain := false
	foundAllowedEmailDomain := false
	for _, d := range allDomainsList.Domains {
		if d.Id == orgDomainId {
			foundOrgDomain = true
			assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, d.DomainType)
		}
		if d.Id == allowedEmailDomainId {
			foundAllowedEmailDomain = true
			assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, d.DomainType)
		}
	}
	assert.True(t, foundOrgDomain, "Organization domain should be in the list")
	assert.True(t, foundAllowedEmailDomain, "Allowed email domain should be in the list")

	// Test listing domains filtered by ORGANIZATION_DOMAIN type
	orgDomainsList, err := client.Domain().ListDomains(context.Background(), testOrg, &scalekit.ListDomainOptions{
		DomainType: scalekit.DomainTypeOrganization,
	})
	assert.NoError(t, err)
	assert.True(t, len(orgDomainsList.Domains) > 0)

	foundOrgDomainInFiltered := false
	foundAllowedEmailDomainInFiltered := false
	for _, d := range orgDomainsList.Domains {
		// All domains in the filtered list should be ORGANIZATION_DOMAIN type
		assert.Equal(t, domains.DomainType_ORGANIZATION_DOMAIN, d.DomainType)
		if d.Id == orgDomainId {
			foundOrgDomainInFiltered = true
		}
		if d.Id == allowedEmailDomainId {
			foundAllowedEmailDomainInFiltered = true
		}
	}
	assert.True(t, foundOrgDomainInFiltered, "Organization domain should be in the filtered list")
	assert.False(t, foundAllowedEmailDomainInFiltered, "Allowed email domain should NOT be in the ORGANIZATION_DOMAIN filtered list")

	// Test listing domains filtered by ALLOWED_EMAIL_DOMAIN type
	allowedEmailDomainsList, err := client.Domain().ListDomains(context.Background(), testOrg, &scalekit.ListDomainOptions{
		DomainType: scalekit.DomainTypeAllowedEmail,
	})
	assert.NoError(t, err)
	assert.True(t, len(allowedEmailDomainsList.Domains) > 0)

	foundOrgDomainInEmailFiltered := false
	foundAllowedEmailDomainInEmailFiltered := false
	for _, d := range allowedEmailDomainsList.Domains {
		// All domains in the filtered list should be ALLOWED_EMAIL_DOMAIN type
		assert.Equal(t, domains.DomainType_ALLOWED_EMAIL_DOMAIN, d.DomainType)
		if d.Id == orgDomainId {
			foundOrgDomainInEmailFiltered = true
		}
		if d.Id == allowedEmailDomainId {
			foundAllowedEmailDomainInEmailFiltered = true
		}
	}
	assert.False(t, foundOrgDomainInEmailFiltered, "Organization domain should NOT be in the ALLOWED_EMAIL_DOMAIN filtered list")
	assert.True(t, foundAllowedEmailDomainInEmailFiltered, "Allowed email domain should be in the filtered list")
}
