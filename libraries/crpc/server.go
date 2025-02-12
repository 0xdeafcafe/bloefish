package crpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/errfuncs"
	"github.com/0xdeafcafe/bloefish/libraries/merr"
	"github.com/0xdeafcafe/bloefish/libraries/mlog"
	"github.com/xeipuuv/gojsonschema"
)

// ResponseWriter is the destination for RPC responses.
type ResponseWriter interface {
	io.Writer
}

// Request contains metadata about the RPC request.
type Request struct {
	Version string
	Method  string

	Body io.ReadCloser

	RemoteAddr    string
	BrowserOrigin string

	originalRequest *http.Request
}

func (r *Request) Context() context.Context {
	return r.originalRequest.Context()
}

type contextKey string

const requestKey contextKey = "crpcrequest"

// GetRequestContext returns the Request from the context object
func GetRequestContext(ctx context.Context) *Request {
	if val, ok := ctx.Value(requestKey).(*Request); ok {
		return val
	}

	return nil
}

func setRequestContext(ctx context.Context, request *Request) context.Context {
	return context.WithValue(ctx, requestKey, request)
}

// HandlerFunc defines a handler for an RPC request. Request and response body
// data will be JSON. Request will immediately io.EOF if there is no request.
type HandlerFunc func(res http.ResponseWriter, req *Request) error

// MiddlewareFunc is a function that wraps HandlerFuncs.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// WrappedFunc contains the wrapped handler, and some additional information
// about the function that was determined during the reflection process
type WrappedFunc struct {
	Handler           HandlerFunc
	HasRequestInput   bool
	HasResponseOutput bool
}

var (
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
	contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
)

// Wrap reflects a HandlerFunc from any function matching the
// following signatures:
//
// func(ctx context.Context, request *T) (response *T, err error)
// func(ctx context.Context, request *T) (err error)
// func(ctx context.Context) (response *T, err error)
// func(ctx context.Context) (err error)
func Wrap(fn any) (*WrappedFunc, error) {
	ctx := context.Background()

	// prevent re-reflection of type that is already a HandlerFunc
	if _, ok := fn.(HandlerFunc); ok {
		return nil, merr.New(ctx, "fn_is_handler", nil)
	} else if _, ok := fn.(*WrappedFunc); ok {
		return nil, merr.New(ctx, "fn_already_wrapped", nil)
	}

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// check the basic type and the number of inputs/outputs
	if fnType.Kind() != reflect.Func {
		return nil, merr.New(ctx, "fn_type_invalid", merr.M{"type": fnType.Kind()})
	}

	inputCount := fnType.NumIn()
	outputCount := fnType.NumOut()
	if inputCount < 1 || inputCount > 2 {
		return nil, merr.New(ctx, "fn_input_params_invalid", merr.M{"count": inputCount})
	} else if outputCount < 1 || outputCount > 2 {
		return nil, merr.New(ctx, "fn_output_params_invalid", merr.M{"count": outputCount})
	}

	firstInput := fnType.In(0)
	lastOutput := fnType.Out(outputCount - 1)
	if !firstInput.Implements(contextType) {
		return nil, merr.New(ctx, "fn_first_input_not_context", merr.M{"type": firstInput})
	} else if !lastOutput.Implements(errorType) {
		return nil, merr.New(ctx, "fn_last_output_not_error", merr.M{"type": lastOutput})
	}

	// resolve function parameter pointers to underlying type for use with
	// reflect.New (which will return pointers).
	var reqType reflect.Type
	var hasResponseOutput bool

	if inputCount == 2 {
		secondInput := fnType.In(1)
		if secondInput.Kind() != reflect.Ptr {
			return nil, merr.New(ctx, "fn_second_input_not_pointer", merr.M{"type": secondInput.Kind()})
		}

		reqType = secondInput.Elem()
		if reqType.Kind() != reflect.Struct {
			return nil, merr.New(ctx, "fn_second_input_not_struct_pointer", merr.M{"type": reqType.Kind()})
		}
	}

	if outputCount == 2 {
		hasResponseOutput = true

		err := checkResponseType(ctx, fnType.Out(0))
		if err != nil {
			return nil, err
		}
	}

	handler := func(w http.ResponseWriter, req *Request) error {
		ctxVal := reflect.ValueOf(req.Context())
		var inputs []reflect.Value

		if reqType == nil {
			if req.Body != nil {
				i, err := req.Body.Read(make([]byte, 1))
				if i != 0 || !errors.Is(err, io.EOF) {
					return cher.New(cher.BadRequest, nil, cher.New("unexpected_request_body", nil))
				}
			}

			inputs = []reflect.Value{ctxVal}
		} else {
			if req.Body == nil {
				return cher.New(cher.BadRequest, nil, cher.New("missing_request_body", nil))
			}

			reqVal := reflect.New(reqType)
			err := json.NewDecoder(req.Body).Decode(reqVal.Interface())
			if errors.Is(err, io.EOF) {
				return cher.New(cher.BadRequest, nil, cher.New("missing_request_body", nil))
			} else if err != nil {
				return merr.New(ctx, "request_body_decode_failed", nil, err)
			}

			inputs = []reflect.Value{ctxVal, reqVal}
		}

		res := fnValue.Call(inputs)

		if errVal := res[len(res)-1]; !errVal.IsNil() {
			return errVal.Interface().(error) //nolint:forcetypeassert // we checked the type above
		}

		if len(res) == 1 {
			w.WriteHeader(http.StatusNoContent)
		} else if len(res) == 2 {
			enc := json.NewEncoder(w)
			enc.SetEscapeHTML(false)
			err := enc.Encode(res[0].Interface())
			if err != nil {
				if strings.Contains(err.Error(), "broken pipe") {
					return nil
				}

				return err
			}
		}

		return nil
	}

	return &WrappedFunc{
		Handler:           handler,
		HasRequestInput:   reqType != nil,
		HasResponseOutput: hasResponseOutput,
	}, nil
}

