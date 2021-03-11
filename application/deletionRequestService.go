package application

import (
	"code-fabrik.com/bend/domain/environment"
	"fmt"
	"net/http"
	"strings"
)

type DeletionService struct {
	Env environment.Environment
}

func (rs DeletionService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	path := strings.TrimPrefix(r.URL.Path, "/delete")

	requestId := getQueryValueOrNil(r.URL.Query(), "requestId")
	if requestId == nil {
		rs.Env.RequestRepository.DeletePath(path)
	} else {
		rs.Env.RequestRepository.DeleteRequestForPath(path, *requestId)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
