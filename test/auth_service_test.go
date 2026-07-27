package test

import (
	"context"
	"os"
	"testing"

	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthServiceUpdateLoginUserDetailsValidation(t *testing.T) {
	authService := client.Auth()
	ctx := context.Background()

	makeReq := func() *scalekit.UpdateLoginUserDetailsRequest {
		return &scalekit.UpdateLoginUserDetailsRequest{
			ConnectionId:   "conn",
			LoginRequestId: "login",
			User: &scalekit.LoggedInUserDetails{
				Sub:   "sub",
				Email: "user@example.com",
			},
		}
	}

	tests := []struct {
		name    string
		req     *scalekit.UpdateLoginUserDetailsRequest
		wantErr bool
	}{
		{name: "nil request", req: nil, wantErr: true},
		{name: "missing connection id", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.ConnectionId = ""; return r }(), wantErr: true},
		{name: "missing login request id", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.LoginRequestId = ""; return r }(), wantErr: true},
		{name: "missing user", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.User = nil; return r }(), wantErr: true},
		{name: "missing sub", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.User.Sub = ""; return r }(), wantErr: true},
		{name: "missing email", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.User.Email = ""; return r }(), wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := authService.UpdateLoginUserDetails(ctx, tc.req)
			assert.Error(t, err)
		})
	}
}

// TestAuthServiceUpdateLoginUserDetailsLive exercises the happy path when a real
// connection id and login request id are supplied via the environment, and
// asserts the typed response carries the auth request id. Skipped otherwise.
func TestAuthServiceUpdateLoginUserDetailsLive(t *testing.T) {
	connectionId := os.Getenv("SCALEKIT_TEST_CONNECTION_ID")
	loginRequestId := os.Getenv("SCALEKIT_TEST_LOGIN_REQUEST_ID")
	if connectionId == "" || loginRequestId == "" {
		t.Skip("set SCALEKIT_TEST_CONNECTION_ID and SCALEKIT_TEST_LOGIN_REQUEST_ID to run this test")
	}

	ctx := context.Background()
	resp, err := client.Auth().UpdateLoginUserDetails(ctx, &scalekit.UpdateLoginUserDetailsRequest{
		ConnectionId:   connectionId,
		LoginRequestId: loginRequestId,
		User: &scalekit.LoggedInUserDetails{
			Sub:   "sub",
			Email: "user@example.com",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.GetAuthRequestId())
}
