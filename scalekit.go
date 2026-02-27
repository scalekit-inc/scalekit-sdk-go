package scalekit

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
)

const authorizeEndpoint = "oauth/authorize"
const logoutEndpoint = "oidc/logout"

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
	GrantTypeClientCredentials GrantType = "client_credentials"
)

var (
	webhookToleranceInSeconds = 5 * time.Minute
	webhookSignatureVersion   = "v1"
)

type GrantType = string

type Scalekit interface {
	Connection() Connection
	Directory() Directory
	Domain() Domain
	Organization() Organization
	User() UserService
	Passwordless() PasswordlessService
	Auth() AuthService
	Client() ClientService
	Session() SessionService
	Role() RoleService
	Permission() PermissionService
	WebAuthn() WebAuthnService
	Token() TokenService
	GetAuthorizationUrl(redirectUri string, options AuthorizationUrlOptions) (*url.URL, error)
	AuthenticateWithCode(ctx context.Context, code string, redirectUri string, options AuthenticationOptions) (*AuthenticationResponse, error)
	GetIdpInitiatedLoginClaims(ctx context.Context, idpInitiateLoginToken string) (*IdpInitiatedLoginClaims, error)
	ValidateAccessToken(ctx context.Context, accessToken string) (bool, error)
	ValidateTokenWithOptions(ctx context.Context, token string, options *ValidateTokenOptions) (bool, error)
	VerifyWebhookPayload(secret string, headers map[string]string, payload []byte) (bool, error)
	VerifyInterceptorPayload(secret string, headers map[string]string, payload []byte) (bool, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	GenerateClientToken(ctx context.Context, options *GenerateClientTokenOptions) (*ClientTokenResponse, error)
	GetLogoutUrl(options LogoutUrlOptions) (*url.URL, error)
	GetAccessTokenClaims(ctx context.Context, accessToken string) (*AccessTokenClaims, error)
	GeneratePKCEConfiguration(options PKCEOptions) (*PKCEConfiguration, error)
	WithSecret(clientSecret string) Scalekit
}

type scalekitClient struct {
	coreClient   *coreClient
	connection   Connection
	domain       Domain
	organization Organization
	directory    Directory
	user         UserService
	passwordless PasswordlessService
	auth         AuthService
	oidcClient   ClientService
	session      SessionService
	role         RoleService
	permission   PermissionService
	webauthn     WebAuthnService
	token        TokenService
}

type AuthorizationUrlOptions struct {
	ConnectionId        string
	OrganizationId      string
	Scopes              []string
	State               string
	Nonce               string
	DomainHint          string
	LoginHint           string
	CodeChallenge       string
	CodeChallengeMethod string
	Provider            string
	Prompt              string
}

type AuthenticationOptions struct {
	CodeVerifier string
}

// GenerateClientTokenOptions defines optional fields for client-credentials token generation.
// The token endpoint does not currently support additional options for this grant type,
// but this struct is intentionally kept for forward compatibility.
type GenerateClientTokenOptions struct {
}

// ValidateTokenOptions defines optional validations for token verification.
type ValidateTokenOptions struct {
	Audience []string
}

type AuthenticationResponse struct {
	User         User
	IdToken      string
	AccessToken  string
	ExpiresIn    int
	RefreshToken string
}

type (
	Claims  map[string]interface{}
	idAlias IdTokenClaims
	atAlias AccessTokenClaims
	tkAlias TokenClaims
)

type IdTokenClaims struct {
	Id                  string     `json:"sub"`
	Username            string     `json:"preferred_username"`
	Name                string     `json:"name"`
	GivenName           string     `json:"given_name"`
	FamilyName          string     `json:"family_name"`
	Email               string     `json:"email"`
	EmailVerified       bool       `json:"email_verified"`
	PhoneNumber         string     `json:"phone_number"`
	PhoneNumberVerified bool       `json:"phone_number_verified"`
	Profile             string     `json:"profile"`
	Picture             string     `json:"picture"`
	Gender              string     `json:"gender"`
	BirthDate           string     `json:"birthdate"`
	ZoneInfo            string     `json:"zoneinfo"`
	Locale              string     `json:"locale"`
	UpdatedAt           string     `json:"updated_at"`
	Identities          []Identity `json:"identities"`
	Metadata            string     `json:"metadata"`
	Claims              Claims     `json:"-"`
}

func (i *IdTokenClaims) UnmarshalJSON(data []byte) error {
	return unmarshalJson(data, (*idAlias)(i), &i.Claims)
}

type Audience []string

type AccessTokenClaims struct {
	Sub      string   `json:"sub"`
	Iss      string   `json:"iss"`
	Audience Audience `json:"aud,omitempty"`
	Iat      int      `json:"iat"`
	Exp      int      `json:"exp"`
	Claims   Claims   `json:"-"`
}

func (a *AccessTokenClaims) UnmarshalJSON(data []byte) error {
	return unmarshalJson(data, (*atAlias)(a), &a.Claims)
}

type TokenClaims struct {
	Sub      string   `json:"sub"`
	Iss      string   `json:"iss"`
	Audience Audience `json:"aud,omitempty"`
	Iat      int      `json:"iat"`
	Exp      int      `json:"exp"`
	Claims   Claims   `json:"-"`
}

func (t *TokenClaims) UnmarshalJSON(data []byte) error {
	return unmarshalJson(data, (*tkAlias)(t), &t.Claims)
}

type User = IdTokenClaims

type Identity struct {
	ConnectionId          string `json:"connection_id"`
	OrganizationId        string `json:"organization_id"`
	ConnectionType        string `json:"connection_type"`
	ProviderName          string `json:"provider_name"`
	Social                bool   `json:"social"`
	ProviderRawAttributes string `json:"provider_raw_attributes"`
}

type IdpInitiatedLoginClaims struct {
	ConnectionID   string  `json:"connection_id"`
	OrganizationID string  `json:"organization_id"`
	LoginHint      string  `json:"login_hint"`
	RelayState     *string `json:"relay_state"`
}

type TokenResponse struct {
	AccessToken  string
	IdToken      string
	RefreshToken string
	ExpiresIn    int
}

type ClientTokenResponse struct {
	AccessToken string
	ExpiresIn   int
}

type LogoutUrlOptions struct {
	IdTokenHint           string
	PostLogoutRedirectUri string
	State                 string
}

// NewScalekitClient creates a new Scalekit client.
//
// For backward compatibility, when a value is provided in opts and opts[0] is a
// string, it is treated as client_secret.
func NewScalekitClient(envUrl, clientId string, opts ...any) Scalekit {
	clientSecret := ""
	if len(opts) > 0 {
		if secret, ok := opts[0].(string); ok {
			clientSecret = secret
		}
	}
	return newScalekitClient(newCoreClient(envUrl, clientId, clientSecret))
}

func newScalekitClient(coreClient *coreClient) *scalekitClient {
	return &scalekitClient{
		coreClient:   coreClient,
		connection:   newConnectionClient(coreClient),
		directory:    newDirectoryClient(coreClient),
		domain:       newDomainClient(coreClient),
		organization: newOrganizationClient(coreClient),
		user:         newUserClient(coreClient),
		passwordless: newPasswordlessClient(coreClient),
		auth:         newAuthService(coreClient),
		oidcClient:   newClientService(coreClient),
		session:      newSessionClient(coreClient),
		role:         newRoleService(coreClient),
		permission:   newPermissionService(coreClient),
		webauthn:     newWebAuthnClient(coreClient),
		token:        newTokenService(coreClient),
	}
}

func (s *scalekitClient) WithSecret(clientSecret string) Scalekit {
	coreCopy := *s.coreClient
	coreCopy.clientSecret = clientSecret

	return newScalekitClient(&coreCopy)
}

func (s *scalekitClient) Connection() Connection {
	return s.connection
}

func (s *scalekitClient) Directory() Directory {
	return s.directory
}

func (s *scalekitClient) Domain() Domain {
	return s.domain
}

func (s *scalekitClient) Organization() Organization {
	return s.organization
}

func (s *scalekitClient) User() UserService {
	return s.user
}

func (s *scalekitClient) Passwordless() PasswordlessService {
	return s.passwordless
}

func (s *scalekitClient) Auth() AuthService {
	return s.auth
}

func (s *scalekitClient) Client() ClientService {
	return s.oidcClient
}

func (s *scalekitClient) Session() SessionService {
	return s.session
}

func (s *scalekitClient) Role() RoleService {
	return s.role
}

func (s *scalekitClient) Permission() PermissionService {
	return s.permission
}

func (s *scalekitClient) WebAuthn() WebAuthnService {
	return s.webauthn
}

func (s *scalekitClient) Token() TokenService {
	return s.token
}

func (s *scalekitClient) GetAuthorizationUrl(redirectUri string, options AuthorizationUrlOptions) (*url.URL, error) {
	scopes := []string{"openid", "profile", "email"}
	if options.Scopes != nil {
		scopes = options.Scopes[:]
	}
	qs := url.Values{}
	qs.Set("response_type", "code")
	qs.Set("client_id", s.coreClient.clientId)
	qs.Set("redirect_uri", redirectUri)
	qs.Set("scope", strings.Join(scopes, " "))
	if options.State != "" {
		qs.Set("state", options.State)
	}
	if options.Nonce != "" {
		qs.Set("nonce", options.Nonce)
	}
	if options.LoginHint != "" {
		qs.Set("login_hint", options.LoginHint)
	}
	if options.DomainHint != "" {
		qs.Set("domain_hint", options.DomainHint)
		qs.Set("domain", options.DomainHint)
	}
	if options.ConnectionId != "" {
		qs.Set("connection_id", options.ConnectionId)
	}
	if options.OrganizationId != "" {
		qs.Set("organization_id", options.OrganizationId)
	}
	if options.CodeChallenge != "" {
		qs.Set("code_challenge", options.CodeChallenge)
	}
	if options.CodeChallengeMethod != "" {
		qs.Set("code_challenge_method", options.CodeChallengeMethod)
	}
	if options.Provider != "" {
		qs.Set("provider", options.Provider)
	}
	if options.Prompt != "" {
		qs.Set("prompt", options.Prompt)
	}

	parsedUrl, err := url.Parse(fmt.Sprintf("%s/%s", s.coreClient.envUrl, authorizeEndpoint))
	if err != nil {
		return nil, err
	}
	parsedUrl.RawQuery = qs.Encode()

	return parsedUrl, nil
}

func (s *scalekitClient) AuthenticateWithCode(
	ctx context.Context,
	code string,
	redirectUri string,
	options AuthenticationOptions,
) (*AuthenticationResponse, error) {
	if code == "" || redirectUri == "" {
		return nil, errors.New("code and redirect uri is required")
	}
	qs := url.Values{}
	qs.Add("code", code)
	qs.Add("redirect_uri", redirectUri)
	qs.Add("grant_type", GrantTypeAuthorizationCode)
	qs.Add("client_id", s.coreClient.clientId)
	if s.coreClient.clientSecret != "" {
		qs.Add("client_secret", s.coreClient.clientSecret)
	}
	if options.CodeVerifier != "" {
		qs.Add("code_verifier", options.CodeVerifier)
	}
	authResp, err := s.coreClient.authenticate(ctx, qs)
	if err != nil {
		return nil, err
	}
	claims, err := ValidateToken[IdTokenClaims](ctx, authResp.IdToken, s.coreClient.GetJwks)
	if err != nil {
		return nil, err
	}

	return &AuthenticationResponse{
		User:         *claims,
		IdToken:      authResp.IdToken,
		AccessToken:  authResp.AccessToken,
		ExpiresIn:    authResp.ExpiresIn,
		RefreshToken: authResp.RefreshToken,
	}, nil
}

func (s *scalekitClient) GetIdpInitiatedLoginClaims(ctx context.Context, idpInitiateLoginToken string) (*IdpInitiatedLoginClaims, error) {
	return ValidateToken[IdpInitiatedLoginClaims](ctx, idpInitiateLoginToken, s.coreClient.GetJwks)
}

func (s *scalekitClient) GetAccessTokenClaims(ctx context.Context, accessToken string) (*AccessTokenClaims, error) {
	at, err := ValidateToken[AccessTokenClaims](ctx, accessToken, s.coreClient.GetJwks)
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (s *scalekitClient) ValidateAccessToken(ctx context.Context, accessToken string) (bool, error) {
	_, err := ValidateToken[AccessTokenClaims](ctx, accessToken, s.coreClient.GetJwks)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ValidateTokenWithOptions validates a signed JWT (access token or ID token)
// and enforces optional checks such as audience validation.
func (s *scalekitClient) ValidateTokenWithOptions(ctx context.Context, token string, options *ValidateTokenOptions) (bool, error) {
	claims, err := ValidateToken[TokenClaims](ctx, token, s.coreClient.GetJwks)
	if err != nil {
		return false, err
	}
	if options == nil {
		return true, nil
	}
	if len(options.Audience) == 0 {
		return true, nil
	}

	audienceSet := map[string]struct{}{}
	for _, audience := range claims.Audience {
		audienceSet[audience] = struct{}{}
	}

	matched := false
	for _, audience := range options.Audience {
		if _, ok := audienceSet[audience]; ok {
			matched = true
			break
		}
	}
	if !matched {
		return false, fmt.Errorf("none of the expected audiences found in token aud claim")
	}

	return true, nil
}

func (s *scalekitClient) VerifyWebhookPayload(
	secret string,
	headers map[string]string,
	payload []byte,
) (bool, error) {
	return s.VerifyPayloadSignature(secret, headers, payload)
}

func (s *scalekitClient) VerifyPayloadSignature(
	secret string,
	headers map[string]string,
	payload []byte,
) (bool, error) {
	webhookId := headers["webhook-id"]
	webhookTimestamp := headers["webhook-timestamp"]
	webhookSignature := headers["webhook-signature"]
	if webhookId == "" || webhookTimestamp == "" || webhookSignature == "" {
		return false, errors.New("Missing required headers")
	}
	secretParts := strings.Split(secret, "_")
	if len(secretParts) < 2 {
		return false, errors.New("Invalid secret")
	}
	secretBytes, err := base64.StdEncoding.DecodeString(secretParts[1])
	if err != nil {
		return false, err
	}
	timestamp, err := verifyTimestamp(webhookTimestamp)
	if err != nil {
		return false, err
	}
	data := fmt.Sprintf("%s.%d.%s", webhookId, timestamp.Unix(), payload)
	computedSignature := computeSignature(secretBytes, data)
	recievedSignatures := strings.Split(webhookSignature, " ")
	for _, versionedSignature := range recievedSignatures {
		signatureParts := strings.Split(versionedSignature, ",")
		if len(signatureParts) < 2 {
			continue
		}
		version := signatureParts[0]
		signature := signatureParts[1]
		if version != webhookSignatureVersion {
			continue
		}
		if hmac.Equal([]byte(signature), []byte(computedSignature)) {
			return true, nil
		}
	}

	return false, errors.New("Invalid signature")
}

func (s *scalekitClient) VerifyInterceptorPayload(
	secret string,
	headers map[string]string,
	payload []byte,
) (bool, error) {
	return s.VerifyPayloadSignature(secret, headers, payload)
}

func verifyTimestamp(timestampStr string) (*time.Time, error) {
	now := time.Now()
	unixTimestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil, err
	}
	timestamp := time.Unix(unixTimestamp, 0)
	if now.Sub(timestamp) > webhookToleranceInSeconds {
		return nil, errors.New("Message timestamp too old")
	}
	if timestamp.Unix() > now.Add(webhookToleranceInSeconds).Unix() {
		return nil, errors.New("Message timestamp too new")
	}
	return &timestamp, nil
}

func ValidateToken[T interface{}](ctx context.Context, token string, jwksFn func(context.Context) (*jose.JSONWebKeySet, error)) (*T, error) {
	var claims T
	keySet, err := jwksFn(ctx)
	if err != nil {
		return nil, err
	}
	jws, err := jose.ParseSigned(token, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return nil, err
	}
	jwt, err := jws.Verify(keySet)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jwt, &claims)
	if err != nil {
		return nil, err
	}

	// Check token expiration
	var rawClaims map[string]interface{}
	err = json.Unmarshal(jwt, &rawClaims)
	if err != nil {
		return nil, err
	}

	if exp, ok := rawClaims["exp"]; ok {
		expFloat, ok := exp.(float64)
		if !ok {
			return nil, ErrInvalidExpClaimFormat
		}

		expTime := int64(expFloat)
		if time.Now().Unix() >= expTime {
			return nil, ErrTokenExpired
		}
	}

	return &claims, nil
}

func computeSignature(secret []byte, data string) string {
	hash := hmac.New(sha256.New, secret)
	hash.Write([]byte(data))
	signature := hash.Sum(nil)

	return base64.StdEncoding.EncodeToString(signature)
}

func (s *scalekitClient) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	if refreshToken == "" {
		return nil, ErrRefreshTokenRequired
	}

	qs := url.Values{}
	qs.Add("refresh_token", refreshToken)
	qs.Add("grant_type", GrantTypeRefreshToken)
	qs.Add("client_id", s.coreClient.clientId)
	if s.coreClient.clientSecret != "" {
		qs.Add("client_secret", s.coreClient.clientSecret)
	}
	authResp, err := s.coreClient.authenticate(ctx, qs)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresIn:    authResp.ExpiresIn,
		IdToken:      authResp.IdToken,
	}, nil
}

