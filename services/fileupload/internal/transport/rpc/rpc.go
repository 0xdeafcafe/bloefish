package rpc

import (
	"context"
	"embed"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/riandyrn/otelchi"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/contexts"
	"github.com/0xdeafcafe/bloefish/libraries/crpc"
	"github.com/0xdeafcafe/bloefish/libraries/crpc/middlewares"
	"github.com/0xdeafcafe/bloefish/libraries/jsonschema"
	"github.com/0xdeafcafe/bloefish/libraries/version"
	"github.com/0xdeafcafe/bloefish/services/fileupload"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/app"
)

//go:embed *.json
var fs embed.FS
var schema = jsonschema.NewFS(fs).LoadJSONExt

// Ensure RPC implements fileupload.Service.
var _ fileupload.Service = (*RPC)(nil)

type RPC struct {
	app *app.App

	httpServer *http.Server
}

func New(ctx context.Context, app *app.App) *RPC {
	rpc := &RPC{app: app}

	svcInfo := contexts.GetServiceInfo(ctx)
	if svcInfo == nil {
		panic("service info not found")
	}

	svr := crpc.NewServer(middlewares.UnsafeNoAuthentication)
	svr.Use(crpc.Logger())

	svr.Register("create_upload", "2025-02-12", schema("create_upload"), rpc.CreateUpload)
	svr.Register("confirm_upload", "2025-02-12", schema("confirm_upload"), rpc.ConfirmUpload)
	svr.Register("get_file", "2025-02-12", schema("get_file"), rpc.GetFile)

	mux := chi.NewRouter()
	mux.Use(version.HeaderMiddleware(svcInfo.ServiceHTTPName))
	mux.Use(otelchi.Middleware(svcInfo.ServiceHTTPName, otelchi.WithChiRoutes(mux)))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mux.Get("/system/health", middlewares.HealthCheck)

	mux.
		With(
			middlewares.StripPrefix("/rpc"),
			middlewares.RequestID,
			middlewares.Logger(clog.Get(ctx)),
		).
		Handle("/rpc/*", svr)

	rpc.httpServer = &http.Server{
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	return rpc
}

func (r *RPC) Run(ctx context.Context, cfg config.Server) (err error) {
	clog.Get(ctx).WithField("addr", cfg.Addr).Info("listening")

	if err = cfg.ListenAndServe(r.httpServer); err != nil {
		err = errors.Wrap(err, "server:")
	}

	return
}
