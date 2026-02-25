package test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Cohesive, prod-like test data constants.
const (
	TestOrgName          = "Acme Corp"
	TestDomainRoot       = "acmecorp.com"
	TestDomainAuth       = "auth.acmecorp.com"
	TestDomainIdp        = "idp.acmecorp.com"
	TestUserEmailJohn    = "john@example.com"
	TestUserEmailJane    = "jane@example.com"
	TestConnectionName   = "Acme Corp SSO"
	TestDirectoryName    = "Acme Corp Directory"
	TestExternalIDPrefix = "acmecorp"
)

// UniqueSuffix returns a short unique suffix for test resources (e.g. acmecorp-20060102150405).
func UniqueSuffix() string {
	return fmt.Sprintf("%s-%d", TestExternalIDPrefix, time.Now().UnixNano()/1e6)
}

// DeleteTestOrganization deletes the org by ID. Idempotent: ignores "not found" errors.
func DeleteTestOrganization(t *testing.T, ctx context.Context, orgID string) {
	t.Helper()
	if orgID == "" {
		return
	}
	if client == nil {
		return
	}
	err := client.Organization().DeleteOrganization(ctx, orgID)
	if err != nil {
		// Idempotent: ignore not-found so defer is safe after partial failure
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return
		}
		t.Logf("deleteTestOrganization %s: %v", orgID, err)
	}
}

// DeleteTestDomain deletes the domain by ID. Idempotent: ignores "not found" errors.
func DeleteTestDomain(t *testing.T, ctx context.Context, orgID, domainID string) {
	t.Helper()
	if orgID == "" || domainID == "" {
		return
	}
	if client == nil {
		return
	}
	err := client.Domain().DeleteDomain(ctx, domainID, orgID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return
		}
		t.Logf("deleteTestDomain %s: %v", domainID, err)
	}
}

// DeleteTestConnection deletes the connection by ID. Idempotent: ignores "not found" errors.
func DeleteTestConnection(t *testing.T, ctx context.Context, orgID, connectionID string) {
	t.Helper()
	if orgID == "" || connectionID == "" {
		return
	}
	if client == nil {
		return
	}
	err := client.Connection().DeleteConnection(ctx, orgID, connectionID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return
		}
		t.Logf("deleteTestConnection %s: %v", connectionID, err)
	}
}

// DeleteTestDirectory deletes the directory by ID. Idempotent: ignores "not found" errors.
func DeleteTestDirectory(t *testing.T, ctx context.Context, orgID, directoryID string) {
	t.Helper()
	if orgID == "" || directoryID == "" {
		return
	}
	if client == nil {
		return
	}
	err := client.Directory().DeleteDirectory(ctx, orgID, directoryID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return
		}
		t.Logf("deleteTestDirectory %s: %v", directoryID, err)
	}
}
