package ws

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"github.com/0xdeafcafe/bloefish/libraries/clog"
	"github.com/0xdeafcafe/bloefish/services/stream/internal/app"
)

type WS struct {
	app *app.App
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO(afr): Change this to validate something later?
	},
}

func New(ctx context.Context, app *app.App, mux *chi.Mux) *WS {
	mux.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			clog.Get(ctx).WithError(err).Error("unable to upgrade websocket connection")
			return
		}

		app.MessageBroker.RegisterConnection(ctx, conn)

		if err := conn.SetReadDeadline(time.Now().Add(10 * time.Minute)); err != nil {
			clog.Get(ctx).WithError(err).Warn("unable to set read deadline on websocket connection")
		}

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				clog.Get(ctx).WithError(err).Info("websocket connection closed due to read error")
				break
			}
		}
	})

	return &WS{app: app}
}
