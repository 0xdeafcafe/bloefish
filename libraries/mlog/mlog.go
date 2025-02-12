package mlog

import (
	"context"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/merr"
	"github.com/0xdeafcafe/bloefish/libraries/mlog/indirect"
	"github.com/0xdeafcafe/bloefish/libraries/stacktrace"
	"github.com/sirupsen/logrus"
)

//nolint:gochecknoinits // currently no way around this
func init() {
	indirect.Debug = Debug
	indirect.Info = Info
	indirect.Warn = Warn
	indirect.Error = Error //nolint:reassign // this is intended, works around import cycle
}

// Debug is to help with tracing the system's behavior
//
// e.g. logic has been evaluated to determine some action
func Debug(ctx context.Context, err merr.Merrer) {
	log(ctx, logrus.DebugLevel, err)
}

// Info is for informational messages which are not errors
//
// e.g. a system is starting up
func Info(ctx context.Context, err merr.Merrer) {
	log(ctx, logrus.InfoLevel, err)
}

// Warn covers issues which were handled and do not require specific action
// from an engineer, but which should be fixed at some point
//
// Only use this for issues that can safely keep occurring indefinitely without
// serious consequences
//
// e.g. a system was unavailable, but failed gracefully
func Warn(ctx context.Context, err merr.Merrer) {
	log(ctx, logrus.WarnLevel, err)
}

// Error represents an issue which requires prompt and individual action from
// an engineer to resolve
//
// e.g. a data integrity issue has been identified which needs to be fixed
func Error(ctx context.Context, err merr.Merrer) {
	log(ctx, logrus.ErrorLevel, err)
}

func log(ctx context.Context, level logrus.Level, err merr.Merrer) {
	if err == nil {
		return
	}

	logger := clog.Get(ctx)
	merr := err.Merr()

	// logrus runs `.String()` on anything implementing `error`
	// so to get proper JSON, we need to copy the merrFields instead
	merrFields := merr.Fields()

	fields := logrus.Fields{
		"error": merr.Error(),
		"merr":  merrFields,
	}

	// if we're doing unstructured/text logging, try to improve readability
	if _, ok := logger.Logger.Formatter.(*logrus.TextFormatter); ok {
		if level <= logrus.InfoLevel {
			delete(merrFields, "stack")
		}
		if merrFields["meta"] == nil {
			delete(merrFields, "meta")
		}
		if merrFields["reason"] == nil {
			delete(merrFields, "reason")
		}

		fields = merrFields
	}

	if merr.Stack != nil && level > logrus.InfoLevel {
		fields["stack_trace"] = stacktrace.FormatFrames(merr.Stack)
	}

	logger.
		WithFields(fields).
		Log(level, merr.Code)
}
