package config

import (
	"flag"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type Config struct {
	RunAddress   string  `env:"RUN_ADDRESS"`
	DatabaseURI  string  `env:"DATABASE_URI"`
	JwtSecretKey string  `env:"JWT_SECRET_KEY"`
	LogLevel     log.Lvl `env:"LOG_LEVEL"`
	LogPath      string  `env:"LOG_PATH"`
	EnableHTTPS  bool    `env:"ENABLE_HTTPS"`
}

func NewConfig() (*Config, error) {
	config := &Config{}

	currentDir, _ := os.Getwd()
	paths := strings.Split(currentDir, "/")
	projectDirIndex := slices.Index(paths, "GophKeeper")
	paths[0] = "/" + paths[0]
	paths = paths[:projectDirIndex+1]
	paths = append(paths, "server.env")
	envFilePath := path.Join(paths...)
	if err := godotenv.Load(envFilePath); err != nil {
		return nil, err
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&config.RunAddress, "a", "localhost:8080", "Run address")
	}

	if flag.Lookup("d") == nil {
		flag.StringVar(&config.DatabaseURI, "d", "", "Database dsn")
	}
	if flag.Lookup("s") == nil {
		flag.StringVar(&config.JwtSecretKey, "s", "", "JWT secret key")
	}
	var logLevel string
	if flag.Lookup("l") == nil {
		flag.StringVar(&logLevel, "l", "", "Log level")
	}
	if flag.Lookup("f") == nil {
		flag.StringVar(&config.LogPath, "f", "", "Log path")
	}
	if flag.Lookup("h") == nil {
		flag.BoolVar(&config.EnableHTTPS, "h", false, "Use https")
	}

	flag.Parse()

	runAddress, exists := os.LookupEnv("RUN_ADDRESS")
	if exists {
		config.RunAddress = runAddress
	}

	databaseURI, exists := os.LookupEnv("DATABASE_URI")
	if exists {
		config.DatabaseURI = databaseURI
	}

	jwtSecretKey, exists := os.LookupEnv("JWT_SECRET_KEY")
	if exists {
		config.JwtSecretKey = jwtSecretKey
	}

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
