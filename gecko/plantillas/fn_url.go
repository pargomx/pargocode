package plantillas

import (
	"fmt"
	"net/url"
)

// addQueryParam agrega un query param a la URL.
func addQueryParam(URL *url.URL, key, value string) string {
	if URL == nil {
		return "/error?msg=url-null"
	}
	q := URL.Query()
	q.Set(key, value)
	URL.RawQuery = q.Encode()
	return URL.String()
}

// addQueryNum agrega un query param a la URL.
func addQueryNum(URL *url.URL, key string, value interface{}) string {
	if URL == nil {
		return "/error?msg=url-null"
	}
	q := URL.Query()
	q.Set(key, fmt.Sprint(value))
	URL.RawQuery = q.Encode()
	return URL.String()
}
