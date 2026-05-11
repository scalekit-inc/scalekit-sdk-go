package test

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	organizationsv1 "github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/organizations"
	"github.com/scalekit-inc/scalekit-sdk-go/v2/pkg/grpc/scalekit/v1/organizations/organizationsconnect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestAuthenticateWithCode(t *testing.T) {
	tests := []struct {
		name     string
		req      func() (string, string, scalekit.AuthenticationOptions)
		mockFn   func(w http.ResponseWriter, r *http.Request)
		assertFn func(t *testing.T, resp *scalekit.AuthenticationResponse, err error)
	}{
		{
			name: "successful_authentication",
			req: func() (string, string, scalekit.AuthenticationOptions) {
				return "test_code", "http://localhost/callback", scalekit.AuthenticationOptions{}
			},
			mockFn: func() func(http.ResponseWriter, *http.Request) {
				// One key pair for both /keys and /oauth/token so JWKS matches signed tokens
				privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					return func(w http.ResponseWriter, _ *http.Request) {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
				keyID := "mock-kid"
				signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privateKey}, (&jose.SignerOptions{}).WithHeader("kid", keyID))
				if err != nil {
					return func(w http.ResponseWriter, _ *http.Request) {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
				jwk := jose.JSONWebKey{
					Key:       privateKey.Public(),
					KeyID:     keyID,
					Algorithm: string(jose.RS256),
					Use:       "sig",
				}
				keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}}

				return func(w http.ResponseWriter, r *http.Request) {
					switch r.URL.Path {
					case "/keys":
						w.Header().Set("Content-Type", "application/json")
						_ = json.NewEncoder(w).Encode(keySet)
					case "/oauth/token":
						now := time.Now()
						iat := now.Unix()
						exp := now.Add(time.Hour).Unix()
						idClaims := map[string]interface{}{
							"sub":            "usr_mock123",
							"name":           "Mock User",
							"email":          "mock@example.com",
							"given_name":     "Mock",
							"family_name":    "User",
							"email_verified": true,
							"iat":            iat,
							"exp":            exp,
							"oid":            "org_mock456",
							"sid":            "ses_mock789",
						}
						idTokenPayload, _ := json.Marshal(idClaims)
						idToken, err := signer.Sign(idTokenPayload)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						idTokenCompact, _ := idToken.CompactSerialize()

						atClaims := map[string]interface{}{
							"sub": "conn_1;user@example.com",
							"iss": "https://mock.example.com",
							"aud": []string{"prd_skc_mock"},
							"iat": iat,
							"exp": exp,
						}
						atPayload, _ := json.Marshal(atClaims)
						atSigned, _ := signer.Sign(atPayload)
						accessToken, _ := atSigned.CompactSerialize()

						w.Header().Set("Content-Type", "application/json")
						resp := map[string]interface{}{
							"access_token": accessToken,
							"id_token":     idTokenCompact,
							"expires_in":   3600,
						}
						_ = json.NewEncoder(w).Encode(resp)
					}
				}
			}(),
			assertFn: func(t *testing.T, resp *scalekit.AuthenticationResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, 3600, resp.ExpiresIn)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.IdToken)
				assert.Equal(t, "usr_mock123", resp.User.Id)
				assert.Equal(t, "Mock User", resp.User.Name)
				assert.Equal(t, "mock@example.com", resp.User.Email)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.mockFn))
			defer server.Close()

			client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
			code, redirectUri, options := tt.req()
			resp, err := client.AuthenticateWithCode(context.Background(), code, redirectUri, options)
			tt.assertFn(t, resp, err)
		})
	}
}

func TestAuthenticateWithCode_ClientSecretBehaviorByClientType(t *testing.T) {
	type testCase struct {
		name                  string
		clientType            string
		usePKCE               bool
		includeSecretAtInit   bool
		expectSecretInTokenRq bool
	}

	tests := []testCase{
		{
			name:                  "WEB_APP with PKCE should not send client_secret",
			clientType:            "WEB_APP",
			usePKCE:               true,
			includeSecretAtInit:   false,
			expectSecretInTokenRq: false,
		},
		{
			name:                  "WEB_APP without PKCE should send client_secret",
			clientType:            "WEB_APP",
			usePKCE:               false,
			includeSecretAtInit:   true,
			expectSecretInTokenRq: true,
		},
		{
			name:                  "SPA with PKCE should not send client_secret",
			clientType:            "SPA",
			usePKCE:               true,
			includeSecretAtInit:   false,
			expectSecretInTokenRq: false,
		},
		{
			name:                  "NTV with PKCE should not send client_secret",
			clientType:            "NTV",
			usePKCE:               true,
			includeSecretAtInit:   false,
			expectSecretInTokenRq: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(t, err)

			keyID := "mock-kid"
			signer, err := jose.NewSigner(
				jose.SigningKey{Algorithm: jose.RS256, Key: privateKey},
				(&jose.SignerOptions{}).WithHeader("kid", keyID),
			)
			require.NoError(t, err)

			jwk := jose.JSONWebKey{
				Key:       privateKey.Public(),
				KeyID:     keyID,
				Algorithm: string(jose.RS256),
				Use:       "sig",
			}
			keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}}

			var expectedVerifier string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/keys":
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(keySet)
				case "/oauth/token":
					require.NoError(t, r.ParseForm())
					assert.Equal(t, "test_code", r.FormValue("code"))
					assert.Equal(t, "http://localhost/callback", r.FormValue("redirect_uri"))
					assert.Equal(t, "authorization_code", r.FormValue("grant_type"))
					assert.Equal(t, "client_id", r.FormValue("client_id"))

					if tt.expectSecretInTokenRq {
						assert.Equal(t, "client_secret", r.FormValue("client_secret"))
					} else {
						assert.Empty(t, r.FormValue("client_secret"))
					}

					if tt.usePKCE {
						assert.Equal(t, expectedVerifier, r.FormValue("code_verifier"))
					} else {
						assert.Empty(t, r.FormValue("code_verifier"))
					}

					now := time.Now()
					iat := now.Unix()
					exp := now.Add(time.Hour).Unix()

					idClaims := map[string]interface{}{
						"sub":            "usr_mock123",
						"name":           "Mock User",
						"email":          "mock@example.com",
						"given_name":     "Mock",
						"family_name":    "User",
						"email_verified": true,
						"iat":            iat,
						"exp":            exp,
					}
					idTokenPayload, _ := json.Marshal(idClaims)
					idToken, signErr := signer.Sign(idTokenPayload)
					require.NoError(t, signErr)
					idTokenCompact, _ := idToken.CompactSerialize()

					atClaims := map[string]interface{}{
						"sub": "conn_1;user@example.com",
						"iss": "https://mock.example.com",
						"aud": []string{"prd_skc_mock"},
						"iat": iat,
						"exp": exp,
					}
					atPayload, _ := json.Marshal(atClaims)
					atSigned, signErr := signer.Sign(atPayload)
					require.NoError(t, signErr)
					accessToken, _ := atSigned.CompactSerialize()

					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"access_token": accessToken,
						"id_token":     idTokenCompact,
						"expires_in":   3600,
					})
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()

			var multiAppClient scalekit.Scalekit
			if tt.includeSecretAtInit {
				multiAppClient = scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
			} else {
				multiAppClient = scalekit.NewScalekitClient(server.URL, "client_id")
			}

			authOptions := scalekit.AuthorizationUrlOptions{
				State: "state-" + strings.ToLower(tt.clientType),
			}
			codeOptions := scalekit.AuthenticationOptions{}

			if tt.usePKCE {
				pkceCfg, pkceErr := multiAppClient.GeneratePKCEConfiguration(scalekit.PKCEOptions{})
				require.NoError(t, pkceErr)
				expectedVerifier = pkceCfg.CodeVerifier
				authOptions.CodeChallenge = pkceCfg.CodeChallenge
				authOptions.CodeChallengeMethod = pkceCfg.CodeChallengeMethod
				codeOptions.CodeVerifier = pkceCfg.CodeVerifier
			}

			authURL, err := multiAppClient.GetAuthorizationUrl("http://localhost/callback", authOptions)
			require.NoError(t, err)
			require.NotNil(t, authURL)

			if tt.usePKCE {
				assert.Equal(t, authOptions.CodeChallenge, authURL.Query().Get("code_challenge"))
				assert.Equal(t, authOptions.CodeChallengeMethod, authURL.Query().Get("code_challenge_method"))
			} else {
				assert.Empty(t, authURL.Query().Get("code_challenge"))
				assert.Empty(t, authURL.Query().Get("code_challenge_method"))
			}

			resp, err := multiAppClient.AuthenticateWithCode(
				context.Background(),
				"test_code",
				"http://localhost/callback",
				codeOptions,
			)
			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, "usr_mock123", resp.User.Id)
			assert.Equal(t, "Mock User", resp.User.Name)
		})
	}
}

