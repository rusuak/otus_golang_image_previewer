package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rusuak/otus_golang_image_previewer/internal/app"
	"github.com/rusuak/otus_golang_image_previewer/internal/resizeproxy"
	"github.com/rusuak/otus_golang_image_previewer/pkg/lrucache"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Logger

	configFile := ".env"
	config, err := app.NewConfig(configFile)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load config: " + err.Error())

		return
	}

	cache := lrucache.NewCache(config.CacheCapacity)
	imageClient := resizeproxy.NewImageClient()
	rp := resizeproxy.NewResizeProxy(&logger, cache, imageClient, config.ImageSupportedTypes)

	newApp, err := app.NewApp(config, rp, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to init app:" + err.Error())

		return
	}

	logger.Info().Msg("Server is starting...")
	if err := newApp.Run(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
