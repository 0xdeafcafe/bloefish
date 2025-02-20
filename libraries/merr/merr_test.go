package merr

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/0xdeafcafe/bloefish/libraries/errfuncs"
	"github.com/0xdeafcafe/bloefish/libraries/stacktrace"
	"github.com/matryer/is"
)

func TestNew(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err := New(ctx, "foo", nil)

	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, nil)
	is.True(strings.HasSuffix(err.Stack[0].File, "/libraries/merr/merr_test.go"))
	is.Equal(err.Reasons, nil)
}

func TestNewMeta(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err := New(ctx, "foo", M{"a": "b"})

	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, M{"a": "b"})
	is.True(strings.HasSuffix(err.Stack[0].File, "/libraries/merr/merr_test.go"))
	is.Equal(err.Reasons, nil)
}

func TestWrap(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err1 := errors.New("underlying error") //nolint:goerr113,forbidigo // needed for testing
	err2 := New(ctx, "foo", nil, err1)

	err, ok := errfuncs.As[E](err2)

	is.True(ok)
	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, nil)
	is.True(strings.HasSuffix(err.Stack[0].File, "/libraries/merr/merr_test.go"))
	is.Equal(len(err.Reasons), 1)
	is.Equal(err.Reasons[0], errors.New("underlying error")) //nolint:goerr113,forbidigo // needed for testing
}

func TestWrapMeta(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err1 := errors.New("underlying error") //nolint:goerr113,forbidigo // needed for testing
	err2 := New(ctx, "foo", M{"a": "b"}, err1)

	err, ok := errfuncs.As[E](err2)

	is.True(ok)
	is.Equal(err.Code, Code("foo"))
	is.Equal(err.Meta, M{"a": "b"})
	is.True(strings.HasSuffix(err.Stack[0].File, "/libraries/merr/merr_test.go"))
	is.Equal(len(err.Reasons), 1)
	is.Equal(err.Reasons[0], errors.New("underlying error")) //nolint:goerr113,forbidigo // needed for testing
}

func TestEqual(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	//nolint:lll // same line to ensure same stack trace 😅
	err1, err2, err3, err4, err5, err6 := New(ctx, "foo", M{"a": "b"}), New(ctx, "foo", M{"a": "b"}), New(ctx, "foo", M{"a": "c"}), New(ctx, "bar", M{"a": "b"}), New(ctx, "foo", nil), New(ctx, "foo", M{"a": "b"})
	err6.Reasons = []error{errors.New("foo")} //nolint:goerr113,forbidigo // needed for testing

	is.True(err1.Equal(err2))
	is.True(!err1.Equal(err3))
	is.True(!err1.Equal(err4))
	is.True(!err1.Equal(err5))
	is.True(!err1.Equal(err6))
}

func TestString(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err := New(ctx, "foo", M{"a": "b"})

	err.Stack = []stacktrace.Frame{
		{
			File:     "/libraries/foo/foo.go",
			Line:     123,
			Function: "github.com/0xdeafcafe/bloefish/libraries/foo.doFoo",
		},
		{
			File:     "/libraries/foo/bar.go",
			Line:     456,
			Function: "github.com/0xdeafcafe/bloefish/libraries/foo.barThing",
		},
	}

	expected := `foo (map[a:b])

github.com/0xdeafcafe/bloefish/libraries/foo.doFoo
	/libraries/foo/foo.go:123
github.com/0xdeafcafe/bloefish/libraries/foo.barThing
	/libraries/foo/bar.go:456
`

	is.Equal(err.String(), expected)
}

func TestEError(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	is.Equal(New(ctx, "foo", nil).Error(), "foo")
	is.Equal(New(ctx, "foo", M{"a": "b"}).Error(), "foo (map[a:b])")
	is.Equal(New(ctx, "bar", nil, errors.New("foo")).Error(), "bar\n- foo") //nolint:goerr113,forbidigo // needed for testing
	is.Equal(New(ctx, "bar", nil, New(ctx, "foo", nil)).Error(), "bar\n- foo")
	is.Equal(New(ctx, "bar", nil, New(ctx, "foo", M{"a": "b"})).Error(), "bar\n- foo (map[a:b])")
	is.Equal(New(ctx, "bar", M{"c": "d"}, New(ctx, "foo", nil)).Error(), "bar (map[c:d])\n- foo")
	is.Equal(New(ctx, "bar", M{"c": "d"}, New(ctx, "foo", M{"a": "b"})).Error(), "bar (map[c:d])\n- foo (map[a:b])")
}

func TestEIsCode(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	errs := []error{
		New(ctx, "foo", nil),
		New(ctx, Code("foo"), nil),
		New(ctx, "foo", M{"a": "b"}),
		New(ctx, Code("foo"), M{"a": "b"}),
		New(ctx, "bar", nil, New(ctx, "foo", nil)),
		New(ctx, "bar", nil, wrappedError{New(ctx, "foo", nil)}),
	}

	for _, err := range errs {
		is.True(errors.Is(err, Code("foo")))
	}
}

func TestEIsE(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	errFoo := New(ctx, "foo", nil)

	errs := []error{
		errFoo,
		New(ctx, "bar", nil, errFoo),
		New(ctx, "bar", nil, wrappedError{errFoo}),
	}

	for _, err := range errs {
		is.True(errors.Is(err, errFoo))
	}
}

func TestEUnwrap(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err1 := errors.New("underlying error") //nolint:goerr113,forbidigo // needed for testing
	err2 := New(ctx, "foo", nil, err1)
	err3 := New(ctx, "bar", nil, err2)

	// Unwrap returns nil for multi-error wrapping
	is.Equal(errors.Unwrap(err3), nil)
	is.Equal(errors.Unwrap(err2), nil)

	is.Equal(err3.Unwrap(), []error{err2})
	is.Equal(err2.Unwrap(), []error{err1})
}

func TestEAs(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	err := New(ctx, "foo", nil)

	var errFoo E
	is.True(errors.As(err, &errFoo))

	_, ok := errfuncs.As[E](err)
	is.True(ok)
}
