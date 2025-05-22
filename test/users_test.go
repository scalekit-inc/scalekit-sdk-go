package test

import (
	"context"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go"
	"github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/commons"
	"github.com/scalekit-inc/scalekit-sdk-go/pkg/grpc/scalekit/v1/users"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	// Test listing users by organization
	usersList, err := client.User().ListUsers(context.Background(), testOrg, &scalekit.ListUsersOptions{
		PageSize:  10,
		PageToken: "",
	})
	assert.NoError(t, err)
	assert.NotNil(t, usersList)
	assert.True(t, len(usersList.Users) > 0)

	// Test getting user by ID
	firstUser := usersList.Users[0]
	user, err := client.User().GetUser(context.Background(), testOrg, firstUser.Id)
	assert.NoError(t, err)
	assert.Equal(t, firstUser.Id, user.User.Id)
	assert.Equal(t, firstUser.Email, user.User.Email)

	// Test updating user
	updateRequest := &users.UpdateUser{
		UserProfile: &commons.UserProfile{
			FirstName: "Test",
			LastName:  "User",
			Name:      "Test User",
			Locale:    "en-US",
		},
	}
	updatedUser, err := client.User().UpdateUser(context.Background(), testOrg, firstUser.Id, updateRequest)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.NotNil(t, updatedUser.User)
	assert.NotNil(t, updatedUser.User.UserProfile)

	// Verify the profile fields were updated
	assert.Equal(t, "Test", updatedUser.User.UserProfile.FirstName)
	assert.Equal(t, "User", updatedUser.User.UserProfile.LastName)
	assert.Equal(t, "Test User", updatedUser.User.UserProfile.Name)
	assert.Equal(t, "en-US", updatedUser.User.UserProfile.Locale)

	// Test listing users with pagination
	paginatedUsers, err := client.User().ListUsers(context.Background(), testOrg, &scalekit.ListUsersOptions{
		PageSize:  5,
		PageToken: usersList.NextPageToken,
	})
	assert.NoError(t, err)
	assert.NotNil(t, paginatedUsers)
	assert.True(t, len(paginatedUsers.Users) > 0)
}

func TestUserOperations(t *testing.T) {
	// Create a new user
	newUser := &users.User{
		Email: "testin.user@example.com",
		Metadata: map[string]string{
			"source": "test",
		},
	}

	var userId string

	// Try to create the user first
	createdUser, err := client.User().CreateUser(context.Background(), testOrg, newUser)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, newUser.Email, createdUser.User.Email)
	userId = createdUser.User.Id

	// Get the user to check their organization membership
	user, err := client.User().GetUser(context.Background(), testOrg, userId)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Delete user
	err = client.User().DeleteUser(context.Background(), testOrg, userId)
	assert.NoError(t, err)

	// Verify user is deleted
	_, err = client.User().GetUser(context.Background(), testOrg, userId)
	assert.Error(t, err)
}
