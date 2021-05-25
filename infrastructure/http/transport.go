package http

import (
	"bytes"
	"code-fabrik.com/bend/domain/request"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Transport struct {
	client *http.Client
}

func addRequestParametersToTargetUrl(requestUri, targetUri string) (string, error) {
	requestURL, err := url.Parse(requestUri)
	if err != nil {
		return "", err
	}
	targetURL, err := url.Parse(targetUri)
	if err != nil {
		return "", err
	}

	q := targetURL.Query()

	for k, vl := range requestURL.Query() {
		if q[k] == nil {
			//only add if target url does not already specify the parameter key
			for _, v := range vl {
				q.Add(k, v)
			}
		}
	}
	targetURL.RawQuery = q.Encode()
	return targetURL.String(), nil

}

func (t Transport) SendRequestToTarget(rr request.Request, targetUrl string) request.Response {
	if t.client == nil {
		t.client = createHttpClient()
	}

	finalUrl, err := addRequestParametersToTargetUrl(rr.Uri, targetUrl)

	result := request.Response{
		Target:             finalUrl,
		ResponseStatusCode: -1,
	}

	req, err := http.NewRequest(rr.Method, finalUrl, bytes.NewBuffer([]byte(rr.Body)))

	if err != nil {
		result.Error = err.Error()
		return result
	}

	for k, vl := range rr.Header {
		for _, v := range vl {
			// If accept encoding is not explicitly set,
			// the compression will be handled handled automatically by transport
			if k == "Accept-Encoding" {
				continue
			}
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
