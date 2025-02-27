package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/libraries/ksuid"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/domain/models"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type persistedFile struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	Size     int64  `bson:"size"`
	MIMEType string `bson:"mime_type"`

	Owner struct {
		Type       string `bson:"type"`
		Identifier string `bson:"identifier"`
	} `bson:"owner"`

	CreatedAt   time.Time  `bson:"created_at"`
	UpdatedAt   *time.Time `bson:"updated_at"`
	ConfirmedAt *time.Time `bson:"confirmed_at"`
	DeletedAt   *time.Time `bson:"deleted_at"`
}

type mgoFile struct {
	c *mongo.Collection
}

func NewMgoFile(db *mongo.Database) ports.FileRepository {
	return &mgoFile{c: db.Collection("files")}
}

func (m *mgoFile) CreateUpload(ctx context.Context, req *models.CreateUploadCommand) (string, error) {
	id := ksuid.Generate(ctx, "file").String()

	result, err := m.c.InsertOne(ctx, bson.M{
		"_id":       id,
		"name":      req.Name,
		"size":      req.Size,
		"mime_type": req.MIMEType,
		"owner": bson.M{
			"type":       req.Owner.Type,
			"identifier": req.Owner.Identifier,
		},
		"created_at":   time.Now(),
		"updated_at":   nil,
		"confirmed_at": nil,
		"deleted_at":   nil,
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

func (m *mgoFile) ConfirmUpload(ctx context.Context, fileID string) (*models.File, error) {
	result := m.c.FindOneAndUpdate(ctx, bson.M{
		"_id": fileID,
	}, bson.M{
		"$currentDate": bson.M{
			"updated_at": true,
		},
		"$set": bson.M{
			"confirmed_at": time.Now(),
		},
	}, options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After))

	var file *persistedFile
	if err := result.Decode(&file); err != nil {
		return nil, err
	}

	return file.ToDomainModel(), nil
}

func (m *mgoFile) Get(ctx context.Context, fileID string) (*models.File, error) {
	result := m.c.FindOne(ctx, bson.M{
		"_id": fileID,
	})
	if err := result.Err(); err != nil {
		return nil, err
	}

	var file *persistedFile
	if err := result.Decode(&file); err != nil {
		return nil, err
	}

	return file.ToDomainModel(), nil
}

func (m *mgoFile) GetMany(ctx context.Context, ids []string) ([]*models.File, error) {
	cursor, err := m.c.Find(ctx, bson.M{
		"_id": bson.M{"$in": ids},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []*persistedFile
	if err := cursor.All(ctx, &files); err != nil {
		return nil, err
	}

	if len(files) != len(ids) {
		foundIDs := make(map[string]bool)
		for _, file := range files {
			foundIDs[file.ID] = true
		}

		missingIDs := make([]string, 0)
		for _, id := range ids {
			if !foundIDs[id] {
				missingIDs = append(missingIDs, id)
			}
		}

		return nil, cher.New("files_not_found", cher.M{
			"missing_file_ids": missingIDs,
		})
	}

	filesDomain := make([]*models.File, len(files))
	for i, file := range files {
		filesDomain[i] = file.ToDomainModel()
	}

	return filesDomain, nil
}

func (p *persistedFile) ToDomainModel() *models.File {
	return &models.File{
		ID:       p.ID,
		Name:     p.Name,
		Size:     p.Size,
		MIMEType: p.MIMEType,
		Owner: &models.Actor{
			Type:       models.ActorType(p.Owner.Type),
			Identifier: p.Owner.Identifier,
		},
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		ConfirmedAt: p.ConfirmedAt,
		DeletedAt:   p.DeletedAt,
	}
}
