package clog

import (
	"context"
	"runtime"

	"github.com/sirupsen/logrus"
)

// HandlePanic structurally logs a panic before optionally propagating.
//
// Propagating a panic can be important for cases like calls made via net/http where the whole process isn't required
// to fail because one request panics. Propagating makes sure we dont disturb upstream panic handling.
func HandlePanic(ctx context.Context, propagate bool) {
	panicVal := recover()
	if panicVal == nil {
		return
	}

	stack := make([]byte, 1<<16) // create a 2 byte stack trace buffer
	stack = stack[:runtime.Stack(stack, false)]

	var logger *logrus.Entry
	ctxLogger := getContextLogger(ctx)
	if ctxLogger != nil {
		logger = ctxLogger.entry
	} else {
		logger = Config{Format: "json", Debug: false}.Configure(ctx)
	}

	logger.WithFields(logrus.Fields{
		"error":       "panic",
		"panic":       panicVal,
		"stack_trace": string(stack),
	}).Error("request")

	if propagate {
		panic(panicVal)
	}
}
