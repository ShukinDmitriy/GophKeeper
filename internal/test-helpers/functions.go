package test_helpers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
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

func createServerConfig(t *testing.T) *config.Config {
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
			resp, err := http.Get(PrepareURL(conf, "/"))
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			return httpServer
		}
	}
}

// GenerateShortKey generate random string
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	shortKey := make([]byte, length)
	for i := range shortKey {
		shortKey[i] = charset[rng.Intn(len(charset))]
	}

	return string(shortKey)
}

// PrepareURL - подготовить URL для запроса
func PrepareURL(conf *config.Config, url string) string {
	return "http://" + conf.RunAddress + url
}

// RunServer - запуск сервера на случайном порту
func RunServer(t *testing.T) (int, *config.Config, *echo.Echo) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	port := random.Intn(60000-1024) + 1024

	// Создаем конфигурацию
	conf := createServerConfig(t)
	conf.RunAddress = fmt.Sprintf("localhost:%d", port)

	// Запускаем сервер
	httpServer := createServer(t, conf)

	return port, conf, httpServer
}

// StopServer - остановить запущенный вебсервер
func StopServer(t *testing.T, httpServer *echo.Echo) {
	ctx := context.TODO()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Failed to shutdown server: %v", err)
	}
}
