package resizeproxy

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetImageSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	imgName := SourceImgName

	imageClient := NewImageClient()

	t.Run("success_download_img", func(t *testing.T) {
		t.Parallel()

		imgRequest := ImageRequest{ImageURL + imgName, nil}
		gotImg, err := imageClient.GetImage(ctx, imgRequest)
		if err != nil {
			t.Errorf("Download() error = %v", err)
			return
		}

		wantImg := loadImage(imgName)
		if !reflect.DeepEqual(gotImg.img, wantImg) {
			t.Errorf("Download() gotImg = %v, want %v", gotImg.img, wantImg)
		}
	})
}

func TestGetImageError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	imgName := ""
	errExpected := ErrNotAcceptableResponseStatus

	imageClient := NewImageClient()

	t.Run("bad_response_case", func(t *testing.T) {
		t.Parallel()

		imgRequest := ImageRequest{ImageURL + imgName, nil}
		_, err := imageClient.GetImage(ctx, imgRequest)
		require.Errorf(t, err, errExpected.Error())
	})
}
