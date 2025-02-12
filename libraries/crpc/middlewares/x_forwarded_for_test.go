package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestInfraClientIP(t *testing.T) {
	tests := []struct {
		Name               string
		Header             string
		ExpectedRemoteAddr string
	}{
		{"NotSet", "", "1.2.3.4"},
		{"Whitespace", " ", "1.2.3.4"},
		{"Set", "8.8.4.4", "8.8.4.4"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			is := is.New(t)

			handlerInvoked := false
			rec := httptest.NewRecorder()
			r := &http.Request{Header: http.Header{"Infra-Client-Ip": []string{test.Header}}, RemoteAddr: "1.2.3.4"}
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

			hn := ClientIP(next)
			is.True(hn != nil)

			hn.ServeHTTP(rec, r)

			is.Equal(test.ExpectedRemoteAddr, r.RemoteAddr)
			is.True(handlerInvoked)
		})
	}
}

func TestXForwardedFor(t *testing.T) {
	tests := []struct {
		Name               string
		Header             string
		ExpectedRemoteAddr string
	}{
		{"NotSet", "", "1.2.3.4"},
		{"Whitespace", " ", "1.2.3.4"},
		{"Set", "8.8.4.4", "8.8.4.4"},
		{"Multiple set", "8.8.4.4, 8.8.8.8", "8.8.4.4"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			is := is.New(t)

			handlerInvoked := false
			rec := httptest.NewRecorder()
			r := &http.Request{Header: http.Header{"X-Forwarded-For": []string{test.Header}}, RemoteAddr: "1.2.3.4"}
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

			hn := ClientIP(next)
			is.True(hn != nil)

			hn.ServeHTTP(rec, r)

			is.Equal(test.ExpectedRemoteAddr, r.RemoteAddr)
			is.True(handlerInvoked)
		})
	}
}
