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

	Title *string `bson:"title"`

	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at"`
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

			"title": nil,

			"created_at": time.Now(),
			"updated_at": nil,
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

func (r *mgoConversation) ListByOwner(ctx context.Context, actor models.Actor) ([]*models.Conversation, error) {
	cursor, err := r.c.Find(ctx, bson.M{
		"owner.type":       actor.Type,
		"owner.identifier": actor.Identifier,
		"deleted_at":       nil,
	}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var persistedConversations []*persistedConversation
	if err := cursor.All(ctx, &persistedConversations); err != nil {
		return nil, err
	}

	conversations := make([]*models.Conversation, len(persistedConversations))
	for i, persistedConversation := range persistedConversations {
		conversations[i] = persistedConversation.ToDomainModel()
	}

	return conversations, nil
}

func (r *mgoConversation) DeleteMany(ctx context.Context, conversationIDs []string) error {
	if _, err := r.c.UpdateMany(ctx, bson.M{
		"_id": bson.M{"$in": conversationIDs},
	}, bson.M{
		"$currentDate": bson.M{
			"deleted_at": true,
			"updated_at": true,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (r *mgoConversation) UpdateTitle(ctx context.Context, conversationID string, title string) error {
	_, err := r.c.UpdateOne(ctx, bson.M{
		"_id": conversationID,
	}, bson.M{
		"$currentDate": bson.M{
			"updated_at": true,
		},
		"$set": bson.M{
			"title": title,
		},
	}, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	return nil
}

func (p *persistedConversation) ToDomainModel() *models.Conversation {
	return &models.Conversation{
		ID:             p.ID,
		IdempotencyKey: p.IdempotencyKey,
		Owner: &models.Actor{
			Type:       models.ActorType(p.Owner.Type),
			Identifier: p.Owner.Identifier,
		},
		AIRelayOptions: &models.AIRelayOptions{
			ProviderID: p.AIRelayOptions.ProviderID,
			ModelID:    p.AIRelayOptions.ModelID,
		},

		Title: p.Title,

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
	}
}
