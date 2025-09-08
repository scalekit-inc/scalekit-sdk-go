package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetConnectedAccountAuth(t *testing.T) {
	// Test with connection name and identifier
	options := &scalekit.GetConnectedAccountAuthOptions{
		Connector:  toPtr("GITHUB"),
		Identifier: toPtr("avinash.kamath@scalekit.com"),
	}

	authResp, err := client.ConnectedAccount().GetConnectedAccountAuth(context.Background(), options)
	assert.NoError(t, err)
	assert.NotNil(t, authResp)
	assert.NotNil(t, authResp.ConnectedAccount)
	fmt.Println("Connected Account Auth Response:", authResp.GetConnectedAccount().GetAuthorizationDetails().GetOauthToken())

	// Verify the connected account details
	connectedAccount := authResp.ConnectedAccount
	assert.NotEmpty(t, connectedAccount.Id)
	assert.Equal(t, "avinash.kamath@scalekit.com", connectedAccount.Identifier)
	assert.Equal(t, "GITHUB", connectedAccount.Connector)
	assert.NotNil(t, connectedAccount.AuthorizationDetails)
	assert.NotNil(t, connectedAccount.UpdatedAt)
}
