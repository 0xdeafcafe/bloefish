package middlewares

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/matryer/is"
)

func TestStripPrefix(t *testing.T) {
	tests := []struct {
		Name   string
		Path   string
		Prefix string
		Result string
	}{
		{"Empty", "/", "", "/"},
		{"HasPrefix", "/foo/bar", "/foo", "/bar"},
		{"NoPrefix", "/foo/bar", "/baz", "/foo/bar"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			is := is.New(t)

			var invoked bool

			fn := StripPrefix(test.Prefix)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				invoked = true

				is.Equal(test.Result, r.URL.Path)
			}))

			fn.ServeHTTP(nil, &http.Request{
				URL: &url.URL{
					Path: test.Path,
				},
			})

			is.True(invoked)
		})
	}
}