func checkResponseType(ctx context.Context, typ reflect.Type) error {
	switch typ.Kind() { //nolint:exhaustive // we only accept a few types
	case reflect.Ptr:
		elem := typ.Elem()

		if elem.Kind() != reflect.Struct && elem.Kind() != reflect.String {
			err := checkResponseType(ctx, elem)
			if err != nil {
				return err
			}
		}

		return nil

	case reflect.Slice:
		elem := typ.Elem()

		if elem.Kind() != reflect.String {
			err := checkResponseType(ctx, elem)
			if err != nil {
				return err
			}
		}

		return nil

	default:
		return merr.New(ctx, "response_type_invalid", merr.M{"type": typ.Kind()})
	}
}

// MustWrap is the same as Wrap, however it panics when passed an
// invalid handler.
func MustWrap(fn any) *WrappedFunc {
	wrapped, err := Wrap(fn)
	if err != nil {
		panic(err)
	}

	return wrapped
}

type wrappedHandler struct {
	v  string
	fn HandlerFunc
}

// Server is an HTTP-compatible crpc handler.
type Server struct {
	// AuthenticationMiddleware applies authentication before any other
	// middleware or request is processed. Servers without middleware must
	// configure the UnsafeNoAuthentication middleware. If no
	// AuthenticationMiddleware is configured, the server will panic.
	AuthenticationMiddleware MiddlewareFunc

	// methods = version -> method -> HandlerFunc
	registeredVersionMethods map[string]map[string]*wrappedHandler
	registeredPreviewMethods map[string]*wrappedHandler

	resolvedMethods map[string]map[string]*wrappedHandler

	mw []MiddlewareFunc
}

// NewServer returns a new RPC Server with an optional exception tracker.
func NewServer(auth MiddlewareFunc) *Server {
	return &Server{
		AuthenticationMiddleware: auth,
	}
}

// Use includes a piece of Middleware to wrap HandlerFuncs.
func (s *Server) Use(mw MiddlewareFunc) {
	s.mw = append(s.mw, mw)
}

const (
	// VersionPreview is used for experimental endpoints in development which
	// are coming but a version identifier has not been decided yet or may
	// be withdrawn at any point.
	VersionPreview = "preview"

	// VersionLatest is used by engineers only to call the latest version
	// of an endpoint in utilities like cURL and Paw.
	VersionLatest = "latest"
)

// expVersion matches valid method versions
var expVersion = regexp.MustCompile(`^(?:preview|20\d{2}-\d{2}-\d{2})$`)

// expMethod matched valid method names
var expMethod = regexp.MustCompile(`^[a-z][a-z\d]*(?:_[a-z\d]+)*$`)

func isValidMethod(method, version string) bool {
	return expMethod.MatchString(method) && expVersion.MatchString(version)
}

// Register reflects a HandlerFunc from fnR and associates it with a
// method name and version. If fnR does not meet the HandlerFunc standard
// defined above, or the presence of the schema doesn't match the presence
// of the input argument, Register will panic. This function is not thread safe
// and must be run in serial if called multiple times.
func (s *Server) Register(method, version string, schema gojsonschema.JSONLoader, fn any, middleware ...MiddlewareFunc) {
	if fn == nil {
		s.RegisterFunc(method, version, schema, nil, middleware...)

		return
	}

	wrapped := MustWrap(fn)
	hasSchema := schema != nil

	if wrapped.HasRequestInput != hasSchema {
		if hasSchema {
			panic("schema validation configured, but handler doesn't accept input")
		} else {
			panic("no schema validation configured")
		}
	}

	s.RegisterFunc(method, version, schema, &wrapped.Handler, middleware...)
}

