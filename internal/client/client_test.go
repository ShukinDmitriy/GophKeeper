package client_test

import (
	"context"
	"fmt"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/models"
	"testing"
	"time"

	"github.com/ShukinDmitriy/GophKeeper/internal/client"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/event"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/http"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/router"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/tui"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
	appLogger "github.com/ShukinDmitriy/GophKeeper/internal/logger"
	"github.com/ShukinDmitriy/GophKeeper/internal/test-helpers"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func createConfig(t *testing.T) *config.Config {
	conf, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	return conf
}

func createLogger(conf *config.Config) appLogger.Logger {
	return appLogger.NewLogger(conf.LogLevel, conf.LogPath)
}

func createHttpClient(t *testing.T, conf *config.Config, appLog appLogger.Logger) *http.Client {
	client := http.NewClient(conf, appLog)
	if client == nil {
		t.Fatal("failed to create http client")
	}

	return client
}

func createEventBus() *event.Observable {
	return event.NewObservable()
}

func createScreen() tcell.Screen {
	return tcell.NewSimulationScreen("UTF-8")
}

func createTUIService(
	t *testing.T,
	appLog appLogger.Logger,
	eventBus *event.Observable,
	screen tcell.Screen,
) *tui.TUIService {
	tuiService := tui.NewTUIService(appLog, eventBus, screen)
	if tuiService == nil {
		t.Fatal("Failed to create tuiService")
	}

	return tuiService
}

func createClient(
	t *testing.T,
	done chan struct{},
	appLog appLogger.Logger,
	conf *config.Config,
	eventBus *event.Observable,
	httpClient *http.Client,
	tuiService *tui.TUIService,
) (context.Context, *client.Client) {
	tClient := client.NewClient(
		appLog,
		conf,
		eventBus,
		httpClient,
		tuiService,
	)
	if tClient == nil {
		t.Fatal("Failed to create tClient")
	}

	ctx := context.Background()

	go func() {
		defer close(done)

		done <- struct{}{}
		err := tClient.Run(ctx)
		if err != nil {
			t.Error(err)
			return
		}

		done <- struct{}{}
	}()

	return ctx, tClient
}

func goToRegisterAndReturnToLoginPage(
	t *testing.T,
	eventBus *event.Observable,
	tuiService *tui.TUIService,
) {
	t.Run("test route to register and return to login page", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	GetStartPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage != "" {
					assert.NotEqual(t, router.RegisterPage, currentPage)
					break GetStartPage
				}
				time.Sleep(1 * time.Second)
			}
		}

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventPressToRegisterButton,
		})

	GetRegisterPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.RegisterPage {
					break GetRegisterPage
				}
				time.Sleep(1 * time.Second)
			}
		}

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventPressToLoginButton,
		})

	GetLoginPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.LoginPage {
					break GetLoginPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})
}

func submitLoginForm(
	t *testing.T,
	eventBus *event.Observable,
	tuiService *tui.TUIService,
) {
	t.Run("test submit login form with errors", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventPressLoginButton,
			Data: requests.UserLogin{
				Login:    "",
				Password: "",
			},
		})

	GetErrorPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.ErrorPage {
					break GetErrorPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test submit valid login form", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventPressLoginButton,
			Data: requests.UserLogin{
				Login:    "login",
				Password: "password",
			},
		})

		time.Sleep(2 * time.Second)
	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})
}

func submitRegisterForm(
	t *testing.T,
	eventBus *event.Observable,
	tuiService *tui.TUIService,
) {
	t.Run("test submit register form with errors", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventPressRegisterButton,
			Data: requests.UserRegister{
				Login:    "",
				Password: "",
			},
		})

	GetErrorPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.ErrorPage {
					break GetErrorPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test submit valid register form", func(t *testing.T) {
		timeout := time.After(30 * time.Second)
	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventPressRegisterButton,
			Data: requests.UserRegister{
				Login:    test_helpers.GenerateRandomString(10),
				Password: test_helpers.GenerateRandomString(10),
			},
		})

		time.Sleep(2 * time.Second)
	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})
}

func selectDataRow(
	t *testing.T,
	eventBus *event.Observable,
	tuiService *tui.TUIService,
) {
	t.Run("test select credentials data row", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventSelectDataRow,
			Data: models.DataInfo{
				ID:          777,
				Type:        models.DataTypeCredentials,
				Description: "test",
				Value:       "ewogICJMb2dpbiI6ICJ0ZXN0IiwKICAiUGFzc3dvcmQiOiAidGVzdCIKfQ==",
			},
		})

	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test select text data row", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventSelectDataRow,
			Data: models.DataInfo{
				ID:          777,
				Type:        models.DataTypeText,
				Description: "test",
				Value:       "ewogICJUZXh0IjogInRlc3QiCn0=",
			},
		})

	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test select binary data row", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventSelectDataRow,
			Data: models.DataInfo{
				ID:          777,
				Type:        models.DataTypeBinary,
				Description: "test",
				Value:       "ewogICJCaW5hcnkiOiAidGVzdCIKfQ==",
			},
		})

	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test select bank data row", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventSelectDataRow,
			Data: models.DataInfo{
				ID:          777,
				Type:        models.DataTypeBankCard,
				Description: "test",
				Value:       "ewogICJOdW1iZXIiOiAiMjIwMDEyMzQ0MzIxOTg3NiIsCiAgIkRhdGUiOiAiMDIvMjUiLAogICJTZWN1cmUiOiAiMTIzIgp9",
			},
		})

	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})
}

