package test

import (
	"context"
	"testing"

	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
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
			err := authService.UpdateLoginUserDetails(ctx, tc.req)
			assert.Error(t, err)
		})
	}
}
