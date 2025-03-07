package errfuncs

import (
	"errors"
)

func As[T error](err error) (target T, ok bool) {
	ok = errors.As(err, &target)
	return
}
