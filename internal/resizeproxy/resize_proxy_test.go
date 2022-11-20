package resizeproxy

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

func TestDefaultServiceFill(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	imageOrigin := loadImage(SourceImgName)
	imageResized := loadImage(DestinationImgName)
	imageResponse := &ImageResponse{img: imageOrigin}

	logger := log.With().Logger()

	mockCache := NewMockCache(ctrl)
	mockImageClient := NewMockImageGetter(ctrl)

	ctx := context.Background()
	wantErr := error(nil)
	resizeRequest := NewResizeRequest(256, 126, ImageURL+SourceImgName, nil)
	expected := &ResizeResponse{imageResized, nil}

	t.Run("success_resized", func(t *testing.T) {
		t.Parallel()

		resizeProxy := NewResizeProxy(
			&logger,
			mockCache,
			mockImageClient,
			getTestImagesSupportedTypes(),
		)

		mockImageClient.EXPECT().GetImage(
			ctx,
			resizeRequest.ImgRequest,
		).Return(imageResponse, nil)

		cacheKey := resizeRequest.getCachingKey()
		mockCache.EXPECT().Get(cacheKey).Return(nil, false)
		mockCache.EXPECT().Set(cacheKey, expected).Return(true)

		_, err := resizeProxy.ResizeImageCached(ctx, resizeRequest)
		if !errors.Is(err, wantErr) {
			t.Errorf("Fill() error = %v, wantErr %v", err, wantErr)
			return
		}
	})
}

func TestResizeImageInCache(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	logger := log.With().Logger()
	ctx := context.Background()
	mockCache := NewMockCache(ctrl)

	resizeProxy := NewResizeProxy(
		&logger,
		mockCache,
		NewImageClient(),
		getTestImagesSupportedTypes(),
	)

	url := ImageURL + SourceImgName
	request := NewResizeRequest(
		256,
		126,
		url,
		nil,
	)

	destImg := loadImage(DestinationImgName)
	resizeResponse := &ResizeResponse{
		Img:     destImg,
		Headers: nil,
	}

	t.Run("get_resized_img_from_cache_case", func(t *testing.T) {
		t.Parallel()

		cacheKey := request.getCachingKey()
		mockCache.EXPECT().Get(cacheKey).Return(resizeResponse, true)

		gotImg, err := resizeProxy.ResizeImageCached(ctx, request)
		if err != nil {
			t.Errorf("Resize() error = %v", err)
			return
		}

		wantImg := loadImage(DestinationImgName)
		if !reflect.DeepEqual(gotImg.Img, wantImg) {
			t.Errorf("Resize() gotImg = %v, want %v", gotImg.Img, wantImg)
		}
	})
}
