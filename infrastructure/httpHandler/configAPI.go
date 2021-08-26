package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/config"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ConfigAPI struct {
	KeyCloakService *keycloak.Service
	ConfigService   application.ConfigService
}

func (cp ConfigAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	id, _ := ParseSubPathAsUUID(r.URL.Path, "/api/configs/")

	switch r.Method {
	case http.MethodPut:
		defer func() {
			_ = r.Body.Close()
		}()
		body, _ := ioutil.ReadAll(r.Body)

		configInput := ConfigInput{}
		err := json.Unmarshal(body, &configInput)
		if err != nil {
			writeResponse(w, "", err)
		} else {
			configFromInput, err := createConfigFromInput(configInput)
			if err != nil {
				writeResponse(w, "", err)
			} else {
				_ = cp.ConfigService.Save(configFromInput)
				writeResponse(w, "ok", err)
			}
		}
	case http.MethodDelete:
		err := cp.ConfigService.Delete(id)
		writeResponse(w, "ok", err)
	}
}

func createConfigFromInput(input ConfigInput) (config.Config, error) {
	path := input.Path
	statusCode, err := strconv.Atoi(input.StatusCode)
	if err != nil {
		return config.Config{}, fmt.Errorf("statuscode is not a number")
	}

	id := uuid.New()
	if input.Id != uuid.Nil {
		id = input.Id
	}

	if strings.Index(path, "/") != 0 {
		path = "/" + path
	}

	return config.Config{
		Path:   path,
		Target: input.Target,
		Response: config.Response{
			StatusCode: statusCode,
			Body:       input.Body,
		},
		Id: id,
	}, nil
}
