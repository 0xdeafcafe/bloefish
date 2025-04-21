package otelopenai

import (
	"net/http"
	"sync"

	"github.com/felixge/httpsnoop"
)

var rrwPool = &sync.Pool{
	New: func() interface{} {
		return &recordingResponseWriter{}
	},
}

type recordingResponseWriter struct {
	writer  http.ResponseWriter
	written bool
	status  int
}

func getRRW(writer http.ResponseWriter) *recordingResponseWriter {
	rrw := rrwPool.Get().(*recordingResponseWriter)
	rrw.written = false
	rrw.status = http.StatusOK
	rrw.writer = httpsnoop.Wrap(writer, httpsnoop.Hooks{
		Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(b []byte) (int, error) {
				if !rrw.written {
					rrw.written = true
				}
				return next(b)
			}
		},
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(statusCode int) {
				if !rrw.written {
					rrw.written = true
					rrw.status = statusCode
					// only call next WriteHeader when header is not written yet
					// this is to prevent superfluous WriteHeader call
					next(statusCode)
				}
			}
		},
	})
	return rrw
}

func putRRW(rrw *recordingResponseWriter) {
	rrw.writer = nil
	rrwPool.Put(rrw)
}
