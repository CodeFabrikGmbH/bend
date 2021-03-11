package application

import (
	"code-fabrik.com/bend/domain/environment"
	"code-fabrik.com/bend/domain/request"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestService struct {
	Env environment.Environment
}

func (rs RequestService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println(rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	req := createRequestObject(r)
	err := rs.Env.RequestRepository.Save(req)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func createRequestObject(r *http.Request) request.Request {
	body, _ := readBody(r.Body)

	return request.Request{
		Timestamp: time.Now().UnixNano(),
		Path:      r.URL.Path,
		Method:    r.Method,
		Body:      string(body),
		Header:    r.Header,
		Host:      r.Host,
		Uri:       r.RequestURI,
	}
}

func readBody(ioBody io.ReadCloser) ([]byte, error) {
	defer func() {
		_ = ioBody.Close()
	}()
	return ioutil.ReadAll(ioBody)
}
