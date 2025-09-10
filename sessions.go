package scalekit

import (
	"context"

	sessionsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/sessions"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/sessions/sessionsconnect"
)

// Type aliases for response types
type SessionDetails = sessionsv1.SessionDetails
type UserSessionDetails = sessionsv1.UserSessionDetails
type RevokeSessionResponse = sessionsv1.RevokeSessionResponse

type SessionService interface {
	GetSession(ctx context.Context, sessionId string) (*SessionDetails, error)
	GetUserSessions(ctx context.Context, userId string) (*UserSessionDetails, error)
	RevokeSession(ctx context.Context, sessionId string) (*RevokeSessionResponse, error)
}

type sessionService struct {
	coreClient *coreClient
	client     sessionsconnect.SessionServiceClient
}

// newSessionClient creates a new session client
func newSessionClient(coreClient *coreClient) SessionService {
	return &sessionService{
		coreClient: coreClient,
		client:     newConnectClient(coreClient, sessionsconnect.NewSessionServiceClient),
	}
}

// GetSession retrieves session details by session ID
func (s *sessionService) GetSession(ctx context.Context, sessionId string) (*SessionDetails, error) {
	return newConnectExecuter(
		s.coreClient,
		s.client.GetSession,
		&sessionsv1.SessionDetailsRequest{
			SessionId: sessionId,
		},
	).exec(ctx)
}

// GetUserSessions retrieves all session details for a user
func (s *sessionService) GetUserSessions(ctx context.Context, userId string) (*UserSessionDetails, error) {
	return newConnectExecuter(
		s.coreClient,
		s.client.GetUserSessions,
		&sessionsv1.UserSessionDetailsRequest{
			UserId: userId,
		},
	).exec(ctx)
}

// RevokeSession revokes a session for a user
func (s *sessionService) RevokeSession(ctx context.Context, sessionId string) (*RevokeSessionResponse, error) {
	return newConnectExecuter(
		s.coreClient,
		s.client.RevokeSession,
		&sessionsv1.RevokeSessionRequest{
			SessionId: sessionId,
		},
	).exec(ctx)
}
