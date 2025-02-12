package stacktrace

import (
	"fmt"
	"runtime"
)

// GetCallerCodePath returns the caller path and line number.
//
// The argument skip is the number of stack frames to skip before identifying
// the frame to use, with 0 identifying the frame for GetCallerCodePath itself
// and 1 identifying the caller of GetCallerCodePath.
func GetCallerCodePath(skip int, fallback string) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok || file == "" {
		return fallback
	}
	return fmt.Sprintf("%s:%d", file, line)
}