// GenerateClientToken generates a client-credentials token using the client_id and
// client_secret configured on the Scalekit client.
//
// The options parameter is reserved for future server-supported fields. For example,
// a future version may support scopes like "usr:read" and "usr:write".
func (s *scalekitClient) GenerateClientToken(
	ctx context.Context,
	_ *GenerateClientTokenOptions,
) (*ClientTokenResponse, error) {
	if s.coreClient.clientSecret == "" {
		return nil, ErrClientSecretRequired
	}

	qs := url.Values{}
	qs.Add("grant_type", GrantTypeClientCredentials)
	qs.Add("client_id", s.coreClient.clientId)
	qs.Add("client_secret", s.coreClient.clientSecret)

	authResp, err := s.coreClient.authenticate(ctx, qs)
	if err != nil {
		return nil, err
	}

	return &ClientTokenResponse{
		AccessToken: authResp.AccessToken,
		ExpiresIn:   authResp.ExpiresIn,
	}, nil
}

func (s *scalekitClient) GetLogoutUrl(options LogoutUrlOptions) (*url.URL, error) {
	qs := url.Values{}

	if options.IdTokenHint != "" {
		qs.Set("id_token_hint", options.IdTokenHint)
	}

	if options.PostLogoutRedirectUri != "" {
		qs.Set("post_logout_redirect_uri", options.PostLogoutRedirectUri)
	}

	if options.State != "" {
		qs.Set("state", options.State)
	}

	parsedUrl, err := url.Parse(fmt.Sprintf("%s/%s", s.coreClient.envUrl, logoutEndpoint))
	if err != nil {
		return nil, err
	}
	parsedUrl.RawQuery = qs.Encode()

	return parsedUrl, nil
}

func unmarshalJson(data []byte, types ...any) error {
	for _, cType := range types {
		err := json.Unmarshal(data, cType)
		if err != nil {
			return err
		}
	}
	return nil
}
