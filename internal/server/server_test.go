package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/ShukinDmitriy/GophKeeper/internal/common/models"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/router"
	"github.com/stretchr/testify/assert"

	"github.com/ShukinDmitriy/GophKeeper/internal/test-helpers"

	"github.com/ShukinDmitriy/GophKeeper/internal/server/config"
)

func userRegister(t *testing.T, conf *config.Config) {
	type args struct {
		login    string
		password string
	}
	type want struct {
		status []int
	}
	userLogin := test_helpers.GenerateRandomString(10)
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Default user",
			args: args{
				login:    "login",
				password: "password",
			},
			want: want{
				status: []int{http.StatusOK, http.StatusConflict},
			},
		},
		{
			name: "Success registration",
			args: args{
				login:    userLogin,
				password: test_helpers.GenerateRandomString(10),
			},
			want: want{
				status: []int{http.StatusOK},
			},
		},
		{
			name: "User exists",
			args: args{
				login:    userLogin,
				password: test_helpers.GenerateRandomString(10),
			},
			want: want{
				status: []int{http.StatusConflict},
			},
		},
		{
			name: "Validation error #1",
			args: args{
				login:    test_helpers.GenerateRandomString(3),
				password: test_helpers.GenerateRandomString(10),
			},
			want: want{
				status: []int{http.StatusBadRequest},
			},
		},
		{
			name: "Validation error #2",
			args: args{
				login:    test_helpers.GenerateRandomString(10),
				password: test_helpers.GenerateRandomString(3),
			},
			want: want{
				status: []int{http.StatusBadRequest},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Отправляем POST-запрос
			body := map[string]string{
				"login":    tt.args.login,
				"password": tt.args.password,
			}
			bodyJson, _ := json.Marshal(body)
			resp, err := http.Post(
				test_helpers.PrepareURL(conf, "/api/user/register"),
				"application/json",
				bytes.NewReader(bodyJson),
			)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Проверяем статус ответа
			if !slices.Contains(tt.want.status, resp.StatusCode) {
				t.Errorf("Expected status code %d, got %d", tt.want.status, resp.StatusCode)
			}
		})
	}
}

func userLogin(t *testing.T, conf *config.Config) {
	type args struct {
		login    string
		password string
	}
	type want struct {
		status []int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Default user",
			args: args{
				login:    "login",
				password: "password",
			},
			want: want{
				status: []int{http.StatusOK},
			},
		},
		{
			name: "invalid password",
			args: args{
				login:    "login",
				password: "invalidPassword",
			},
			want: want{
				status: []int{http.StatusUnauthorized},
			},
		},
		{
			name: "User not exists",
			args: args{
				login:    test_helpers.GenerateRandomString(10),
				password: test_helpers.GenerateRandomString(10),
			},
			want: want{
				status: []int{http.StatusUnauthorized},
			},
		},
		{
			name: "Validation error",
			args: args{
				login:    test_helpers.GenerateRandomString(3),
				password: test_helpers.GenerateRandomString(3),
			},
			want: want{
				status: []int{http.StatusBadRequest},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Отправляем POST-запрос
			body := map[string]string{
				"login":    tt.args.login,
				"password": tt.args.password,
			}
			bodyJson, _ := json.Marshal(body)
			resp, err := http.Post(
				test_helpers.PrepareURL(conf, "/api/user/login"),
				"application/json",
				bytes.NewReader(bodyJson),
			)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Проверяем статус ответа
			if !slices.Contains(tt.want.status, resp.StatusCode) {
				t.Errorf("Expected status code %d, got %d", tt.want.status, resp.StatusCode)
			}
		})
	}
}

func userRefresh(t *testing.T, conf *config.Config) {
	// Авторизация
	body := map[string]string{
		"login":    "login",
		"password": "password",
	}
	bodyJson, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", test_helpers.PrepareURL(conf, router.ApiLoginPath), bytes.NewReader(bodyJson))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	// Сохранение cookie в переменную
	var cookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "refresh-token" {
			cookie = c
			break
		}
	}

	assert.NotNil(t, cookie)

	t.Run("Get list data with only refresh token", func(t *testing.T) {
		// Запрос для получения данных
		req, err = http.NewRequest("GET", test_helpers.PrepareURL(conf, router.ApiDataListPath), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Добавление cookie к запросу
		req.AddCookie(cookie)

		// Выполнение запроса получения данных
		resp, err = client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		// Проверяем статус ответа
		successStatuses := []int{http.StatusOK, http.StatusNoContent}
		if !slices.Contains(successStatuses, resp.StatusCode) {
			t.Errorf("Expected status code %d, got %d", successStatuses, resp.StatusCode)
		}
	})
}

