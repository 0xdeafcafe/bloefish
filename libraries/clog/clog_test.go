package clog

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/matryer/is"
	"github.com/sirupsen/logrus"
)

func TestContextLogger(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		is := is.New(t)

		log := logrus.New().WithField("foo", "bar")

		r := &http.Request{}
		r = r.WithContext(Set(r.Context(), log))

		l := Get(r.Context())

		is.Equal(log.Data, l.Data)
	})

	t.Run("SetFields", func(t *testing.T) {
		is := is.New(t)

		log := logrus.New().WithField("foo", "bar")

		req := &http.Request{}
		req = req.WithContext(Set(req.Context(), log))

		SetFields(req.Context(), Fields{
			"foo2": "bar2",
		})

		cl := getContextLogger(req.Context()).GetLogger()
		is.Equal("bar", cl.Data["foo"])
		is.Equal("bar2", cl.Data["foo2"])
	})

	t.Run("SetField", func(t *testing.T) {
		is := is.New(t)

		log := logrus.New().WithField("foo", "bar")

		req := &http.Request{}
		req = req.WithContext(Set(req.Context(), log))

		SetField(req.Context(), "foo2", "bar2")

		cl := getContextLogger(req.Context()).GetLogger()
		is.Equal("bar", cl.Data["foo"])
		is.Equal("bar2", cl.Data["foo2"])
	})

	t.Run("SetError", func(t *testing.T) {
		is := is.New(t)

		log := logrus.New().WithField("foo", "bar")

		req := &http.Request{}
		req = req.WithContext(Set(req.Context(), log))

		testError := errors.New("test error") //nolint:forbidigo,goerr113 // required for test

		SetError(req.Context(), testError)

		cl := getContextLogger(req.Context()).GetLogger()
		is.Equal(testError, cl.Data["error"])
	})

	t.Run("Logger when no clog is set", func(t *testing.T) {
		is := is.New(t)

		r := &http.Request{}
		l := Get(r.Context())

		is.True(l != nil)
	})

	t.Run("SetField when no logger is set", func(t *testing.T) {
		defer func() { _ = recover() }()

		SetField(context.Background(), "foo", "bar")

		t.Error("should have panicked")
	})

	t.Run("SetFields when no logger is set", func(t *testing.T) {
		defer func() { _ = recover() }()

		SetFields(context.Background(), Fields{"foo": "bar"})

		t.Error("should have panicked")
	})

	t.Run("SetError when no logger is set", func(t *testing.T) {
		defer func() { _ = recover() }()

		SetError(context.Background(), errors.New("foo")) //nolint:forbidigo,goerr113 // required for test

		t.Error("should have panicked")
	})
}

func TestDetermineLevel(t *testing.T) {
	type testCase struct {
		name             string
		err              error
		timeoutsAsErrors bool
		expected         logrus.Level
	}

	tests := []testCase{
		{
			name:             "bad request",
			err:              cher.New("bad_request", nil),
			timeoutsAsErrors: false,
			expected:         logrus.WarnLevel,
		},
		{
			name:             "context cancelled",
			err:              cher.New(cher.ContextCanceled, nil),
			timeoutsAsErrors: false,
			expected:         logrus.InfoLevel,
		},
		{
			name:             "context cancelled with timeouts as errors",
			err:              cher.New(cher.ContextCanceled, nil),
			timeoutsAsErrors: true,
			expected:         logrus.ErrorLevel,
		},
		{
			name:             "unknown",
			err:              cher.New(cher.Unknown, nil),
			timeoutsAsErrors: false,
			expected:         logrus.ErrorLevel,
		},
		{
			name:             "postgres context cancelled",
			err:              errors.New("pq: canceling statement due to user request"), //nolint:forbidigo,goerr113 // required for test
			timeoutsAsErrors: false,
			expected:         logrus.InfoLevel,
		},
		{
			name:             "postgres context cancelled with timeouts as errors",
			err:              errors.New("pq: canceling statement due to user request"), //nolint:forbidigo,goerr113 // required for test
			timeoutsAsErrors: true,
			expected:         logrus.ErrorLevel,
		},
		{
			name:             "other error",
			err:              errors.New("something, something darkside"), //nolint:forbidigo,goerr113 // required for test
			timeoutsAsErrors: false,
			expected:         logrus.ErrorLevel,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			result := DetermineLevel(tc.err, tc.timeoutsAsErrors)
			is.Equal(tc.expected, result)
		})
	}
}
