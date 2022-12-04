package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigSuccess(t *testing.T) {
	t.Run("ok_case", func(t *testing.T) {
		config, err := NewConfig("./testdata/testing_config.env")

		expected := &Config{
			ServerAddr:              "0.0.0.0:8080",
			ServerReadTimeout:       15,
			ServerReadHeaderTimeout: 5,
			CacheCapacity:           3,
			ImageSupportedTypes:     []string{"image/jpeg", "image/jpg"},
		}

		require.Equal(t, expected, config)
		require.Nil(t, err)
	})
}

func TestNotFoundFile(t *testing.T) {
	t.Run("not_found_file_case", func(t *testing.T) {
		config, err := NewConfig("123")

		require.Nil(t, config)
		require.Errorf(t, err, ErrFailedToLoadEnvFile.Error())
	})
}

func TestNewConfigInvalid(t *testing.T) {
	expectedError := `your have to set SERVER_IP;
your have to set positive SERVER_PORT;
your have to set positive SERVER_READ_TIMEOUT;
your have to set positive SERVER_READ_HEADER_TIMEOUT;
your have to set positive CACHE_CAPACITY;
your have to set IMAGE_SUPPORTED_TYPES`

	t.Run("new_config_invalid_case", func(t *testing.T) {
		config, err := NewConfig("./testdata/testing_config_invalid.env")

		require.Nil(t, config)
		require.Equal(t, err.Error(), expectedError)
	})
}
