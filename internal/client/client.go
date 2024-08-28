package client

import (
	"context"

	commonRequests "github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
	"github.com/go-playground/validator/v10"

	"github.com/ShukinDmitriy/GophKeeper/internal/client/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/event"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/http"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/tui"
	"github.com/ShukinDmitriy/GophKeeper/internal/logger"
)

// Client - основная структура для работы с клиентом
type Client struct {
	appLog     logger.Logger
	config     *config.Config
	eventBus   *event.Observable
	http       *http.Client
	tuiService *tui.TUIService
}

// NewClient - создаёт клиента с заданным конфигом
func NewClient(
	appLog logger.Logger,
	config *config.Config,
	eventBus *event.Observable,
	http *http.Client,
	tuiService *tui.TUIService,
) *Client {
	c := &Client{
		appLog:     appLog,
		config:     config,
		eventBus:   eventBus,
		http:       http,
		tuiService: tuiService,
	}

	return c
}

// Run - Запускает клиента
func (c *Client) Run(ctx context.Context) error {
	c.tuiService.LoginPage()

	c.eventBus.Subscribe(func(e *event.Event) {
		switch e.Name {
		case event.ClientEventPressToRegisterButton:
			c.tuiService.RegisterPage()
		case event.ClientEventPressToLoginButton:
			c.tuiService.LoginPage()
		case event.ClientEventPressLoginButton:
			loginFormData, ok := e.Data.(commonRequests.UserLogin)
			if !ok {
				c.appLog.Error("Не удалось получить данные формы")
			}

			validate := validator.New(validator.WithRequiredStructEnabled())
			err := validate.Struct(loginFormData)
			if err != nil {
				c.appLog.Error("error login %v", err)
				c.tuiService.LoginError("Необходимо ввести корректные логин и пароль")
				return
			}

			err = c.http.Login(ctx, loginFormData)
			if err != nil {
				c.appLog.Error("error login %v", err)
				c.tuiService.LoginError(err.Error())
			}

			// TODO переход к списку данных
		case event.ClientEventPressRegisterButton:
			registerFormData, ok := e.Data.(commonRequests.UserRegister)
			if !ok {
				c.appLog.Error("Не удалось получить данные формы")
			}

			validate := validator.New(validator.WithRequiredStructEnabled())
			err := validate.Struct(registerFormData)
			if err != nil {
				c.appLog.Error("error register %v", err)
				c.tuiService.RegisterError("Необходимо ввести корректные логин и пароль")
				return
			}

			err = c.http.Register(ctx, registerFormData)
			if err != nil {
				c.appLog.Error("error register %v", err)
				c.tuiService.RegisterError(err.Error())
			}

			// TODO переход к списку данных
		}
	})

	// Должен работать в основном потоке...
	err := c.tuiService.Run()
	if err != nil {
		c.appLog.Error(err)
	}

	return nil
}

// Shutdown - остановить приложение
func (c *Client) Shutdown(_ context.Context) error {
	c.tuiService.Stop()
	return nil
}