func TestGetAccessToken(t *testing.T) {
	type testCase struct {
		name     string
		token    string
		mockFn   func(w http.ResponseWriter, r *http.Request)
		assertFn func(t *testing.T, token *scalekit.AccessTokenClaims, err error)
	}

	tests := []testCase{
		{
			name:  "successful token validation",
			token: "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vYWlyZGV2LmxvY2FsaG9zdDo4ODg4Iiwic3ViIjoiY29ubl83NTQxNjU3OTA0MjQ3NDIwNDtzcmluaXZhcy5rYXJyZUBzY2FsZWtpdC5jb20iLCJhdWQiOlsicHJkX3NrY18xNzAwMjMzNDIyNzg1NzUwOCJdLCJleHAiOjE5MDY4MDQ4MzcsImlhdCI6MTc0OTAyMDA3NywibmJmIjoxNzQ5MDIwMDc3LCJjbGllbnRfaWQiOiJwcmRfc2tjXzE3MDAyMzM0MjI3ODU3NTA4IiwianRpIjoidGtuXzc1NDE4NDE0MTAwODA1ODUyIn0.SxlKHr1EFBAvfm3Zm7CliKcSWZ8LUFWx8Cs3_3bf1SVouVvRu-zE2_ghB4iAmarsxErurU0kHDEX-Fpx6euemiWXN3Z-mECB4clmb1PF8RThh7bbHx1zxqp3z_MIcDbO4ZKTXMSRx39JbcWyThQSTbeAo50TEFpIT7RsWhNYrBnhsZNibrfZXWUVDBYB930LZMzhdKPRUXBhA-HuKIjggg2jWEAv2leJ3UPbLVccbKrdq2qSzGaxLpvlPoX6RpcrA2Cbuig4vJ7bCy46M-DUg73NO91arPpl5BOnHHx2Oappk_i2S4cMOGdSyX3s50owX1xRDyELNMEIo-VoQ7rfww\n",
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/keys":
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, token *scalekit.AccessTokenClaims, err error) {
				require.NoError(t, err)
				require.NotNil(t, token)
				assert.Equal(t, 1906804837, token.Exp)

				// Verify parsed claims
				assert.Equal(t, "conn_75416579042474204;srinivas.karre@scalekit.com", token.Sub)
				assert.Equal(t, scalekit.Audience{"prd_skc_17002334227857508"}, token.Audience)

				// Verify custom claim
				rawClaims := token.Claims
				assert.Equal(t, "conn_75416579042474204;srinivas.karre@scalekit.com", rawClaims["sub"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.mockFn))
			defer server.Close()

			client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
			token, err := client.GetAccessTokenClaims(context.Background(), tt.token)
			tt.assertFn(t, token, err)
		})
	}
}

func TestVerifyWebhookPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
	createSignature := func(webhookID string, timestamp int64, payload []byte, secretStr string) string {
		data := fmt.Sprintf("%s.%d.%s", webhookID, timestamp, payload)
		secretBytes, _ := base64.StdEncoding.DecodeString(secretStr)
		hash := hmac.New(sha256.New, secretBytes)
		hash.Write([]byte(data))
		return base64.StdEncoding.EncodeToString(hash.Sum(nil))
	}

	tests := []struct {
		name          string
		secret        string
		headers       map[string]string
		payload       []byte
		expectedValid bool
		expectedError string
	}{
		{
			name:    "valid signature",
			secret:  "whsec_dGVzdHNlY3JldA==",
			payload: []byte(`{"event": "user.created", "data": {"id": "123"}}`),
			headers: func() map[string]string {
				timestamp := time.Now().Unix()
				webhookID := "webhook_123"
				signature := createSignature(webhookID, timestamp, []byte(`{"event": "user.created", "data": {"id": "123"}}`), "dGVzdHNlY3JldA==")
				return map[string]string{
					"webhook-id":        webhookID,
					"webhook-timestamp": fmt.Sprintf("%d", timestamp),
					"webhook-signature": fmt.Sprintf("v1,%s", signature),
				}
			}(),
			expectedValid: true,
			expectedError: "",
		},
		{
			name:   "missing headers",
			secret: "whsec_dGVzdHNlY3JldA==",
			headers: map[string]string{
				"webhook-id": "webhook_123",
			},
			payload:       []byte("{}"),
			expectedValid: false,
			expectedError: "missing required headers",
		},
		{
			name:   "invalid secret",
			secret: "invalid_secret",
			headers: map[string]string{
				"webhook-id":        "webhook_123",
				"webhook-timestamp": fmt.Sprintf("%d", time.Now().Unix()),
				"webhook-signature": "v1,invalid",
			},
			payload:       []byte("{}"),
			expectedValid: false,
			expectedError: "illegal base64 data at input byte 4",
		},
		{
			name:   "expired timestamp",
			secret: "whsec_dGVzdHNlY3JldA==",
			headers: map[string]string{
				"webhook-id":        "webhook_123",
				"webhook-timestamp": fmt.Sprintf("%d", time.Now().Add(-10*time.Minute).Unix()),
				"webhook-signature": "v1,somesignature",
			},
			payload:       []byte("{}"),
			expectedValid: false,
			expectedError: "message timestamp too old",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := client.VerifyWebhookPayload(tt.secret, tt.headers, tt.payload)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedValid, valid)
		})
	}
}