// RegisterFunc associates a method name and version with a HandlerFunc,
// and optional middleware. This function is not thread safe and must be run in
// serial if called multiple times.
func (s *Server) RegisterFunc(method, version string, schema gojsonschema.JSONLoader, fn *HandlerFunc, middleware ...MiddlewareFunc) {
	if s.registeredVersionMethods == nil {
		s.registeredVersionMethods = make(map[string]map[string]*wrappedHandler)
	}

	if s.registeredPreviewMethods == nil {
		s.registeredPreviewMethods = make(map[string]*wrappedHandler)
	}

	if fn == nil && schema != nil {
		panic("schema validation configured, but handler is nil")
	}

	if !isValidMethod(method, version) {
		panic("invalid method/version")
	} else if s.AuthenticationMiddleware == nil {
		panic("no authentication configured")
	}

	if fn == nil && version == VersionPreview {
		panic(fmt.Sprintf("cannot set preview method '%s' as nil", method))
	}

	if s.isRouteDefined(method, version) {
		panic(fmt.Sprintf("cannot set '%s' on version '%s', it's already defined", method, version))
	}

	if fn == nil {
		s.setRoute(version, method, nil)
	} else {
		if schema != nil {
			compiledSchema, err := gojsonschema.NewSchemaLoader().Compile(schema)
			if err != nil {
				panic(fmt.Sprintf("json schema error in %s: %s", method, err))
			}

			middleware = append([]MiddlewareFunc{s.AuthenticationMiddleware, Validate(compiledSchema)}, middleware...)
		} else {
			middleware = append([]MiddlewareFunc{s.AuthenticationMiddleware}, middleware...)
		}

		// This wraps the middleware funcs inside each one in reverse order
		for i := range middleware {
			p := middleware[len(middleware)-1-i](*fn)
			fn = &p
		}

		for i := range s.mw {
			p := s.mw[len(s.mw)-1-i](*fn)
			fn = &p
		}

		s.setRoute(version, method, &wrappedHandler{version, *fn})
	}

	s.buildRoutes()
}

func (s Server) isRouteDefined(method, version string) bool {
	if version == VersionPreview {
		_, ok := s.registeredPreviewMethods[method]
		return ok
	}

	if methodSet, ok := s.registeredVersionMethods[version]; ok {
		_, ok := methodSet[method]
		return ok
	}

	return false
}

func (s *Server) setRoute(version, method string, handler *wrappedHandler) {
	if version == VersionPreview {
		s.registeredPreviewMethods[method] = handler

		return
	}

	versions, ok := s.registeredVersionMethods[version]
	if !ok {
		versions = make(map[string]*wrappedHandler)
		s.registeredVersionMethods[version] = versions
	}

	versions[method] = handler
}

func (s *Server) buildRoutes() {
	knownVersions := sort.StringSlice{}
	resolvedMethods := make(map[string]map[string]*wrappedHandler)

	// build known versions
	for version, methodSet := range s.registeredVersionMethods {
		if methodSet == nil {
			continue
		}

		if _, ok := resolvedMethods[version]; !ok {
			knownVersions = append(knownVersions, version)
			resolvedMethods[version] = make(map[string]*wrappedHandler)
		}
	}

	// We must ensure that the earliest version is done first
	sort.Sort(knownVersions)

	var previousVersion string

	//	loop over versions, earliest first
	// build up each version by copying the previous version as a base
	// then setting on the version any explicitly defined method
	for _, version := range knownVersions {
		if previousVersion != "" {
			for mn, fn := range resolvedMethods[previousVersion] {
				resolvedMethods[version][mn] = fn
			}
		}

		for mn, fn := range s.registeredVersionMethods[version] {
			if fn == nil {
				delete(resolvedMethods[version], mn)
			} else {
				resolvedMethods[version][mn] = fn
			}
		}

		previousVersion = version
	}

	// build up latest methodSet if previous methodSet has been made
	if previousVersion != "" {
		resolvedMethods[VersionLatest] = resolvedMethods[previousVersion]
	}

	// Handle preview methods
	if len(s.registeredPreviewMethods) > 0 {
		resolvedMethods[VersionPreview] = make(map[string]*wrappedHandler)
	}

	for mn, fn := range s.registeredPreviewMethods {
		if fn == nil {
			panic("cannot set preview method as nil")
		}

		resolvedMethods[VersionPreview][mn] = fn
	}

	s.resolvedMethods = resolvedMethods
}