func dataCRUD(t *testing.T, conf *config.Config) {
	// Авторизация
	body := map[string]string{
		"login":    "login",
		"password": "password",
	}
	bodyJson, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", test_helpers.PrepareURL(conf, router.ApiLoginPath), bytes.NewReader(bodyJson))
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	// Сохранение cookie в переменную
	var cookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "access-token" {
			cookie = c
			break
		}
	}

	assert.NotNil(t, cookie)

	t.Run("Get empty list data", func(t *testing.T) {
		// Запрос для получения данных
		req, err = http.NewRequest("GET", test_helpers.PrepareURL(conf, router.ApiDataListPath), nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Добавление cookie к запросу
		req.AddCookie(cookie)

		// Выполнение запроса получения данных
		resp, err = client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		// Чтение ответа
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		var list []models.DataInfo
		err = json.Unmarshal(resBody, &list)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, 0, len(list))
	})

	var lastID uint

	t.Run("Success create data", func(t *testing.T) {
		// Запрос для создания данных
		data := requests.DataModel{
			Type:        models.DataTypeText,
			Description: "test data",
			Value:       "test value",
		}
		dataJson, _ := json.Marshal(data)
		req, err = http.NewRequest("POST", test_helpers.PrepareURL(conf, router.ApiDataCreatePath), bytes.NewReader(dataJson))
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Добавление cookie к запросу
		req.AddCookie(cookie)

		// Выполнение запроса получения данных
		resp, err = client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		// Чтение ответа
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		var resData models.DataInfo
		err = json.Unmarshal(resBody, &resData)
		if err != nil {
			t.Error(err)
			return
		}

		assert.NotEqual(t, resData.ID, 0)
		assert.Equal(t, resData.Type, data.Type)
		assert.Equal(t, resData.Description, data.Description)
		assert.Equal(t, resData.Value, data.Value)

		lastID = resData.ID
	})

	t.Run("Success get data", func(t *testing.T) {
		// Запрос для получения данных
		url := test_helpers.PrepareURL(conf, router.ApiDataReadPath)
		url = strings.Replace(url, ":id", strconv.Itoa(int(lastID)), 1)
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Добавление cookie к запросу
		req.AddCookie(cookie)

		// Выполнение запроса получения данных
		resp, err = client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		// Чтение ответа
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		var reqData models.DataInfo
		err = json.Unmarshal(resBody, &reqData)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, reqData.ID, lastID)
	})

	t.Run("Getting non-existent data", func(t *testing.T) {
		// Запрос для получения данных
		url := test_helpers.PrepareURL(conf, router.ApiDataReadPath)
		url = strings.Replace(url, ":id", strconv.Itoa(int(lastID+1)), 1)
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Добавление cookie к запросу
		req.AddCookie(cookie)

		// Выполнение запроса получения данных
		resp, err = client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusNotFound)
	})

	t.Run("Success update data", func(t *testing.T) {
		// Запрос для создания данных
		data := requests.DataModel{
			Type:        models.DataTypeText,
			Description: "test data updated",
			Value:       "test value updated",
		}
		dataJson, _ := json.Marshal(data)
		url := test_helpers.PrepareURL(conf, router.ApiDataUpdatePath)
		url = strings.Replace(url, ":id", strconv.Itoa(int(lastID)), 1)
		req, err = http.NewRequest("PUT", url, bytes.NewReader(dataJson))
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Добавление cookie к запросу
		req.AddCookie(cookie)

		// Выполнение запроса получения данных
		resp, err = client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		// Чтение ответа
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		var resData models.DataInfo
		err = json.Unmarshal(resBody, &resData)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, resData.ID, lastID)
		assert.Equal(t, resData.Type, data.Type)
		assert.Equal(t, resData.Description, data.Description)
		assert.Equal(t, resData.Value, data.Value)
	})

	t.Run("Success delete data", func(t *testing.T) {
		// Запрос для создания данных
		url := test_helpers.PrepareURL(conf, router.ApiDataUpdatePath)
		url = strings.Replace(url, ":id", strconv.Itoa(int(lastID)), 1)
		req, err = http.NewRequest("DELETE", url, nil)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		// Добавление cookie к запросу
		req.AddCookie(cookie)

		// Выполнение запроса получения данных
		resp, err = client.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusAccepted)
	})
}

func TestServer(t *testing.T) {
	// Запуск сервера
	_, conf, httpServer := test_helpers.RunServer(t)

	// Запуск тестов сервера
	userRegister(t, conf)
	userLogin(t, conf)
	userRefresh(t, conf)
	dataCRUD(t, conf)

	// Отключаем сервер
	test_helpers.StopServer(t, httpServer)
}
