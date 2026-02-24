package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/commons"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_ListOrganizationUsers(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	// Create 2 users so we can verify list, get, and update
	email1 := fmt.Sprintf("john.list.%d@example.com", time.Now().UnixNano()/1e6)
	email2 := fmt.Sprintf("jane.list.%d@example.com", time.Now().UnixNano()/1e6)
	user1, err := client.User().CreateUserAndMembership(ctx, orgId, &users.CreateUser{Email: email1, Metadata: map[string]string{"source": "list_test"}}, false)
	require.NoError(t, err)
	require.NotNil(t, user1)
	require.NotEmpty(t, user1.GetUser().GetId())
	defer func() { _ = client.User().DeleteUser(ctx, user1.GetUser().GetId()) }()
	user2, err := client.User().CreateUserAndMembership(ctx, orgId, &users.CreateUser{Email: email2, Metadata: map[string]string{"source": "list_test"}}, false)
	require.NoError(t, err)
	require.NotNil(t, user2)
	require.NotEmpty(t, user2.GetUser().GetId())
	defer func() { _ = client.User().DeleteUser(ctx, user2.GetUser().GetId()) }()

	usersList, err := client.User().ListOrganizationUsers(ctx, orgId, &scalekit.ListUsersOptions{
		PageSize:  10,
		PageToken: "",
	})
	require.NoError(t, err)
	require.NotNil(t, usersList)
	require.GreaterOrEqual(t, len(usersList.GetUsers()), 2, "expected at least 2 users from list")

	// Verify both created users appear in list
	ids := make(map[string]bool)
	for _, u := range usersList.GetUsers() {
		ids[u.GetId()] = true
		require.NotEmpty(t, u.GetId())
		require.NotEmpty(t, u.GetEmail())
	}
	require.True(t, ids[user1.GetUser().GetId()], "user1 should be in list")
	require.True(t, ids[user2.GetUser().GetId()], "user2 should be in list")

	// Get and verify first user
	firstUser := usersList.GetUsers()[0]
	user, err := client.User().GetUser(ctx, firstUser.GetId())
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, user.GetUser())
	assert.Equal(t, firstUser.GetId(), user.GetUser().GetId())
	assert.Equal(t, firstUser.GetEmail(), user.GetUser().GetEmail())
	assert.NotEmpty(t, user.GetUser().GetEnvironmentId())

	// Update first user profile (use GivenName/FamilyName; FirstName/LastName are deprecated)
	givenName := "John"
	familyName := "Doe"
	name := "John Doe"
	locale := "en-US"
	updateRequest := &users.UpdateUser{
		UserProfile: &users.UpdateUserProfile{
			GivenName:  &givenName,
			FamilyName: &familyName,
			Name:       &name,
			Locale:     &locale,
		},
	}
	updatedUser, err := client.User().UpdateUser(ctx, firstUser.GetId(), updateRequest)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.NotNil(t, updatedUser.GetUser().GetUserProfile())
	assert.Equal(t, "John", updatedUser.GetUser().GetUserProfile().GetGivenName())
	assert.Equal(t, "Doe", updatedUser.GetUser().GetUserProfile().GetFamilyName())
	assert.Equal(t, "John Doe", updatedUser.GetUser().GetUserProfile().GetName())
	assert.Equal(t, "en-US", updatedUser.GetUser().GetUserProfile().GetLocale())

	// If pagination is supported, exercise it
	if usersList.GetNextPageToken() != "" {
		paginatedUsers, err := client.User().ListOrganizationUsers(ctx, orgId, &scalekit.ListUsersOptions{
			PageSize:  5,
			PageToken: usersList.GetNextPageToken(),
		})
		require.NoError(t, err)
		require.NotNil(t, paginatedUsers)
		for _, u := range paginatedUsers.GetUsers() {
			assert.NotEmpty(t, u.GetId())
			assert.NotEmpty(t, u.GetEmail())
		}
	}
}

