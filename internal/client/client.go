package client

import (
	"context"

	"github.com/ShukinDmitriy/GophKeeper/internal/common/models"

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
	http       http.ClientInterface
	tuiService *tui.TUIService
}

// NewClient - создаёт клиента с заданным конфигом
func NewClient(
	appLog logger.Logger,
	config *config.Config,
	eventBus *event.Observable,
	http http.ClientInterface,
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

			c.tuiService.DataPage()
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

			c.tuiService.DataPage()
		case event.ClientEventSelectDataType:
			dataType, ok := e.Data.(models.DataType)
			if !ok {
				c.appLog.Error("error type data %t", e.Data)
				return
			}

			dataList, err := c.http.GetList(ctx, dataType)
			if err != nil {
				c.appLog.Error("error get data %v", err)
				c.tuiService.DataError(err.Error())
				return
			}

			c.tuiService.DrawDataList(dataType, dataList)
		case event.ClientEventSelectDataRow:
			data, ok := e.Data.(models.DataInfo)
			if !ok {
				c.appLog.Error("error type data %t", e.Data)
				return
			}

			c.tuiService.DrawDataRow(data)
		case event.ClientEventPressToCreateFormButton:
			c.tuiService.DrawCreateForm()
		case event.ClientEventSelectCreateDataType:
			dataType, ok := e.Data.(string)
			if !ok {
				c.appLog.Error("error type data %t", e.Data)
				return
			}

			switch dataType {
			case "Учетные данные":
				c.tuiService.DrawCreateCredentialsForm()
			case "Текстовые данные":
				c.tuiService.DrawCreateTextForm()
			case "Бинарные данные":
				c.tuiService.DrawCreateBinaryForm()
			case "Данные банковских карт":
				c.tuiService.DrawCreateBankForm()
			}
		case event.ClientEventCreateData:
			data, ok := e.Data.(commonRequests.DataModel)
			if !ok {
				c.appLog.Error("error type data %t", e.Data)
				return
			}

			dataInfo, err := c.http.CreateData(ctx, data)
			if err != nil {
				c.appLog.Error("error create data %v", err)
				return
			}

			c.eventBus.Next(&event.Event{
				Name: event.ClientEventCreatedData,
				Data: *dataInfo,
			})

			c.eventBus.Next(&event.Event{
				Name: event.ClientEventSelectDataType,
				Data: data.Type,
			})
		case event.ClientEventUpdateData:
			data, ok := e.Data.(models.DataInfo)
			if !ok {
				c.appLog.Error("error type data %t", e.Data)
				return
			}

			dataInfo, err := c.http.UpdateData(ctx, data)
			if err != nil {
				c.appLog.Error("error update data %v", err)
				return
			}

			c.eventBus.Next(&event.Event{
				Name: event.ClientEventUpdatedData,
				Data: *dataInfo,
			})

			c.eventBus.Next(&event.Event{
				Name: event.ClientEventSelectDataType,
				Data: data.Type,
			})
		case event.ClientEventDeleteData:
			data, ok := e.Data.(models.DataInfo)
			if !ok {
				c.appLog.Error("error type data %t", e.Data)
				return
			}

			err := c.http.DeleteData(ctx, data.ID)
			if err != nil {
				c.appLog.Error("error delete data %v", err)
				return
			}

			c.eventBus.Next(&event.Event{
				Name: event.ClientEventDeletedData,
				Data: data,
			})

			c.eventBus.Next(&event.Event{
				Name: event.ClientEventSelectDataType,
				Data: data.Type,
			})
		}
	})

	c.tuiService.LoginPage()
	// c.tuiService.DataPage()

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
