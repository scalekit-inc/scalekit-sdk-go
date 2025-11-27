package test

import (
	"context"
	"testing"

	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	authv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/auth"
	"github.com/stretchr/testify/assert"
)

func TestAuthServiceUpdateLoginUserDetailsValidation(t *testing.T) {
	authService := client.Auth()
	ctx := context.Background()

	makeReq := func() *scalekit.UpdateLoginUserDetailsRequest {
		return &scalekit.UpdateLoginUserDetailsRequest{
			ConnectionId:   "conn",
			LoginRequestId: "login",
			User: &authv1.User{
				Sub:   "sub",
				Email: "user@example.com",
			},
		}
	}

	tests := []struct {
		name    string
		req     *scalekit.UpdateLoginUserDetailsRequest
		wantErr string
	}{
		{name: "nil request", req: nil, wantErr: "update login user details request is required"},
		{name: "missing connection id", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.ConnectionId = ""; return r }(), wantErr: "connectionId is required"},
		{name: "missing login request id", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.LoginRequestId = ""; return r }(), wantErr: "loginRequestId is required"},
		{name: "missing user", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.User = nil; return r }(), wantErr: "user details are required"},
		{name: "missing sub", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.User.Sub = ""; return r }(), wantErr: "user sub is required"},
		{name: "missing email", req: func() *scalekit.UpdateLoginUserDetailsRequest { r := makeReq(); r.User.Email = ""; return r }(), wantErr: "user email is required"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := authService.UpdateLoginUserDetails(ctx, tc.req)
			assert.EqualError(t, err, tc.wantErr)
		})
	}
}
