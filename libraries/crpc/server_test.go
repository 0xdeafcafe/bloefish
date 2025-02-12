package crpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xdeafcafe/bloefish/libraries/merr"
	"github.com/matryer/is"
	"github.com/xeipuuv/gojsonschema"
)

type testInput struct{}

type testOutput struct{}

func TestWrap(t *testing.T) {
	tests := []struct {
		Name  string
		Fn    any
		Error merr.Code
	}{
		{
			"HandlerFunc",
			HandlerFunc(func(w http.ResponseWriter, r *Request) error { return nil }),
			"fn_is_handler",
		},
		{
			"WrappedFunc", &WrappedFunc{},
			"fn_already_wrapped",
		},
		{
			"NotFunc", "string",
			"fn_type_invalid",
		},
		{
			"NoInput", func() {},
			"fn_input_params_invalid",
		},
		{
			"LongInput", func(ctx context.Context, foo string, bar string) {},
			"fn_input_params_invalid",
		},
		{
			"NoOutput", func(ctx context.Context) {},
			"fn_output_params_invalid",
		},
		{
			"LongOutput", func(ctx context.Context) (foo, bar string, err error) { return },
			"fn_output_params_invalid",
		},
		{
			"ContextRequired", func(foo string) error { return nil },
			"fn_first_input_not_context",
		},
		{
			"ErrorRequired", func(ctx context.Context) string { return "" },
			"fn_last_output_not_error",
		},
		{
			"InputNotPointer", func(ctx context.Context, in testInput) error { return nil },
			"fn_second_input_not_pointer",
		},
		{
			"InputNotStruct", func(ctx context.Context, in *string) error { return nil },
			"fn_second_input_not_struct_pointer",
		},
		{
			"OutputNotPointer", func(ctx context.Context) (out testOutput, err error) { return },
			"response_type_invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			is := is.New(t)

			_, err := Wrap(test.Fn)
			is.Equal(test.Error, err.(merr.E).Code) //nolint:errorlint,forcetypeassert // required for test
		})
	}
}

func TestMethodsAreBroughtForward(t *testing.T) {
	is := is.New(t)

	foov1 := &wrappedHandler{v: "2019-01-01"}
	barv1 := &wrappedHandler{v: "2019-02-02"}

	rpc := Server{
		registeredVersionMethods: map[string]map[string]*wrappedHandler{
			"2019-01-01": {
				"foo": foov1,
			},
			"2019-02-02": {
				"bar": barv1,
			},
		},
	}

	rpc.buildRoutes()

	expected := map[string]map[string]*wrappedHandler{
		"2019-01-01": {
			"foo": foov1,
		},
		"2019-02-02": {
			"foo": foov1,
			"bar": barv1,
		},
		"latest": {
			"foo": foov1,
			"bar": barv1,
		},
	}

	is.Equal(expected, rpc.resolvedMethods)
}

func TestMethodsAreBroughtForwardComplex(t *testing.T) {
	is := is.New(t)

	foov1 := &wrappedHandler{v: "2019-01-01"}
	foov2 := &wrappedHandler{v: "2019-02-02"}
	barv1 := &wrappedHandler{v: "2019-02-02"}
	barv2 := &wrappedHandler{v: "2019-03-03"}
	barv3 := &wrappedHandler{v: "2019-04-04"}

	rpc := Server{
		registeredVersionMethods: map[string]map[string]*wrappedHandler{
			"2019-01-01": {
				"foo": foov1,
			},
			"2019-02-02": {
				"foo": foov2,
				"bar": barv1,
			},
			"2019-03-03": {
				"bar": barv2,
			},
			"2019-04-04": {
				"bar": barv3,
			},
		},
	}

	rpc.buildRoutes()

	expected := map[string]map[string]*wrappedHandler{
		"2019-01-01": {
			"foo": foov1,
		},
		"2019-02-02": {
			"foo": foov2,
			"bar": barv1,
		},
		"2019-03-03": {
			"foo": foov2,
			"bar": barv2,
		},
		"2019-04-04": {
			"foo": foov2,
			"bar": barv3,
		},
		"latest": {
			"foo": foov2,
			"bar": barv3,
		},
	}

	is.Equal(expected, rpc.resolvedMethods)
}

