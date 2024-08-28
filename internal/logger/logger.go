package logger

import (
	"os"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

func NewLogger(
	logLevel log.Lvl,
	logPath string,
) Logger {
	logger := logrus.New()

	// установим уровень логирования
	var level logrus.Level
	switch logLevel {
	case log.DEBUG:
		level = logrus.DebugLevel
	case log.INFO:
		level = logrus.InfoLevel
	case log.WARN:
		level = logrus.WarnLevel
	case log.ERROR:
		level = logrus.ErrorLevel
	case log.OFF:
		level = logrus.PanicLevel
	}
	logger.SetLevel(level)

	// установим форматирование логов в джейсоне
	logger.SetFormatter(&logrus.JSONFormatter{})

	// установим вывод логов в файл
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err == nil {
		logger.SetOutput(file)
	} else {
		logger.Info("Не удалось открыть файл логов, используется стандартный stderr")
	}

	return logger
}
