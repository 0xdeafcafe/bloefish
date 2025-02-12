package version

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestHeader(t *testing.T) {
	tests := []struct {
		Name          string
		App, Revision string
		Expected      string
	}{
		{"Full", "Test", "dev", "Test/dev"},
		{"Truncated", "Test", "e51d8e3c9eca72a41f205a11a2698373bbebb447", "Test/e51d8e3"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			is := is.New(t)

			handlerInvoked := false
			res := httptest.NewRecorder()
			req := &http.Request{}
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { handlerInvoked = true })

			Revision = test.Revision
			Truncated = genTruncated()
			mw := HeaderMiddleware(test.App)
			is.True(mw != nil)

			hn := mw(next)
			is.True(hn != nil)

			hn.ServeHTTP(res, req)

			is.Equal(test.Expected, res.Header().Get("Server"))
			is.True(handlerInvoked)
		})
	}
}
