package middlewares

import (
	"net/http"

	"github.com/0xdeafcafe/bloefish/libraries/crpc"
)

func UnsafeNoAuthentication(next crpc.HandlerFunc) crpc.HandlerFunc {
	return func(res http.ResponseWriter, req *crpc.Request) error {
		return next(res, req)
	}
}
