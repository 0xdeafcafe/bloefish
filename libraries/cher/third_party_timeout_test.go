package cher

import (
	"errors"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

type testError struct {
	error   string
	timeout bool
}

func (e testError) Error() string { return e.error }
func (e testError) Timeout() bool { return e.timeout }

func TestCoerceThirdPartyTimeout(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect func(is *is.I, err error)
	}{
		{
			name: "wrapped timeout error",
			err: fmt.Errorf("wrap: %w", testError{ //nolint:forbidigo // required for test
				error:   "foobar",
				timeout: true,
			}),
			expect: func(is *is.I, err error) {
				var cErr E
				is.True(errors.As(err, &cErr))
				is.Equal("third_party_timeout", cErr.Code)
				is.Equal("foobar", cErr.Meta["error"])
			},
		},
		{
			name: "other error",
			err:  errors.New("any error"), //nolint:forbidigo,goerr113 // required for test
			expect: func(is *is.I, err error) {
				is.Equal(err.Error(), "any error")
			},
		},
		{
			name: "no error",
			err:  nil,
			expect: func(is *is.I, err error) {
				is.NoErr(err)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			err := CoerceThirdPartyTimeout(tc.err)
			tc.expect(is, err)
		})
	}
}