// TestValidateTokenViaInterface verifies that ValidateToken is callable directly
// through the Scalekit interface type, as described in SK-2598.
func TestValidateTokenViaInterface(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	keyID := "test-kid"
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: privateKey},
		(&jose.SignerOptions{}).WithHeader("kid", keyID),
	)
	require.NoError(t, err)

	jwk := jose.JSONWebKey{
		Key:       privateKey.Public(),
		KeyID:     keyID,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}
	keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/keys" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(keySet)
		}
	}))
	defer server.Close()

	// NewScalekitClient returns the Scalekit interface directly (SK-2598)
	c := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")

	// Build a signed token using our test key
	now := time.Now()
	rawClaims := map[string]interface{}{
		"sub": "usr_test123",
		"iss": server.URL,
		"exp": now.Add(time.Hour).Unix(),
		"iat": now.Unix(),
	}
	claimsBytes, _ := json.Marshal(rawClaims)
	signed, err := signer.Sign(claimsBytes)
	require.NoError(t, err)
	token, err := signed.CompactSerialize()
	require.NoError(t, err)

	// Callable directly through the interface — no package-level function required
	claims, err := c.ValidateToken(context.Background(), token)
	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, "usr_test123", claims["sub"])
}

func TestGenerateClientToken(t *testing.T) {
	tests := []struct {
		name     string
		options  scalekit.GenerateClientTokenOptions
		mockFn   func(w http.ResponseWriter, r *http.Request)
		assertFn func(t *testing.T, resp *scalekit.ClientTokenResponse, err error)
	}{
		{
			name: "successful token generation",
			options: scalekit.GenerateClientTokenOptions{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				Scopes:       []string{"usr:read", "usr:write"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/oauth/token" {
					_ = r.ParseForm()
					assert.Equal(t, "client_credentials", r.FormValue("grant_type"))
					assert.Equal(t, "test_client_id", r.FormValue("client_id"))
					assert.Equal(t, "test_client_secret", r.FormValue("client_secret"))
					assert.Equal(t, "usr:read usr:write", r.FormValue("scope"))
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"access_token":"mock_token_123","expires_in":3600}`))
				}
			},
			assertFn: func(t *testing.T, resp *scalekit.ClientTokenResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, "mock_token_123", resp.AccessToken)
				assert.Equal(t, 3600, resp.ExpiresIn)
			},
		},
		{
			// authenticate() does not inspect HTTP status codes; it errors only on
			// network failure or JSON decode failure. Return malformed JSON so the
			// decode step fails and a real error is propagated to the caller.
			name: "server error",
			options: scalekit.GenerateClientTokenOptions{
				ClientID:     "bad_id",
				ClientSecret: "bad_secret",
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/oauth/token" {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write([]byte(`not valid json`))
				}
			},
			assertFn: func(t *testing.T, resp *scalekit.ClientTokenResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
		{
			name: "missing client secret",
			options: scalekit.GenerateClientTokenOptions{
				ClientID: "test_client_id",
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {},
			assertFn: func(t *testing.T, resp *scalekit.ClientTokenResponse, err error) {
				assert.ErrorIs(t, err, scalekit.ErrClientSecretRequired)
				assert.Nil(t, resp)
			},
		},
		{
			name:    "missing client id",
			options: scalekit.GenerateClientTokenOptions{},
			mockFn:  func(w http.ResponseWriter, r *http.Request) {},
			assertFn: func(t *testing.T, resp *scalekit.ClientTokenResponse, err error) {
				assert.ErrorIs(t, err, scalekit.ErrClientIdRequired)
				assert.Nil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.mockFn))
			defer server.Close()

			c := scalekit.NewScalekitClient(server.URL, "stored_id", "stored_secret")
			resp, err := c.GenerateClientToken(context.Background(), tt.options)
			tt.assertFn(t, resp, err)
		})
	}
}

func TestGetClientAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/token" {
			_ = r.ParseForm()
			// GetClientAccessToken should use the stored credentials
			assert.Equal(t, "client_credentials", r.FormValue("grant_type"))
			assert.Equal(t, "stored_id", r.FormValue("client_id"))
			assert.Equal(t, "stored_secret", r.FormValue("client_secret"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"stored_cred_token","expires_in":3600}`))
		}
	}))
	defer server.Close()

	c := scalekit.NewScalekitClient(server.URL, "stored_id", "stored_secret")
	token, err := c.GetClientAccessToken(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "stored_cred_token", token)
}

func TestGetAuthorizationUrl(t *testing.T) {
	type testCase struct {
		name        string
		redirectUri string
		options     scalekit.AuthorizationUrlOptions
		wantUrl     string
		wantErr     bool
	}

	tests := []testCase{
		{
			name:        "basic authorization url with default scopes",
			redirectUri: "http://localhost/callback",
			options:     scalekit.AuthorizationUrlOptions{},
			wantUrl:     "http://test.com/oauth/authorize?client_id=client_id&redirect_uri=http%3A%2F%2Flocalhost%2Fcallback&response_type=code&scope=openid+profile+email",
			wantErr:     false,
		},
		{
			name:        "authorization url with custom scopes",
			redirectUri: "http://localhost/callback",
			options: scalekit.AuthorizationUrlOptions{
				Scopes: []string{"openid", "profile", "email", "offline_access"},
			},
			wantUrl: "http://test.com/oauth/authorize?client_id=client_id&redirect_uri=http%3A%2F%2Flocalhost%2Fcallback&response_type=code&scope=openid+profile+email+offline_access",
			wantErr: false,
		},
		{
			name:        "authorization url with all options",
			redirectUri: "http://localhost/callback",
			options: scalekit.AuthorizationUrlOptions{
				ConnectionId:        "conn_123",
				OrganizationId:      "org_123",
				State:               "state123",
				Nonce:               "nonce123",
				DomainHint:          "example.com",
				LoginHint:           "user@example.com",
				CodeChallenge:       "challenge123",
				CodeChallengeMethod: "S256",
				Provider:            "google",
			},
			wantUrl: "http://test.com/oauth/authorize?client_id=client_id&code_challenge=challenge123&code_challenge_method=S256&connection_id=conn_123&domain=example.com&domain_hint=example.com&login_hint=user%40example.com&nonce=nonce123&organization_id=org_123&provider=google&redirect_uri=http%3A%2F%2Flocalhost%2Fcallback&response_type=code&scope=openid+profile+email&state=state123",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := scalekit.NewScalekitClient("http://test.com", "client_id", "client_secret")
			got, err := client.GetAuthorizationUrl(tt.redirectUri, tt.options)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.wantUrl, got.String())
		})
	}
}

