package httpHandler

import (
	"code-fabrik.com/bend/domain/authentication"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"context"
	"log/slog"
	"net/http"
)

type LoginViewData struct {
	Name   string
	Origin string
	Error  string
}

type LoginPage struct {
	KeyCloakService *keycloak.Service
}

func (lp LoginPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("panic in LoginPage", "recover", rec)
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

		htmlTemplate.PresentHtmlTemplate(w, "login.html", createLoginViewData(user, origin, ""))
	}
}

func (lp LoginPage) login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	origin := r.PostFormValue("origin")

	user, err := lp.KeyCloakService.Login(ctx, w, username, password)
	errorString := ""
	if err != nil {
		errorString = "bad credentials"
	}

	htmlTemplate.PresentHtmlTemplate(w, "login.html", createLoginViewData(user, origin, errorString))
}

func (lp LoginPage) logout(ctx context.Context, w http.ResponseWriter) {
	lp.KeyCloakService.Logout(ctx, w)
	htmlTemplate.PresentHtmlTemplate(w, "login.html", createLoginViewData(nil, "", ""))
}

func (lp LoginPage) isLogoutRequest(r *http.Request) bool {
	return r.Method == http.MethodPost && len(r.PostFormValue("logout")) != 0
}

func (lp LoginPage) isLoginRequest(r *http.Request) bool {
	if r.Method != http.MethodPost {
		return false
	}
	return len(r.PostFormValue("username")) != 0 && len(r.PostFormValue("password")) != 0
}

func createLoginViewData(userInfo *authentication.User, origin string, err string) LoginViewData {
	name := ""
	if userInfo != nil {
		name = userInfo.GivenName + " " + userInfo.FamilyName
	}
	return LoginViewData{
		Name:   name,
		Origin: origin,
		Error:  err,
	}
}
