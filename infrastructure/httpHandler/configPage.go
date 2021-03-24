package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"fmt"
	"net/http"
	"strings"
)

type ConfigPage struct {
	KeyCloakService *keycloak.Service
	ConfigService   application.ConfigService
}

type ConfigInput struct {
	OriginalPath string `json:"originalPath"`
	Path         string `json:"path"`
	Target       string `json:"target"`
	StatusCode   string `json:"statusCode"`
	Body         string `json:"body"`
}

func (cp ConfigPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
		}
	}()

	_, err := cp.KeyCloakService.Authenticate(w, r)
	if err != nil {
		http.Redirect(w, r, "/login?origin="+r.RequestURI, http.StatusTemporaryRedirect)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/configs")

	switch r.Method {
	case http.MethodGet:
		configData := cp.ConfigService.GetConfigData(path)
		htmlTemplate.PresentHtmlTemplate(w, "resources/config.html", configData)
	}
}
