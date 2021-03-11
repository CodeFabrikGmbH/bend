package dashboardpage

import "net/http"

type Page interface {
	Present(w http.ResponseWriter, board DashBoard)
}
