package services

import (
	"context"
	"fmt"
	"mime"
	"time"

	"github.com/0xdeafcafe/bloefish/libraries/cher"
	"github.com/0xdeafcafe/bloefish/services/fileupload/internal/domain/ports"
	"github.com/minio/minio-go/v7"
)

var (
	minioPresignedUploadExpiry = 15 * time.Minute
)

type MinioFileObject struct {
	client     *minio.Client
	bucketName string
}

func NewMinioFileObject(client *minio.Client, bucketName string) ports.FileObjectService {
	return &MinioFileObject{
		client:     client,
		bucketName: bucketName,
	}
}

func (m *MinioFileObject) CreatePresignedUploadURL(ctx context.Context, fileID string) (string, error) {
	presignedURL, err := m.client.PresignedPutObject(ctx, m.bucketName, fileID, minioPresignedUploadExpiry)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (m *MinioFileObject) CreatePresignedDownloadURL(ctx context.Context, fileID string, expiry time.Duration) (string, error) {
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, fileID, expiry, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (m *MinioFileObject) EnsureFileUploadValidity(ctx context.Context, fileID, mimeType string, size int) error {
	object, err := m.client.StatObject(ctx, m.bucketName, fileID, minio.StatObjectOptions{})
	if err != nil {
		return err
	}

	if object.Size != int64(size) {
		return cher.New("file_size_mismatch", cher.M{
			"expected_size": size,
			"actual_size":   object.Size,
		})
	}

	mediaType, _, err := mime.ParseMediaType(object.ContentType)
	if err != nil {
		return fmt.Errorf("failed to parse content type: %w", err)
	}

	if mediaType != mimeType {
		return cher.New("file_mime_type_mismatch", cher.M{
			"expected_mime_type": mimeType,
			"actual_mime_type":   mediaType,
		})
	}

	return nil
}
