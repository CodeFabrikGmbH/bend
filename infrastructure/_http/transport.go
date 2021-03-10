package _http

import (
	"bytes"
	"code-fabrik.com/bend/domain/request"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

type Transport struct {
	client *http.Client
}

func (t Transport) SendRequestToTarget(rr request.Request, targetUrl string) request.Response {
	if t.client == nil {
		t.client = createHttpClient()
	}

	result := request.Response{}

	req, err := http.NewRequest(rr.Method, targetUrl, bytes.NewBuffer([]byte(rr.Body)))

	if err != nil {
		result.Error = err.Error()
		return result
	}

	for k, vl := range rr.Header {
		for _, v := range vl {
			req.Header.Add(k, v)
		}
	}

	response, err := t.client.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	defer func() {
		_ = response.Body.Close()
	}()

	responseBody, err := ioutil.ReadAll(response.Body)

	return request.Response{
		ResponseBody:       string(responseBody),
		ResponseStatusCode: response.StatusCode,
	}
}

func createHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 2 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
}