// NoLongerSupported is a HandlerFunc which always returns the error
// `no_longer_supported` to the requester to indicate methods which have
// been withdrawn.
func NoLongerSupported(_ http.ResponseWriter, _ *Request) error {
	return cher.New(cher.NoLongerSupported, nil)
}

// Serve executes an RPC request.
func (s *Server) Serve(res http.ResponseWriter, req *Request) error {
	if s.AuthenticationMiddleware == nil {
		return cher.New(cher.AccessDenied, nil)
	}

	if s.resolvedMethods == nil {
		return cher.New("no_methods_registered", nil)
	}

	methodSet, ok := s.resolvedMethods[req.Version]
	if !ok {
		return cher.New(cher.NotFound, cher.M{"version": req.Version})
	}

	handler, ok := methodSet[req.Method]
	if !ok || handler == nil {
		return cher.New(cher.NotFound, cher.M{"method": req.Method, "version": req.Version})
	}

	// append latest version to Infra-Endpoint-Status
	appendInfraEndpointStatus(res, req.Version, handler.v)

	fn := handler.fn

	return fn(res, req)
}

const (
	// InfraEndpointStatus is the header appended to the response indicating the
	// usability status of the endpoint.
	InfraEndpointStatus = `Infra-Endpoint-Status`

	// PreviewNotice is the contents of the `Infra-Endpoint-Status` header when an
	// endpoint is called with the preview version.
	PreviewNotice = `preview; msg="endpoint is experimental and may change/be withdrawn without notice"`

	// LatestNotice is the contents of the `Infra-Endpoint-Status` header when an
	// endpoint is called to request the latest version.
	LatestNotice = `latest; msg="subject to change without notice"`

	// StableNotice is the contents of the `Infra-Endpoint-Status` header when an
	// endpoint is called and is not expected to change.
	StableNotice = `stable`
)

// appendInfraEndpointStatus applies the appropriate `Infra-Endpoint-Status`
// header for the method version requested by the client.
func appendInfraEndpointStatus(w http.ResponseWriter, requestedVersion, resolvedVersion string) {
	switch requestedVersion {
	case VersionPreview:
		w.Header().Set(InfraEndpointStatus, PreviewNotice)

	case VersionLatest:
		message := fmt.Sprintf(`%s; v="%s"`, LatestNotice, resolvedVersion)
		w.Header().Set(InfraEndpointStatus, message)

	default:
		w.Header().Set(InfraEndpointStatus, StableNotice)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if strings.ToUpper(r.Method) != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.RawQuery != "" {
		s.writeError(ctx, w, cher.New("unexpected_input", nil))
		return
	}

	req := &Request{
		Body: r.Body,

		RemoteAddr:    r.RemoteAddr,
		BrowserOrigin: r.Header.Get("Origin"),
	}

	ctx = setRequestContext(ctx, req)
	r = r.WithContext(ctx)
	req.originalRequest = r

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var ok bool
	req.Method, req.Version, ok = requestPath(r.URL.Path)
	if !ok {
		s.writeError(ctx, w, cher.New(cher.NotFound, nil))
		return
	}

	s.writeError(ctx, w, s.Serve(w, req))
}

// expRequestPath only matches HTTP Paths formed of /<version date>/<method name>
var expRequestPath = regexp.MustCompile(`^/(preview|latest|20\d{2}-\d{2}-\d{2})/([a-z0-9\_]+)$`)

func requestPath(path string) (method, version string, ok bool) {
	match := expRequestPath.FindStringSubmatch(path)
	if len(match) != 3 {
		return
	}

	version = match[1]
	method = match[2]
	ok = true
	return
}

func (s *Server) writeError(ctx context.Context, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var body cher.E

	if err, ok := errfuncs.As[cher.E](err); ok {
		body = err
	} else if err, ok := errfuncs.As[*json.SyntaxError](err); ok {
		body = cher.New(
			"invalid_json",
			cher.M{
				"error":  err.Error(),
				"offset": err.Offset,
			},
		)
	} else if err, ok := errfuncs.As[*json.UnmarshalTypeError](err); ok {
		body = cher.New(
			"invalid_json",
			cher.M{
				"expected": err.Type.Kind().String(),
				"actual":   err.Value,
				"name":     err.Field,
			},
		)
	} else {
		body = cher.New(cher.Unknown, nil)
	}

	w.WriteHeader(body.StatusCode())

	werr := json.NewEncoder(w).Encode(body)
	if werr != nil {
		mlog.Warn(ctx, merr.New(ctx, "crpc_write_error_failed", nil, werr))
	}
}
