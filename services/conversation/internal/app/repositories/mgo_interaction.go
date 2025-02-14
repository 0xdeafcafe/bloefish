package repositories

import (
	"context"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/ksuid"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/models"
	"github.com/0xdeafcafe/bloefish/services/conversation/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type persistedInteraction struct {
	ID             string   `bson:"_id"`
	IdempotencyKey string   `bson:"idempotency_key"`
	ConversationID string   `bson:"conversation_id"`
	MessageContent string   `bson:"message_content"`
	FileIDs        []string `bson:"file_ids"`

	Owner struct {
		Type       string `bson:"type"`
		Identifier string `bson:"identifier"`
	} `bson:"owner"`
	AIRelayOptions struct {
		ProviderID string `bson:"provider_id"`
		ModelID    string `bson:"model_id"`
	} `bson:"ai_relay_options"`

	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at"`
}

type mgoInteraction struct {
	c *mongo.Collection
}

func NewMgoInteraction(db *mongo.Database) ports.InteractionRepository {
	return &mgoInteraction{c: db.Collection("interactions")}
}

func (r *mgoInteraction) Create(ctx context.Context, cmd *models.CreateInteractionCommand) (*models.Interaction, error) {
	result := r.c.FindOneAndUpdate(ctx, bson.M{
		"idempotency_key":  cmd.IdempotencyKey,
		"conversation_id":  cmd.ConversationID,
		"owner.type":       cmd.Owner.Type,
		"owner.identifier": cmd.Owner.Identifier,
	}, bson.M{
		"$currentDate": bson.M{
			"updated_at": true,
		},
		"$setOnInsert": bson.M{
			"_id":             ksuid.Generate(ctx, "interaction").String(),
			"idempotency_key": cmd.IdempotencyKey,
			"conversation_id": cmd.ConversationID,
			"message_content": cmd.MessageContent,
			"file_ids":        cmd.FileIDs,

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

	var interaction *persistedInteraction
	if err := result.Decode(&interaction); err != nil {
		return nil, err
	}

	return interaction.ToDomainModel(), nil
}

func (r *mgoInteraction) GetByID(ctx context.Context, interactionID string) (*models.Interaction, error) {
	result := r.c.FindOne(ctx, bson.M{
		"_id":        interactionID,
		"deleted_at": nil,
	})

	var interaction *persistedInteraction
	if err := result.Decode(&interaction); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, cher.New("interaction_not_found", cher.M{"interaction_id": interactionID})
		}

		return nil, err
	}

	return interaction.ToDomainModel(), nil
}

func (r *mgoInteraction) GetAllByConversationID(ctx context.Context, conversationID string) ([]*models.Interaction, error) {
	cursor, err := r.c.Find(ctx, bson.M{
		"conversation_id": conversationID,
		"deleted_at":      nil,
	}, options.Find().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var persistedInteractions []*persistedInteraction
	if err := cursor.All(ctx, &persistedInteractions); err != nil {
		return nil, err
	}

	interactions := make([]*models.Interaction, len(persistedInteractions))
	for i, persistedInteraction := range persistedInteractions {
		interactions[i] = persistedInteraction.ToDomainModel()
	}

	return interactions, nil
}

func (p *persistedInteraction) ToDomainModel() *models.Interaction {
	return &models.Interaction{
		ID:             p.ID,
		IdempotencyKey: p.IdempotencyKey,
		ConversationID: p.ConversationID,
		MessageContent: p.MessageContent,
		FileIDs:        p.FileIDs,

		Owner: &models.Actor{
			Type:       models.ActorType(p.Owner.Type),
			Identifier: p.Owner.Identifier,
		},
		AIRelayOptions: &models.AIRelayOptions{
			ProviderID: p.AIRelayOptions.ProviderID,
			ModelID:    p.AIRelayOptions.ModelID,
		},

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
	}
}
