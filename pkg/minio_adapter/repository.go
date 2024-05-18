package minio_adapter

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
	"io"
)

var ErrSendingObject = fmt.Errorf("error while sending object to storage")

// New creates a new MinIO file storage Adapter.
func New(logger *zerolog.Logger, config *MinioConfig, secure bool) (*Adapter, error) {
	adapter := &Adapter{
		logger: logger,
		config: config,
	}

	client, err := minio.New(config.ServerURL, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: secure,
	})

	if err != nil {
		return nil, err
	}
	adapter.client = client

	return adapter, nil
}

type Adapter struct {
	logger *zerolog.Logger
	config *MinioConfig
	client *minio.Client
}

// Save saves file in the MinIO and returns created object's name.
func (a *Adapter) Save(ctx context.Context, objName, contentType string, data io.Reader, size int64) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := a.client.PutObject(ctx, a.config.BucketName, objName, data, size, opts)
	if err != nil {
		a.logger.Error().Err(err).Msg("Error saving a file in the MinIO!")
		return "", ErrSendingObject
	}
	a.logger.Debug().Any("created obj info", info)

	return info.Key, nil

}

func (a *Adapter) Remove(ctx context.Context, objName string) error {
	err := a.client.RemoveObject(ctx, a.config.BucketName, objName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
