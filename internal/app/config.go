package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddr              string
	ServerReadTimeout       int
	ServerReadHeaderTimeout int
	CacheCapacity           int
	ImageSupportedTypes     []string
}

var ErrFailedToLoadEnvFile = errors.New("failed to load env file")

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, ErrFailedToLoadEnvFile
	}

	configErrors := make([]string, 0)

	serverIP := viper.GetString("SERVER_IP")
	if serverIP == "" {
		configErrors = append(configErrors, "your have to set SERVER_IP")
	}
	serverPort := viper.GetInt("SERVER_PORT")
	if serverPort <= 0 {
		configErrors = append(configErrors, "your have to set positive SERVER_PORT")
	}
	serverReadTimeout := viper.GetInt("SERVER_READ_TIMEOUT")
	if serverReadTimeout <= 0 {
		configErrors = append(configErrors, "your have to set positive SERVER_READ_TIMEOUT")
	}
	serverReadHeaderTimeout := viper.GetInt("SERVER_READ_HEADER_TIMEOUT")
	if serverReadHeaderTimeout <= 0 {
		configErrors = append(configErrors, "your have to set positive SERVER_READ_HEADER_TIMEOUT")
	}
	cacheCapacity := viper.GetInt("CACHE_CAPACITY")
	if cacheCapacity <= 0 {
		configErrors = append(configErrors, "your have to set positive CACHE_CAPACITY")
	}
	imageSupportedTypes := viper.GetStringSlice("IMAGE_SUPPORTED_TYPES")
	if len(imageSupportedTypes) == 0 {
		configErrors = append(configErrors, "your have to set IMAGE_SUPPORTED_TYPES")
	}

	if len(configErrors) > 0 {
		errorText := strings.Join(configErrors, ";\n")

		return nil, errors.New(errorText)
	}

	return &Config{
		ServerAddr:              fmt.Sprintf("%s:%d", serverIP, serverPort),
		ServerReadTimeout:       serverReadTimeout,
		ServerReadHeaderTimeout: serverReadHeaderTimeout,
		CacheCapacity:           cacheCapacity,
		ImageSupportedTypes:     imageSupportedTypes,
	}, nil
}