func TestGetIdpInitiatedLoginClaims(t *testing.T) {
	type testCase struct {
		name     string
		token    string
		mockFn   func(w http.ResponseWriter, r *http.Request)
		assertFn func(t *testing.T, claims *scalekit.IdpInitiatedLoginClaims, err error)
	}

	tests := []testCase{
		{
			name:  "successful idp initiated login claims",
			token: "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJvcmdhbml6YXRpb25faWQiOiJvcmdfNzIyODk4OTcwMDc4NzQxNTEiLCJjb25uZWN0aW9uX2lkIjoiY29ubl83NTQxNjU3OTA0MjQ3NDIwNCIsImxvZ2luX2hpbnQiOiJzcmluaXZhcy5rYXJyZUBzY2FsZWtpdC5jb20iLCJleHAiOjE5MDY4MDUwNDcsIm5iZiI6MTc0OTAyMDI4NywiaWF0IjoxNzQ5MDIwMjg3LCJpc3MiOiJodHRwOi8vYWlyZGV2LmxvY2FsaG9zdDo4ODg4In0.kMoDVBUlPx-FXKklZnM2ceawYxuL2ccHh9V8lSgf1GncdTVfMlUHgVHvF839JK3b5UiMH0ZIOx2ELTpmCZYjNYP9RNsTmn9JAxxW-K-Mu-tKM7y9k4ZIDpq2MuYrCk_hHgVhdgSDNVnol78PPL8SuLBdZenFNuRBrq4kV9B0x9Mn31QcXL3zoZ4mKV5IRX6ArO7tNT77seXQNSTzF0iMaswri86GP7NfXXBYABRmsULdUKCzn5raWLbqrqiLoIa8ieO81XBYOJiMBvqReUeNfe4hBC2-XJ9txvPBIPlAfT_-9ysOWRpXFGZ4WwNKGSusndxIp103slYAoP2IiCIXRg",
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/keys":
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, claims *scalekit.IdpInitiatedLoginClaims, err error) {
				require.NoError(t, err)
				require.NotNil(t, claims)

				// Verify parsed claims
				assert.Equal(t, "conn_75416579042474204", claims.ConnectionID)
				assert.Equal(t, "org_72289897007874151", claims.OrganizationID)
				assert.Equal(t, "srinivas.karre@scalekit.com", claims.LoginHint)
			},
		},
		{
			name:  "invalid token format",
			token: "invalid.token.format",
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/keys":
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, claims *scalekit.IdpInitiatedLoginClaims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.mockFn))
			defer server.Close()

			client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
			claims, err := client.GetIdpInitiatedLoginClaims(context.Background(), tt.token)
			tt.assertFn(t, claims, err)
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		mockFn   func(w http.ResponseWriter, r *http.Request)
		assertFn func(t *testing.T, isValid bool, err error)
	}{
		{
			name:  "valid access token",
			token: "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vYWlyZGV2LmxvY2FsaG9zdDo4ODg4Iiwic3ViIjoiY29ubl83NTQxNjU3OTA0MjQ3NDIwNDtzcmluaXZhcy5rYXJyZUBzY2FsZWtpdC5jb20iLCJhdWQiOlsicHJkX3NrY18xNzAwMjMzNDIyNzg1NzUwOCJdLCJleHAiOjE5MDY4MDQ4MzcsImlhdCI6MTc0OTAyMDA3NywibmJmIjoxNzQ5MDIwMDc3LCJjbGllbnRfaWQiOiJwcmRfc2tjXzE3MDAyMzM0MjI3ODU3NTA4IiwianRpIjoidGtuXzc1NDE4NDE0MTAwODA1ODUyIn0.SxlKHr1EFBAvfm3Zm7CliKcSWZ8LUFWx8Cs3_3bf1SVouVvRu-zE2_ghB4iAmarsxErurU0kHDEX-Fpx6euemiWXN3Z-mECB4clmb1PF8RThh7bbHx1zxqp3z_MIcDbO4ZKTXMSRx39JbcWyThQSTbeAo50TEFpIT7RsWhNYrBnhsZNibrfZXWUVDBYB930LZMzhdKPRUXBhA-HuKIjggg2jWEAv2leJ3UPbLVccbKrdq2qSzGaxLpvlPoX6RpcrA2Cbuig4vJ7bCy46M-DUg73NO91arPpl5BOnHHx2Oappk_i2S4cMOGdSyX3s50owX1xRDyELNMEIo-VoQ7rfww",
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/keys":
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.NoError(t, err)
				assert.True(t, isValid)
			},
		},
		{
			name:  "invalid token format",
			token: "invalid.token.format",
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/keys":
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.Error(t, err)
				assert.False(t, isValid)
			},
		},
		{
			name:  "empty token",
			token: "",
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/keys":
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.Error(t, err)
				assert.False(t, isValid)
			},
		},
		{
			name:  "jwks endpoint error",
			token: "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vYWlyZGV2LmxvY2FsaG9zdDo4ODg4Iiwic3ViIjoiY29ubl83NTQxNjU3OTA0MjQ3NDIwNDtzcmluaXZhcy5rYXJyZUBzY2FsZWtpdC5jb20iLCJhdWQiOlsicHJkX3NrY18xNzAwMjMzNDIyNzg1NzUwOCJdLCJleHAiOjE5MDY4MDQ4MzcsImlhdCI6MTc0OTAyMDA3NywibmJmIjoxNzQ5MDIwMDc3LCJjbGllbnRfaWQiOiJwcmRfc2tjXzE3MDAyMzM0MjI3ODU3NTA4IiwianRpIjoidGtuXzc1NDE4NDE0MTAwODA1ODUyIn0.SxlKHr1EFBAvfm3Zm7CliKcSWZ8LUFWx8Cs3_3bf1SVouVvRu-zE2_ghB4iAmarsxErurU0kHDEX-Fpx6euemiWXN3Z-mECB4clmb1PF8RThh7bbHx1zxqp3z_MIcDbO4ZKTXMSRx39JbcWyThQSTbeAo50TEFpIT7RsWhNYrBnhsZNibrfZXWUVDBYB930LZMzhdKPRUXBhA-HuKIjggg2jWEAv2leJ3UPbLVccbKrdq2qSzGaxLpvlPoX6RpcrA2Cbuig4vJ7bCy46M-DUg73NO91arPpl5BOnHHx2Oappk_i2S4cMOGdSyX3s50owX1xRDyELNMEIo-VoQ7rfww\n",
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/keys":
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.Error(t, err)
				assert.False(t, isValid)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.mockFn))
			defer server.Close()

			client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
			isValid, err := client.ValidateAccessToken(context.Background(), tt.token)
			tt.assertFn(t, isValid, err)
		})
	}
}

