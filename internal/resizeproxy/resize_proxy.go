package resizeproxy

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

func NewResizeProxy(
	logger *zerolog.Logger,
	cache Cache,
	imageGetter ImageGetter,
	imageSupportedTypes []string,
) *ResizeProxy {
	return &ResizeProxy{
		logger:              logger,
		cache:               cache,
		imageGetter:         imageGetter,
		imageResizer:        newImageResizer(logger, imageSupportedTypes),
		imageSupportedTypes: imageSupportedTypes,
	}
}

type ResizeProxy struct {
	logger              *zerolog.Logger
	cache               Cache
	imageGetter         ImageGetter
	imageResizer        *imageResizer
	imageSupportedTypes []string
}

func (rp *ResizeProxy) ResizeImageCached(ctx context.Context, request *ResizeRequest) (*ResizeResponse, error) {
	cacheKey := request.getCachingKey()
	if result, ok := rp.cache.Get(cacheKey); ok {
		return result.(*ResizeResponse), nil
	}

	result, err := rp.ResizeImage(ctx, request)
	if err != nil {
		return nil, err
	}

	rp.cache.Set(cacheKey, result)

	return result, nil
}

func (rp *ResizeProxy) ResizeImage(ctx context.Context, request *ResizeRequest) (*ResizeResponse, error) {
	imageResponse, err := rp.imageGetter.GetImage(ctx, request.ImgRequest)
	if err != nil {
		return nil, err
	}

	resizedImg, err := rp.imageResizer.Resize(imageResponse.img, request.Width, request.Height)
	if err != nil {
		return nil, err
	}

	return &ResizeResponse{resizedImg, imageResponse.headers}, nil
}

type ResizeRequest struct {
	Width      int
	Height     int
	ImgRequest ImageRequest
}

type ResizeResponse struct {
	Img     []byte
	Headers map[string][]string
}

type ImageRequest struct {
	URL     string
	Headers map[string][]string
}

type ImageResponse struct {
	img     []byte
	headers map[string][]string
}

func NewResizeRequest(width int, height int, url string, headers map[string][]string) *ResizeRequest {
	return &ResizeRequest{
		Width:  width,
		Height: height,
		ImgRequest: ImageRequest{
			URL:     url,
			Headers: headers,
		},
	}
}

func (rp *ResizeRequest) getCachingKey() string {
	return fmt.Sprintf("%d%d%s", rp.Width, rp.Height, rp.ImgRequest.URL)
}
