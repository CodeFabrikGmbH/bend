package httpHandler

import (
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
