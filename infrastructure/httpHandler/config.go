package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

	path := strings.TrimPrefix(r.URL.Path, "/config")

	switch r.Method {
	case http.MethodGet:
		configData := cp.ConfigService.GetConfigData(path)
		htmlTemplate.PresentHtmlTemplate(w, "resources/config.html", configData)

	case http.MethodPut:
		defer func() {
			_ = r.Body.Close()
		}()
		body, _ := ioutil.ReadAll(r.Body)

		configInput := ConfigInput{}
		err := json.Unmarshal(body, &configInput)
		if err != nil {
			cp.writeResponse(w, err)
		} else {
			config, err := createConfigFromInput(configInput)
			if err != nil {
				cp.writeResponse(w, err)
			} else {
				if len(configInput.OriginalPath) != 0 {
					_ = cp.ConfigService.Delete(configInput.OriginalPath)
				}

				cp.ConfigService.Save(config)
				cp.writeResponse(w, err)
			}
		}
	case http.MethodDelete:
		err := cp.ConfigService.Delete(path)
		cp.writeResponse(w, err)
	}
}

func (cp ConfigPage) writeResponse(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

func createConfigFromInput(input ConfigInput) (config.Config, error) {
	path := input.Path
	statusCode, err := strconv.Atoi(input.StatusCode)
	if err != nil {
		return config.Config{}, fmt.Errorf("statuscode is not a number")
	}

	return config.Config{
		Path:   path,
		Target: input.Target,
		Response: config.Response{
			StatusCode: statusCode,
			Body:       input.Body,
		},
	}, nil
}
