package tui

import (
	"github.com/ShukinDmitriy/GophKeeper/internal/client/event"
	"github.com/ShukinDmitriy/GophKeeper/internal/client/router"
	commonRequests "github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
	"github.com/ShukinDmitriy/GophKeeper/internal/logger"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TUIService сервис для работы с консолью
type TUIService struct {
	application *tview.Application
	appLog      logger.Logger
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
func (t *TUIService) errorPage(err string, returnToPage string) {
	t.appLog.Debug("Create error page")

	form := tview.NewForm().
		AddTextView("Ошибка", err, 50, 5, true, true).
		AddButton("Понятно", func() {
			t.pages.SwitchToPage(returnToPage)
		})

	t.pages.AddAndSwitchToPage(router.ErrorPage, form, true)

	if t.running {
		t.application.Draw()
	}
}

// Run - запустить консольное приложение
func (t *TUIService) Run() error {
	t.running = true
	err := t.application.Run()
	if err != nil {
		return err
	}

	return nil
}

// LoginPage - отобразить страницу авторизации
func (t *TUIService) LoginPage() {
	if !t.pages.HasPage(router.LoginPage) {
		t.appLog.Debug("Create login page")
		data := commonRequests.UserLogin{}
		form := tview.NewForm().
			AddInputField("Логин", "", 20, nil, func(text string) {
				data.Login = text
			}).
			AddPasswordField("Пароль", "", 20, '*', func(text string) {
				data.Password = text
			}).
			AddButton("Авторизоваться", func() {
				t.appLog.Debug("Press Login button")
				t.eventBus.Next(&event.Event{
					Name: event.ClientEventPressLoginButton,
					Data: data,
				})
			}).
			AddButton("Перейти к регистрации", func() {
				t.appLog.Debug("Press Register button")
				t.eventBus.Next(&event.Event{
					Name: event.ClientEventPressToRegisterButton,
				})
			}).
			AddButton("Закончить", func() {
				t.appLog.Debug("Press Stop button")
				t.application.Stop()
			})

		t.pages.AddPage(router.LoginPage, form, true, false)
	}

	t.pages.SwitchToPage(router.LoginPage)

	if t.running {
		t.application.Draw()
	}
}

// LoginError - отобразить ошибку авторизации
func (t *TUIService) LoginError(err string) {
	t.errorPage(err, router.LoginPage)
}

// RegisterPage - отобразить страницу регистрации
func (t *TUIService) RegisterPage() {
	if !t.pages.HasPage(router.RegisterPage) {
		t.appLog.Debug("Create register page")
		data := commonRequests.UserRegister{}
		form := tview.NewForm().
			AddInputField("Логин", "", 20, nil, func(text string) {
				data.Login = text
			}).
			AddPasswordField("Пароль", "", 20, '*', func(text string) {
				data.Password = text
			}).
			AddButton("Зарегистрироваться", func() {
				t.appLog.Debug("Press Register button")
				t.eventBus.Next(&event.Event{
					Name: event.ClientEventPressRegisterButton,
					Data: data,
				})
			}).
			AddButton("Перейти к авторизации", func() {
				t.appLog.Debug("Press \"To login\" button")
				t.eventBus.Next(&event.Event{
					Name: event.ClientEventPressToLoginButton,
				})
			}).
			AddButton("Закончить", func() {
				t.appLog.Debug("Press Stop button")
				t.Stop()
			})

		t.pages.AddPage(router.RegisterPage, form, true, false)
	}

	t.pages.SwitchToPage(router.RegisterPage)

	if t.running {
		t.application.Draw()
	}
}

// RegisterError - отобразить ошибку регистрации
func (t *TUIService) RegisterError(err string) {
	t.errorPage(err, router.RegisterPage)
}

// GetCurrentPage - получить имя текущей страницы
func (t *TUIService) GetCurrentPage() string {
	currentPage, _ := t.pages.GetFrontPage()

	return currentPage
}

// Stop - остановить консольное приложение
func (t *TUIService) Stop() {
	t.application.Stop()
	t.running = false
}
