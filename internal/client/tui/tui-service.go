package tui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ShukinDmitriy/GophKeeper/internal/client/event"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/router"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/models"
	commonRequests "github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
	"github.com/ShukinDmitriy/GophKeeper/internal/logger"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TUIService сервис для работы с консолью
type TUIService struct {
	application *tview.Application
	appLog      logger.Logger
	dataForm    *tview.Form
	dataList    *tview.List
	dataTypes   *tview.List
	eventBus    *event.Observable
	pages       *tview.Pages
	running     bool
}

// NewTUIService конструктор для TUIService
func NewTUIService(
	appLog logger.Logger,
	eventBus *event.Observable,
	screen tcell.Screen,
) *TUIService {
	pages := tview.NewPages()

	pages.SetBorder(true).
		SetTitle(router.ApplicationName).
		SetTitleAlign(tview.AlignLeft)

	application := tview.NewApplication().
		SetScreen(screen).
		SetRoot(pages, true).
		SetFocus(pages).
		EnableMouse(true)

	return &TUIService{
		application: application,
		appLog:      appLog,
		eventBus:    eventBus,
		pages:       pages,
	}
}

// errorPage - отобразить страницу ошибки
func (tuiService *TUIService) errorPage(err string, returnToPage string) {
	tuiService.appLog.Debug("Create error page")

	form := tview.NewForm().
		AddTextView("Ошибка", err, 50, 5, true, true).
		AddButton("Понятно", func() {
			tuiService.pages.SwitchToPage(returnToPage)
		})

	tuiService.pages.AddAndSwitchToPage(router.ErrorPage, form, true)

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// dataToBase64 - преобразовать данные в base64
func (tuiService *TUIService) dataToBase64(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		tuiService.appLog.Error(fmt.Sprintf("error data to json: %v", err))
		return ""
	}

	return base64.StdEncoding.EncodeToString(jsonData)
}

// drawDataTypes - отрисовать типы данных
func (tuiService *TUIService) drawDataTypes() {
	type dataTypeStruct struct {
		dataType models.DataType
		title    string
		shortcut rune
	}
	dataTypes := []dataTypeStruct{
		{
			dataType: models.DataTypeCredentials,
			title:    "Учетные данные",
			shortcut: 0,
		},
		{
			dataType: models.DataTypeText,
			title:    "Текстовые данные",
			shortcut: 0,
		},
		{
			dataType: models.DataTypeBinary,
			title:    "Бинарные данные",
			shortcut: 0,
		},
		{
			dataType: models.DataTypeBankCard,
			title:    "Данные банковских карт",
			shortcut: 0,
		},
	}

	for _, dataType := range dataTypes {
		tuiService.dataTypes.AddItem(dataType.title, fmt.Sprintf("%v", dataType.dataType), dataType.shortcut, func() {
			tuiService.application.SetFocus(tuiService.dataList)
		})
	}

	tuiService.dataTypes.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		tuiService.eventBus.Next(&event.Event{
			Name: event.ClientEventSelectDataType,
			Data: dataTypes[index].dataType,
		})
	})

	tuiService.eventBus.Next(&event.Event{
		Name: event.ClientEventSelectDataType,
		Data: dataTypes[0].dataType,
	})
}

