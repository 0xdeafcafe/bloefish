package repositories

import (
	"context"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/ksuid"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/models"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type persistedConversation struct {
	ID             string `bson:"_id"`
	IdempotencyKey string `bson:"idempotency_key"`

	Owner struct {
		Type       string `bson:"type"`
		Identifier string `bson:"identifier"`
	} `bson:"owner"`
	AIRelayOptions struct {
		ProviderID string `bson:"provider_id"`
		ModelID    string `bson:"model_id"`
	} `bson:"ai_relay_options"`

	CreatedAt time.Time  `bson:"created_at"`
	DeletedAt *time.Time `bson:"deleted_at"`
}

type mgoConversation struct {
	c *mongo.Collection
}

func NewMgoConversation(db *mongo.Database) ports.ConversationRepository {
	return &mgoConversation{c: db.Collection("conversations")}
}

func (r *mgoConversation) Create(ctx context.Context, cmd *models.CreateConversationCommand) (*models.Conversation, error) {
	result := r.c.FindOneAndUpdate(ctx, bson.M{
		"idempotency_key":  cmd.IdempotencyKey,
		"owner.type":       cmd.Owner.Type,
		"owner.identifier": cmd.Owner.Identifier,
	}, bson.M{
		"$setOnInsert": bson.M{
			"_id":             ksuid.Generate(ctx, "conversation").String(),
			"idempotency_key": cmd.IdempotencyKey,
			"owner": bson.M{
				"type":       cmd.Owner.Type,
				"identifier": cmd.Owner.Identifier,
			},
			"ai_relay_options": bson.M{
				"provider_id": cmd.AIRelayOptions.ProviderID,
				"model_id":    cmd.AIRelayOptions.ModelID,
			},
			"created_at": time.Now(),
			"deleted_at": nil,
		},
	}, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))

	var conversation *persistedConversation
	if err := result.Decode(&conversation); err != nil {
		return nil, err
	}

	return conversation.ToDomainModel(), nil
}

func (r *mgoConversation) GetByID(ctx context.Context, conversationID string) (*models.Conversation, error) {
	result := r.c.FindOne(ctx, bson.M{"_id": conversationID})
	if err := result.Err(); err != nil {
		return nil, err
	}

	var conversation *persistedConversation
	if err := result.Decode(&conversation); err != nil {
		return nil, err
	}

	return conversation.ToDomainModel(), nil
}

func (p *persistedConversation) ToDomainModel() *models.Conversation {
	return &models.Conversation{
		ID:             p.ID,
		IdempotencyKey: p.IdempotencyKey,
		Owner: &models.Actor{
			Type:       models.ActorType(p.Owner.Type),
			Identifier: p.Owner.Identifier,
		},
		AIRelayOptions: &models.ConversationAIRelayOptions{
			ProviderID: p.AIRelayOptions.ProviderID,
			ModelID:    p.AIRelayOptions.ModelID,
		},
		CreatedAt: p.CreatedAt,
		DeletedAt: p.DeletedAt,
	}
}
