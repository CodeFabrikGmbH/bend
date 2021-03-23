package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"fmt"
	"net/http"
	"strings"
)

type DashboardPage struct {
	KeyCloakService  *keycloak.Service
	DashboardService application.DashboardService
}

func (dp DashboardPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
		}
	}()

	_, err := dp.KeyCloakService.Authenticate(w, r)
	if err != nil {
		http.Redirect(w, r, "/login?origin="+r.RequestURI, http.StatusTemporaryRedirect)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/dashboard")
	requestId := getQueryValueOrNil(r.URL.Query(), "requestId")

	dashBoardData := dp.DashboardService.GenerateDashboardViewData(path, requestId)

	htmlTemplate.PresentHtmlTemplate(w, "resources/dashboard.html", dashBoardData)
}