// drawDataRowCredentials - отрисовать форму просмотра "Учетные данные"
func (tuiService *TUIService) drawDataRowCredentials(data models.DataInfo) {
	decodedValue, err := base64.StdEncoding.DecodeString(data.Value)
	if err != nil {
		tuiService.appLog.Error(fmt.Sprintf("error decode value: %v", err))
		return
	}
	value := &struct {
		Login    string
		Password string
	}{}
	err = json.Unmarshal(decodedValue, value)
	if err != nil {
		tuiService.appLog.Error(fmt.Sprintf("error decode value: %v", err))
		return
	}

	tuiService.dataForm.Clear(true)

	tuiService.dataForm.SetTitle("Подробно")

	tuiService.dataForm.
		AddTextView("Идентификатор", fmt.Sprintf("%d", data.ID), 50, 1, true, true).
		AddTextView("Тип", "Учетные данные", 50, 1, true, true).
		AddInputField("Описание", data.Description, 50, nil, func(text string) {
			data.Description = text
		}).
		AddInputField("Логин", value.Login, 50, nil, func(text string) {
			value.Login = text
		}).
		AddInputField("Пароль", value.Password, 50, nil, func(text string) {
			value.Password = text
		}).
		AddButton("Изменить", func() {
			base64Data := tuiService.dataToBase64(value)
			if base64Data == "" {
				return
			}
			data.Value = base64Data

			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventUpdateData,
				Data: data,
			})
		}).
		AddButton("Удалить", func() {
			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventDeleteData,
				Data: data,
			})
		})

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// drawDataRowText - отрисовать форму просмотра "Текстовые данные"
func (tuiService *TUIService) drawDataRowText(data models.DataInfo) {
	decodedValue, err := base64.StdEncoding.DecodeString(data.Value)
	if err != nil {
		tuiService.appLog.Error(fmt.Sprintf("error decode value: %v", err))
		return
	}
	value := &struct {
		Text string
	}{}
	err = json.Unmarshal(decodedValue, value)
	if err != nil {
		tuiService.appLog.Error(fmt.Sprintf("error decode value: %v", err))
		return
	}

	tuiService.dataForm.Clear(true)

	tuiService.dataForm.SetTitle("Подробно")

	tuiService.dataForm.
		AddTextView("Идентификатор", fmt.Sprintf("%d", data.ID), 50, 1, true, true).
		AddTextView("Тип", "Текстовые данные", 50, 1, true, true).
		AddInputField("Описание", data.Description, 50, nil, func(text string) {
			data.Description = text
		}).
		AddInputField("Текст", value.Text, 50, nil, func(text string) {
			value.Text = text
		}).
		AddButton("Изменить", func() {
			base64Data := tuiService.dataToBase64(value)
			if base64Data == "" {
				return
			}
			data.Value = base64Data

			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventUpdateData,
				Data: data,
			})
		}).
		AddButton("Удалить", func() {
			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventDeleteData,
				Data: data,
			})
		})

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// drawDataRowBinary - отрисовать форму просмотра "Бинарные данные"
func (tuiService *TUIService) drawDataRowBinary(data models.DataInfo) {
	decodedValue, err := base64.StdEncoding.DecodeString(data.Value)
	if err != nil {
		tuiService.appLog.Error(fmt.Sprintf("error decode value: %v", err))
		return
	}
	value := &struct {
		Binary string
	}{}
	err = json.Unmarshal(decodedValue, value)
	if err != nil {
		tuiService.appLog.Error(fmt.Sprintf("error decode value: %v", err))
		return
	}

	tuiService.dataForm.Clear(true)

	tuiService.dataForm.SetTitle("Подробно")

	tuiService.dataForm.
		AddTextView("Идентификатор", fmt.Sprintf("%d", data.ID), 50, 1, true, true).
		AddTextView("Тип", "Текстовые данные", 50, 1, true, true).
		AddInputField("Описание", data.Description, 50, nil, func(text string) {
			data.Description = text
		}).
		AddInputField("Данные", value.Binary, 50, nil, func(text string) {
			value.Binary = text
		}).
		AddButton("Изменить", func() {
			base64Data := tuiService.dataToBase64(value)
			if base64Data == "" {
				return
			}
			data.Value = base64Data

			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventUpdateData,
				Data: data,
			})
		}).
		AddButton("Удалить", func() {
			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventDeleteData,
				Data: data,
			})
		})

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// drawDataRowBank - отрисовать форму просмотра "Банковские данные"
func (tuiService *TUIService) drawDataRowBank(data models.DataInfo) {
	decodedValue, err := base64.StdEncoding.DecodeString(data.Value)
	if err != nil {
		tuiService.appLog.Error(fmt.Sprintf("error decode value: %v", err))
		return
	}
	value := &struct {
		Number string
		Date   string
		Secure string
	}{}
	err = json.Unmarshal(decodedValue, value)
	if err != nil {
		tuiService.appLog.Error(fmt.Sprintf("error decode value: %v", err))
		return
	}

	tuiService.dataForm.Clear(true)

	tuiService.dataForm.SetTitle("Подробно")

	tuiService.dataForm.
		AddTextView("Идентификатор", fmt.Sprintf("%d", data.ID), 50, 1, true, true).
		AddTextView("Тип", "Текстовые данные", 50, 1, true, true).
		AddInputField("Описание", data.Description, 50, nil, func(text string) {
			data.Description = text
		}).
		AddInputField("Номер карты", value.Number, 50, nil, func(text string) {
			value.Number = text
		}).
		AddInputField("Срок действия", value.Date, 50, nil, func(text string) {
			value.Date = text
		}).
		AddInputField("Секретный код", value.Secure, 50, nil, func(text string) {
			value.Secure = text
		}).
		AddButton("Изменить", func() {
			base64Data := tuiService.dataToBase64(value)
			if base64Data == "" {
				return
			}
			data.Value = base64Data

			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventUpdateData,
				Data: data,
			})
		}).
		AddButton("Удалить", func() {
			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventDeleteData,
				Data: data,
			})
		})

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// Run - запустить консольное приложение
func (tuiService *TUIService) Run() error {
	tuiService.running = true
	err := tuiService.application.Run()
	if err != nil {
		return err
	}

	return nil
}

