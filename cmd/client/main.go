package main

import (
	"context"

	"github.com/gdamore/tcell/v2"

	"github.com/ShukinDmitriy/GophKeeper/internal/client"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/event"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/http"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/tui"
	appLogger "github.com/ShukinDmitriy/GophKeeper/internal/logger"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	appLog := appLogger.NewLogger(conf.LogLevel, conf.LogPath)
	httpClient := http.NewClient(conf, appLog)
	eventBus := event.NewObservable()
	screen, err := tcell.NewScreen()
	if err != nil {
		appLog.Fatal("Failed to initialize screen", err)
	}
	tuiService := tui.NewTUIService(appLog, eventBus, screen)
	if tuiService == nil {
		appLog.Fatal("Failed to create tuiService")
		return
	}
	tClient := client.NewClient(
		appLog,
		conf,
		eventBus,
		httpClient,
		tuiService,
	)

	ctx := context.Background()
	appLog.Info("Running GophKeeper client")
	// Здесь будет зависание основного потока
	err = tClient.Run(ctx)
	if err != nil {
		appLog.Fatal("shutting down the GophKeeper client ", err.Error())
		return
	}

	err = tClient.Shutdown(ctx)
	if err != nil {
		appLog.Error(err.Error())
	}

	appLog.Info("Stopping the GophKeeper client")
}
