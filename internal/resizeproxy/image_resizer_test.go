package resizeproxy

import (
	"context"
	"reflect"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		ctx         context.Context
		originalImg []byte
		resizedImg  []byte
		width       int
		height      int
		err         error
	}{
		{
			name:        "success resize 256x126",
			ctx:         ctx,
			width:       256,
			height:      126,
			originalImg: loadImage("_gopher_original_1024x504.jpg"),
			resizedImg:  loadImage("gopher_256x126_resized.jpg"),
			err:         nil,
		},
		{
			name:        "not_allowed_type_img_case",
			ctx:         ctx,
			width:       256,
			height:      126,
			originalImg: loadImage("gopher.png"),
			resizedImg:  nil,
			err:         ErrNotSupportedType,
		},
	}

	logger := log.With().Logger()
	ir := newImageResizer(&logger, getTestImagesSupportedTypes())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotImg, err := ir.Resize(tt.originalImg, tt.width, tt.height)
			if err != nil {
				require.Nil(t, gotImg)
				require.Errorf(t, err, tt.err.Error())
				return
			}

			if !reflect.DeepEqual(gotImg, tt.resizedImg) {
				t.Errorf("Resize()  gotImg = %v, want %v", gotImg, tt.resizedImg)
			}
		})
	}
}