// LoginPage - отобразить страницу авторизации
func (tuiService *TUIService) LoginPage() {
	if !tuiService.pages.HasPage(router.LoginPage) {
		tuiService.appLog.Debug("Create login page")
		data := commonRequests.UserLogin{}
		form := tview.NewForm().
			AddInputField("Логин", "", 20, nil, func(text string) {
				data.Login = text
			}).
			AddPasswordField("Пароль", "", 20, '*', func(text string) {
				data.Password = text
			}).
			AddButton("Авторизоваться", func() {
				tuiService.appLog.Debug("Press Login button")
				tuiService.eventBus.Next(&event.Event{
					Name: event.ClientEventPressLoginButton,
					Data: data,
				})
			}).
			AddButton("Перейти к регистрации", func() {
				tuiService.appLog.Debug("Press Register button")
				tuiService.eventBus.Next(&event.Event{
					Name: event.ClientEventPressToRegisterButton,
				})
			}).
			AddButton("Закончить", func() {
				tuiService.appLog.Debug("Press Stop button")
				tuiService.application.Stop()
			})

		tuiService.pages.AddPage(router.LoginPage, form, true, false)
	}

	tuiService.pages.SwitchToPage(router.LoginPage)

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// LoginError - отобразить ошибку авторизации
func (tuiService *TUIService) LoginError(err string) {
	tuiService.errorPage(err, router.LoginPage)
}

// RegisterPage - отобразить страницу регистрации
func (tuiService *TUIService) RegisterPage() {
	if !tuiService.pages.HasPage(router.RegisterPage) {
		tuiService.appLog.Debug("Create register page")
		data := commonRequests.UserRegister{}
		form := tview.NewForm().
			AddInputField("Логин", "", 20, nil, func(text string) {
				data.Login = text
			}).
			AddPasswordField("Пароль", "", 20, '*', func(text string) {
				data.Password = text
			}).
			AddButton("Зарегистрироваться", func() {
				tuiService.appLog.Debug("Press Register button")
				tuiService.eventBus.Next(&event.Event{
					Name: event.ClientEventPressRegisterButton,
					Data: data,
				})
			}).
			AddButton("Перейти к авторизации", func() {
				tuiService.appLog.Debug("Press \"To login\" button")
				tuiService.eventBus.Next(&event.Event{
					Name: event.ClientEventPressToLoginButton,
				})
			}).
			AddButton("Закончить", func() {
				tuiService.appLog.Debug("Press Stop button")
				tuiService.Stop()
			})

		tuiService.pages.AddPage(router.RegisterPage, form, true, false)
	}

	tuiService.pages.SwitchToPage(router.RegisterPage)

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// RegisterError - отобразить ошибку регистрации
func (tuiService *TUIService) RegisterError(err string) {
	tuiService.errorPage(err, router.RegisterPage)
}

// DataPage - отобразить страницу данных
func (tuiService *TUIService) DataPage() {
	if !tuiService.pages.HasPage(router.DataPage) {
		tuiService.appLog.Debug("Create data page")

		tuiService.dataTypes = tview.NewList().ShowSecondaryText(false)
		tuiService.dataTypes.SetBorder(true).SetTitle("Типы данных")

		form := tview.NewForm().AddButton("Создать", func() {
			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventPressToCreateFormButton,
			})
		})

		mainMenu := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tuiService.dataTypes, 0, 3, true).
			AddItem(form, 0, 1, false)

		tuiService.dataList = tview.NewList().ShowSecondaryText(false)
		tuiService.dataList.SetBorder(true).SetTitle("Записи")

		tuiService.dataForm = tview.NewForm()
		tuiService.dataForm.SetBorder(true).SetTitle("Подробно")

		flex := tview.NewFlex().
			AddItem(mainMenu, 0, 1, true).
			AddItem(tuiService.dataList, 0, 1, false).
			AddItem(tuiService.dataForm, 0, 3, false)

		flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyLeft:
				if tuiService.dataList.HasFocus() {
					tuiService.application.SetFocus(tuiService.dataTypes)

					return event
				}

				if tuiService.dataForm.HasFocus() {
					tuiService.application.SetFocus(tuiService.dataList)

					return event
				}
			case tcell.KeyRight:
				if tuiService.dataTypes.HasFocus() {
					tuiService.application.SetFocus(tuiService.dataList)

					return event
				}

				if tuiService.dataList.HasFocus() {
					tuiService.application.SetFocus(tuiService.dataForm)

					return event
				}
			}

			return event
		})

		tuiService.pages.AddPage(router.DataPage, flex, true, false)

		tuiService.drawDataTypes()
	}

	tuiService.pages.SwitchToPage(router.DataPage)

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// DrawDataList - отрисовать список записей по типу
func (tuiService *TUIService) DrawDataList(dataType models.DataType, dataList []models.DataInfo) {
	switch dataType {
	case models.DataTypeCredentials:
		tuiService.dataTypes.SetCurrentItem(0)
	case models.DataTypeText:
		tuiService.dataTypes.SetCurrentItem(1)
	case models.DataTypeBinary:
		tuiService.dataTypes.SetCurrentItem(2)
	case models.DataTypeBankCard:
		tuiService.dataTypes.SetCurrentItem(3)
	}
	tuiService.dataList.Clear()
	tuiService.dataForm.Clear(true)
	tuiService.dataForm.SetTitle("Подробно")

	for _, data := range dataList {
		tuiService.dataList.AddItem(fmt.Sprintf("%d. %s", data.ID, data.Description), "", 0, func() {
			tuiService.application.SetFocus(tuiService.dataForm)
		})
	}

	if len(dataList) > 0 {
		tuiService.dataList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventSelectDataRow,
				Data: dataList[index],
			})
		})
	}

	if tuiService.running {
		tuiService.application.Draw()
	}

	if len(dataList) > 0 {
		tuiService.eventBus.Next(&event.Event{
			Name: event.ClientEventSelectDataRow,
			Data: dataList[0],
		})
	}

	tuiService.application.SetFocus(tuiService.dataTypes)
}

