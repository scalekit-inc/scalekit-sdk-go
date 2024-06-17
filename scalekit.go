package scalekit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-jose/go-jose/v4"
)

const authorizeEndpoint = "oauth/authorize"

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeRefreshToken      GrantType = "refresh_token"
	GrantTypeClientCredentials GrantType = "client_credentials"
)

type GrantType = string

type Scalekit interface {
	Connection() Connection
	Domain() Domain
	Organization() Organization
	GetAuthorizationUrl(redirectUri string, options AuthorizationUrlOptions) (*url.URL, error)
	AuthenticateWithCode(
		code string,
		redirectUri string,
		options AuthenticationOptions,
	) (*AuthenticationResponse, error)
	ValidateAccessToken(accessToken string) (bool, error)
}

type scalekitClient struct {
	coreClient   *coreClient
	connection   Connection
	domain       Domain
	organization Organization
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

func NewScalekitClient(envUrl, clientId, clientSecret string) Scalekit {
	coreClient := newCoreClient(envUrl, clientId, clientSecret)
	return &scalekitClient{
		coreClient:   coreClient,
		connection:   newConnectionClient(coreClient),
		domain:       newDomainClient(coreClient),
		organization: newOrganizationClient(coreClient),
	}
}

func (s *scalekitClient) Connection() Connection {
	return s.connection
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

	var claims IdTokenClaims
	jws, err := jose.ParseSigned(authResp.IdToken, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jws.UnsafePayloadWithoutVerification(), &claims)
	if err != nil {
		return nil, err
	}

	return &AuthenticationResponse{
		User:        claims,
		IdToken:     authResp.IdToken,
		AccessToken: authResp.AccessToken,
		ExpiresIn:   authResp.ExpiresIn,
	}, nil
}

func (s *scalekitClient) ValidateAccessToken(accessToken string) (bool, error) {
	err := s.coreClient.getJwks()
	if err != nil {
		return false, err
	}
	jws, err := jose.ParseSigned(accessToken, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return false, err
	}
	_, err = jws.Verify(s.coreClient.jsonWebKeySet)
	if err != nil {
		return false, err
	}

	return true, nil
}
