package test

import (
	"context"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationSessionPolicy_GetDefaultPolicy(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	policy, err := client.Organization().GetOrganizationSessionPolicy(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, policy)
	assert.Equal(t, scalekit.SessionPolicySourceApplication, policy.PolicySource,
		"new org should inherit APPLICATION policy by default")
}

func TestOrganizationSessionPolicy_SetCustomPolicy(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	timeout := int32(360)
	idleTimeout := int32(60)
	idleEnabled := true

	policy, err := client.Organization().UpdateOrganizationSessionPolicy(ctx, orgId, scalekit.OrganizationSessionPolicy{
		PolicySource:               scalekit.SessionPolicySourceCustom,
		AbsoluteSessionTimeout:     &timeout,
		AbsoluteSessionTimeoutUnit: scalekit.TimeUnitMinutes,
		IdleSessionTimeoutEnabled:  &idleEnabled,
		IdleSessionTimeout:         &idleTimeout,
		IdleSessionTimeoutUnit:     scalekit.TimeUnitMinutes,
	})
	require.NoError(t, err)
	require.NotNil(t, policy)
	assert.Equal(t, scalekit.SessionPolicySourceCustom, policy.PolicySource)

	fetched, err := client.Organization().GetOrganizationSessionPolicy(ctx, orgId)
	require.NoError(t, err)
	require.NotNil(t, fetched)
	assert.Equal(t, scalekit.SessionPolicySourceCustom, fetched.PolicySource)
	require.NotNil(t, fetched.AbsoluteSessionTimeout)
	assert.Equal(t, int32(360), fetched.AbsoluteSessionTimeout.GetValue())
	require.NotNil(t, fetched.IdleSessionTimeoutEnabled)
	assert.True(t, fetched.IdleSessionTimeoutEnabled.GetValue())
}

func TestOrganizationSessionPolicy_RevertToApplication(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	timeout := int32(120)
	_, err := client.Organization().UpdateOrganizationSessionPolicy(ctx, orgId, scalekit.OrganizationSessionPolicy{
		PolicySource:               scalekit.SessionPolicySourceCustom,
		AbsoluteSessionTimeout:     &timeout,
		AbsoluteSessionTimeoutUnit: scalekit.TimeUnitMinutes,
	})
	require.NoError(t, err)

	reverted, err := client.Organization().UpdateOrganizationSessionPolicy(ctx, orgId, scalekit.OrganizationSessionPolicy{
		PolicySource: scalekit.SessionPolicySourceApplication,
	})
	require.NoError(t, err)
	require.NotNil(t, reverted)
	assert.Equal(t, scalekit.SessionPolicySourceApplication, reverted.PolicySource)
}

func TestOrganizationSessionPolicy_SetIdleTimeoutDisabled(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	timeout := int32(480)
	idleEnabled := false

	policy, err := client.Organization().UpdateOrganizationSessionPolicy(ctx, orgId, scalekit.OrganizationSessionPolicy{
		PolicySource:               scalekit.SessionPolicySourceCustom,
		AbsoluteSessionTimeout:     &timeout,
		AbsoluteSessionTimeoutUnit: scalekit.TimeUnitMinutes,
		IdleSessionTimeoutEnabled:  &idleEnabled,
	})
	require.NoError(t, err)
	require.NotNil(t, policy)
	assert.Equal(t, scalekit.SessionPolicySourceCustom, policy.PolicySource)
	require.NotNil(t, policy.IdleSessionTimeoutEnabled)
	assert.False(t, policy.IdleSessionTimeoutEnabled.GetValue())
}