func TestMethodsAreBroughtForwardAndRemoved(t *testing.T) {
	is := is.New(t)

	foov1 := &wrappedHandler{v: "2019-01-01"}
	barv1 := &wrappedHandler{v: "2019-01-01"}

	rpc := Server{
		registeredVersionMethods: map[string]map[string]*wrappedHandler{
			"2019-01-01": {
				"foo": foov1,
				"bar": barv1,
			},
			"2019-02-02": {
				"bar": nil,
			},
		},
	}

	rpc.buildRoutes()

	expected := map[string]map[string]*wrappedHandler{
		"2019-01-01": {
			"foo": foov1,
			"bar": barv1,
		},
		"2019-02-02": {
			"foo": foov1,
		},
		"latest": {
			"foo": foov1,
		},
	}

	is.Equal(expected, rpc.resolvedMethods)
}

func TestMethodsAreDefinedRemovedMultiple(t *testing.T) {
	is := is.New(t)

	foov1 := &wrappedHandler{v: "2019-01-01"}
	barv1 := &wrappedHandler{v: "2019-01-01"}
	foov2 := &wrappedHandler{v: "2019-02-02"}
	foov3 := &wrappedHandler{v: "2019-03-03"}

	rpc := Server{
		registeredVersionMethods: map[string]map[string]*wrappedHandler{
			"2019-01-01": {
				"foo": foov1,
				"bar": barv1,
			},
			"2019-02-02": {
				"foo": nil,
			},
			"2019-03-03": {
				"foo": foov2,
			},
			"2019-04-04": {
				"foo": nil,
			},
			"2019-05-05": {
				"foo": foov3,
			},
		},
	}

	rpc.buildRoutes()

	expected := map[string]map[string]*wrappedHandler{
		"2019-01-01": {
			"foo": foov1,
			"bar": barv1,
		},
		"2019-02-02": {
			"bar": barv1,
		},
		"2019-03-03": {
			"foo": foov2,
			"bar": barv1,
		},
		"2019-04-04": {
			"bar": barv1,
		},
		"2019-05-05": {
			"foo": foov3,
			"bar": barv1,
		},
		"latest": {
			"foo": foov3,
			"bar": barv1,
		},
	}

	is.Equal(expected, rpc.resolvedMethods)
}

func TestPreviewMethodsAreRegistered(t *testing.T) {
	is := is.New(t)

	barv1 := &wrappedHandler{v: "2019-01-01"}
	fooPrev := &wrappedHandler{v: "preview"}

	rpc := Server{
		registeredPreviewMethods: map[string]*wrappedHandler{
			"foo": fooPrev,
		},
		registeredVersionMethods: map[string]map[string]*wrappedHandler{
			"2019-01-01": {
				"bar": barv1,
			},
		},
	}

	rpc.buildRoutes()

	expected := map[string]map[string]*wrappedHandler{
		"preview": {
			"foo": fooPrev,
		},
		"2019-01-01": {
			"bar": barv1,
		},
		"latest": {
			"bar": barv1,
		},
	}

	is.Equal(expected, rpc.resolvedMethods)
}

func TestNilPreviewMethodsPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic when setting preview methods to nil")
		}
	}()

	zs := NewServer(UnsafeNoAuthentication)

	zs.Register("foo", "preview", nil, nil)
}

func TestPanicIfMethodDeclaredTwice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic when declaring same method, same version, twice")
		}
	}()

	zs := NewServer(UnsafeNoAuthentication)

	zs.Register("foo", "2019-01-01", nil, func(context.Context) error { return nil })

	zs.Register("foo", "2019-01-01", nil, func(context.Context) error { return nil })
}

