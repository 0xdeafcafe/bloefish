package middlewares

import (
	"context"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/matryer/is"
)

func TestParseVersionHeader(t *testing.T) {
	t.Run("empty should error", func(t *testing.T) {
		is := is.New(t)

		parsed, err := parseVersionHeader(context.Background(), "")
		is.True(err != nil)
		is.Equal(parsed, nil)
	})

	t.Run("invalid should error", func(t *testing.T) {
		is := is.New(t)

		parsed, err := parseVersionHeader(context.Background(), "a-1-")
		is.True(err != nil)
		is.Equal(parsed, nil)
	})

	t.Run("ios should work", func(t *testing.T) {
		is := is.New(t)

		parsed, err := parseVersionHeader(context.Background(), "ios-3.6.8-1337")
		is.NoErr(err)
		is.Equal(parsed.Platform, ClientPlatformIOS)
		is.Equal(parsed.Version, semver.MustParse("3.6.8"))
		is.Equal(parsed.Build, 1337)
	})

	t.Run("android should work", func(t *testing.T) {
		is := is.New(t)

		parsed, err := parseVersionHeader(context.Background(), "android-0.0.1-1")
		is.NoErr(err)
		is.Equal(parsed.Platform, ClientPlatformAndroid)
		is.Equal(parsed.Version, semver.MustParse("0.0.1"))
		is.Equal(parsed.Build, 1)
	})
}
