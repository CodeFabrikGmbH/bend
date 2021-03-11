package keycloak

import (
	"code-fabrik.com/bend/domain/authentication"
	"code-fabrik.com/bend/infrastructure/env"
	"code-fabrik.com/bend/infrastructure/jwt/jwks"
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v7"
	"net/http"
	"time"
)

type Config struct {
	HostName     string
	ClientId     string
	ClientSecret string
	Realm        string
}

type Service struct {
	config Config
}

var (
	defaultUser = authentication.User{}
)

func New() *Service {
	config := Config{
		HostName:     env.KEYCLOAK_HOST,
		ClientId:     env.KEYCLOAK_CLIENT_ID,
		ClientSecret: "",
		Realm:        env.KEYCLOAK_REALM,
	}

	return &Service{config}
}

func (k *Service) Authenticate(w http.ResponseWriter, r *http.Request) (*authentication.User, error) {
	if len(env.KEYCLOAK_HOST) == 0 {
		return &defaultUser, nil
	}

	accessCookie, err := r.Cookie("access")
	if err != nil {
		return nil, err
	}

	jwksService := jwks.JwksService{}

	if user, err := jwksService.BuildUserFromJWT(accessCookie.Value); err == nil {
		return user, nil
	}

	refreshCookie, err := r.Cookie("refresh")
	if err != nil {
		return nil, err
	}

	return k.BuildUserFromOIDCResponse(r.Context(), w, accessCookie.Value, refreshCookie.Value)
}

func (k *Service) BuildUserFromOIDCResponse(ctx context.Context, w http.ResponseWriter, accessToken, refreshToken string) (*authentication.User, error) {
	if len(accessToken) == 0 {
		return nil, fmt.Errorf("access token undefined")
	}

	client := gocloak.NewClient(k.config.HostName)

	if userInfo, err := client.GetUserInfo(ctx, accessToken, k.config.Realm); err == nil {
		return k.newUser(userInfo)
	}

	jwt, err := client.RefreshToken(ctx, refreshToken, k.config.ClientId, k.config.ClientSecret, k.config.Realm)
	if err != nil {
		return nil, err
	}

	return k.handleNewJWT(ctx, w, err, client, jwt)
}

func (k *Service) Login(ctx context.Context, w http.ResponseWriter, username, password string) (*authentication.User, error) {
	client := gocloak.NewClient(k.config.HostName)

	jwt, err := client.Login(ctx, k.config.ClientId, k.config.ClientSecret, k.config.Realm, username, password)
	if err != nil {
		return nil, err
	}

	return k.handleNewJWT(ctx, w, err, client, jwt)
}

func (k *Service) handleNewJWT(ctx context.Context, w http.ResponseWriter, err error, client gocloak.GoCloak, jwt *gocloak.JWT) (*authentication.User, error) {
	userInfo, err := client.GetUserInfo(ctx, jwt.AccessToken, k.config.Realm)
	if err != nil {
		return nil, err
	}

	err = k.saveTokenAsCookie(jwt, w)
	if err != nil {
		return nil, err
	}

	return k.newUser(userInfo)
}

func (k *Service) Logout(_ context.Context, w http.ResponseWriter) {
	k.deleteCookies(w)
}

func (k *Service) saveTokenAsCookie(token *gocloak.JWT, w http.ResponseWriter) error {
	if token == nil {
		return fmt.Errorf("token undefined")
	}

	accessCookie := http.Cookie{Path: "/", Name: "access", Value: token.AccessToken, Expires: time.Now().Add(365 * 24 * time.Hour)}
	http.SetCookie(w, &accessCookie)

	refreshCookie := http.Cookie{Path: "/", Name: "refresh", Value: token.RefreshToken, Expires: time.Now().Add(365 * 24 * time.Hour)}
	http.SetCookie(w, &refreshCookie)
	return nil
}

func (k *Service) deleteCookies(w http.ResponseWriter) {
	accessCookie := http.Cookie{Path: "/", Name: "access", Value: "", Expires: time.Unix(0, 0)}
	http.SetCookie(w, &accessCookie)

	refreshCookie := http.Cookie{Path: "/", Name: "refresh", Value: "", Expires: time.Unix(0, 0)}
	http.SetCookie(w, &refreshCookie)
}

func (k *Service) newUser(userInfo *gocloak.UserInfo) (*authentication.User, error) {
	if userInfo == nil {
		return nil, fmt.Errorf("user info undefined")
	}
	return &authentication.User{
		Sub:               safeString(userInfo.Sub),
		Email:             safeString(userInfo.Email),
		FamilyName:        safeString(userInfo.FamilyName),
		GivenName:         safeString(userInfo.GivenName),
		PreferredUsername: safeString(userInfo.PreferredUsername),
	}, nil
}

func safeString(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}
