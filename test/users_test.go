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
)

func TestUsers(t *testing.T) {
	// Test listing users by organization
	usersList, err := client.User().ListOrganizationUsers(context.Background(), testOrg, &scalekit.ListUsersOptions{
		PageSize:  10,
		PageToken: "",
	})
	assert.NoError(t, err)
	assert.NotNil(t, usersList)
	assert.True(t, len(usersList.Users) > 0)

	// Test getting user by ID
	firstUser := usersList.Users[0]
	assert.NotEmpty(t, firstUser.Id)
	assert.NotEmpty(t, firstUser.Email)

	user, err := client.User().GetUser(context.Background(), firstUser.Id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, user.User)
	assert.Equal(t, firstUser.Id, user.User.Id)
	assert.Equal(t, firstUser.Email, user.User.Email)
	assert.NotEmpty(t, user.User.EnvironmentId)

	// Test updating user
	firstName := "Test"
	lastName := "User"
	name := "Test User"
	locale := "en-US"
	updateRequest := &users.UpdateUser{
		UserProfile: &users.UpdateUserProfile{
			FirstName: &firstName,
			LastName:  &lastName,
			Name:      &name,
			Locale:    &locale,
		},
	}
	updatedUser, err := client.User().UpdateUser(context.Background(), firstUser.Id, updateRequest)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.NotNil(t, updatedUser.User)
	assert.NotNil(t, updatedUser.User.UserProfile)
	assert.Equal(t, firstUser.Id, updatedUser.User.Id)
	assert.Equal(t, firstUser.Email, updatedUser.User.Email)

	// Verify the profile fields were updated
	assert.Equal(t, "Test", updatedUser.User.UserProfile.FirstName)
	assert.Equal(t, "User", updatedUser.User.UserProfile.LastName)
	assert.Equal(t, "Test User", updatedUser.User.UserProfile.Name)
	assert.Equal(t, "en-US", updatedUser.User.UserProfile.Locale)

	// Test listing users with pagination
	paginatedUsers, err := client.User().ListOrganizationUsers(context.Background(), testOrg, &scalekit.ListUsersOptions{
		PageSize:  5,
		PageToken: usersList.NextPageToken,
	})
	assert.NoError(t, err)
	assert.NotNil(t, paginatedUsers)
	assert.True(t, len(paginatedUsers.Users) > 0)

	// Assert basic attributes for paginated users
	for _, u := range paginatedUsers.Users {
		assert.NotEmpty(t, u.Id)
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.EnvironmentId)
	}
}

func TestUserOperations(t *testing.T) {
	// Create a new user
	newUser := &users.CreateUser{
		Email: "testin.user@example.com",
		Metadata: map[string]string{
			"source": "test",
		},
	}

	var userId string

	// Try to create the user first
	createdUser, err := client.User().CreateUserAndMembership(context.Background(), testOrg, newUser, false)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.NotNil(t, createdUser.User)
	assert.NotEmpty(t, createdUser.User.Id)
	assert.Equal(t, newUser.Email, createdUser.User.Email)
	assert.NotEmpty(t, createdUser.User.EnvironmentId)
	assert.NotNil(t, createdUser.User.CreateTime)
	assert.NotNil(t, createdUser.User.UpdateTime)
	userId = createdUser.User.Id

	// Get the user to check their organization membership
	user, err := client.User().GetUser(context.Background(), userId)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, user.User)
	assert.Equal(t, userId, user.User.Id)
	assert.Equal(t, newUser.Email, user.User.Email)
	assert.NotEmpty(t, user.User.EnvironmentId)
	assert.NotNil(t, user.User.CreateTime)
	assert.NotNil(t, user.User.UpdateTime)

	// Delete user
	err = client.User().DeleteUser(context.Background(), userId)
	assert.NoError(t, err)

	// Verify user is deleted
	_, err = client.User().GetUser(context.Background(), userId)
	assert.Error(t, err)
}

