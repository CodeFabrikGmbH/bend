package dashboard

import "net/http"

type Presenter interface {
	Present(w http.ResponseWriter, board DashBoard)
}
