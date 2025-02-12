package cher

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/matryer/is"
	"github.com/pkg/errors"
)

func TestE(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		is := is.New(t)

		m := M{"foo": "bar"}
		e := New(NotFound, m, E{Code: "foo"})

		is.Equal(e, E{
			Code: NotFound,
			Meta: m,
			Reasons: []E{
				{Code: "foo"},
			},
		})
	})

	t.Run("Errorf", func(t *testing.T) {
		is := is.New(t)

		m := M{"foo": "bar"}
		e := Errorf(NotFound, m, "foo %s", "bar")

		is.Equal(e, E{
			Code: NotFound,
			Meta: M{
				"foo":     "bar",
				"message": "foo bar",
			},
		})
	})

	t.Run("StatusCode", func(t *testing.T) {
		tests := []struct {
			Name       string
			E          E
			StatusCode int
		}{
			{"BadRequest", E{Code: BadRequest}, http.StatusBadRequest},
			{"Unauthorized", E{Code: Unauthorized}, http.StatusUnauthorized},
			{"AccessDenied", E{Code: AccessDenied}, http.StatusForbidden},
			{"NotFound", E{Code: NotFound}, http.StatusNotFound},
			{"Unknown", E{Code: Unknown}, http.StatusInternalServerError},
			{"Handled", E{Code: "some_developer_code"}, http.StatusBadRequest},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				is := is.New(t)

				sc := test.E.StatusCode()
				is.Equal(test.StatusCode, sc)
			})
		}
	})

	t.Run("Error", func(t *testing.T) {
		is := is.New(t)

		e := E{Code: NotFound}
		is.Equal(NotFound, e.Error())
	})
}

func TestCoerce(t *testing.T) {
	tests := []struct {
		Name   string
		Src    any
		Result E
	}{
		{"E", E{Code: "foo"}, E{Code: "foo"}},
		{"String", "foo", E{Code: "foo"}},
		{"JSON", []byte(`{"code":"foo"}`), E{Code: "foo"}},
		{"BadJSON", []byte(`{"code":0}`), E{Code: CoercionError, Meta: M{"message": "json: cannot unmarshal number into Go struct field E.code of type string"}}},
		{"Error", errors.New("foo"), E{Code: Unknown, Meta: M{"message": "foo"}}}, //nolint:forbidigo // required for test
		{"Unknown", nil, E{Code: CoercionError}},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			is := is.New(t)

			e := Coerce(test.Src)

			is.Equal(test.Result, e)
		})
	}
}

//nolint:dupl // verbosity isn't an issue here
func TestWrapIfNotCher(t *testing.T) {
	type testCase struct {
		name   string
		msg    string
		err    error
		expect func(*is.I, error)
	}

	tests := []testCase{
		{
			name: "nil",
			msg:  "foo",
			err:  nil,
			expect: func(is *is.I, err error) {
				is.NoErr(err)
			},
		},
		{
			name: "err",
			msg:  "foo",
			err:  errors.New("nope"), //nolint:forbidigo // required for test
			expect: func(is *is.I, err error) {
				is.Equal(err.Error(), "foo: nope")
			},
		},
		{
			name: "cher",
			msg:  "foo",
			err:  New("nope", nil),
			expect: func(is *is.I, err error) {
				cErr, ok := err.(E) //nolint:errorlint // required for test
				is.True(ok)
				is.Equal(cErr.Code, "nope")
			},
		},
		{
			name: "cher unknown",
			msg:  "foo",
			err:  New("unknown", nil),
			expect: func(is *is.I, err error) {
				is.Equal(err.Error(), "foo: unknown")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			result := WrapIfNotCher(tc.err, tc.msg)
			tc.expect(is, result)
		})
	}
}

//nolint:dupl // verbosity isn't an issue here
func TestWrapIfNotCherCodes(t *testing.T) {
	type testCase struct {
		name   string
		msg    string
		err    error
		expect func(*is.I, error)
	}

	tests := []testCase{
		{
			name: "nil",
			msg:  "foo",
			err:  nil,
			expect: func(is *is.I, err error) {
				is.NoErr(err)
			},
		},
		{
			name: "err",
			msg:  "foo",
			err:  errors.New("nope"), //nolint:forbidigo // required for test
			expect: func(is *is.I, err error) {
				is.Equal(err.Error(), "foo: nope")
			},
		},
		{
			name: "cher specified code",
			msg:  "foo",
			err:  New("code_1", nil),
			expect: func(is *is.I, err error) {
				cErr, ok := err.(E) //nolint:errorlint // required for test
				is.True(ok)
				is.Equal(cErr.Code, "code_1")
			},
		},
		{
			name: "cher other code",
			msg:  "foo",
			err:  New("unknown", nil),
			expect: func(is *is.I, err error) {
				is.Equal(err.Error(), "foo: unknown")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			result := WrapIfNotCher(tc.err, tc.msg)
			tc.expect(is, result)
		})
	}
}

func TestAsCherWithCode(t *testing.T) {
	type testCase struct {
		name   string
		err    error
		codes  []string
		expect func(is *is.I, cErr E, ok bool)
	}

	tests := []testCase{
		{
			name:  "nil",
			err:   nil,
			codes: []string{"code_1"},
			expect: func(is *is.I, cErr E, ok bool) {
				is.Equal(ok, false)
			},
		},
		{
			name:  "normal error",
			err:   fmt.Errorf("nope"), //nolint:forbidigo,err113 // required for test
			codes: []string{"code_1"},
			expect: func(is *is.I, cErr E, ok bool) {
				is.Equal(ok, false)
			},
		},
		{
			name:  "normal error with same string",
			err:   fmt.Errorf("code_1"), //nolint:forbidigo,err113 // required for test
			codes: []string{"code_1"},
			expect: func(is *is.I, cErr E, ok bool) {
				is.Equal(ok, false)
			},
		},
		{
			name:  "cher specified code",
			err:   New("code_1", nil),
			codes: []string{"code_1"},
			expect: func(is *is.I, cErr E, ok bool) {
				is.True(ok)
				is.Equal(cErr.Code, "code_1")
			},
		},
		{
			name:  "cher other code",
			err:   New("unknown", nil),
			codes: []string{"code_1"},
			expect: func(is *is.I, cErr E, ok bool) {
				is.Equal(ok, false)
			},
		},
		{
			name:  "wrapped cher",
			err:   errors.Wrap(New("code_1", nil), "wrapped"),
			codes: []string{"code_1"},
			expect: func(is *is.I, cErr E, ok bool) {
				is.True(ok)
				is.Equal(cErr.Code, "code_1")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			cErr, ok := AsCherWithCode(tc.err, tc.codes...)
			tc.expect(is, cErr, ok)
		})
	}
}
