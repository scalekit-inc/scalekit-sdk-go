package scalekit

import (
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

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
	GrantTypeClientCredentials GrantType = "client_credentials"
)

var webhookToleranceInSeconds = 5 * time.Minute
var webhookSignatureVersion = "v1"

type GrantType = string

type Scalekit interface {
	Connection() Connection
	Directory() Directory
	Domain() Domain
	Organization() Organization
	GetAuthorizationUrl(redirectUri string, options AuthorizationUrlOptions) (*url.URL, error)
	AuthenticateWithCode(
		code string,
		redirectUri string,
		options AuthenticationOptions,
	) (*AuthenticationResponse, error)
	GetIdpInitiatedLoginClaims(idpInitiateLoginToken string) (*IdpInitiatedLoginClaims, error)
	ValidateAccessToken(accessToken string) (bool, error)
	VerifyWebhookPayload(secret string, headers map[string]string, payload []byte) (bool, error)
	GetAccessToken(accessToken string) (*AccessTokenClaims, error)
}

type scalekitClient struct {
	coreClient   *coreClient
	connection   Connection
	domain       Domain
	organization Organization
	directory    Directory
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
}

type AuthenticationOptions struct {
	CodeVerifier string
}

type AuthenticationResponse struct {
	User        User
	IdToken     string
	AccessToken string
	ExpiresIn   int
}

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
	// Alias is used to avoid infinite recursion during unmarshalling
	type Alias AccessTokenClaims
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	a.Claims = temp
	return nil
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

type Claims map[string]interface{}

func (i *IdTokenClaims) UnmarshalJSON(data []byte) error {
	// Alias is used to avoid infinite recursion during unmarshalling
	type Alias IdTokenClaims
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Claims = temp
	return nil
}

func NewScalekitClient(envUrl, clientId, clientSecret string) Scalekit {
	coreClient := newCoreClient(envUrl, clientId, clientSecret)
	return &scalekitClient{
		coreClient:   coreClient,
		connection:   newConnectionClient(coreClient),
		directory:    newDirectoryClient(coreClient),
		domain:       newDomainClient(coreClient),
		organization: newOrganizationClient(coreClient),
	}
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

	parsedUrl, err := url.Parse(fmt.Sprintf("%s/%s", s.coreClient.envUrl, authorizeEndpoint))
	if err != nil {
		return nil, err
	}
	parsedUrl.RawQuery = qs.Encode()

	return parsedUrl, nil
}

func (s *scalekitClient) AuthenticateWithCode(
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
	qs.Add("client_secret", s.coreClient.clientSecret)
	if options.CodeVerifier != "" {
		qs.Add("code_verifier", options.CodeVerifier)
	}
	authResp, err := s.coreClient.authenticate(qs)
	if err != nil {
		return nil, err
	}
	claims, err := validateToken[IdTokenClaims](authResp.IdToken, s.coreClient.getJwks)
	if err != nil {
		return nil, err
	}

	return &AuthenticationResponse{
		User:        *claims,
		IdToken:     authResp.IdToken,
		AccessToken: authResp.AccessToken,
		ExpiresIn:   authResp.ExpiresIn,
	}, nil
}

func (s *scalekitClient) GetIdpInitiatedLoginClaims(idpInitiateLoginToken string) (*IdpInitiatedLoginClaims, error) {
	return validateToken[IdpInitiatedLoginClaims](idpInitiateLoginToken, s.coreClient.getJwks)
}

func (s *scalekitClient) GetAccessToken(accessToken string) (*AccessTokenClaims, error) {
	at, err := validateToken[AccessTokenClaims](accessToken, s.coreClient.getJwks)
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (s *scalekitClient) ValidateAccessToken(accessToken string) (bool, error) {
	_, err := validateToken[AccessTokenClaims](accessToken, s.coreClient.getJwks)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *scalekitClient) VerifyWebhookPayload(
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

func validateToken[T interface{}](token string, jwksFn func() (*jose.JSONWebKeySet, error)) (*T, error) {
	var claims T
	keySet, err := jwksFn()
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
	return &claims, nil
}

func computeSignature(secret []byte, data string) string {
	hash := hmac.New(sha256.New, secret)
	hash.Write([]byte(data))
	signature := hash.Sum(nil)

	return base64.StdEncoding.EncodeToString(signature)
}
