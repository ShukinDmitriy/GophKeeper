package server_test

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"

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
				strings.NewReader(string(bodyJson)),
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
				strings.NewReader(string(bodyJson)),
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

func TestServer(t *testing.T) {
	// Запуск сервера
	_, conf, httpServer := test_helpers.RunServer(t)

	// Запуск тестов сервера
	userRegister(t, conf)
	userLogin(t, conf)

	// Отключаем сервер
	test_helpers.StopServer(t, httpServer)
}