func TestUser_EndToEndIntegration(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	uniqueEmail := fmt.Sprintf("testin.user.%d@example.com", time.Now().UnixNano()/1e6)
	newUser := &users.CreateUser{
		Email: uniqueEmail,
		Metadata: map[string]string{
			"source": "test",
		},
	}

	createdUser, err := client.User().CreateUserAndMembership(ctx, orgId, newUser, false)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotNil(t, createdUser.GetUser())
	require.NotEmpty(t, createdUser.GetUser().GetId())
	assert.Equal(t, uniqueEmail, createdUser.GetUser().GetEmail())
	userId := createdUser.GetUser().GetId()

	user, err := client.User().GetUser(ctx, userId)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, user.GetUser())
	assert.Equal(t, userId, user.GetUser().GetId())
	assert.Equal(t, uniqueEmail, user.GetUser().GetEmail())

	err = client.User().DeleteUser(ctx, userId)
	require.NoError(t, err)

	_, err = client.User().GetUser(ctx, userId)
	assert.Error(t, err)
}

func TestUser_MembershipOperations(t *testing.T) {
	ctx := context.Background()
	orgId1 := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId1)
	orgId2 := createOrg(t, ctx, "Acme Corp 2", UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId2)

	uniqueEmail := fmt.Sprintf("membership.test.%d@example.com", time.Now().UnixNano()/1e6)
	newUser := &users.CreateUser{
		Email: uniqueEmail,
		Metadata: map[string]string{
			"source": "membership_test",
		},
	}

	createdUser, err := client.User().CreateUserAndMembership(ctx, orgId1, newUser, false)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotEmpty(t, createdUser.GetUser().GetId())
	userId := createdUser.GetUser().GetId()
	defer func() {
		_ = client.User().DeleteUser(ctx, userId)
	}()

	membership := &users.CreateMembership{
		Roles: []*commons.Role{
			{Name: "admin"},
		},
		Metadata: map[string]string{
			"membership_type": "test",
		},
	}
	membershipResponse, err := client.User().CreateMembership(ctx, orgId2, userId, membership, false)
	require.NoError(t, err)
	require.NotNil(t, membershipResponse)
	require.NotNil(t, membershipResponse.GetUser())
	assert.Equal(t, userId, membershipResponse.GetUser().GetId())
	assert.Equal(t, uniqueEmail, membershipResponse.GetUser().GetEmail())

	updateMembership := &users.UpdateMembership{
		Roles: []*commons.Role{
			{Name: "member"},
		},
		Metadata: map[string]string{
			"membership_type": "updated_test",
		},
	}
	updatedMembershipResponse, err := client.User().UpdateMembership(ctx, orgId2, userId, updateMembership)
	require.NoError(t, err)
	require.NotNil(t, updatedMembershipResponse)

	err = client.User().DeleteMembership(ctx, orgId2, userId, false)
	require.NoError(t, err)

	userAfterDelete, err := client.User().GetUser(ctx, userId)
	require.NoError(t, err)
	require.NotNil(t, userAfterDelete)
	require.NotNil(t, userAfterDelete.GetUser())
	assert.Equal(t, userId, userAfterDelete.GetUser().GetId())
}

func TestUser_ResendInvite(t *testing.T) {
	ctx := context.Background()
	orgId := createOrg(t, ctx, TestOrgName, UniqueSuffix())
	defer DeleteTestOrganization(t, ctx, orgId)

	uniqueEmail := fmt.Sprintf("resend.invite.test.%d@example.com", time.Now().UnixNano()/1e6)
	newUser := &users.CreateUser{
		Email: uniqueEmail,
		Metadata: map[string]string{
			"source": "resend_invite_test",
		},
	}

	createdUser, err := client.User().CreateUserAndMembership(ctx, orgId, newUser, true)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotNil(t, createdUser.GetUser())
	require.NotEmpty(t, createdUser.GetUser().GetId())
	userId := createdUser.GetUser().GetId()
	defer func() {
		_ = client.User().DeleteUser(ctx, userId)
	}()

	resendResponse, err := client.User().ResendInvite(ctx, orgId, userId)
	require.NoError(t, err)
	require.NotNil(t, resendResponse)
	require.NotNil(t, resendResponse.GetInvite())
	assert.Equal(t, userId, resendResponse.GetInvite().GetUserId())
	assert.Equal(t, orgId, resendResponse.GetInvite().GetOrganizationId())
	assert.Equal(t, "PENDING_INVITE", resendResponse.GetInvite().GetStatus())
	require.NotNil(t, resendResponse.GetInvite().GetCreatedAt())
	require.NotNil(t, resendResponse.GetInvite().GetExpiresAt())
}
