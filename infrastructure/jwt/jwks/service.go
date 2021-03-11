package jwks

import (
	"code-fabrik.com/bend/domain/authentication"
	"code-fabrik.com/bend/infrastructure/env"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/s12v/go-jwks"
)

type FetchPublicKey func(key string) (*rsa.PublicKey, error)

type JwksService struct{}

// Returns user id from JWT token.
// Retrieves public key from OpenID Connect server and checks token signature, expiry and claims.
// Returns error if user object cannot be created for whatever reason (caller unauthorized, mandatory values not present).
func (a *JwksService) BuildUserFromJWT(token string) (*authentication.User, error) {
	user, err := parseAndValidate(token, fetchJwksPublicKey)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func getClaim(key string, claims jwt.MapClaims) (string, error) {
	if v, ok := claims[key].(string); ok && len(v) > 0 {
		return v, nil
	}
	return "", fmt.Errorf("missing claim %v", key)
}

func parseAndValidate(tokenString string, fetchPublicKey FetchPublicKey) (*authentication.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		keyId := token.Header["kid"]
		if keyId == nil {
			return nil, fmt.Errorf("missing key id")
		}

		publicKey, err := fetchPublicKey(keyId.(string))
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, err := getClaim("sub", claims)
		email, err := getClaim("email", claims)
		familyName, err := getClaim("family_name", claims)
		givenName, err := getClaim("given_name", claims)
		preferredUsername, _ := getClaim("preferred_username", claims)

		if err != nil {
			return nil, err
		}

		user := authentication.User{
			Sub:               sub,
			Email:             email,
			FamilyName:        familyName,
			GivenName:         givenName,
			PreferredUsername: preferredUsername,
		}
		return &user, nil
	} else {
		return nil, fmt.Errorf("claims not found or token invalid")
	}
}

func jwksUrl() string {
	return fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/certs", env.KEYCLOAK_HOST, env.KEYCLOAK_REALM)
}

// Calls OIDC certs endpoint and retrieves RSA public key ("use":"sig") corresponding to given key id.
// Does not cache known keys.
func fetchJwksPublicKey(keyId string) (*rsa.PublicKey, error) {
	webKeySet, err := jwks.NewWebSource(jwksUrl()).JSONWebKeySet()
	if err != nil {
		return nil, err
	}

	key := webKeySet.Key(keyId)
	if len(key) != 1 {
		return nil, fmt.Errorf("error retrieving public key")
	}

	return key[0].Key.(*rsa.PublicKey), nil
}
