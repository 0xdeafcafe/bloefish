package merr

import (
	"context"
	"errors"
	"testing"

	"github.com/matryer/is"
)

type wrappedError struct {
	reason error
}

func (e wrappedError) Error() string {
	return "wrappedError"
}

func (e wrappedError) Unwrap() error {
	return e.reason
}

func TestCodeError(t *testing.T) {
	is := is.New(t)

	err := Code("foo")

	is.Equal(err.Error(), "foo")
}

func TestCodeComparable(t *testing.T) {
	is := is.New(t)

	err1 := Code("foo")
	err2 := Code("foo")

	is.True(err1 == err2)
}

func TestIsCode(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	errs := []error{
		Code("foo"),
		New(ctx, "foo", nil),
		New(ctx, Code("foo"), nil),
		New(ctx, "foo", M{"a": "b"}),
		New(ctx, Code("foo"), M{"a": "b"}),
		New(ctx, "bar", nil, New(ctx, "foo", nil)),
		New(ctx, "bar", nil, wrappedError{New(ctx, "foo", nil)}),
	}

	for _, err := range errs {
		is.True(IsCode(err, "foo"))
		is.True(IsCode(err, Code("foo")))
		is.True(errors.Is(err, Code("foo")))
	}
}
