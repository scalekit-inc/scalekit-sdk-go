package scalekit

import (
	"context"

	authv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/auth"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/auth/authconnect"
)

// Type aliases for WebAuthn request/response types
type ListCredentialsRequest = authv1.ListCredentialsRequest
type ListCredentialsResponse = authv1.ListCredentialsResponse
type UpdateCredentialRequest = authv1.UpdateCredentialRequest
type UpdateCredentialResponse = authv1.UpdateCredentialResponse
type DeleteCredentialRequest = authv1.DeleteCredentialRequest
type DeleteCredentialResponse = authv1.DeleteCredentialResponse
type WebAuthnCredential = authv1.WebAuthnCredential
type AllAcceptedCredentialsOptions = authv1.AllAcceptedCredentialsOptions
type UnknownCredentialOptions = authv1.UnknownCredentialOptions

// WebAuthnService interface defines the methods for WebAuthn/passkey operations
type WebAuthnService interface {
	// ListCredentials retrieves all registered passkeys for a user.
	// If userId is empty, it will list credentials for the current authenticated user.
	ListCredentials(ctx context.Context, userId string) (*ListCredentialsResponse, error)
	// UpdateCredential updates the display name of a passkey credential
	UpdateCredential(ctx context.Context, credentialId string, displayName string) (*UpdateCredentialResponse, error)
	// DeleteCredential deletes a specific passkey credential
	DeleteCredential(ctx context.Context, credentialId string) (*DeleteCredentialResponse, error)
}

type webAuthnService struct {
	coreClient *coreClient
	client     authconnect.WebAuthnServiceClient
}

// newWebAuthnClient creates a new WebAuthn client
func newWebAuthnClient(coreClient *coreClient) WebAuthnService {
	return &webAuthnService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, authconnect.NewWebAuthnServiceClient),
	}
}

// ListCredentials retrieves all registered passkeys for a user
func (w *webAuthnService) ListCredentials(ctx context.Context, userId string) (*ListCredentialsResponse, error) {
	request := &authv1.ListCredentialsRequest{
		UserId: userId,
	}

	return newConnectExecuter(
		w.coreClient,
		w.client.ListCredentials,
		request,
	).exec(ctx)
}

// UpdateCredential updates the display name of a passkey credential
func (w *webAuthnService) UpdateCredential(ctx context.Context, credentialId string, displayName string) (*UpdateCredentialResponse, error) {
	request := &authv1.UpdateCredentialRequest{
		CredentialId: credentialId,
		DisplayName:  displayName,
	}

	return newConnectExecuter(
		w.coreClient,
		w.client.UpdateCredential,
		request,
	).exec(ctx)
}

// DeleteCredential deletes a specific passkey credential
func (w *webAuthnService) DeleteCredential(ctx context.Context, credentialId string) (*DeleteCredentialResponse, error) {
	request := &authv1.DeleteCredentialRequest{
		CredentialId: credentialId,
	}

	return newConnectExecuter(
		w.coreClient,
		w.client.DeleteCredential,
		request,
	).exec(ctx)
}

