package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_SOURCE", "user=postgres password=postgres dbname=simple_bank sslmode=disable")
	os.Setenv("SERVER_ADDRESS", "localhost:8080")

	// Invalid path should return error
	config, err := LoadConfig("../../../../../../")
	require.Error(t, err)

	// Load configuration
	config, err = LoadConfig("../../")
	require.NoError(t, err)

	// Assert values
	require.Equal(t, "postgres", config.DBDriver)
	require.Equal(t, "user=postgres password=postgres dbname=simple_bank sslmode=disable", config.DBSource)
	require.Equal(t, "localhost:8080", config.ServerAddress)

}