func TestValidateTokenWithOptions(t *testing.T) {
	signedToken := func(claims map[string]interface{}, keyID string) (string, string) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)

		signer, err := jose.NewSigner(
			jose.SigningKey{Algorithm: jose.RS256, Key: privateKey},
			(&jose.SignerOptions{}).WithHeader("kid", keyID),
		)
		require.NoError(t, err)

		tokenPayload, err := json.Marshal(claims)
		require.NoError(t, err)
		idToken, err := signer.Sign(tokenPayload)
		require.NoError(t, err)
		idTokenCompact, err := idToken.CompactSerialize()
		require.NoError(t, err)

		jwk := jose.JSONWebKey{
			Key:       privateKey.Public(),
			KeyID:     keyID,
			Algorithm: string(jose.RS256),
			Use:       "sig",
		}
		keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}}
		keySetBytes, err := json.Marshal(keySet)
		require.NoError(t, err)

		return idTokenCompact, string(keySetBytes)
	}

	validIDToken, validIDTokenJWKS := func() (string, string) {
		now := time.Now()
		return signedToken(map[string]interface{}{
			"sub":            "usr_mock123",
			"name":           "Mock User",
			"email":          "mock@example.com",
			"given_name":     "Mock",
			"family_name":    "User",
			"email_verified": true,
			"iat":            now.Unix(),
			"exp":            now.Add(time.Hour).Unix(),
		}, "mock-id-token-kid")
	}()

	validScopedToken, validScopedTokenJWKS := func() (string, string) {
		now := time.Now()
		return signedToken(map[string]interface{}{
			"sub":   "usr_mock123",
			"iss":   "http://test.com",
			"aud":   []string{"prd_skc_17002334227857508"},
			"iat":   now.Unix(),
			"exp":   now.Add(time.Hour).Unix(),
			"scope": "usr:read usr:write",
		}, "mock-scoped-token-kid")
	}()

	tokenWithoutScopeClaim, tokenWithoutScopeClaimJWKS := func() (string, string) {
		now := time.Now()
		return signedToken(map[string]interface{}{
			"sub": "usr_mock123",
			"iss": "http://test.com",
			"aud": []string{"prd_skc_17002334227857508"},
			"iat": now.Unix(),
			"exp": now.Add(time.Hour).Unix(),
		}, "mock-no-scope-token-kid")
	}()

	tokenWithInvalidScopeClaim, tokenWithInvalidScopeClaimJWKS := func() (string, string) {
		now := time.Now()
		return signedToken(map[string]interface{}{
			"sub":   "usr_mock123",
			"iss":   "http://test.com",
			"aud":   []string{"prd_skc_17002334227857508"},
			"iat":   now.Unix(),
			"exp":   now.Add(time.Hour).Unix(),
			"scope": []string{"usr:read"},
		}, "mock-invalid-scope-token-kid")
	}()

	tests := []struct {
		name     string
		token    string
		options  *scalekit.ValidateTokenOptions
		mockFn   func(w http.ResponseWriter, r *http.Request)
		assertFn func(t *testing.T, isValid bool, err error)
	}{
		{
			name:  "valid access token with matching audience",
			token: "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vYWlyZGV2LmxvY2FsaG9zdDo4ODg4Iiwic3ViIjoiY29ubl83NTQxNjU3OTA0MjQ3NDIwNDtzcmluaXZhcy5rYXJyZUBzY2FsZWtpdC5jb20iLCJhdWQiOlsicHJkX3NrY18xNzAwMjMzNDIyNzg1NzUwOCJdLCJleHAiOjE5MDY4MDQ4MzcsImlhdCI6MTc0OTAyMDA3NywibmJmIjoxNzQ5MDIwMDc3LCJjbGllbnRfaWQiOiJwcmRfc2tjXzE3MDAyMzM0MjI3ODU3NTA4IiwianRpIjoidGtuXzc1NDE4NDE0MTAwODA1ODUyIn0.SxlKHr1EFBAvfm3Zm7CliKcSWZ8LUFWx8Cs3_3bf1SVouVvRu-zE2_ghB4iAmarsxErurU0kHDEX-Fpx6euemiWXN3Z-mECB4clmb1PF8RThh7bbHx1zxqp3z_MIcDbO4ZKTXMSRx39JbcWyThQSTbeAo50TEFpIT7RsWhNYrBnhsZNibrfZXWUVDBYB930LZMzhdKPRUXBhA-HuKIjggg2jWEAv2leJ3UPbLVccbKrdq2qSzGaxLpvlPoX6RpcrA2Cbuig4vJ7bCy46M-DUg73NO91arPpl5BOnHHx2Oappk_i2S4cMOGdSyX3s50owX1xRDyELNMEIo-VoQ7rfww",
			options: &scalekit.ValidateTokenOptions{
				Audience: []string{"prd_skc_17002334227857508"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.NoError(t, err)
				assert.True(t, isValid)
			},
		},
		{
			name:  "valid access token with none of the expected audiences",
			token: "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vYWlyZGV2LmxvY2FsaG9zdDo4ODg4Iiwic3ViIjoiY29ubl83NTQxNjU3OTA0MjQ3NDIwNDtzcmluaXZhcy5rYXJyZUBzY2FsZWtpdC5jb20iLCJhdWQiOlsicHJkX3NrY18xNzAwMjMzNDIyNzg1NzUwOCJdLCJleHAiOjE5MDY4MDQ4MzcsImlhdCI6MTc0OTAyMDA3NywibmJmIjoxNzQ5MDIwMDc3LCJjbGllbnRfaWQiOiJwcmRfc2tjXzE3MDAyMzM0MjI3ODU3NTA4IiwianRpIjoidGtuXzc1NDE4NDE0MTAwODA1ODUyIn0.SxlKHr1EFBAvfm3Zm7CliKcSWZ8LUFWx8Cs3_3bf1SVouVvRu-zE2_ghB4iAmarsxErurU0kHDEX-Fpx6euemiWXN3Z-mECB4clmb1PF8RThh7bbHx1zxqp3z_MIcDbO4ZKTXMSRx39JbcWyThQSTbeAo50TEFpIT7RsWhNYrBnhsZNibrfZXWUVDBYB930LZMzhdKPRUXBhA-HuKIjggg2jWEAv2leJ3UPbLVccbKrdq2qSzGaxLpvlPoX6RpcrA2Cbuig4vJ7bCy46M-DUg73NO91arPpl5BOnHHx2Oappk_i2S4cMOGdSyX3s50owX1xRDyELNMEIo-VoQ7rfww",
			options: &scalekit.ValidateTokenOptions{
				Audience: []string{"non_matching_audience"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.Error(t, err)
				assert.False(t, isValid)
				assert.EqualError(t, err, "none of the expected audiences found in token aud claim")
			},
		},
		{
			name:  "valid access token when any expected audience matches",
			token: "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vYWlyZGV2LmxvY2FsaG9zdDo4ODg4Iiwic3ViIjoiY29ubl83NTQxNjU3OTA0MjQ3NDIwNDtzcmluaXZhcy5rYXJyZUBzY2FsZWtpdC5jb20iLCJhdWQiOlsicHJkX3NrY18xNzAwMjMzNDIyNzg1NzUwOCJdLCJleHAiOjE5MDY4MDQ4MzcsImlhdCI6MTc0OTAyMDA3NywibmJmIjoxNzQ5MDIwMDc3LCJjbGllbnRfaWQiOiJwcmRfc2tjXzE3MDAyMzM0MjI3ODU3NTA4IiwianRpIjoidGtuXzc1NDE4NDE0MTAwODA1ODUyIn0.SxlKHr1EFBAvfm3Zm7CliKcSWZ8LUFWx8Cs3_3bf1SVouVvRu-zE2_ghB4iAmarsxErurU0kHDEX-Fpx6euemiWXN3Z-mECB4clmb1PF8RThh7bbHx1zxqp3z_MIcDbO4ZKTXMSRx39JbcWyThQSTbeAo50TEFpIT7RsWhNYrBnhsZNibrfZXWUVDBYB930LZMzhdKPRUXBhA-HuKIjggg2jWEAv2leJ3UPbLVccbKrdq2qSzGaxLpvlPoX6RpcrA2Cbuig4vJ7bCy46M-DUg73NO91arPpl5BOnHHx2Oappk_i2S4cMOGdSyX3s50owX1xRDyELNMEIo-VoQ7rfww",
			options: &scalekit.ValidateTokenOptions{
				Audience: []string{"non_matching_audience", "prd_skc_17002334227857508"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.NoError(t, err)
				assert.True(t, isValid)
			},
		},
		{
			name:  "invalid token should return original validation error",
			token: "invalid.token.format",
			options: &scalekit.ValidateTokenOptions{
				Audience: []string{"prd_skc_17002334227857508"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.Error(t, err)
				assert.False(t, isValid)
			},
		},
		{
			name:    "valid id token with no audience checks",
			token:   validIDToken,
			options: &scalekit.ValidateTokenOptions{},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(validIDTokenJWKS))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.NoError(t, err)
				assert.True(t, isValid)
			},
		},
		{
			name:    "valid id token with nil options skips audience checks",
			token:   validIDToken,
			options: nil,
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(validIDTokenJWKS))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.NoError(t, err)
				assert.True(t, isValid)
			},
		},
		{
			name:    "invalid token with nil options returns validation error",
			token:   "invalid.token.format",
			options: nil,
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
					_, _ = w.Write([]byte(resp))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.Error(t, err)
				assert.False(t, isValid)
			},
		},
		{
			name:  "valid token with matching scopes",
			token: validScopedToken,
			options: &scalekit.ValidateTokenOptions{
				Scopes: []string{"usr:read", "usr:write"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(validScopedTokenJWKS))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.NoError(t, err)
				assert.True(t, isValid)
			},
		},
		{
			name:  "valid token with missing expected scope",
			token: validScopedToken,
			options: &scalekit.ValidateTokenOptions{
				Scopes: []string{"usr:read", "usr:delete"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(validScopedTokenJWKS))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.Error(t, err)
				assert.False(t, isValid)
				assert.EqualError(t, err, `missing expected scope "usr:delete" in token scope claim`)
			},
		},
		{
			name:  "token missing scope claim",
			token: tokenWithoutScopeClaim,
			options: &scalekit.ValidateTokenOptions{
				Scopes: []string{"usr:read"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(tokenWithoutScopeClaimJWKS))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.Error(t, err)
				assert.False(t, isValid)
				assert.EqualError(t, err, "token missing scope claim")
			},
		},
		{
			name:  "token scope claim must be string",
			token: tokenWithInvalidScopeClaim,
			options: &scalekit.ValidateTokenOptions{
				Scopes: []string{"usr:read"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(tokenWithInvalidScopeClaimJWKS))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.Error(t, err)
				assert.False(t, isValid)
				assert.EqualError(t, err, "token scope claim must be a string")
			},
		},
		{
			name:  "valid token with matching audience and scopes",
			token: validScopedToken,
			options: &scalekit.ValidateTokenOptions{
				Audience: []string{"prd_skc_17002334227857508"},
				Scopes:   []string{"usr:read"},
			},
			mockFn: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/keys" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(validScopedTokenJWKS))
				}
			},
			assertFn: func(t *testing.T, isValid bool, err error) {
				assert.NoError(t, err)
				assert.True(t, isValid)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.mockFn))
			defer server.Close()

			client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
			isValid, err := client.ValidateTokenWithOptions(context.Background(), tt.token, tt.options)
			tt.assertFn(t, isValid, err)
		})
	}
}

