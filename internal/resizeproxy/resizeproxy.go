package resizeproxy

import (
	"context"
)

type ImageGetter interface {
	GetImage(ctx context.Context, imgRequest ImageRequest) (*ImageResponse, error)
}

type ImageResizer interface {
	ResizeImageCached(ctx context.Context, request *ResizeRequest) (*ResizeResponse, error)
}

type Cache interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
}
