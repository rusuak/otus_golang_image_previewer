package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rusuak/otus_golang_image_previewer/internal/resizeproxy"
)

type Handler struct {
	logger  *zerolog.Logger
	resizer resizeproxy.ImageResizer
}

func NewHandler(
	logger *zerolog.Logger,
	resizer resizeproxy.ImageResizer,
) *Handler {
	return &Handler{logger: logger, resizer: resizer}
}

var (
	ErrInvalidWidth  = fmt.Errorf("width is invalid")
	ErrInvalidHeight = fmt.Errorf("height is invalid")
)

func (h *Handler) ResizeHandler(w http.ResponseWriter, r *http.Request) {
	request, err := h.createResizeRequest(mux.Vars(r), r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("validation error"))
		h.logger.Err(err).Msg(err.Error())

		return
	}

	resizeResponse, err := h.resizer.ResizeImageCached(r.Context(), request)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("resize image failed"))
		h.logger.Err(err).Msg(err.Error())

		return
	}

	for name, values := range resizeResponse.Headers {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(resizeResponse.Img)))
	if _, err := w.Write(resizeResponse.Img); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Err(err).Msg(err.Error())
	}
}

func (h *Handler) createResizeRequest(
	vars map[string]string,
	headers map[string][]string,
) (r *resizeproxy.ResizeRequest, err error) {
	width, err := strconv.Atoi(vars["width"])
	if err != nil {
		return nil, ErrInvalidWidth
	}

	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		return nil, ErrInvalidHeight
	}

	imageURL, err := url.Parse(vars["imageURL"])
	if err != nil {
		return nil, err
	}
	imageURL.Scheme = "https"

	return resizeproxy.NewResizeRequest(width, height, imageURL.String(), headers), nil
}