// DrawDataRow - отрисовать конкретную запись
func (tuiService *TUIService) DrawDataRow(data models.DataInfo) {
	tuiService.dataForm.SetTitle("Подробно")

	switch data.Type {
	case models.DataTypeCredentials:
		tuiService.drawDataRowCredentials(data)
	case models.DataTypeText:
		tuiService.drawDataRowText(data)
	case models.DataTypeBinary:
		tuiService.drawDataRowBinary(data)
	case models.DataTypeBankCard:
		tuiService.drawDataRowBank(data)
	}

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// DrawCreateForm - отрисовать форму создания записи
func (tuiService *TUIService) DrawCreateForm() {
	tuiService.dataForm.Clear(true)

	tuiService.dataForm.SetTitle("Создание записи")
	tuiService.application.SetFocus(tuiService.dataForm)

	dropdown := tview.NewDropDown().
		SetLabel("Выберите тип данных: ").
		SetOptions([]string{"Учетные данные", "Текстовые данные", "Бинарные данные", "Данные банковских карт"}, func(text string, index int) {
			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventSelectCreateDataType,
				Data: text,
			})
		})
	tuiService.dataForm.AddFormItem(dropdown)

	tuiService.dataForm.SetFocus(0)

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// DrawCreateCredentialsForm - отрисовать форму создания учетных данных
func (tuiService *TUIService) DrawCreateCredentialsForm() {
	tuiService.dataForm.Clear(true)

	tuiService.dataForm.SetTitle("Создание записи")
	tuiService.application.SetFocus(tuiService.dataForm)

	dropdown := tview.NewDropDown().
		SetLabel("Выберите тип данных: ").
		SetOptions([]string{"Учетные данные", "Текстовые данные", "Бинарные данные", "Данные банковских карт"}, func(text string, index int) {
			if text == "Учетные данные" {
				return
			}

			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventSelectCreateDataType,
				Data: text,
			})
		}).
		SetCurrentOption(0)

	data := struct {
		Login    string
		Password string
	}{}
	var description string
	descriptionInput := tview.NewInputField().
		SetLabel("Описание").
		SetChangedFunc(func(text string) {
			description = text
		})
	loginInput := tview.NewInputField().
		SetLabel("Логин").
		SetChangedFunc(func(text string) {
			data.Login = text
		})
	passwordInput := tview.NewInputField().
		SetLabel("Пароль").
		SetChangedFunc(func(text string) {
			data.Password = text
		})

	tuiService.dataForm.AddFormItem(dropdown)
	tuiService.dataForm.AddFormItem(descriptionInput)
	tuiService.dataForm.AddFormItem(loginInput)
	tuiService.dataForm.AddFormItem(passwordInput)
	tuiService.dataForm.AddButton("Сохранить", func() {
		base64Data := tuiService.dataToBase64(data)
		if base64Data == "" {
			return
		}

		tuiService.eventBus.Next(&event.Event{
			Name: event.ClientEventCreateData,
			Data: commonRequests.DataModel{
				Type:        models.DataTypeCredentials,
				Description: description,
				Value:       base64Data,
			},
		})
	})
	tuiService.dataForm.SetFocus(0)

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// DrawCreateTextForm - отрисовать форму создания текстовых данных
func (tuiService *TUIService) DrawCreateTextForm() {
	tuiService.dataForm.Clear(true)

	tuiService.dataForm.SetTitle("Создание записи")
	tuiService.application.SetFocus(tuiService.dataForm)

	dropdown := tview.NewDropDown().
		SetLabel("Выберите тип данных: ").
		SetOptions([]string{"Учетные данные", "Текстовые данные", "Бинарные данные", "Данные банковских карт"}, func(text string, index int) {
			if text == "Текстовые данные" {
				return
			}

			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventSelectCreateDataType,
				Data: text,
			})
		}).
		SetCurrentOption(1)

	data := struct {
		Text string
	}{}
	var description string
	descriptionInput := tview.NewInputField().
		SetLabel("Описание").
		SetChangedFunc(func(text string) {
			description = text
		})
	textInput := tview.NewInputField().
		SetLabel("Текст").
		SetChangedFunc(func(text string) {
			data.Text = text
		})

	tuiService.dataForm.AddFormItem(dropdown)
	tuiService.dataForm.AddFormItem(descriptionInput)
	tuiService.dataForm.AddFormItem(textInput)
	tuiService.dataForm.AddButton("Сохранить", func() {
		base64Data := tuiService.dataToBase64(data)
		if base64Data == "" {
			return
		}

		tuiService.eventBus.Next(&event.Event{
			Name: event.ClientEventCreateData,
			Data: commonRequests.DataModel{
				Type:        models.DataTypeText,
				Description: description,
				Value:       base64Data,
			},
		})
	})
	tuiService.dataForm.SetFocus(0)

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// DrawCreateBinaryForm - отрисовать форму создания бинарных данных
func (tuiService *TUIService) DrawCreateBinaryForm() {
	tuiService.dataForm.Clear(true)

	tuiService.dataForm.SetTitle("Создание записи")
	tuiService.application.SetFocus(tuiService.dataForm)

	dropdown := tview.NewDropDown().
		SetLabel("Выберите тип данных: ").
		SetOptions([]string{"Учетные данные", "Текстовые данные", "Бинарные данные", "Данные банковских карт"}, func(text string, index int) {
			if text == "Бинарные данные" {
				return
			}

			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventSelectCreateDataType,
				Data: text,
			})
		}).
		SetCurrentOption(2)

	data := struct {
		Binary string
	}{}
	var description string
	descriptionInput := tview.NewInputField().
		SetLabel("Описание").
		SetChangedFunc(func(text string) {
			description = text
		})
	binaryInput := tview.NewInputField().
		SetLabel("Данные").
		SetChangedFunc(func(text string) {
			data.Binary = text
		})

	tuiService.dataForm.AddFormItem(dropdown)
	tuiService.dataForm.AddFormItem(descriptionInput)
	tuiService.dataForm.AddFormItem(binaryInput)
	tuiService.dataForm.AddButton("Сохранить", func() {
		base64Data := tuiService.dataToBase64(data)
		if base64Data == "" {
			return
		}

		tuiService.eventBus.Next(&event.Event{
			Name: event.ClientEventCreateData,
			Data: commonRequests.DataModel{
				Type:        models.DataTypeBinary,
				Description: description,
				Value:       base64Data,
			},
		})
	})
	tuiService.dataForm.SetFocus(0)

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// DrawCreateBankForm - отрисовать форму создания данных банковских карт
func (tuiService *TUIService) DrawCreateBankForm() {
	tuiService.dataForm.Clear(true)

	tuiService.dataForm.SetTitle("Создание записи")
	tuiService.application.SetFocus(tuiService.dataForm)

	dropdown := tview.NewDropDown().
		SetLabel("Выберите тип данных: ").
		SetOptions([]string{"Учетные данные", "Текстовые данные", "Бинарные данные", "Данные банковских карт"}, func(text string, index int) {
			if text == "Данные банковских карт" {
				return
			}

			tuiService.eventBus.Next(&event.Event{
				Name: event.ClientEventSelectCreateDataType,
				Data: text,
			})
		}).
		SetCurrentOption(3)

	data := struct {
		Number string
		Date   string
		Secure string
	}{}
	var description string
	descriptionInput := tview.NewInputField().
		SetLabel("Описание").
		SetChangedFunc(func(text string) {
			description = text
		})
	numberInput := tview.NewInputField().
		SetLabel("Номер карты").
		SetChangedFunc(func(text string) {
			data.Number = text
		})
	dateInput := tview.NewInputField().
		SetLabel("Срок действия").
		SetChangedFunc(func(text string) {
			data.Date = text
		})
	secureInput := tview.NewInputField().
		SetLabel("Секретный код").
		SetChangedFunc(func(text string) {
			data.Secure = text
		})

	tuiService.dataForm.AddFormItem(dropdown)
	tuiService.dataForm.AddFormItem(descriptionInput)
	tuiService.dataForm.AddFormItem(numberInput)
	tuiService.dataForm.AddFormItem(dateInput)
	tuiService.dataForm.AddFormItem(secureInput)
	tuiService.dataForm.AddButton("Сохранить", func() {
		base64Data := tuiService.dataToBase64(data)
		if base64Data == "" {
			return
		}

		tuiService.eventBus.Next(&event.Event{
			Name: event.ClientEventCreateData,
			Data: commonRequests.DataModel{
				Type:        models.DataTypeBankCard,
				Description: description,
				Value:       base64Data,
			},
		})
	})
	tuiService.dataForm.SetFocus(0)

	if tuiService.running {
		tuiService.application.Draw()
	}
}

// DataError - отобразить ошибку получения данных
func (tuiService *TUIService) DataError(err string) {
	tuiService.errorPage(err, router.DataPage)
}

// GetCurrentPage - получить имя текущей страницы
func (tuiService *TUIService) GetCurrentPage() string {
	currentPage, _ := tuiService.pages.GetFrontPage()

	return currentPage
}

// Stop - остановить консольное приложение
func (tuiService *TUIService) Stop() {
	tuiService.application.Stop()
	tuiService.running = false
}
