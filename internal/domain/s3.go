package domain

import (
	"context"
	"io"
)

type S3Adapter interface {
	Save(ctx context.Context, objName, contentType string, data io.Reader, size int64) (string, error)
	GetPublicLink(objName string) string
}
