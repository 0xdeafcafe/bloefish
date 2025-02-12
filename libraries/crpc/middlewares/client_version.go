package middlewares

import (
	"context"
	"net/http"
	"regexp"
	"strconv"

	"github.com/0xdeafcafe/bloefish/libraries/merr"
	"github.com/blang/semver/v4"
)

const (
	ClientPlatformIOS     = "ios"
	ClientPlatformAndroid = "android"
)

// keep in sync with the CDN function which does this check globally
var clientHeaderRegexp = regexp.MustCompile(`^(ios|android)-(\d+\.\d+\.\d+)-(\d+)$`)

// ClientVersion represents a parsed client version header string
type ClientVersion struct {
	Platform string
	Version  semver.Version
	Build    int
}

// GetClientVersion retrieves the parsed version from a request
func GetClientVersion(r *http.Request) *ClientVersion {
	return GetClientVersionContext(r.Context())
}

// GetClientVersionContext retrieves the parsed version from a request context
func GetClientVersionContext(ctx context.Context) *ClientVersion {
	if clientVersion, ok := ctx.Value(clientVersionContextKey).(*ClientVersion); ok {
		return clientVersion
	}

	return nil
}

// ParseClientVersion attempts to parse the infra-client-version HTTP header and add
// it as a struct to the context
func ParseClientVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		parsedClientVersion, err := parseVersionHeader(ctx, r.Header.Get("Infra-Client-Version"))
		if err == nil {
			r = r.WithContext(context.WithValue(ctx, clientVersionContextKey, parsedClientVersion))
		}

		next.ServeHTTP(w, r)
	})
}

func parseVersionHeader(ctx context.Context, clientVersionHeader string) (*ClientVersion, error) {
	if clientVersionHeader == "" {
		return nil, merr.New(ctx, "client_version_missing", nil)
	}

	versionParts := clientHeaderRegexp.FindStringSubmatch(clientVersionHeader)
	if len(versionParts) != 4 {
		return nil, merr.New(ctx, "client_version_invalid", merr.M{"header": clientVersionHeader})
	}

	platform := versionParts[1]
	version := versionParts[2]
	buildStr := versionParts[3]

	clientSemver, err := semver.Parse(version)
	if err != nil {
		return nil, merr.New(ctx, "client_version_invalid", merr.M{"header": clientVersionHeader}, err)
	}

	build, err := strconv.Atoi(buildStr)
	if err != nil {
		return nil, merr.New(ctx, "client_version_invalid", merr.M{"header": clientVersionHeader}, err)
	}

	return &ClientVersion{
		Platform: platform,
		Version:  clientSemver,
		Build:    build,
	}, nil
}
