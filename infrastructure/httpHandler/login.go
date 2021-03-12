package httpHandler

import (
	"code-fabrik.com/bend/domain/authentication"
	"code-fabrik.com/bend/domain/loginpage"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"context"
	"fmt"
	"net/http"
)

type LoginPage struct {
	KeyCloakService *keycloak.Service
}

func (lp LoginPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
		}
	}()
	ctx := context.Background()

	if lp.isLogoutRequest(r) {
		lp.logout(ctx, w)
	} else if lp.isLoginRequest(r) {
		lp.login(ctx, w, r)
	} else {
		query := r.URL.Query()
		origin := query.Get("origin")

		user, _ := lp.KeyCloakService.Authenticate(w, r)
		htmlTemplate.PresentLoginPage(w, createLoginData(user, origin, ""))
	}
}

func (lp LoginPage) login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	username := query.Get("username")
	password := query.Get("password")
	origin := query.Get("origin")

	user, err := lp.KeyCloakService.Login(ctx, w, username, password)
	errorString := ""
	if err != nil {
		errorString = "bad credentials"
	}

	htmlTemplate.PresentLoginPage(w, createLoginData(user, origin, errorString))
}

func (lp LoginPage) logout(ctx context.Context, w http.ResponseWriter) {
	lp.KeyCloakService.Logout(ctx, w)
	htmlTemplate.PresentLoginPage(w, createLoginData(nil, "", ""))
}

func (lp LoginPage) isLogoutRequest(r *http.Request) bool {
	query := r.URL.Query()
	logout := query.Get("logout")
	return len(logout) != 0
}

func (lp LoginPage) isLoginRequest(r *http.Request) bool {
	query := r.URL.Query()
	username := query.Get("username")
	password := query.Get("password")

	return len(username) != 0 && len(password) != 0
}

func createLoginData(userInfo *authentication.User, origin string, err string) loginpage.Login {
	name := ""
	if userInfo != nil {
		name = userInfo.GivenName + " " + userInfo.FamilyName
	}
	return loginpage.Login{
		Name:   name,
		Origin: origin,
		Error:  err,
	}
}
