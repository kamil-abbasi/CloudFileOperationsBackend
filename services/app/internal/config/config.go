package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           uint32
	RootPath       string
	MaxUploadBytes uint64
	UsageCacheTTL  uint32
	LogLevel       string
	GinMode        string
	PgDb           string
	PgUser         string
	PgPassword     string
	PgHost         string
}

func Load() (*Config, error) {
	config := &Config{}

	portStr := os.Getenv("PORT")
	rootPath := os.Getenv("ROOT_PATH")
	maxUploadBytesStr := os.Getenv("MAX_UPLOAD_BYTES")
	usageCacheTTLStr := os.Getenv("USAGE_CACHE_TTL")
	logLevel := os.Getenv("LOG_LEVEL")
	ginMode := os.Getenv("GIN_MODE")
	pgDb := os.Getenv("PG_DB")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgHost := os.Getenv("PG_HOST")

	port, err := strconv.Atoi(portStr)

	if err != nil {
		return nil, fmt.Errorf("PORT must be an integer")
	}

	maxUploadBytes, err := strconv.Atoi(maxUploadBytesStr)

	if err != nil {
		return nil, fmt.Errorf("MAX_UPLOAD_BYTES must be an integer")
	}

	usageCacheTTL, err := strconv.Atoi(usageCacheTTLStr)

	if err != nil {
		return nil, fmt.Errorf("USAGE_CACHE_TTL must be an integer")
	}

	config.MaxUploadBytes = uint64(maxUploadBytes)
	config.Port = uint32(port)
	config.RootPath = rootPath
	config.UsageCacheTTL = uint32(usageCacheTTL)
	config.GinMode = ginMode
	config.LogLevel = logLevel
	config.PgDb = pgDb
	config.PgUser = pgUser
	config.PgPassword = pgPassword
	config.PgHost = pgHost

	return config, nil
}
