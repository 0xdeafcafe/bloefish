package jsonschema

import (
	"testing"

	"github.com/matryer/is"
)

func TestKSUIDFormatChecker(t *testing.T) {
	is := is.New(t)

	c := ksuidFormatChecker{}

	is.Equal(c.IsFormat("foo"), false)
	is.Equal(c.IsFormat("user_000000CBEvdtGRrnrcQKCsSDNNKmR"), true)
	is.Equal(c.IsFormat("test_user_000000CBEvefIcrYXsZObXFKZBQrh"), true)
}
