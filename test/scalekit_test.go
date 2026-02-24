package test

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
							"family_name":   "User",
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
			expectedError: "Missing required headers",
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
			expectedError: "Message timestamp too old",
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
