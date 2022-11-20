package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rusuak/otus_golang_image_previewer/internal/resizeproxy"
	"github.com/rusuak/otus_golang_image_previewer/internal/server"
)

type App struct {
	config       *Config
	imageResizer resizeproxy.ImageResizer
	logger       *zerolog.Logger
}

func NewApp(config *Config, ir resizeproxy.ImageResizer, logger *zerolog.Logger) (*App, error) {
	return &App{
		config:       config,
		imageResizer: ir,
		logger:       logger,
	}, nil
}

func (app App) Run() error {
	srv := server.NewServer(server.Config{
		Addr:              app.config.ServerAddr,
		ReadTimeout:       app.config.ServerReadTimeout,
		ReadHeaderTimeout: app.config.ServerReadHeaderTimeout,
	})

	handler := server.NewHandler(app.logger, app.imageResizer)
	router := app.newRouter(handler)

	return srv.Run(router)
}

func (app App) newRouter(handler *server.Handler) *mux.Router {
	resizeURLTemplate := "/fill/{width:[0-9]+}/{height:[0-9]+}/{imageURL:.*}"

	router := mux.NewRouter()
	router.HandleFunc(resizeURLTemplate, handler.ResizeHandler)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Please use resize url by template:" + resizeURLTemplate))
	})

	return router
}
