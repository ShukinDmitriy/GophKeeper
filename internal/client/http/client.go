package http

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	commonRequests "github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/router"

	"github.com/ShukinDmitriy/GophKeeper/internal/client/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/logger"
	"github.com/go-resty/resty/v2"
)

var (
	ErrInvalidAuth   = errors.New(`неправильные имя пользователя и пароль`)
	ErrUserExist     = errors.New(`пользователь существует`)
	ErrServerProblem = errors.New(`попробуйте позже`)
)

// Client - http client
type Client struct {
	config *config.Config
	client *resty.Client
	appLog logger.Logger
}

// NewClient - Создаёт клиента для подключения к серверу по http/HTTPS
func NewClient(
	c *config.Config,
	appLog logger.Logger,
) *Client {
	r := resty.New()
	if c.EnableHTTPS {
		appLog.Info("Use TLS")
		certPath := "ssl/localhost.crt"
		caCert, err := os.ReadFile(certPath)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			appLog.Error("Can't add cert to certpool")
		}
		r = r.SetRootCertificate(certPath)
	}
	return &Client{
		config: c,
		client: r,
		appLog: appLog,
	}
}

// Login - авторизация пользователя
func (hc *Client) Login(ctx context.Context, data commonRequests.UserLogin) error {
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(data).
		Post(fmt.Sprintf("%s%s", hc.config.ServerAddress, router.ApiLoginPath))
	if err != nil {
		return fmt.Errorf("Не удалось выполнить запрос: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return fmt.Errorf("Не удалось авторизоваться: %v", data)
		case http.StatusUnauthorized:
			return fmt.Errorf("Не удалось авторизоваться: %w", ErrInvalidAuth)
		case http.StatusInternalServerError:
			return fmt.Errorf("Не удалось авторизоваться %w", ErrServerProblem)
		}
	}

	cookies := resp.Cookies()
	hc.client.SetCookies(cookies)
	hc.appLog.Debug(fmt.Sprintf("Auth on server, cookies=%v", cookies))

	return nil
}

// Register - регистрация пользователя
func (hc *Client) Register(ctx context.Context, data commonRequests.UserRegister) error {
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(data).
		Post(fmt.Sprintf("%s%s", hc.config.ServerAddress, router.ApiRegisterPath))
	if err != nil {
		return fmt.Errorf("Не удалось выполнить запрос: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return fmt.Errorf("Не удалось зарегистрироваться: %v", data)
		case http.StatusConflict:
			return fmt.Errorf("Не удалось зарегистрироваться: %w", ErrUserExist)
		case http.StatusInternalServerError:
			return fmt.Errorf("Не удалось зарегистрироваться %w", ErrServerProblem)
		}
	}

	cookies := resp.Cookies()
	hc.client.SetCookies(cookies)
	hc.appLog.Debug(fmt.Sprintf("User %s successfuly register", data.Login))

	return nil
}
