package httpHandler

import (
	"code-fabrik.com/bend/application"
	"fmt"
	"net/http"
	"strings"
)

type Deletion struct {
	RequestService application.RequestService
}

func (rs Deletion) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	path := strings.TrimPrefix(r.URL.Path, "/delete")
	requestId := getQueryValueOrNil(r.URL.Query(), "requestId")

	err := rs.RequestService.Delete(path, requestId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
