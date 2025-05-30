package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scalekit-inc/scalekit-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateWithCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/keys":
			{
				// Return mock token response {
				w.Header().Set("Content-Type", "application/json")
				resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
				_, err := w.Write([]byte(resp))
				if err != nil {
					return
				}
			}
		case "/oauth/token":
			{
				// Return mock token response
				w.Header().Set("Content-Type", "application/json")
				resp := `{
				"access_token": "mock_access_token",
				"id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJhbXIiOlsiY29ubl82MDM1ODU1NTczOTM4ODIxMCJdLCJhdF9oYXNoIjoid2hCTHlyWVJFdGtXaHY2ekM2T09hdyIsImF1ZCI6WyJwcmRfc2tjXzE3MDAyMzM0MjI3ODU3NTA4Il0sImF6cCI6InByZF9za2NfMTcwMDIzMzQyMjc4NTc1MDgiLCJjX2hhc2giOiJHY2NRZW9tSG1JNmNqNTUyOUtnenFRIiwiY2xpZW50X2lkIjoicHJkX3NrY18xNzAwMjMzNDIyNzg1NzUwOCIsImVtYWlsIjoic3JpbnZhc2thcnJhQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJleHAiOjE3NDg1NDE4NzIsImdpdmVuX25hbWUiOiJzcmludmFza2FycmEiLCJpYXQiOjE3NDg1NDAwNzIsImlzcyI6Imh0dHA6Ly9haXJkZXYubG9jYWxob3N0Ojg4ODgiLCJuYW1lIjoiZ21haWwiLCJvaWQiOiJvcmdfNzIyOTgzMDI5ODAyMjcxNzUiLCJzaWQiOiJzZXNfNzQ2MTI5NDQ4OTMxMTY1MTgiLCJzdWIiOiJ1c3JfNzIyOTc4NTM0NjgyNzg4ODcifQ.Arti6kfBAjJI2sxy97bTGJwANKOdjfxfIBAdpEeL931pG-Rc89iN9vyyKK6V2W4CSAIF1qsWYJwVeSg0yKBC-w94n-79x5D1f3AydVE_Pp-YSN_8asLJlWQrbnQPOI6SSlItVQdV_1ag2D_CcpQpkYNhrv_AHC9fmIhlabMWCYx-vRFKqr0Jj9BWVjkynIG6wb3m7lbijt2_bnF135-3ob7dRJ0B_f0ZdIBli_numj6ik5Q-PpHrUP5UcZHO0ieE2jqC_z9sF-Msmn2xUYPhJCd2JkFOaEKDULI5k_-01Gyk-1zFWNBDJjKiFu8SjIQDU5nGVc2Hrbptxu7Aoqx8BA",
				"expires_in": 3600
				}`
				_, err := w.Write([]byte(resp))
				if err != nil {
					return
				}
			}
		}
	}))

	client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")

	resp, err := client.AuthenticateWithCode("test_code", "http://localhost/callback", scalekit.AuthenticationOptions{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "mock_access_token", resp.AccessToken)
	assert.Equal(t, 3600, resp.ExpiresIn)

	// Verify parsed claims
	assert.Equal(t, "usr_72297853468278887", resp.User.Id)
	assert.Equal(t, "gmail", resp.User.Name)
	assert.Equal(t, "srinvaskarra@gmail.com", resp.User.Email)

	// Verify custom claim
	rawClaims := resp.User.Claims
	assert.Equal(t, "usr_72297853468278887", rawClaims["sub"])
	assert.Equal(t, "org_72298302980227175", rawClaims["oid"])
	assert.Equal(t, "ses_74612944893116518", rawClaims["sid"])
}

func TestGetAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/keys":
			{
				// Return mock token response {
				w.Header().Set("Content-Type", "application/json")
				resp := `{"keys":[{"use":"sig","kty":"RSA","kid":"snk_17002334227791972","alg":"RS256","n":"8HgCyscnWpT78Jscy7GOSrdK30R8AkBu7BSsXPnWNTCBMmdoRYa2kJf4al9XXW28FNYwM9oHAxCFsiRQna_ouClsRyW1_rYXxqQeeW4GvI1uRpq-3kgRvDm1cjekXH4a0bu_cGNcdTVherrUiBH3WoHxnIMTO0i__BD0qbyh4teUfYaoRgE8T-zsBB_QGdDfMl7EfGLIFgI8eTZFGn_-ONpV9Z9HvVefnyr4Oibyu58z77cOytd6r4lCF0dErAUkjiPNk-cTUDv-QRBNLG4uNcLEqgKL-nvNW-7JrUMiWCcrkHKUlwUncuMvbwWrLlT_dJp7XRjN8RampGUEQUbzGw","e":"AQAB"}]}`
				_, err := w.Write([]byte(resp))
				if err != nil {
					return
				}
			}
		}
	}))
	
	client := scalekit.NewScalekitClient(server.URL, "client_id", "client_secret")

	accessToken, err := client.GetAccessToken("eyJhbGciOiJSUzI1NiIsImtpZCI6InNua18xNzAwMjMzNDIyNzc5MTk3MiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsicHJkX3NrY18xNzAwMjMzNDIyNzg1NzUwOCJdLCJjbGllbnRfaWQiOiJwcmRfc2tjXzE3MDAyMzM0MjI3ODU3NTA4IiwiZXhwIjoxNzQ4NjAwMzAwLCJpYXQiOjE3NDg1OTk0MDAsImlzcyI6Imh0dHA6Ly9haXJkZXYubG9jYWxob3N0Ojg4ODgiLCJqdGkiOiJ0a25fNzQ3MTI2MzQyMjE0OTIzODEiLCJuYmYiOjE3NDg1OTk0MDAsInNpZCI6InNlc183NDcxMjYzMzI4MTkwMjc0OSIsInN1YiI6InVzcl83MjI5Nzg1MzQ2ODI3ODg4NyJ9.Dj2Le9PEFd5hcGub1QlPu5oa58gVbOKwZXedVa_AJl_4bxdXPAB7iZqOL_eGFWCpwXe6yrA8cuZCOQvdL3EDLEMHpLkFMo0LeXPf2ukuAN8VrRNlVG6rAdGxxwvDdVP4bC1-4m6atqAtIk8sYYFb1Hd8CP4B42VVr3oV6RoOGlVhuIFdKKO6Sin2hoVnZsLsm6Q6u3nc17GNf6wmskKCwktnooEAv7L1Mp_SJYNNDOyTBtsXoAkeK4CzNwVwHivys_dSd4euMDkMIyEQTX-rYzWLnd4iemJAzQLCKfmTVYGagh1Cnxr92T4wShWT5rwi5XorDCJA-LrEmNPeT4OKHg")

	assert.NoError(t, err)
	assert.NotNil(t, accessToken)
	assert.Equal(t, 1748600300, accessToken.Exp)

	// Verify parsed claims
	assert.Equal(t, "usr_72297853468278887", accessToken.Sub)
	assert.Equal(t, scalekit.Audience{"prd_skc_17002334227857508"}, accessToken.Audience)

	// Verify custom claim
	rawClaims := accessToken.Claims
	assert.Equal(t, "usr_72297853468278887", rawClaims["sub"])
	assert.Equal(t, "ses_74712633281902749", rawClaims["sid"])
}