func TestMembershipOperations(t *testing.T) {
	// Create a new user first with unique email using timestamp
	timestamp := time.Now().Unix()
	uniqueEmail := fmt.Sprintf("membership.test.%d@example.com", timestamp)

	newUser := &users.CreateUser{
		Email: uniqueEmail,
		Metadata: map[string]string{
			"source": "membership_test",
		},
	}

	// Use testOrg for user creation
	createdUser, err := client.User().CreateUserAndMembership(context.Background(), testOrg, newUser, false)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.NotNil(t, createdUser.User)
	assert.NotEmpty(t, createdUser.User.Id)
	assert.Equal(t, uniqueEmail, createdUser.User.Email)
	assert.NotEmpty(t, createdUser.User.EnvironmentId)
	assert.NotNil(t, createdUser.User.CreateTime)
	assert.NotNil(t, createdUser.User.UpdateTime)
	userId := createdUser.User.Id

	// Use testOrg2 for membership operations
	membership := &users.CreateMembership{
		Roles: []*commons.Role{
			{
				Name: "admin",
			},
		},
		Metadata: map[string]string{
			"membership_type": "test",
		},
	}

	membershipResponse, err := client.User().CreateMembership(context.Background(), testOrg2, userId, membership, false)
	assert.NoError(t, err)
	assert.NotNil(t, membershipResponse)
	if membershipResponse != nil {
		assert.NotNil(t, membershipResponse.User)
		assert.Equal(t, userId, membershipResponse.User.Id)
		assert.Equal(t, uniqueEmail, membershipResponse.User.Email)
		assert.NotEmpty(t, membershipResponse.User.EnvironmentId)
	}

	updateMembership := &users.UpdateMembership{
		Roles: []*commons.Role{
			{
				Name: "member",
			},
		},
		Metadata: map[string]string{
			"membership_type": "updated_test",
		},
	}

	updatedMembershipResponse, err := client.User().UpdateMembership(context.Background(), testOrg2, userId, updateMembership)
	assert.NoError(t, err)
	assert.NotNil(t, updatedMembershipResponse)
	if updatedMembershipResponse != nil {
		assert.NotNil(t, updatedMembershipResponse.User)
		assert.Equal(t, userId, updatedMembershipResponse.User.Id)
		assert.Equal(t, uniqueEmail, updatedMembershipResponse.User.Email)
		assert.NotEmpty(t, updatedMembershipResponse.User.EnvironmentId)
	}

	err = client.User().DeleteMembership(context.Background(), testOrg2, userId, false)
	assert.NoError(t, err)

	userAfterDelete, err := client.User().GetUser(context.Background(), userId)
	assert.NoError(t, err)
	assert.NotNil(t, userAfterDelete)
	assert.NotNil(t, userAfterDelete.User)
	assert.Equal(t, userId, userAfterDelete.User.Id)
	assert.Equal(t, uniqueEmail, userAfterDelete.User.Email)
	assert.NotEmpty(t, userAfterDelete.User.EnvironmentId)

	err = client.User().DeleteUser(context.Background(), userId)
	assert.NoError(t, err)
}

func TestResendInvite(t *testing.T) {
	// Create a new user first with unique email using timestamp
	timestamp := time.Now().Unix()
	uniqueEmail := fmt.Sprintf("resend.invite.test.%d@example.com", timestamp)

	newUser := &users.CreateUser{
		Email: uniqueEmail,
		Metadata: map[string]string{
			"source": "resend_invite_test",
		},
	}

	// Create user with invitation email
	createdUser, err := client.User().CreateUserAndMembership(context.Background(), testOrg, newUser, true)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.NotNil(t, createdUser.User)
	assert.NotEmpty(t, createdUser.User.Id)
	assert.Equal(t, uniqueEmail, createdUser.User.Email)
	userId := createdUser.User.Id

	// Resend invite
	resendResponse, err := client.User().ResendInvite(context.Background(), testOrg, userId)
	assert.NoError(t, err)
	assert.NotNil(t, resendResponse)
	assert.NotNil(t, resendResponse.Invite)
	assert.Equal(t, userId, resendResponse.Invite.UserId)
	assert.Equal(t, testOrg, resendResponse.Invite.OrganizationId)
	assert.Equal(t, "PENDING_INVITE", resendResponse.Invite.Status)
	assert.NotNil(t, resendResponse.Invite.CreatedAt)
	assert.NotNil(t, resendResponse.Invite.ExpiresAt)
	assert.Equal(t, int32(1), resendResponse.Invite.ResentCount)

	// Clean up
	err = client.User().DeleteUser(context.Background(), userId)
	assert.NoError(t, err)
}
