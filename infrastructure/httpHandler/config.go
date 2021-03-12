package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ConfigPage struct {
	KeyCloakService *keycloak.Service
	ConfigService   application.ConfigService
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

	path := strings.TrimPrefix(r.URL.Path, "/config")

	switch r.Method {
	case "PUT":
		defer func() {
			_ = r.Body.Close()
		}()
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println(body)

	case "DELETE":
		err := cp.ConfigService.Delete(path)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}
	default:
		configData := cp.ConfigService.GetConfigData(path)
		htmlTemplate.PresentHtmlTemplate(w, "resources/config.html", configData)
	}
}
