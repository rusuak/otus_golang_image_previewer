package server

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/rusuak/otus_golang_image_previewer/internal/resizeproxy"
	"github.com/stretchr/testify/require"
)

var defaultImgURL = "https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/"

func TestResizeHandlerSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	mockImageResizer := resizeproxy.NewMockImageResizer(ctrl)
	logger := log.With().Logger()

	image := loadImage("_gopher_original_1024x504.jpg")

	width := 500
	height := 600
	url := defaultImgURL + "_gopher_original_1024x504.jpg"
	response := string(image)
	resizeResponse := &resizeproxy.ResizeResponse{Img: image}

	t.Run("ok_case", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
		req = mux.SetURLVars(req, map[string]string{
			"width":    strconv.Itoa(width),
			"height":   strconv.Itoa(height),
			"imageURL": url,
		})

		request := resizeproxy.NewResizeRequest(width, height, url, req.Header)
		mockImageResizer.EXPECT().ResizeImageCached(req.Context(), request).Return(resizeResponse, nil)

		h := NewHandler(&logger, mockImageResizer)

		w := httptest.NewRecorder()

		h.ResizeHandler(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		if w.Body.String() != response {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), response)
		}
	})
}

func TestResizeHandlerNegative(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	mockImageResizer := resizeproxy.NewMockImageResizer(ctrl)
	logger := log.With().Logger()

	tests := []struct {
		name           string
		width          string
		height         string
		url            string
		response       string
		resizeResponse *resizeproxy.ResizeResponse
		err            error
		httpStatus     int
	}{
		{
			name:       "bad_request_case",
			width:      "foo",
			height:     "bar",
			url:        defaultImgURL + "_gopher_original_1024x504.jpg",
			response:   "validation error",
			httpStatus: http.StatusBadRequest,
		},
		{
			name:           "bad_gateway_case",
			width:          "300",
			height:         "400",
			url:            defaultImgURL + "_gopher_original_1024x504.jpg",
			response:       "resize image failed",
			resizeResponse: nil,
			httpStatus:     http.StatusBadGateway,
			err:            errors.New("error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
			req = mux.SetURLVars(req, map[string]string{
				"width":    tt.width,
				"height":   tt.height,
				"imageURL": tt.url,
			})

			if tt.resizeResponse != nil || tt.err != nil {
				width, err := strconv.Atoi(tt.width)
				if err != nil {
					t.Errorf("error converted width to int")
				}

				height, err := strconv.Atoi(tt.height)
				if err != nil {
					t.Errorf("error converted height to int")
				}

				request := resizeproxy.NewResizeRequest(width, height, tt.url, req.Header)
				mockImageResizer.EXPECT().ResizeImageCached(req.Context(), request).Return(tt.resizeResponse, tt.err)
			}

			h := NewHandler(&logger, mockImageResizer)

			w := httptest.NewRecorder()

			h.ResizeHandler(w, req)

			if status := w.Code; status == http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			if w.Body.String() != tt.response {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.response)
			}
		})
	}
}

func TestResizeHandlerProxyHeaders(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	mockImageResizer := resizeproxy.NewMockImageResizer(ctrl)
	logger := log.With().Logger()

	image := loadImage("gopher_256x126_resized.jpg")

	headers := map[string][]string{
		"Content-Length": {0: "6495"},
		"Content-Type":   {0: "image/jpeg"},
	}

	width := 200
	height := 300
	url := defaultImgURL + "_gopher_original_1024x504.jpg"
	fillResponse := &resizeproxy.ResizeResponse{Img: image, Headers: headers}

	t.Run("good headers", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
		req = mux.SetURLVars(req, map[string]string{
			"width":    strconv.Itoa(width),
			"height":   strconv.Itoa(height),
			"imageURL": url,
		})

		fillParams := resizeproxy.NewResizeRequest(
			width,
			height,
			url,
			req.Header,
		)
		mockImageResizer.EXPECT().
			ResizeImageCached(req.Context(), fillParams).
			Return(fillResponse, nil)

		h := NewHandler(&logger, mockImageResizer)

		w := httptest.NewRecorder()

		h.ResizeHandler(w, req)

		for name, values := range fillResponse.Headers {
			for _, value := range values {
				require.Equal(t, value, w.Header().Get(name))
			}
		}
	})
}

func loadImage(imgName string) []byte {
	fileToBeUploaded := "../../img_example/" + imgName
	file, err := os.Open(fileToBeUploaded)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileInfo, _ := file.Stat()
	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	return bytes
}