func TestGeneratePKCEConfiguration(t *testing.T) {
	type testCase struct {
		name     string
		options  scalekit.PKCEOptions
		assertFn func(t *testing.T, cfg *scalekit.PKCEConfiguration, err error)
	}

	validVerifier := strings.Repeat("a", 43)

	tests := []testCase{
		{
			name:    "defaults to S256 with generated verifier",
			options: scalekit.PKCEOptions{},
			assertFn: func(t *testing.T, cfg *scalekit.PKCEConfiguration, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, "S256", cfg.CodeChallengeMethod)
				assert.Len(t, cfg.CodeVerifier, 64)

				hash := sha256.Sum256([]byte(cfg.CodeVerifier))
				expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
				assert.Equal(t, expectedChallenge, cfg.CodeChallenge)
			},
		},
		{
			name: "uses provided verifier with S256 method",
			options: scalekit.PKCEOptions{
				CodeChallengeMethod: "S256",
				CodeVerifier:        validVerifier,
			},
			assertFn: func(t *testing.T, cfg *scalekit.PKCEConfiguration, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, "S256", cfg.CodeChallengeMethod)
				assert.Equal(t, validVerifier, cfg.CodeVerifier)

				hash := sha256.Sum256([]byte(validVerifier))
				expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
				assert.Equal(t, expectedChallenge, cfg.CodeChallenge)
			},
		},
		{
			name: "fails for plain method",
			options: scalekit.PKCEOptions{
				CodeChallengeMethod: "plain",
				CodeVerifier:        validVerifier,
			},
			assertFn: func(t *testing.T, cfg *scalekit.PKCEConfiguration, err error) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			},
		},
		{
			name: "fails for unsupported method",
			options: scalekit.PKCEOptions{
				CodeChallengeMethod: "sha512",
			},
			assertFn: func(t *testing.T, cfg *scalekit.PKCEConfiguration, err error) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			},
		},
		{
			name: "fails for invalid verifier length",
			options: scalekit.PKCEOptions{
				VerifierLength: 42,
			},
			assertFn: func(t *testing.T, cfg *scalekit.PKCEConfiguration, err error) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			},
		},
		{
			name: "fails for invalid verifier characters",
			options: scalekit.PKCEOptions{
				CodeVerifier: strings.Repeat("a", 42) + "+",
			},
			assertFn: func(t *testing.T, cfg *scalekit.PKCEConfiguration, err error) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			},
		},
	}

	client := scalekit.NewScalekitClient("http://test.com", "client_id", "client_secret")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := client.GeneratePKCEConfiguration(tt.options)
			tt.assertFn(t, cfg, err)
		})
	}
}

