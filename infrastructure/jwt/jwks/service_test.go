package jwks

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const tokenString = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIzVHVzam9EQkhnZWVXSHdGaVNuWktXZ3ZsRGt4QXRreU9GZDdVeEVEQ0lnIn0.eyJleHAiOjE1OTk4MjQ2MjgsImlhdCI6MTU5OTgyNDU2OCwiYXV0aF90aW1lIjoxNTk5ODI0NTY3LCJqdGkiOiI4MzAwY2I4ZS1iNGY5LTQ2ZWEtODcyNi1hMjJlZGM3NzU3ODQiLCJpc3MiOiJodHRwczovL2tleWNsb2FrLnRlc3Qta2V5Y2xvYWsuY29kZS1mYWJyaWsuY29tL2F1dGgvcmVhbG1zL21hc3RlciIsImF1ZCI6WyJodHRwczovL3N0YXJ0dXAtYmF0dGxlIiwibWFzdGVyLXJlYWxtIiwiYWNjb3VudCJdLCJzdWIiOiJhOGFhNDU3YS0xYzcyLTRmMDgtYTBkNS1hM2Y2ODVjNTVhM2YiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ2dWUtYXBwIiwibm9uY2UiOiJlZjFkMjRmNy00ZGE4LTRlZTItOGRhYS1hMGMwNWVlZGFlOTQiLCJzZXNzaW9uX3N0YXRlIjoiZDlmZWJjODQtZTBhNS00MjFjLWIyZDctN2Q3Yzk1MTQ3MzI5IiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwOi8vbG9jYWxob3N0OjgwODAiLCJodHRwOi8vbG9jYWxob3N0OjgwODEiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbImNyZWF0ZS1yZWFsbSIsIm9mZmxpbmVfYWNjZXNzIiwiYWRtaW4iLCJqZW5raW5zX2FkbWluIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJtYXN0ZXItcmVhbG0iOnsicm9sZXMiOlsidmlldy1pZGVudGl0eS1wcm92aWRlcnMiLCJ2aWV3LXJlYWxtIiwibWFuYWdlLWlkZW50aXR5LXByb3ZpZGVycyIsImltcGVyc29uYXRpb24iLCJjcmVhdGUtY2xpZW50IiwibWFuYWdlLXVzZXJzIiwicXVlcnktcmVhbG1zIiwidmlldy1hdXRob3JpemF0aW9uIiwicXVlcnktY2xpZW50cyIsInF1ZXJ5LXVzZXJzIiwibWFuYWdlLWV2ZW50cyIsIm1hbmFnZS1yZWFsbSIsInZpZXctZXZlbnRzIiwidmlldy11c2VycyIsInZpZXctY2xpZW50cyIsIm1hbmFnZS1hdXRob3JpemF0aW9uIiwibWFuYWdlLWNsaWVudHMiLCJxdWVyeS1ncm91cHMiXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsIm5hbWUiOiJhZG1pbkZOIGFkbWluTE4iLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJhZG1pbiIsImdpdmVuX25hbWUiOiJhZG1pbkZOIiwiZmFtaWx5X25hbWUiOiJhZG1pbkxOIiwicGljdHVyZSI6Ii9pbWcvbG9nby44MmI5YzdhNS5wbmciLCJlbWFpbCI6ImFkbWluQHRlc3QuZGUifQ.qcvmLKTHhMcLKL1RNhMctVRKQbms82z_qLcMnvmJtJEcUzOdR7zQ23MWEcw90tNdQYjKXRDeW1HtcFAEaSj_YeVCSiyNfW0vIk3HIBk9zB-yzE2kFfhuJOQYoLKIndbVPaq8jH2baof0oDi_TOISTyC7h8iR3vO3LqeWlmv4jwOGgpqgP6Cw_Ujtku7FI-nFdXIqYfPgB1l1HM8hMndwQl-w2vrliu1nQAYeDQMQZLxWHQmOr0VjtrqUDrtZ9eMlUm7KzLvS7TZSAX0yf2Dfnz2Xr1byx1VuJLBRw3ipt56uraRZC4_h2I1LuSxoYY5actbSW32rArt3Fxd9tmhnyw"

func Test_CallerUnauthorizedTokenExpired(t *testing.T) {
	_, err := parseAndValidate(tokenString, fetchPublicKeyForTest)

	assert.NotNil(t, err, fmt.Sprintf("create user should return error"))

	assert.Equal(t, "Token is expired", err.Error(), "unexpected root cause")
}

const publicKey = `
-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArUd/eupoAY5sSpslzrZV
LxEJCIx9twWx2m0F2hEi+P6x8fXTZoP5rFr6m0u4b9X7uqoN2aTpaRshxRGBhb77
BFH1XtDlg+j1HzECa4DwUOoPUNJjfaWpDUtvoDhKBdeJcdVhFBh5QDpWtUbn0DTs
gqUlK9ht5yHZiUHmSxTMOvkXbzEkv0SKeGQhSHoSI8XWxluRRLHbklIVhwT7t+XJ
y/OcrIXJzXUHhiGV76MqFpv+f7snGciOL8+U2tAgynpzLVLlUklEz6+xihPQQgPj
thvsGjPp0MgkPNX5MWR4cTqFabzrrsZqlkwZkT1Ai5tiu9hHiNrg4GqUudQdSLS9
wwIDAQAB
-----END RSA PUBLIC KEY-----
`

var fetchPublicKeyForTest = func(string) (*rsa.PublicKey, error) {
	publicKey, err := parseRsaPublicKeyFromPemStr(publicKey)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

func parseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}
