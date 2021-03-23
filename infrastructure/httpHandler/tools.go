package httpHandler

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func getQueryValueOrNil(v url.Values, key string) *string {
	value := v.Get(key)
	if len(value) == 0 {
		value = v.Get(strings.ToLower(key))
		if len(value) == 0 {
			return nil
		}
	}
	return &value
}

func writeResponse(w http.ResponseWriter, okResponse string, err error) {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(okResponse))
	}
}

func readRequestBody(r *http.Request) ([]byte, error) {
	defer func() {
		_ = r.Body.Close()
	}()
	return ioutil.ReadAll(r.Body)
}
