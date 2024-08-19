package server_test

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ShukinDmitriy/GophKeeper/internal/server"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/auth"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/controllers"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/repositories"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GenerateShortKey generate random string
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	shortKey := make([]byte, length)
	for i := range shortKey {
		shortKey[i] = charset[rng.Intn(len(charset))]
	}

	return string(shortKey)
}

func prepareURL(conf *config.Config, url string) string {
	return "http://" + conf.RunAddress + url
}

func createConfig(t *testing.T) *config.Config {
	conf, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	return conf
}

func createServer(t *testing.T, conf *config.Config) *echo.Echo {
	postgresqlURL := conf.DatabaseURI

	if postgresqlURL == "" {
		t.Fatal("no DATABASE_URI in .env")
	}
	db, err := gorm.Open(postgres.Open(postgresqlURL), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	userRepository := repositories.NewUserRepository(db)
	authUser := auth.NewAuthUser(userRepository)
	authService := auth.NewAuthService(*authUser)
	userController := controllers.NewUserController(
		authService,
		userRepository,
	)

	httpServer := server.NewHTTPServer(
		conf,
		authService,
		userController,
	)

	go func() {
		_ = httpServer.Start(conf.RunAddress)
	}()

	// Ждем запуска сервера
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatal("timed out waiting for server to start")
		default:
			resp, err := http.Get(prepareURL(conf, "/"))
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			return httpServer
		}
	}
}

func userRegister(t *testing.T, conf *config.Config) {
	type args struct {
		login    string
		password string
	}
	type want struct {
		status []int
	}
	userLogin := generateRandomString(10)
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
				password: generateRandomString(10),
			},
			want: want{
				status: []int{http.StatusOK},
			},
		},
		{
			name: "User exists",
			args: args{
				login:    userLogin,
				password: generateRandomString(10),
			},
			want: want{
				status: []int{http.StatusConflict},
			},
		},
		{
			name: "Validation error #1",
			args: args{
				login:    generateRandomString(3),
				password: generateRandomString(10),
			},
			want: want{
				status: []int{http.StatusBadRequest},
			},
		},
		{
			name: "Validation error #2",
			args: args{
				login:    generateRandomString(10),
				password: generateRandomString(3),
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
				prepareURL(conf, "/api/user/register"),
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
				login:    generateRandomString(10),
				password: generateRandomString(10),
			},
			want: want{
				status: []int{http.StatusUnauthorized},
			},
		},
		{
			name: "Validation error",
			args: args{
				login:    generateRandomString(3),
				password: generateRandomString(3),
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
				prepareURL(conf, "/api/user/login"),
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
	// Создаем конфигурацию
	conf := createConfig(t)
	// Запускаем сервер
	httpServer := createServer(t, conf)

	// Отключаем сервер
	defer func() {
		ctx := context.TODO()

		_ = httpServer.Shutdown(ctx)
	}()

	// Запуск тестов
	userRegister(t, conf)
	userLogin(t, conf)
}
