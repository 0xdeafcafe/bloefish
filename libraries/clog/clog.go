package clog

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/otel"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/contexts"
	"github.com/0xdeafcafe/bloefish/libraries/errfuncs"
	"github.com/0xdeafcafe/bloefish/libraries/version"
)

type contextKey string

type Format string

type Fields map[string]any

// LoggerKey is the key used for request-scoped loggers in a requests context map
const loggerKey contextKey = "clog"

const (
	// TextFormat is the logrus text format
	TextFormat Format = "text"

	// JSONFormat is the logrus json format
	JSONFormat Format = "json"

	// ServiceKey is the log entry key for the name of the crpc service
	ServiceKey = "_service"

	// HostKey is the log entry key for the hostname / container ID
	HostKey = "_hostname"

	// VersionKey is the log entry key for the current version of the codebase
	VersionKey = "_commit_hash"

	// LevelKey is the log entry key for the log level
	LevelKey = "_level"

	// MessageKey is the log entry key for a generic message
	MessageKey = "_message"

	// TimestampKey is the log entry key for the timestamp
	TimestampKey = "_timestamp"
)

// Config allows services to configure the logging format, level and storage options
// for Logrus logging.
type Config struct {
	// Format configures the output format. Possible options:
	//   - text - logrus default text output, good for local development
	//   - json - fields and message encoded as json, good for storage in e.g. cloudwatch
	Format Format `env:"format"`

	// Debug enables debug level logging, otherwise INFO level
	Debug bool `env:"debug"`
}

// Configure applies standard Logging structure options to a logrus Entry.
func (c Config) Configure(ctx context.Context) *logrus.Entry {
	var serviceName string
	if svc := contexts.GetServiceInfo(ctx); svc != nil {
		serviceName = svc.Service
	}

	log := logrus.WithFields(logrus.Fields{
		ServiceKey: serviceName,
		VersionKey: version.Revision,
	})

	if otel.GetTracerProvider() != nil {
		log.Logger.Hooks.Add(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		)))
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.WithError(err).Warn("logger hostname configuration failed")
		hostname = "unknown"
	}

	log = log.WithField(HostKey, hostname)

	switch c.Format {
	case "json", "logstash":
		log.Logger.Formatter = &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyLevel: LevelKey,
				logrus.FieldKeyMsg:   MessageKey,
				logrus.FieldKeyTime:  TimestampKey,
			},
		}

	default:
		log.Logger.Formatter = &logrus.TextFormatter{}
	}

	if c.Debug {
		log.Logger.Level = logrus.DebugLevel
		log.Debug("debug logging enabled")
	} else {
		log.Logger.Level = logrus.InfoLevel
	}

	return log
}

// ContextLogger wraps logrus Entry to allow field mutation, which means the
// context itself can store a pointer to a ContextLogger, so it doesn't need
// replacing each time new fields are added to the logger
type ContextLogger struct {
	entry            *logrus.Entry
	timeoutsAsErrors bool
}

// NewContextLogger creates a new (mutable) ContextLogger instance from an (immutable) logrus Entry
func NewContextLogger(log *logrus.Entry) *ContextLogger {
	return &ContextLogger{entry: log}
}

// GetLogger returns (an immutable) logrus entry from a (mutable) ContextLogger instance
func (l *ContextLogger) GetLogger() *logrus.Entry {
	return l.entry
}

// SetField updates the internal field map
func (l *ContextLogger) SetField(field string, value any) *ContextLogger {
	l.entry = l.entry.WithField(field, value)
	return l
}

// SetFields updates the internal field map with multiple fields at a time
func (l *ContextLogger) SetFields(fields logrus.Fields) *ContextLogger {
	l.entry = l.entry.WithFields(fields)
	return l
}

// SetError updates the internal error
func (l *ContextLogger) SetError(err error) *ContextLogger {
	l.entry = l.entry.WithError(err)
	return l
}

// getContextLogger retrieves the ContextLogger instance from the context
func getContextLogger(ctx context.Context) *ContextLogger {
	ctxLogger, _ := ctx.Value(loggerKey).(*ContextLogger)
	return ctxLogger
}

func mustGetContextLogger(ctx context.Context) *ContextLogger {
	ctxLogger := getContextLogger(ctx)
	if ctxLogger != nil {
		return ctxLogger
	}

	panic("no clog exists in the context")
}

// Set sets a persistent, mutable logger for the request context.
func Set(ctx context.Context, log *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, NewContextLogger(log))
}

// Get retrieves the logrus Entry from the ContextLogger in a context
// and returns a new logrus Entry if none is found
func Get(ctx context.Context) *logrus.Entry {
	ctxLogger := getContextLogger(ctx)
	if ctxLogger != nil {
		return ctxLogger.GetLogger().WithContext(ctx)
	}

	logger := logrus.NewEntry(logrus.New()).WithContext(ctx)

	logger.Warn("no clog exists in the context")

	return logger
}

// SetField adds or updates a field to the ContextLogger in a context
func SetField(ctx context.Context, field string, value any) {
	mustGetContextLogger(ctx).SetField(field, value)
}

// SetFields adds or updates fields to the ContextLogger in a context
func SetFields(ctx context.Context, fields Fields) {
	mustGetContextLogger(ctx).SetFields(logrus.Fields(fields))
}

// SetError adds or updates an error to the ContextLogger in a context
func SetError(ctx context.Context, err error) {
	ctxLogger := mustGetContextLogger(ctx)

	ctxLogger.SetError(err)

	cherErr := cher.E{}
	if errors.As(err, &cherErr) {
		ctxLogger.SetField("error_code", cherErr.Code)
		if len(cherErr.Reasons) > 0 {
			ctxLogger.SetField("error_reasons", cherErr.Reasons)
		}

		if cherErr.Meta != nil {
			ctxLogger.SetField("error_meta", cherErr.Meta)
		}
	}
}

// ConfigureTimeoutsAsErrors changes to default behaviour of logging timeouts as info, to log them as errors
func ConfigureTimeoutsAsErrors(ctx context.Context) {
	mustGetContextLogger(ctx).timeoutsAsErrors = true
}

// TimeoutsAsErrors determines whether ConfigureTimeoutsAsErrors was called on the context
func TimeoutsAsErrors(ctx context.Context) bool {
	return mustGetContextLogger(ctx).timeoutsAsErrors
}

// DetermineLevel returns a suggested logrus Level type for a given error
func DetermineLevel(err error, timeoutsAsErrors bool) logrus.Level {
	if cherError, ok := errfuncs.As[cher.E](err); ok {
		if cherError.StatusCode() >= 500 {
			return logrus.ErrorLevel
		}

		switch cherError.Code {
		case cher.ContextCanceled:
			return levelForContextCancelation(timeoutsAsErrors)

		default:
			return logrus.WarnLevel
		}
	}

	if strings.Contains(err.Error(), "canceling statement due to user request") {
		return levelForContextCancelation(timeoutsAsErrors)
	}

	// non-cher errors are "unhandled" so warrant an error
	return logrus.ErrorLevel
}

func levelForContextCancelation(timeoutsAsErrors bool) logrus.Level {
	if timeoutsAsErrors {
		return logrus.ErrorLevel
	}

	return logrus.InfoLevel
}
