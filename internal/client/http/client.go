package http

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ShukinDmitriy/GophKeeper/internal/common/models"

	commonRequests "github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/router"

	"github.com/ShukinDmitriy/GophKeeper/internal/client/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/logger"
	"github.com/go-resty/resty/v2"
)

var (
	ErrInvalidAuth      = errors.New(`неправильные имя пользователя и пароль`)
	ErrUserExist        = errors.New(`пользователь существует`)
	ErrUserUnauthorized = errors.New(`пользователь не авторизован`)
	ErrServerProblem    = errors.New(`попробуйте позже`)
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
	hc.appLog.Debug(fmt.Sprintf("User %s successfully register", data.Login))

	return nil
}

// GetList - получить список данных по типу
func (hc *Client) GetList(ctx context.Context, dataType models.DataType) ([]models.DataInfo, error) {
	var dataList []models.DataInfo

	resp, err := hc.client.R().
		SetContext(ctx).
		Get(fmt.Sprintf("%s%s?type=%v", hc.config.ServerAddress, router.ApiDataListPath, dataType))
	if err != nil {
		return dataList, fmt.Errorf("Не удалось выполнить запрос: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return dataList, fmt.Errorf("Не удалось получить данные: %v", dataType)
		case http.StatusUnauthorized:
			return dataList, fmt.Errorf("Не удалось получить данные: %w", ErrUserUnauthorized)
		case http.StatusInternalServerError:
			return dataList, fmt.Errorf("Не удалось получить данные %w", ErrServerProblem)
		}
	}

	err = json.Unmarshal(resp.Body(), &dataList)
	if err != nil {
		return dataList, fmt.Errorf("Не удалось разобрать ответ: %w", err)
	}

	hc.appLog.Debug(fmt.Sprintf("Data %v successfully getting:", dataList))

	return dataList, nil
}

// CreateData - создать новую запись
func (hc *Client) CreateData(ctx context.Context, data commonRequests.DataModel) (*models.DataInfo, error) {
	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(data).
		Post(fmt.Sprintf("%s%s", hc.config.ServerAddress, router.ApiDataCreatePath))
	if err != nil {
		return nil, fmt.Errorf("Не удалось выполнить запрос: %w", err)
	}
	if resp.StatusCode() != http.StatusCreated {
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return nil, fmt.Errorf("Не удалось создать запись: %v", data)
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("Не удалось создать запись: %w", ErrUserUnauthorized)
		case http.StatusInternalServerError:
			return nil, fmt.Errorf("Не удалось создать запись %w", ErrServerProblem)
		}
	}

	resData := &models.DataInfo{}
	err = json.Unmarshal(resp.Body(), &resData)
	if err != nil {
		return resData, fmt.Errorf("Не удалось разобрать ответ: %w", err)
	}

	hc.appLog.Debug(fmt.Sprintf("Создана запись: %v", data))

	return resData, nil
}

// UpdateData - изменить данные
func (hc *Client) UpdateData(ctx context.Context, data models.DataInfo) (*models.DataInfo, error) {
	url := fmt.Sprintf("%s%s", hc.config.ServerAddress, router.ApiDataUpdatePath)
	url = strings.Replace(url, ":id", fmt.Sprintf("%d", data.ID), 1)

	resp, err := hc.client.R().
		SetContext(ctx).
		SetBody(data).
		Put(url)
	if err != nil {
		return nil, fmt.Errorf("Не удалось выполнить запрос: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return nil, fmt.Errorf("Не удалось изменить запись: %v", data)
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("Не удалось изменить запись: %w", ErrUserUnauthorized)
		case http.StatusInternalServerError:
			return nil, fmt.Errorf("Не удалось изменить запись %w", ErrServerProblem)
		}
	}

	resData := &models.DataInfo{}
	err = json.Unmarshal(resp.Body(), &resData)
	if err != nil {
		return resData, fmt.Errorf("Не удалось разобрать ответ: %w", err)
	}

	hc.appLog.Debug(fmt.Sprintf("Изменена запись: %v", data))

	return resData, nil
}

// DeleteData - удалить данные
func (hc *Client) DeleteData(ctx context.Context, id uint) error {
	url := fmt.Sprintf("%s%s", hc.config.ServerAddress, router.ApiDataDeletePath)
	url = strings.Replace(url, ":id", fmt.Sprintf("%d", id), 1)

	resp, err := hc.client.R().
		SetContext(ctx).
		Delete(url)
	if err != nil {
		return fmt.Errorf("Не удалось выполнить запрос: %w", err)
	}
	if resp.StatusCode() != http.StatusAccepted {
		switch resp.StatusCode() {
		case http.StatusBadRequest:
			return fmt.Errorf("Не удалось удалить запись: %v", id)
		case http.StatusUnauthorized:
			return fmt.Errorf("Не удалось удалить запись: %w", ErrUserUnauthorized)
		case http.StatusInternalServerError:
			return fmt.Errorf("Не удалось удалить запись %w", ErrServerProblem)
		}
	}

	hc.appLog.Debug(fmt.Sprintf("Удалена запись: %v", id))

	return nil
}
