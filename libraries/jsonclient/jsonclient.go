package jsonclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	pathlib "path"
	"strings"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/errfuncs"
	"github.com/0xdeafcafe/bloefish/libraries/version"
)

// ErrNoResponse is returned when a client request is given a body to
// unmarshal to however the server does not return any content (HTTP 204).
var ErrNoResponse = ClientRequestError{"no response to unmarshal to body", nil}

// DefaultUserAgent is the default HTTP User-Agent Header that is presented to the server.
var DefaultUserAgent = "jsonclient/" + version.Truncated + " (+https://github.com/0xdeafcafe/bloefish/tree/main/lib/jsonclient)"

// Client represents a json-client HTTP client.
type Client struct {
	Scheme string
	Host   string
	Prefix string

	UserAgent string

	Client *http.Client
}

// NewClient returns a client configured with a transport scheme, remote host
// and URL prefix supplied as a URL <scheme>://<host></prefix>
func NewClient(baseURL string, client *http.Client) *Client {
	remote, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &Client{
		Scheme: remote.Scheme,
		Host:   remote.Host,
		Prefix: remote.Path,

		UserAgent: DefaultUserAgent,

		Client: client,
	}
}

// Do executes an HTTP request against the configured server.
func (c *Client) Do(ctx context.Context, method, path string, params url.Values, src, dst any, requestModifiers ...func(r *http.Request)) error {
	return c.DoWithHeaders(ctx, method, path, nil, params, src, dst, requestModifiers...)
}

// DoWithHeaders executes an HTTP request against the configured server with custom headers.
func (c *Client) DoWithHeaders(ctx context.Context, method, path string, headers http.Header, params url.Values, src, dst any, requestModifiers ...func(r *http.Request)) error {
	fullPath := pathlib.Join("/", c.Prefix, path)
	req := &http.Request{
		Method: method,
		URL: &url.URL{
			Scheme: c.Scheme,
			Host:   c.Host,
			Path:   fullPath,
		},
		Header: http.Header{
			"Accept":     []string{"application/json"},
			"User-Agent": []string{c.UserAgent},
		},
		Host: c.Host,
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	for key, value := range headers {
		req.Header[key] = value
	}

	for _, modifier := range requestModifiers {
		modifier(req)
	}

	err := setRequestBody(req, src)
	if err != nil {
		return ClientRequestError{"could not marshal", err}
	}

	res, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		if netErr, ok := errfuncs.As[net.Error](err); ok {
			if netErr.Timeout() {
				return cher.New(cher.RequestTimeout, cher.M{"method": method, "path": fullPath, "host": c.Host, "scheme": c.Scheme, "timeout_error": netErr.Error()})
			}

			return ClientTransportError{method, path, "request failed", netErr}
		}

		return ClientTransportError{method, path, "unknown error", err}
	}

	defer res.Body.Close()

	return handleResponse(res, method, path, dst)
}

func setRequestBody(req *http.Request, src any) error {
	if src != nil {
		data, err := json.Marshal(src)
		if err != nil {
			return err
		}

		req.Body = io.NopCloser(bytes.NewReader(data))
		req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(data)), nil }
		req.ContentLength = int64(len(data))

		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return nil
}

func handleResponse(res *http.Response, method, path string, dst any) error {
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		if dst == nil {
			return nil
		}

		if res.StatusCode == http.StatusNoContent || res.Body == nil {
			return ErrNoResponse
		}

		err := json.NewDecoder(res.Body).Decode(dst)
		if errors.Is(err, io.EOF) {
			return ErrNoResponse
		} else if err != nil {
			return ClientTransportError{method, path, "could not unmarshal", err}
		}

		return nil
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return ClientTransportError{method, path, "could not read response body stream", err}
	}

	var body cher.E
	if err := json.Unmarshal(resBody, &body); err == nil && body.Code != "" {
		return body
	}

	var errorResBody any
	if err := json.Unmarshal(resBody, &errorResBody); err != nil {
		errorResBody = string(resBody)
	}

	statusText := http.StatusText(res.StatusCode)
	if statusText == "" {
		statusText = "unknown"
	}

	statusParts := strings.Fields(statusText)

	for i := range statusParts {
		statusParts[i] = strings.ToLower(statusParts[i])
	}

	newErrorMessage := strings.Join(statusParts, "_")

	return cher.New(newErrorMessage, cher.M{
		"httpStatus": res.StatusCode,
		"data":       errorResBody,
		"method":     res.Request.Method,
		"url":        res.Request.URL.String(),
	})
}

// ClientRequestError is returned when an error related to
// constructing a client request occurs.
type ClientRequestError struct {
	ErrorString string

	cause error
}

// Cause returns the causal error (if wrapped) or nil
func (cre ClientRequestError) Cause() error {
	return cre.cause
}

func (cre ClientRequestError) Error() string {
	if cre.cause != nil {
		return cre.ErrorString + ": " + cre.cause.Error()
	}

	return cre.ErrorString
}

// ClientTransportError is returned when an error related to
// executing a client request occurs.
type ClientTransportError struct {
	Method, Path, ErrorString string

	cause error
}

// Cause returns the causal error (if wrapped) or nil
func (cte ClientTransportError) Cause() error {
	return cte.cause
}

func (cte ClientTransportError) Error() string {
	if cte.cause != nil {
		return fmt.Sprintf("%s %s %s: %s", cte.Method, cte.Path, cte.ErrorString, cte.cause.Error())
	}

	return fmt.Sprintf("%s %s %s", cte.Method, cte.Path, cte.ErrorString)
}
