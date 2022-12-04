package resizeproxy

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/rs/zerolog"
)

var ErrNotSupportedType = fmt.Errorf("not supported image type")

func newImageResizer(logger *zerolog.Logger, imageSupportedTypes []string) *imageResizer {
	return &imageResizer{logger: logger, imageSupportedTypes: imageSupportedTypes}
}

type imageResizer struct {
	logger              *zerolog.Logger
	imageSupportedTypes []string
}

func (ir *imageResizer) Resize(sourceImg []byte, width int, height int) ([]byte, error) {
	if !ir.isImageTypeSupported(http.DetectContentType(sourceImg)) {
		return nil, ErrNotSupportedType
	}

	newImgName := fmt.Sprintf("image_%d_%d.jpg", width, height)

	file, err := os.Create(newImgName)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			ir.logger.Err(err).Msg(err.Error())
		}
	}(file)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			ir.logger.Err(err).Msg(err.Error())
		}
	}(newImgName)

	_, err = io.Copy(file, bytes.NewReader(sourceImg))
	if err != nil {
		return nil, err
	}

	src, err := imaging.Open(newImgName)
	if err != nil {
		return nil, err
	}

	img := imaging.Resize(src, width, height, imaging.Lanczos)

	imgBuffer := new(bytes.Buffer)
	err = jpeg.Encode(imgBuffer, img, nil)
	if err != nil {
		return nil, err
	}

	return imgBuffer.Bytes(), nil
}

func (ir *imageResizer) isImageTypeSupported(imgType string) bool {
	for _, supported := range ir.imageSupportedTypes {
		if supported == imgType {
			return true
		}
	}
	return false
}
