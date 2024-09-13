package config

import (
	"os"
	"path"
	"slices"
	"strings"

	"github.com/labstack/gommon/log"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string  `env:"SERVER_ADDRESS"`
	LogLevel      log.Lvl `env:"LOG_LEVEL"`
	LogPath       string  `env:"LOG_PATH"`
	EnableHTTPS   bool    `env:"ENABLE_HTTPS"`
}

func NewConfig() (*Config, error) {
	config := &Config{}

	currentDir, _ := os.Getwd()
	paths := strings.Split(currentDir, "/")
	projectDirIndex := slices.Index(paths, "GophKeeper")
	paths[0] = "/" + paths[0]
	paths = paths[:projectDirIndex+1]
	paths = append(paths, "client.env")
	envFilePath := path.Join(paths...)
	if err := godotenv.Load(envFilePath); err != nil {
		return nil, err
	}

	serverAddress, exists := os.LookupEnv("SERVER_ADDRESS")
	if exists {
		config.ServerAddress = serverAddress
	}

	var logLevel string
	envLogLevel, exists := os.LookupEnv("LOG_LEVEL")
	if exists {
		logLevel = envLogLevel
	}

	logPath, exists := os.LookupEnv("LOG_PATH")
	if exists {
		config.LogPath = logPath
	}

	enableHTTPS, exists := os.LookupEnv("ENABLE_HTTPS")
	if exists {
		config.EnableHTTPS = enableHTTPS == "1"
	}

	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		config.LogLevel = log.DEBUG
	case "INFO":
		config.LogLevel = log.INFO
	case "WARN":
		config.LogLevel = log.WARN
	case "ERROR":
		config.LogLevel = log.ERROR
	case "OFF":
		config.LogLevel = log.OFF
	}

	return config, nil
}
