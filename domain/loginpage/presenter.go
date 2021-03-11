package loginpage

import "net/http"

type Page interface {
	Present(w http.ResponseWriter, login Login)
}