func TestNewScalekitClientSecretCompatibilityAndWithSecret(t *testing.T) {
	tests := []struct {
		name           string
		expectedSecret string
		clientFn       func(serverURL string) scalekit.Scalekit
		assertFn       func(t *testing.T, serverURL string)
	}{
		{
			name:           "uses variadic string as backward-compatible client_secret",
			expectedSecret: "client_secret",
			clientFn: func(serverURL string) scalekit.Scalekit {
				return scalekit.NewScalekitClient(serverURL, "client_id", "client_secret")
			},
			assertFn: func(t *testing.T, serverURL string) {
				client := scalekit.NewScalekitClient(serverURL, "client_id", "client_secret")
				resp, err := client.RefreshAccessToken(context.Background(), "refresh_token")
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, "it", resp.IdToken)
			},
		},
		{
			name:           "WithSecret overwrites client_secret",
			expectedSecret: "new_secret",
			clientFn: func(serverURL string) scalekit.Scalekit {
				return scalekit.NewScalekitClient(serverURL, "client_id").WithSecret("new_secret")
			},
			assertFn: func(t *testing.T, serverURL string) {
				client := scalekit.NewScalekitClient(serverURL, "client_id").WithSecret("new_secret")
				resp, err := client.RefreshAccessToken(context.Background(), "refresh_token")
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, "it", resp.IdToken)
			},
		},
		{
			name:           "WithSecret returns isolated client and does not mutate original",
			expectedSecret: "base_secret",
			clientFn: func(serverURL string) scalekit.Scalekit {
				return scalekit.NewScalekitClient(serverURL, "client_id", "base_secret")
			},
			assertFn: func(t *testing.T, serverURL string) {
				baseClient := scalekit.NewScalekitClient(serverURL, "client_id", "base_secret")
				derivedClient := baseClient.WithSecret("new_secret")

				_, err := baseClient.RefreshAccessToken(context.Background(), "refresh_token")
				require.NoError(t, err)
				_, err = derivedClient.RefreshAccessToken(context.Background(), "refresh_token")
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientSecrets := []string{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/oauth/token", r.URL.Path)
				require.NoError(t, r.ParseForm())
				clientSecrets = append(clientSecrets, r.FormValue("client_secret"))

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"access_token":"at","id_token":"it","refresh_token":"rt","expires_in":3600}`))
			}))
			defer server.Close()

			if tt.assertFn != nil {
				tt.assertFn(t, server.URL)
			} else {
				testClient := tt.clientFn(server.URL)
				resp, err := testClient.RefreshAccessToken(context.Background(), "refresh_token")
				require.NoError(t, err)
				require.NotNil(t, resp)
			}

			if tt.name == "WithSecret returns isolated client and does not mutate original" {
				require.Equal(t, []string{"base_secret", "new_secret"}, clientSecrets)
				return
			}
			require.NotEmpty(t, clientSecrets)
			require.Equal(t, tt.expectedSecret, clientSecrets[0])
		})
	}
}

func TestAuthenticateWithCode_Non2xx(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		body           string
		wantErrContain string
	}{
		{
			name:           "401 unauthorized",
			statusCode:     http.StatusUnauthorized,
			body:           `{"error":"invalid_client"}`,
			wantErrContain: "authentication failed: HTTP 401",
		},
		{
			name:           "400 bad request",
			statusCode:     http.StatusBadRequest,
			body:           `{"error":"invalid_grant"}`,
			wantErrContain: "authentication failed: HTTP 400",
		},
		{
			name:           "503 unavailable with empty body",
			statusCode:     http.StatusServiceUnavailable,
			body:           "",
			wantErrContain: "authentication failed: HTTP 503",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
			_, err := client.AuthenticateWithCode(context.Background(), "code", "http://localhost/cb", scalekit.AuthenticationOptions{})
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErrContain)
			// Non-2xx responses must be extractable as *scalekit.Error with StatusCode set.
			var sdkErr *scalekit.Error
			require.True(t, errors.As(err, &sdkErr), "err should be *scalekit.Error for non-2xx")
			assert.Equal(t, tt.statusCode, sdkErr.StatusCode)
		})
	}
}

func TestValidateAccessToken_JwksError(t *testing.T) {
	tests := []struct {
		name           string
		jwksStatus     int
		wantErrContain string
	}{
		{
			name:           "JWKS 500 internal server error",
			jwksStatus:     http.StatusInternalServerError,
			wantErrContain: "failed to fetch JWKS: HTTP 500",
		},
		{
			name:           "JWKS 403 forbidden",
			jwksStatus:     http.StatusForbidden,
			wantErrContain: "failed to fetch JWKS: HTTP 403",
		},
	}

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vYWlyZGV2LmxvY2FsaG9zdDo4ODg4Iiwic3ViIjoiY29ubl83NTQxNjU3OTA0MjQ3NDIwNDtzcmluaXZhcy5rYXJyZUBzY2FsZWtpdC5jb20iLCJhdWQiOlsicHJkX3NrY18xNzAwMjMzNDIyNzg1NzUwOCJdLCJleHAiOjE5MDY4MDQ4MzcsImlhdCI6MTc0OTAyMDA3NywibmJmIjoxNzQ5MDIwMDc3LCJjbGllbnRfaWQiOiJwcmRfc2tjXzE3MDAyMzM0MjI3ODU3NTA4IiwianRpIjoidGtuXzc1NDE4NDE0MTAwODA1ODUyIn0.SxlKHr1EFBAvfm3Zm7CliKcSWZ8LUFWx8Cs3_3bf1SVouVvRu-zE2_ghB4iAmarsxErurU0kHDEX-Fpx6euemiWXN3Z-mECB4clmb1PF8RThh7bbHx1zxqp3z_MIcDbO4ZKTXMSRx39JbcWyThQSTbeAo50TEFpIT7RsWhNYrBnhsZNibrfZXWUVDBYB930LZMzhdKPRUXBhA-HuKIjggg2jWEAv2leJ3UPbLVccbKrdq2qSzGaxLpvlPoX6RpcrA2Cbuig4vJ7bCy46M-DUg73NO91arPpl5BOnHHx2Oappk_i2S4cMOGdSyX3s50owX1xRDyELNMEIo-VoQ7rfww"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.jwksStatus)
			}))
			defer server.Close()

			client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
			_, err := client.ValidateAccessToken(context.Background(), token)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErrContain)
			var sdkErr *scalekit.Error
			require.True(t, errors.As(err, &sdkErr), "JWKS non-2xx should be *scalekit.Error")
			assert.Equal(t, tt.jwksStatus, sdkErr.StatusCode)
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	c := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
	ctx := context.Background()

	tests := []struct {
		name         string
		fn           func() error
		wantSentinel error
	}{
		{
			name:         "ErrTokenRequired on empty ValidateToken",
			fn:           func() error { _, err := c.Token().ValidateToken(ctx, ""); return err },
			wantSentinel: scalekit.ErrTokenRequired,
		},
		{
			name:         "ErrTokenRequired on empty InvalidateToken",
			fn:           func() error { return c.Token().InvalidateToken(ctx, "") },
			wantSentinel: scalekit.ErrTokenRequired,
		},
		{
			name:         "ErrRefreshTokenRequired on empty refresh token",
			fn:           func() error { _, err := c.RefreshAccessToken(ctx, ""); return err },
			wantSentinel: scalekit.ErrRefreshTokenRequired,
		},
		{
			name:         "ErrCodeOrLinkTokenRequired on nil options",
			fn:           func() error { _, err := c.Passwordless().VerifyPasswordlessEmail(ctx, nil); return err },
			wantSentinel: scalekit.ErrCodeOrLinkTokenRequired,
		},
		{
			name: "ErrCodeOrLinkTokenRequired on non-nil empty options",
			fn: func() error {
				_, err := c.Passwordless().VerifyPasswordlessEmail(ctx, &scalekit.VerifyPasswordlessOptions{})
				return err
			},
			wantSentinel: scalekit.ErrCodeOrLinkTokenRequired,
		},
		{
			name:         "ErrOrganizationIdRequired on CreateToken with empty orgId",
			fn:           func() error { _, err := c.Token().CreateToken(ctx, "", scalekit.CreateTokenOptions{}); return err },
			wantSentinel: scalekit.ErrOrganizationIdRequired,
		},
		{
			name: "ErrCodeRequired on AuthenticateWithCode with empty code",
			fn: func() error {
				_, err := c.AuthenticateWithCode(ctx, "", "http://localhost/cb", scalekit.AuthenticationOptions{})
				return err
			},
			wantSentinel: scalekit.ErrCodeRequired,
		},
		{
			name: "ErrRedirectUriRequired on AuthenticateWithCode with empty redirectUri",
			fn: func() error {
				_, err := c.AuthenticateWithCode(ctx, "code", "", scalekit.AuthenticationOptions{})
				return err
			},
			wantSentinel: scalekit.ErrRedirectUriRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantSentinel)
		})
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	noopJwks := func(_ context.Context) (*jose.JSONWebKeySet, error) {
		return &jose.JSONWebKeySet{}, nil
	}
	_, err := scalekit.ValidateToken[map[string]interface{}](context.Background(), "", noopJwks)
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrTokenRequired)
	// Backward compatibility: empty token also matches ErrTokenValidationFailed (joined by SDK).
	assert.ErrorIs(t, err, scalekit.ErrTokenValidationFailed)
}

func TestValidateToken_MissingExpClaim(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	kid := "test-key"
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: privKey},
		(&jose.SignerOptions{}).WithHeader("kid", kid),
	)
	require.NoError(t, err)

	// JWT with no exp field
	payload, _ := json.Marshal(map[string]interface{}{"sub": "user123", "iss": "test"})
	jws, err := signer.Sign(payload)
	require.NoError(t, err)
	token, err := jws.CompactSerialize()
	require.NoError(t, err)

	jwksFn := func(_ context.Context) (*jose.JSONWebKeySet, error) {
		return &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
			{Key: &privKey.PublicKey, KeyID: kid, Algorithm: string(jose.RS256)},
		}}, nil
	}

	_, err = scalekit.ValidateToken[map[string]interface{}](context.Background(), token, jwksFn)
	require.Error(t, err)
	assert.ErrorIs(t, err, scalekit.ErrMissingExpClaim)
	// Backward compat: ErrInvalidExpClaimFormat is an alias for ErrMissingExpClaim
	assert.ErrorIs(t, err, scalekit.ErrInvalidExpClaimFormat) //nolint:staticcheck // testing deprecated alias
}

// TestConnectRetryOn401 verifies that a 401 on a Connect RPC triggers re-auth and one retry.
// The success response uses manual gRPC wire framing; this is fragile against Connect protocol
// changes. For stability, consider using connect/connecttest or a real connectrpc handler.
func TestConnectRetryOn401(t *testing.T) {
	listOrganizationPath := organizationsconnect.OrganizationServiceListOrganizationProcedure
	var rpcCallCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oauth/token" && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"new_token","expires_in":3600}`))
		case r.URL.Path == listOrganizationPath:
			n := rpcCallCount.Add(1)
			if n == 1 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// Second call: success with empty list (manual gRPC wire: 5-byte prefix + marshaled message).
			// Connect expects Grpc-Status as a trailer (after the body); declare it and set it after writing.
			msg, err := proto.Marshal(&organizationsv1.ListOrganizationsResponse{})
			require.NoError(t, err)
			w.Header().Set("Content-Type", "application/grpc")
			w.Header().Set("Trailer", "Grpc-Status")
			w.WriteHeader(http.StatusOK)
			prefix := make([]byte, 5)
			prefix[0] = 0 // no compression
			binary.BigEndian.PutUint32(prefix[1:5], uint32(len(msg)))
			_, _ = w.Write(prefix)
			_, _ = w.Write(msg)
			w.Header().Set("Grpc-Status", "0")
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")
	ctx := context.Background()
	resp, err := client.Organization().ListOrganization(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int32(2), rpcCallCount.Load(), "ListOrganization should be called twice (401 then retry)")
}
