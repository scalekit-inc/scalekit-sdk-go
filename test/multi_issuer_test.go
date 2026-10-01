package test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	scalekit "github.com/scalekit-inc/scalekit-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	multiIssuerBase     = "https://acme.scalekit.cloud"
	multiIssuerResource = "https://acme.scalekit.cloud/resources/res_123"
)

// newMultiIssuerFixture returns a client whose JWKS endpoint is a local server
// (no network beyond loopback) and a helper that signs tokens with the matching key.
func newMultiIssuerFixture(t *testing.T) (scalekit.Scalekit, func(claims map[string]any) string) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
		Key: &privateKey.PublicKey, KeyID: "test-key-1", Algorithm: string(jose.RS256), Use: "sig",
	}}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	t.Cleanup(server.Close)

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: jose.JSONWebKey{Key: privateKey, KeyID: "test-key-1"}},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	require.NoError(t, err)

	sign := func(claims map[string]any) string {
		if _, ok := claims["exp"]; !ok {
			claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
		}
		payload, err := json.Marshal(claims)
		require.NoError(t, err)
		jws, err := signer.Sign(payload)
		require.NoError(t, err)
		token, err := jws.CompactSerialize()
		require.NoError(t, err)
		return token
	}

	return scalekit.NewScalekitClient(server.URL, "client_id", "client_secret"), sign
}

func TestMultiIssuerValidation(t *testing.T) {
	c, sign := newMultiIssuerFixture(t)
	ctx := context.Background()

	withIss := func(iss string) string { return sign(map[string]any{"iss": iss, "sub": "user_1"}) }

	cases := []struct {
		name    string
		token   string
		options *scalekit.ValidateTokenOptions
		valid   bool
	}{
		// existing single-string behavior (unchanged)
		{"Issuer string matches", withIss(multiIssuerBase), &scalekit.ValidateTokenOptions{Issuer: multiIssuerBase}, true},
		{"Issuer string mismatch", withIss(multiIssuerBase), &scalekit.ValidateTokenOptions{Issuer: multiIssuerResource}, false},

		// Issuers: valid if iss equals ANY entry
		{"Issuers matches 1st", withIss(multiIssuerBase), &scalekit.ValidateTokenOptions{Issuers: []string{multiIssuerBase, multiIssuerResource}}, true},
		{"Issuers matches 2nd (resource-scoped token)", withIss(multiIssuerResource), &scalekit.ValidateTokenOptions{Issuers: []string{multiIssuerBase, multiIssuerResource}}, true},
		{"Issuers matches none", withIss("https://acme.scalekit.cloud/resources/res_other"), &scalekit.ValidateTokenOptions{Issuers: []string{multiIssuerBase, multiIssuerResource}}, false},
		{"Issuers list of one behaves like a string", withIss(multiIssuerBase), &scalekit.ValidateTokenOptions{Issuers: []string{multiIssuerResource}}, false},
		{"exact match only, no trailing-slash normalization", withIss(multiIssuerBase + "/"), &scalekit.ValidateTokenOptions{Issuers: []string{multiIssuerBase}}, false},

		// Issuer and Issuers combine: any-of across both
		{"Issuer + Issuers: matches Issuer", withIss(multiIssuerBase), &scalekit.ValidateTokenOptions{Issuer: multiIssuerBase, Issuers: []string{multiIssuerResource}}, true},
		{"Issuer + Issuers: matches Issuers entry", withIss(multiIssuerResource), &scalekit.ValidateTokenOptions{Issuer: multiIssuerBase, Issuers: []string{multiIssuerResource}}, true},
		{"Issuer + Issuers: matches neither", withIss("https://evil.example.com"), &scalekit.ValidateTokenOptions{Issuer: multiIssuerBase, Issuers: []string{multiIssuerResource}}, false},

		// issuer check skipped
		{"nil options skips", withIss(multiIssuerResource), nil, true},
		{"empty options skips", withIss(multiIssuerResource), &scalekit.ValidateTokenOptions{}, true},
		{"nil Issuers skips", withIss(multiIssuerResource), &scalekit.ValidateTokenOptions{Issuers: nil}, true},
		{"empty Issuers skips", withIss(multiIssuerResource), &scalekit.ValidateTokenOptions{Issuers: []string{}}, true},

		// non-empty Issuers is always enforced, even with blank entries (fails closed)
		{"blank-only Issuers fails closed", withIss(multiIssuerResource), &scalekit.ValidateTokenOptions{Issuers: []string{""}}, false},
		{"blank entries harmless next to a real one", withIss(multiIssuerResource), &scalekit.ValidateTokenOptions{Issuers: []string{"", multiIssuerResource}}, true},
		{"blank entries do not widen the match", withIss(multiIssuerResource), &scalekit.ValidateTokenOptions{Issuers: []string{"", multiIssuerBase}}, false},
		{"token without iss never matches a blank entry", sign(map[string]any{"sub": "user_1"}), &scalekit.ValidateTokenOptions{Issuers: []string{""}}, false},
		{"token without iss never matches Issuers", sign(map[string]any{"sub": "user_1"}), &scalekit.ValidateTokenOptions{Issuers: []string{"", multiIssuerBase}}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := c.ValidateTokenWithOptions(ctx, tc.token, tc.options)
			if tc.valid {
				require.NoError(t, err)
				assert.True(t, valid)
				return
			}
			require.Error(t, err)
			assert.False(t, valid)
			assert.ErrorIs(t, err, scalekit.ErrIssuerMismatch)
		})
	}
}
