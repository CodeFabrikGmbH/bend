package httpHandler

import (
	"code-fabrik.com/bend/application"
	"code-fabrik.com/bend/domain/request"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type TrackRequest struct {
	RequestService application.RequestService
}

func (rs TrackRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	defer func() {
		_ = r.Body.Close()
	}()
	body, _ := ioutil.ReadAll(r.Body)

	req := request.Request{
		Timestamp: time.Now().UnixNano(),
		Path:      r.URL.Path,
		Method:    r.Method,
		Body:      string(body),
		Header:    r.Header,
		Host:      r.Host,
		Uri:       r.RequestURI,
	}

	if req.Path == "/" {
		writeResponse(w, "", fmt.Errorf("root is not server"))
	} else {
		response := rs.RequestService.TrackRequest(req)
		addHeaders(w, response.ResponseHeader)
		w.WriteHeader(response.ResponseStatusCode)
		_, _ = w.Write([]byte(response.ResponseBody))
	}
}

func addHeaders(w http.ResponseWriter, header map[string][]string) {
	for key, values := range header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
}