func TestMiddlewareIsLoadedInOrder(t *testing.T) {
	ctx := context.Background()
	rpc := NewServer(UnsafeNoAuthentication)

	rpc.Register("foo", "preview", nil, makeRPCCall("called foo!"))
	rpc.Use(addHeaderMiddleware("X-Is-Test", "win!"))
	rpc.Register("bar", "preview", nil, makeRPCCall("called bar!"))

	rec := httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/preview/foo", nil)

	rpc.ServeHTTP(rec, r)

	if _, ok := rec.Result().Header["X-Is-Test"]; ok {
		t.Error("was expecting 'X-Is-Test' header to not be present")
	}

	rec = httptest.NewRecorder()
	r, _ = http.NewRequestWithContext(ctx, http.MethodPost, "/preview/bar", nil)

	rpc.ServeHTTP(rec, r)

	if _, ok := rec.Result().Header["X-Is-Test"]; !ok {
		t.Error("was expecting 'X-Is-Test' header to be present")
	}
}

func TestMiddlewareRunsGlobalInOrderAndRequestSpecific(t *testing.T) {
	ctx := context.Background()
	rpc := NewServer(UnsafeNoAuthentication)

	rpc.Use(addHeaderMiddleware("X-Present-On-Both", "win!"))
	rpc.Register("foo", "preview", nil, makeRPCCall("called foo!"))
	rpc.Use(addHeaderMiddleware("X-Present-On-Bar", "win!"))
	rpc.Register("bar", "preview", nil, makeRPCCall("called bar!"), addHeaderMiddleware("X-Also-On-Bar", "wat?"))

	rec1 := httptest.NewRecorder()
	rec2 := httptest.NewRecorder()
	r1, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/preview/foo", nil)
	r2, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/preview/bar", nil)

	rpc.ServeHTTP(rec1, r1)
	rpc.ServeHTTP(rec2, r2)

	if _, ok := rec1.Result().Header["X-Present-On-Both"]; !ok {
		t.Error("was expecting 'X-Present-On-Both' header to be present")
	}

	if _, ok := rec2.Result().Header["X-Present-On-Both"]; !ok {
		t.Error("was expecting 'X-Present-On-Both' header to be present")
	}

	if _, ok := rec1.Result().Header["X-Present-On-Bar"]; ok {
		t.Error("was expecting 'X-Present-On-Bar' header to NOT be present")
	}

	if _, ok := rec2.Result().Header["X-Present-On-Bar"]; !ok {
		t.Error("was expecting 'X-Present-On-Bar' header to be present")
	}

	if _, ok := rec1.Result().Header["X-Also-On-Bar"]; ok {
		t.Error("was expecting 'X-Also-On-Bar' header to NOT be present")
	}

	if _, ok := rec2.Result().Header["X-Also-On-Bar"]; !ok {
		t.Error("was expecting 'X-Also-On-Bar' header to be present")
	}
}

type testResponse struct {
	Message string `json:"message"`
}

func addHeaderMiddleware(headerToAdd, value string) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(res http.ResponseWriter, req *Request) error {
			res.Header().Add(headerToAdd, value)

			return next(res, req)
		}
	}
}

func makeRPCCall(messageToReturn string) func(context.Context) (*testResponse, error) {
	return func(_ context.Context) (*testResponse, error) {
		return &testResponse{
			Message: messageToReturn,
		}, nil
	}
}

func TestSchemasAreCompiled(t *testing.T) {
	brokenSchema := gojsonschema.NewStringLoader(`{
		"type": "object",
		"properties":}
	}`)
	validSchema := gojsonschema.NewStringLoader(`{
		"type": "object",
		"properties": {
			"foo": {
				"type": "string"
			}
		}
	}`)

	handler := func(_ context.Context, _ *struct{}) error {
		return nil
	}

	rpc := NewServer(UnsafeNoAuthentication)

	rpc.Register("should_pass", "2019-01-01", validSchema, handler)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should_crash method should panic")
		}
	}()

	rpc.Register("should_crash", "2019-01-01", brokenSchema, handler)
}

func UnsafeNoAuthentication(next HandlerFunc) HandlerFunc {
	return func(res http.ResponseWriter, req *Request) error {
		return next(res, req)
	}
}
