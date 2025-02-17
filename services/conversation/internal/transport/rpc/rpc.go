package rpc

import (
	"context"
	"embed"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/libraries/contexts"
	"github.com/0xdeafcafe/bloefish/libraries/crpc"
	"github.com/0xdeafcafe/bloefish/libraries/crpc/middlewares"
	"github.com/0xdeafcafe/bloefish/libraries/jsonschema"
	"github.com/0xdeafcafe/bloefish/libraries/version"
	"github.com/0xdeafcafe/bloefish/services/conversation"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/app"
)

//go:embed *.json
var fs embed.FS
var schema = jsonschema.NewFS(fs).LoadJSONExt

// Ensure RPC implements fileupload.Service.
var _ conversation.Service = (*RPC)(nil)

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

	svr.Register("create_conversation", "2025-02-12", schema("create_conversation"), rpc.CreateConversation)
	svr.Register("create_conversation_message", "2025-02-12", schema("create_conversation_message"), rpc.CreateConversationMessage)
	svr.Register("get_interaction", "2025-02-12", schema("get_interaction"), rpc.GetInteraction)
	svr.Register("get_conversation_with_interactions", "2025-02-12", schema("get_conversation_with_interactions"), rpc.GetConversationWithInteractions)
	svr.Register("list_conversations_with_interactions", "2025-02-12", schema("list_conversations_with_interactions"), rpc.ListConversationsWithInteractions)

	mux := chi.NewRouter()
	mux.Use(version.HeaderMiddleware(svcInfo.ServiceHTTPName))
	mux.Get("/system/health", middlewares.HealthCheck)

	mux.
		With(
			cors.Handler(cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"POST", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300,
			}),
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
