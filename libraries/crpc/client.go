package crpc

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/0xdeafcafe/bloefish/libraries/contexts"
	"github.com/0xdeafcafe/bloefish/libraries/errfuncs"
	"github.com/0xdeafcafe/bloefish/libraries/jsonclient"
	"github.com/0xdeafcafe/bloefish/libraries/version"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	userAgentTemplate            = "crpc/%s (+https://github.com/0xdeafcafe/bloefish/tree/main/lib/crpc)"
	userAgentTemplateWithService = "crpc/%s (+https://github.com/0xdeafcafe/bloefish/tree/main/lib/crpc) [%s/%s]"
)

// Client represents a crpc client. It builds on top of jsonclient, so error
// variables/structs and the authenticated round tripper live there.
type Client struct {
	client *jsonclient.Client
}

// NewClient returns a client configured with a transport scheme, remote host
// and URL prefix supplied as a URL <scheme>://<host></prefix>
func NewClient(ctx context.Context, baseURL string, c *http.Client) *Client {
	if c == nil {
		c = &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport, otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return fmt.Sprintf("crpc %s%s", r.URL.Hostname(), r.URL.Path)
			})),
		}
	}

	jcc := jsonclient.NewClient(baseURL, c)
	svc := contexts.GetServiceInfo(ctx)
	if svc != nil {
		jcc.UserAgent = fmt.Sprintf(userAgentTemplateWithService, version.Truncated, svc.Service, svc.Environment)
	} else {
		jcc.UserAgent = fmt.Sprintf(userAgentTemplate, version.Truncated)
	}

	return &Client{jcc}
}

// Do executes an RPC request against the configured server.
func (c *Client) Do(ctx context.Context, method, version string, src, dst any, requestModifiers ...func(r *http.Request)) error {
	err := c.client.Do(ctx, "POST", path.Join(version, method), nil, src, dst, requestModifiers...)
	if err == nil {
		return nil
	}

	if err, ok := errfuncs.As[jsonclient.ClientTransportError](err); ok {
		return ClientTransportError{method, version, err.ErrorString, err.Cause()}
	}

	return err
}

// ClientTransportError is returned when an error related to
// executing a client request occurs.
type ClientTransportError struct {
	Method, Version, ErrorString string

	cause error
}

// Cause returns the causal error (if wrapped) or nil
func (cte ClientTransportError) Cause() error {
	return cte.cause
}

func (cte ClientTransportError) Error() string {
	if cte.cause != nil {
		return fmt.Sprintf("%s/%s %s: %s", cte.Version, cte.Method, cte.ErrorString, cte.cause.Error())
	}

	return fmt.Sprintf("%s/%s %s", cte.Version, cte.Method, cte.ErrorString)
}
