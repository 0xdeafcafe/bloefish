package rpc

import (
	"context"
	"embed"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/contexts"
	"github.com/0xdeafcafe/bloefish/libraries/crpc"
	"github.com/0xdeafcafe/bloefish/libraries/crpc/middlewares"
	"github.com/0xdeafcafe/bloefish/libraries/jsonschema"
	"github.com/0xdeafcafe/bloefish/libraries/version"
	"github.com/0xdeafcafe/bloefish/services/stream"
	"github.com/0xdeafcafe/bloefish/services/stream/internal/app"
)

//go:embed *.json
var fs embed.FS
var schema = jsonschema.NewFS(fs).LoadJSONExt

// Ensure RPC implements stream.Service.
var _ stream.Service = (*RPC)(nil)

type RPC struct {
	app *app.App
}

func New(ctx context.Context, app *app.App, mux *chi.Mux) *RPC {
	rpc := &RPC{app: app}

	svcInfo := contexts.GetServiceInfo(ctx)
	if svcInfo == nil {
		panic("service info not found")
	}

	svr := crpc.NewServer(middlewares.UnsafeNoAuthentication)
	svr.Use(crpc.Logger())

	svr.Register("send_message_full", "2025-02-12", schema("send_message_full"), rpc.SendMessageFull)
	svr.Register("send_message_fragment", "2025-02-12", schema("send_message_fragment"), rpc.SendMessageFragment)
	svr.Register("send_error_message", "2025-02-12", schema("send_error_message"), rpc.SendErrorMessage)

	mux.Use(version.HeaderMiddleware(svcInfo.ServiceHTTPName))
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

	return rpc
}