func clickCreateButton(
	t *testing.T,
	eventBus *event.Observable,
	tuiService *tui.TUIService,
) {
	t.Run("test click create button", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventPressToCreateFormButton,
		})

	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test select credentials create form", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventSelectCreateDataType,
			Data: "Учетные данные",
		})

	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test select text create form", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventSelectCreateDataType,
			Data: "Текстовые данные",
		})

	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test select binary create form", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventSelectCreateDataType,
			Data: "Бинарные данные",
		})

	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test select bank create form", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		eventBus.Next(&event.Event{
			Name: event.ClientEventSelectCreateDataType,
			Data: "Данные банковских карт",
		})

	GetDataPage:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				currentPage := tuiService.GetCurrentPage()
				if currentPage == router.DataPage {
					break GetDataPage
				}
				time.Sleep(1 * time.Second)
			}
		}
	})
}

func crudData(
	t *testing.T,
	eventBus *event.Observable,
	tuiService *tui.TUIService,
) {
	var dataInfo models.DataInfo

	// Подписываемся на событие, чтобы отловить созданную запись
	subscription := eventBus.Subscribe(func(e *event.Event) {
		if e.Name != event.ClientEventCreatedData {
			return
		}

		eventDataInfo, ok := e.Data.(models.DataInfo)
		if ok {
			dataInfo = eventDataInfo
		}
	})

	t.Run("test create data", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		dataCreated := false
		// Подписка, чтобы отследить успешное выполнение
		subscription := eventBus.Subscribe(func(e *event.Event) {
			if e.Name != event.ClientEventCreatedData {
				return
			}

			dataCreated = true
		})

		eventBus.Next(&event.Event{
			Name: event.ClientEventCreateData,
			Data: requests.DataModel{
				Type:        models.DataTypeCredentials,
				Description: "test",
				Value:       "ewogICJMb2dpbiI6ICJ0ZXN0IiwKICAiUGFzc3dvcmQiOiAidGVzdCIKfQ==",
			},
		})

	AwaitDataCreated:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				if dataCreated {
					subscription.Unsubscribe()
					break AwaitDataCreated
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test update data", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		dataUpdated := false
		// Подписка, чтобы отследить успешное выполнение
		subscription := eventBus.Subscribe(func(e *event.Event) {
			if e.Name != event.ClientEventUpdatedData {
				return
			}

			dataUpdated = true
		})

		eventBus.Next(&event.Event{
			Name: event.ClientEventUpdateData,
			Data: models.DataInfo{
				ID:          dataInfo.ID,
				Type:        dataInfo.Type,
				Description: "test updated",
				Value:       "ewogICJMb2dpbiI6ICJ0ZXN0IHVwZGF0ZWQiLAogICJQYXNzd29yZCI6ICJ0ZXN0IHVwZGF0ZWQiCn0=",
			},
		})

	AwaitDataUpdated:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				if dataUpdated {
					subscription.Unsubscribe()
					break AwaitDataUpdated
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	t.Run("test delete data", func(t *testing.T) {
		timeout := time.After(30 * time.Second)

	_:
		dataDeleted := false
		// Подписка, чтобы отследить успешное выполнение
		subscription := eventBus.Subscribe(func(e *event.Event) {
			if e.Name != event.ClientEventDeletedData {
				return
			}

			dataDeleted = true
		})

		eventBus.Next(&event.Event{
			Name: event.ClientEventDeleteData,
			Data: models.DataInfo{
				ID:          dataInfo.ID,
				Type:        dataInfo.Type,
				Description: "test updated",
				Value:       "ewogICJMb2dpbiI6ICJ0ZXN0IHVwZGF0ZWQiLAogICJQYXNzd29yZCI6ICJ0ZXN0IHVwZGF0ZWQiCn0=",
			},
		})

	AwaitDataDeleted:
		for {
			select {
			case <-timeout:
				t.Fatal("timed out waiting")
			default:
				if dataDeleted {
					subscription.Unsubscribe()
					break AwaitDataDeleted
				}
				time.Sleep(1 * time.Second)
			}
		}
	})

	subscription.Unsubscribe()
}

func TestClient(t *testing.T) {
	// Запуск сервера
	_, serverConf, httpServer := test_helpers.RunServer(t)

	conf := createConfig(t)
	if serverConf.EnableHTTPS {
		conf.ServerAddress = fmt.Sprintf("https://%s", serverConf.RunAddress)
	} else {
		conf.ServerAddress = fmt.Sprintf("http://%s", serverConf.RunAddress)
	}

	appLog := createLogger(conf)
	httpClient := createHttpClient(t, conf, appLog)
	eventBus := createEventBus()
	screen := createScreen()
	tuiService := createTUIService(t, appLog, eventBus, screen)
	done := make(chan struct{})
	ctx, tClient := createClient(
		t,
		done,
		appLog,
		conf,
		eventBus,
		httpClient,
		tuiService,
	)

	// Ждем отработки горутины с запуском приложения
	<-done

	// Запускаем тесты
	goToRegisterAndReturnToLoginPage(t, eventBus, tuiService)
	submitLoginForm(t, eventBus, tuiService)
	submitRegisterForm(t, eventBus, tuiService)
	selectDataRow(t, eventBus, tuiService)
	clickCreateButton(t, eventBus, tuiService)
	crudData(t, eventBus, tuiService)

	// Отработали, останавливаем приложение
	err := tClient.Shutdown(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// Проверка, что приложение успешно завершилось
	select {
	case <-done:
		// Горутина завершилась без ошибок
		break
	case <-time.After(time.Second * 2):
		t.Fatal("Горутина не завершилась в течение 2 секунд")
	}

	// Отключаем сервер
	test_helpers.StopServer(t, httpServer)
}
