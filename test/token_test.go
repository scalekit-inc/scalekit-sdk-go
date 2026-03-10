package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	ctx := context.Background()

	created, err := client.Token().CreateToken(ctx, testOrg, scalekit.CreateTokenOptions{
		Description: "test token",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	t.Cleanup(func() {
		_ = client.Token().InvalidateToken(ctx, created.Token)
	})

	assert.NotEmpty(t, created.Token)
	assert.NotEmpty(t, created.TokenId)
	assert.NotNil(t, created.TokenInfo)
	assert.Equal(t, testOrg, created.TokenInfo.OrganizationId)
}

func TestCreateTokenWithCustomClaims(t *testing.T) {
	ctx := context.Background()

	claims := map[string]string{
		"team":        "engineering",
		"environment": "test",
	}

	created, err := client.Token().CreateToken(ctx, testOrg, scalekit.CreateTokenOptions{
		Description:  "token with claims",
		CustomClaims: claims,
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	t.Cleanup(func() {
		_ = client.Token().InvalidateToken(ctx, created.Token)
	})

	assert.NotEmpty(t, created.Token)
	assert.NotEmpty(t, created.TokenId)
	assert.NotNil(t, created.TokenInfo)
	assert.Equal(t, "engineering", created.TokenInfo.CustomClaims["team"])
	assert.Equal(t, "test", created.TokenInfo.CustomClaims["environment"])
}

func TestCreateUserScopedToken(t *testing.T) {
	ctx := context.Background()

	// Create a user with sendInvitationEmail=false to get an active membership
	timestamp := time.Now().Unix()
	uniqueEmail := fmt.Sprintf("token.test.%d@example.com", timestamp)
	newUser := &users.CreateUser{
		Email: uniqueEmail,
		Metadata: map[string]string{
			"source": "token_test",
		},
	}

	createdUser, err := client.User().CreateUserAndMembership(ctx, testOrg, newUser, false)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	t.Cleanup(func() {
		_ = client.User().DeleteUser(ctx, createdUser.User.Id)
	})

	userId := createdUser.User.Id

	// Create a user-scoped token
	created, err := client.Token().CreateToken(ctx, testOrg, scalekit.CreateTokenOptions{
		UserId:      userId,
		Description: "user scoped token",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	t.Cleanup(func() {
		_ = client.Token().InvalidateToken(ctx, created.Token)
	})

	assert.NotEmpty(t, created.Token)
	assert.NotEmpty(t, created.TokenId)
	assert.NotNil(t, created.TokenInfo)
	assert.Equal(t, testOrg, created.TokenInfo.OrganizationId)
	require.NotNil(t, created.TokenInfo.UserId)
	assert.Equal(t, userId, *created.TokenInfo.UserId)
}

func TestValidateTokenByOpaqueToken(t *testing.T) {
	ctx := context.Background()

	created, err := client.Token().CreateToken(ctx, testOrg, scalekit.CreateTokenOptions{
		Description: "validate test token",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	t.Cleanup(func() {
		_ = client.Token().InvalidateToken(ctx, created.Token)
	})

	validated, err := client.Token().ValidateToken(ctx, created.Token)
	assert.NoError(t, err)
	assert.NotNil(t, validated)
	assert.NotNil(t, validated.TokenInfo)
	assert.Equal(t, created.TokenId, validated.TokenInfo.TokenId)
	assert.Equal(t, testOrg, validated.TokenInfo.OrganizationId)
}

func TestValidateTokenByTokenId(t *testing.T) {
	ctx := context.Background()

	created, err := client.Token().CreateToken(ctx, testOrg, scalekit.CreateTokenOptions{
		Description: "validate by id test",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	t.Cleanup(func() {
		_ = client.Token().InvalidateToken(ctx, created.Token)
	})

	validated, err := client.Token().ValidateToken(ctx, created.TokenId)
	assert.NoError(t, err)
	assert.NotNil(t, validated)
	assert.NotNil(t, validated.TokenInfo)
	assert.Equal(t, created.TokenId, validated.TokenInfo.TokenId)
	assert.Equal(t, testOrg, validated.TokenInfo.OrganizationId)
}

func TestListTokens(t *testing.T) {
	ctx := context.Background()

	// Create a token to ensure there's at least one
	created, err := client.Token().CreateToken(ctx, testOrg, scalekit.CreateTokenOptions{
		Description: "list test token",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	t.Cleanup(func() {
		_ = client.Token().InvalidateToken(ctx, created.Token)
	})

	listed, err := client.Token().ListTokens(ctx, testOrg, scalekit.ListTokensOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, listed)
	assert.True(t, len(listed.Tokens) > 0)
	assert.True(t, listed.TotalCount > 0)

	// Verify token fields
	for _, token := range listed.Tokens {
		assert.NotEmpty(t, token.TokenId)
		assert.Equal(t, testOrg, token.OrganizationId)
	}
}

func TestListTokensWithPagination(t *testing.T) {
	ctx := context.Background()

	// Create 3 tokens
	for i := 0; i < 3; i++ {
		created, err := client.Token().CreateToken(ctx, testOrg, scalekit.CreateTokenOptions{
			Description: fmt.Sprintf("pagination test token %d", i),
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = client.Token().InvalidateToken(ctx, created.Token)
		})
	}

	// List with page size 1
	page1, err := client.Token().ListTokens(ctx, testOrg, scalekit.ListTokensOptions{
		PageSize: 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, page1)
	assert.Equal(t, 1, len(page1.Tokens))
	assert.NotEmpty(t, page1.NextPageToken)

	// Get next page
	page2, err := client.Token().ListTokens(ctx, testOrg, scalekit.ListTokensOptions{
		PageSize:  1,
		PageToken: page1.NextPageToken,
	})
	assert.NoError(t, err)
	assert.NotNil(t, page2)
	assert.Equal(t, 1, len(page2.Tokens))

	// Ensure different tokens on different pages
	assert.NotEqual(t, page1.Tokens[0].TokenId, page2.Tokens[0].TokenId)
}

func TestInvalidateToken(t *testing.T) {
	ctx := context.Background()

	created, err := client.Token().CreateToken(ctx, testOrg, scalekit.CreateTokenOptions{
		Description: "invalidate test token",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	t.Cleanup(func() {
		// Safe to call even if already invalidated (idempotent)
		_ = client.Token().InvalidateToken(ctx, created.Token)
	})

	// Invalidate the token
	err = client.Token().InvalidateToken(ctx, created.Token)
	assert.NoError(t, err)

	// Validate should fail for invalidated token
	_, err = client.Token().ValidateToken(ctx, created.Token)
	assert.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrTokenValidationFailed)
}

func TestInvalidateTokenIdempotent(t *testing.T) {
	ctx := context.Background()

	created, err := client.Token().CreateToken(ctx, testOrg, scalekit.CreateTokenOptions{
		Description: "idempotent invalidate test",
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	// Invalidate twice - second should not error
	err = client.Token().InvalidateToken(ctx, created.Token)
	assert.NoError(t, err)

	err = client.Token().InvalidateToken(ctx, created.Token)
	assert.NoError(t, err)
}
