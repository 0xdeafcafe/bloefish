package middlewares

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/matryer/is"
	"github.com/sirupsen/logrus"
)

func TestResponseWriter(t *testing.T) {
	t.Run("WriteHeader", func(t *testing.T) {
		is := is.New(t)

		w := httptest.NewRecorder()
		rw := &responseWriter{ResponseWriter: w}

		rw.WriteHeader(http.StatusTeapot)

		is.Equal(http.StatusTeapot, w.Code)
		is.Equal(http.StatusTeapot, rw.Status)
	})

	t.Run("Write", func(t *testing.T) {
		is := is.New(t)

		rec := httptest.NewRecorder()
		res := &responseWriter{ResponseWriter: rec}

		data := []byte("hello")

		n, err := res.Write(data)
		is.NoErr(err)

		is.Equal(n, 5)
		is.Equal(http.StatusOK, rec.Code)
		is.Equal(http.StatusOK, res.Status)
		is.Equal(data, rec.Body.Bytes())
	})
}

func TestLogger(t *testing.T) {
	tests := []struct {
		Name     string
		Status   int
		Contains string
	}{
		{"Error", http.StatusInternalServerError, "request"},
		{"Warning", http.StatusBadRequest, "request"},
		{"Success", http.StatusOK, "request"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctx := context.Background()
			is := is.New(t)

			log := logrus.New().WithField("foo", "bar")
			ctx = clog.Set(ctx, log)

			var buf bytes.Buffer
			log.Logger.Out = &buf

			data := []byte("hello")

			handlerInvoked := false
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				handlerInvoked = true
				w.WriteHeader(test.Status)
				_, err := w.Write(data)
				is.NoErr(err)
			})

			mw := Logger(log)
			is.True(mw != nil)

			fn := mw(next)
			is.True(fn != nil)

			rec := httptest.NewRecorder()
			r := (&http.Request{
				Method:     http.MethodGet,
				URL:        &url.URL{Path: "/"},
				Proto:      "HTTP/1.1",
				RemoteAddr: "127.0.0.1",
				Header: http.Header{
					"User-Agent": []string{"FooBar"},
					"Referer":    []string{"FooBar"},
				},
			}).WithContext(ctx)

			fn.ServeHTTP(rec, r)

			is.Equal(test.Status, rec.Code)
			is.Equal(data, rec.Body.Bytes())
			is.True(handlerInvoked)

			// TODO(jc): compare whole log entry, currently contains timestamp so
			// plain comparison will not work
			is.True(strings.Contains(buf.String(), test.Contains))
		})
	}
}
