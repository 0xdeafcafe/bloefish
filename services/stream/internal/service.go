package internal

import (
	"context"
	"net"
	"net/http"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/libraries/config"
	"github.com/0xdeafcafe/bloefish/services/stream/internal/app"
	"github.com/0xdeafcafe/bloefish/services/stream/internal/app/services"
	"github.com/0xdeafcafe/bloefish/services/stream/internal/transport/rpc"
	"github.com/0xdeafcafe/bloefish/services/stream/internal/transport/ws"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Config struct {
	Server  config.Server `env:"SERVER"`
	Logging clog.Config   `env:"LOGGING"`
}

func defaultConfig() Config {
	return Config{
		Server: config.Server{
			Addr: ":4005",
		},

		Logging: clog.Config{
			Format: clog.TextFormat,
			Debug:  true,
		},
	}
}

func Run(ctx context.Context) error {
	cfg := defaultConfig()
	config.MustHydrateFromEnvironment(ctx, &cfg)
	ctx = clog.Set(ctx, cfg.Logging.Configure(ctx))

	app := &app.App{
		MessageBroker: services.NewWebSocketMessageBroker(),
	}

	mux := chi.NewRouter()
	_ = rpc.New(ctx, app, mux)
	_ = ws.New(ctx, app, mux)

	clog.Get(ctx).WithField("addr", cfg.Server.Addr).Info("listening")
	if err := cfg.Server.ListenAndServe(&http.Server{
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}); err != nil {
		return errors.Wrap(err, "server:")
	}

	return nil
}
