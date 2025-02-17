package app

import (
	"github.com/0xdeafcafe/bloefish/services/airelay"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/ports"
	"github.com/0xdeafcafe/bloefish/services/stream"
	"github.com/0xdeafcafe/bloefish/services/user"
)

type App struct {
	ConversationRepository ports.ConversationRepository
	InteractionRepository  ports.InteractionRepository

	AIRelayService airelay.Service
	StreamService  stream.Service
	UserService    user.Service
}
