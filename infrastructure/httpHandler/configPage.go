package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type ConfigPage struct {
	KeyCloakService *keycloak.Service
	ConfigService   application.ConfigService
}

type ConfigInput struct {
	Path       string    `json:"path"`
	Target     string    `json:"target"`
	StatusCode string    `json:"statusCode"`
	Body       string    `json:"body"`
	Id         uuid.UUID `json:"id"`
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

	id, _ := ParseSubPathAsUUID(r.URL.Path, "/configs/")

	switch r.Method {
	case http.MethodGet:
		configData := cp.ConfigService.GetConfigData(id)
		htmlTemplate.PresentHtmlTemplate(w, "resources/config.html", configData)
	}
}

func ParseSubPathAsUUID(path string, prefix string) (uuid.UUID, error) {
	idAsString := strings.TrimPrefix(path, prefix)
	return uuid.Parse(idAsString)
}
