package stacktrace

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestGetCallerCodePath(t *testing.T) {
	is := is.New(t)

	codePath := GetCallerCodePath(1, "")

	is.True(strings.HasSuffix(codePath, "/lib/stacktrace/codepath_test.go:13"))
}
