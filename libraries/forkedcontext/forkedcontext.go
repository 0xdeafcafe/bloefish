package forkedcontext

import (
	"context"
	"fmt"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
)

// forkedContext allows cloning context values
// while discarding deadlines and cancellations
type forkedContext struct {
	ctx context.Context
}

func (forkedContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (forkedContext) Done() <-chan struct{} {
	return nil
}

func (forkedContext) Err() error {
	return nil
}

func (d forkedContext) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

func cloneContext(ctx context.Context) context.Context {
	return forkedContext{ctx}
}

// ContextKey maps are type aware, define a custom string type for context keys
// to prevent collisions with third-party context that uses the same key.
type ContextKey string

// ForkContext provides a callback function with a new context inheriting values from the request context, and will log any error returned by the callback
func ForkContext(ctx context.Context, fn func(context.Context) error) {
	newCtx := cloneContext(ctx)
	go func() {
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}

			if err != nil {
				clog.Get(ctx).WithError(err).Log(clog.DetermineLevel(err, true), "forked context errored")
			}
		}()

		err = fn(newCtx)
	}()
}

// ForkContextWithTimeout provides a callback function with a new context inheriting values from the request context with a timeout, and will log any error returned by the callback
func ForkContextWithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) {
	newCtx, cancel := context.WithTimeout(cloneContext(ctx), timeout)

	go func() {
		var err error
		defer func() {
			cancel()
			r := recover()
			if r != nil {
				err = fmt.Errorf("panic: %v", r)
			}

			if err != nil {
				clog.Get(ctx).WithError(err).Log(clog.DetermineLevel(err, true), "forked context errored")
			}
		}()

		err = fn(newCtx)
	}()

}
