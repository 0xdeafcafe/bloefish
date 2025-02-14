package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/ksuid"
	"github.com/0xdeafcafe/bloefish/services/user/internal/domain/models"
	"github.com/0xdeafcafe/bloefish/services/user/internal/domain/ports"
)

var (
	ErrCodeUserNotFound = "user_not_found"

	userKSUIDResource = "user"
)

type persistedUser struct {
	ID string `bson:"_id"`

	DefaultUser bool `bson:"default_user"`

	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at"`
}

type mgoUser struct {
	c *mongo.Collection
}

func NewMgoUser(db *mongo.Database) ports.UserRepository {
	return &mgoUser{
		c: db.Collection("users"),
	}
}

func (r *mgoUser) GetByUserID(ctx context.Context, userID string) (*models.User, error) {
	var user *persistedUser
	if err := r.c.FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, cher.New(ErrCodeUserNotFound, cher.M{"user_id": userID})
		}

		return nil, err
	}

	return user.ToUser(), nil
}

func (r *mgoUser) GetOrCreateDefaultUser(ctx context.Context) (*models.User, error) {
	result := r.c.FindOneAndUpdate(ctx, bson.M{
		"default_user": true,
	}, bson.M{
		"$setOnInsert": bson.M{
			"_id":          ksuid.Generate(ctx, userKSUIDResource).String(),
			"default_user": true,
			"created_at":   time.Now(),
			"deleted_at":   nil,
		},
	}, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))

	var user *persistedUser
	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	return user.ToUser(), nil
}

func (u *persistedUser) ToUser() *models.User {
	return &models.User{
		ID:          u.ID,
		DefaultUser: u.DefaultUser,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
	}
}
