package logger_test

import (
	"testing"

	"github.com/ShukinDmitriy/GophKeeper/internal/logger"

	"github.com/ShukinDmitriy/GophKeeper/internal/server/config"
	"github.com/labstack/gommon/log"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		conf *config.Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success #1",
			args: args{
				conf: &config.Config{
					LogLevel: log.DEBUG,
				},
			},
		},
		{
			name: "success #2",
			args: args{
				conf: &config.Config{
					LogLevel: log.INFO,
				},
			},
		},
		{
			name: "success #3",
			args: args{
				conf: &config.Config{
					LogLevel: log.WARN,
				},
			},
		},
		{
			name: "success #4",
			args: args{
				conf: &config.Config{
					LogLevel: log.ERROR,
				},
			},
		},
		{
			name: "success #5",
			args: args{
				conf: &config.Config{
					LogLevel: log.OFF,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.NewLogger(tt.args.conf.LogLevel, "")
		})
	}
}
