package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/ksuid"
	"github.com/0xdeafcafe/bloefish/services/skillset/internal/domain/models"
	"github.com/0xdeafcafe/bloefish/services/skillset/internal/domain/ports"
)

type persistedSkillSet struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Icon        string `bson:"icon"`
	Description string `bson:"description"`
	Prompt      string `bson:"prompt"`

	Owner struct {
		Type       string `bson:"type"`
		Identifier string `bson:"identifier"`
	} `bson:"owner"`

	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at"`
}

type mgoSkillSet struct {
	c *mongo.Collection
}

func NewMgoSkillSet(db *mongo.Database) ports.SkillSetRepository {
	return &mgoSkillSet{c: db.Collection("skill_sets")}
}

func (m *mgoSkillSet) CreateSkillSet(ctx context.Context, req models.CreateSkillSetCommand) (string, error) {
	id := ksuid.Generate(ctx, "skillset").String()

	result, err := m.c.InsertOne(ctx, bson.M{
		"_id":         id,
		"name":        req.Name,
		"icon":        req.Icon,
		"description": req.Description,
		"prompt":      req.Prompt,

		"owner": bson.M{
			"type":       req.Owner.Type,
			"identifier": req.Owner.Identifier,
		},

		"created_at": time.Now(),
		"updated_at": nil,
		"deleted_at": nil,
	})
	if err != nil {
		return "", err
	}

	id, ok := result.InsertedID.(string)
	if !ok {
		return "", errors.New("failed to cast InsertedID to string")
	}

	return id, nil
}

func (m *mgoSkillSet) GetSkillSet(ctx context.Context, id string) (*models.SkillSet, error) {
	var p persistedSkillSet
	if err := m.c.FindOne(ctx, bson.M{"_id": id}).Decode(&p); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, cher.New("skill_set_not_found", cher.M{
				"skill_set_id": id,
			})
		}
	}

	return p.ToDomainModel(), nil
}

func (m *mgoSkillSet) ListSkillSetsByOwner(ctx context.Context, ownerType, ownerIdentifier string) ([]*models.SkillSet, error) {
	cursor, err := m.c.Find(ctx, bson.M{
		"owner.type":       ownerType,
		"owner.identifier": ownerIdentifier,
		"deleted_at":       nil,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var skillSets []*persistedSkillSet
	if err := cursor.All(ctx, &skillSets); err != nil {
		return nil, err
	}

	skillSetsDomain := make([]*models.SkillSet, len(skillSets))
	for i, skillSet := range skillSets {
		skillSetsDomain[i] = skillSet.ToDomainModel()
	}

	return skillSetsDomain, nil
}

func (m *mgoSkillSet) GetManySkillSets(ctx context.Context, ids []string) ([]*models.SkillSet, error) {
	cursor, err := m.c.Find(ctx, bson.M{
		"_id": bson.M{"$in": ids},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var skillSets []*persistedSkillSet
	if err := cursor.All(ctx, &skillSets); err != nil {
		return nil, err
	}

	if len(skillSets) != len(ids) {
		foundIDs := make(map[string]bool)
		for _, skillSet := range skillSets {
			foundIDs[skillSet.ID] = true
		}

		missingIDs := make([]string, 0)
		for _, id := range ids {
			if !foundIDs[id] {
				missingIDs = append(missingIDs, id)
			}
		}

		return nil, cher.New("skill_sets_not_found", cher.M{
			"missing_skill_set_ids": missingIDs,
		})
	}

	skillSetsDomain := make([]*models.SkillSet, len(skillSets))
	for i, skillSet := range skillSets {
		skillSetsDomain[i] = skillSet.ToDomainModel()
	}

	return skillSetsDomain, nil
}

func (p *persistedSkillSet) ToDomainModel() *models.SkillSet {
	return &models.SkillSet{
		ID:          p.ID,
		Name:        p.Name,
		Icon:        p.Icon,
		Description: p.Description,
		Prompt:      p.Prompt,

		Owner: &models.Actor{
			Type:       models.ActorType(p.Owner.Type),
			Identifier: p.Owner.Identifier,
		},

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
	}
}
