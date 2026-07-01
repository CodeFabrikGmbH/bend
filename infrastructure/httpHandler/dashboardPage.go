package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/infrastructure/htmlTemplate"
	"code-fabrik.com/bend/infrastructure/jwt/keycloak"
	"log/slog"
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
			slog.Error("panic in DashboardPage", "recover", rec)
		}
	}()

	_, err := dp.KeyCloakService.Authenticate(w, r)
	if err != nil {
		http.Redirect(w, r, "/login?origin="+r.RequestURI, http.StatusTemporaryRedirect)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/dashboard")

	requestPath, requestId := getRequestPathAndId(path)

	dashBoardData := dp.DashboardService.GenerateDashboardViewData(requestPath, requestId)

	htmlTemplate.PresentHtmlTemplate(w, "dashboard.html", dashBoardData)
}
