package merr

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xdeafcafe/bloefish/libraries/stacktrace"
	"github.com/google/go-cmp/cmp"
	"go.opentelemetry.io/otel/trace"
)

var _ interface {
	error
	Is(error) bool
	Unwrap() []error
} = E{}

type StackError interface {
	GetStack() []stacktrace.Frame
}

// Merrer (merr-er) represents a merr-compatible error
//
// It primarily exists to allow `Wrap` to return nil without forcing us to use
// pointers for `E`, but also allows other structs to offer a merr.E option
type Merrer interface {
	error
	Merr() E
}

type E struct {
	Code Code `json:"code"`
	Meta M    `json:"meta"`

	TraceID trace.TraceID `json:"trace_id"`
	SpanID  trace.SpanID  `json:"span_id"`

	Stack []stacktrace.Frame `json:"stack"`

	Reasons []error `json:"reasons"`
}

type M map[string]any

func New(ctx context.Context, code Code, meta M, reasons ...error) E {
	spanContext := trace.SpanContextFromContext(ctx)

	for _, reason := range reasons {
		if reason == nil {
			panic("merr.New: nil reason provided for " + code.String())
		}
	}
	stack := stacktrace.GetCallerFrames(2)
	// For errors with only one reason, we can merge stack traces so we can actually see where the error originally came from
	if len(reasons) == 1 {
		if stackErr, ok := reasons[0].(StackError); ok {
			rootStack := stackErr.GetStack()
			stack = stacktrace.MergeStacks(rootStack, stack)
		}
	}
	return E{
		Code: code,
		Meta: meta,

		TraceID: spanContext.TraceID(),
		SpanID:  spanContext.SpanID(),

		Stack: stack,

		Reasons: reasons,
	}
}

func (e E) Merr() E {
	return e
}

func (e E) Fields() map[string]any {
	return map[string]any{
		"code":    e.Code,
		"meta":    e.Meta,
		"stack":   e.Stack,
		"reasons": e.Reasons,
	}
}

func (e E) Equal(e2 E) bool {
	return e.Code == e2.Code &&
		cmp.Equal(e.Meta, e2.Meta) &&
		cmp.Equal(e.Stack, e2.Stack) &&
		cmp.Equal(e.Reasons, e2.Reasons)
}

func (e E) String() string {
	return e.Error() + "\n\n" + stacktrace.FormatFrames(e.Stack)
}

// Error implements the error interface
//
// Provides a simple string representation of the error, but cannot include the
// complete data contained in the error
//
// No compatibility guarantees are made with its output - it may change at any time
func (e E) Error() string {
	str := string(e.Code)

	if len(e.Meta) > 0 {
		str += fmt.Sprintf(" (%v)", e.Meta)
	}

	for _, reason := range e.Reasons {
		str += fmt.Sprintf("\n- %v", reason)
	}

	return str
}

// Is enables the use of `errors.Is`
func (e E) Is(err error) bool {
	if errors.Is(e.Code, err) {
		return true
	}

	if merr, ok := err.(E); ok {
		return e.Equal(merr)
	}

	return false
}

// Unwrap enables the use of `errors.Unwrap`
func (e E) Unwrap() []error {
	return e.Reasons
}

func (e E) GetStack() []stacktrace.Frame {
	return e.Stack
}
