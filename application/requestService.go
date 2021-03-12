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
	req.Response = rs.getOrCreateResponse(req)

	err := rs.Env.RequestRepository.Add(req)
	if err != nil {
		panic(err)
	}

	writeResponse(w, req.Response)
}

func writeResponse(w http.ResponseWriter, response request.Response) {
	w.WriteHeader(response.ResponseStatusCode)
	w.Write([]byte(response.ResponseBody))
}

func (rs RequestService) getOrCreateResponse(req request.Request) request.Response {
	config := rs.Env.ConfigRepository.Find(req.Path)

	if config == nil {
		return request.Response{
			Target:             "no target - mocked response",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       "ok",
		}
	}
	if len(config.Target) != 0 {
		return rs.Env.Transport.SendRequestToTarget(req, config.Target)
	}

	return request.Response{
		Target:             "no target - mocked response",
		ResponseStatusCode: config.Response.StatusCode,
		ResponseBody:       config.Response.Body,
	}
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
