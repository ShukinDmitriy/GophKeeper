package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"path"

	"github.com/ShukinDmitriy/GophKeeper/internal/server"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/auth"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/controllers"
	appLogger "github.com/ShukinDmitriy/GophKeeper/internal/server/logger"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/repositories"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// main Запуск сервера
func main() {
	fx.New(
		fx.Provide(
			// Конфигурация
			func() *config.Config {
				conf, err := config.NewConfig()
				if err != nil {
					panic(err)
				}
				return conf
			},
			// логгер
			func(conf *config.Config) appLogger.Logger {
				return appLogger.NewLogger(conf)
			},
			// база данных
			func(
				conf *config.Config,
				appLog appLogger.Logger,
			) *gorm.DB {
				postgresqlURL := conf.DatabaseURI

				if postgresqlURL == "" {
					appLog.Fatal("no DATABASE_URI in .env")
					return nil
				}

				var level logger.LogLevel
				switch conf.LogLevel {
				case log.DEBUG:
				case log.INFO:
					level = logger.Info
				case log.WARN:
					level = logger.Warn
				case log.ERROR:
					level = logger.Error
				case log.OFF:
					level = logger.Silent
				}

				db, err := gorm.Open(postgres.Open(postgresqlURL), &gorm.Config{
					Logger: logger.Default.LogMode(level),
				})
				if err != nil {
					appLog.Fatal(err)
					return nil
				}

				return db
			},
			// Репозиторий пользователя
			func(DB *gorm.DB) *repositories.UserRepository {
				return repositories.NewUserRepository(DB)
			},
			// Аутентификация
			func(userRepository *repositories.UserRepository) *auth.AuthUser {
				return auth.NewAuthUser(userRepository)
			},
			// Сервис работы с аутентификацией
			func(authUser *auth.AuthUser) *auth.AuthService {
				return auth.NewAuthService(*authUser)
			},
			// Контроллер пользователя
			func(
				authService *auth.AuthService,
				userRepository *repositories.UserRepository,
			) *controllers.UserController {
				return controllers.NewUserController(
					authService,
					userRepository,
				)
			},
			// http сервер
			func(
				lc fx.Lifecycle,
				conf *config.Config,
				appLog appLogger.Logger,
				authService *auth.AuthService,
				userController *controllers.UserController,
			) *echo.Echo {
				httpServer := server.NewHTTPServer(
					conf,
					authService,
					userController,
				)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := httpServer.Start(conf.RunAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
								appLog.Fatal("shutting down the GophKeeper ", err.Error())
							}

							appLog.Info("Running GophKeeper")
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						return httpServer.Shutdown(ctx)
					},
				})

				return httpServer
			},
		),
		// Запускаем сервер
		fx.Invoke(func(*echo.Echo) {}),
		// Запускаем миграции
		fx.Invoke(func(
			appLog appLogger.Logger,
			config *config.Config,
		) error {
			db, err := sql.Open("postgres", config.DatabaseURI)
			if err != nil {
				appLog.Error("can't connect to db", err.Error())
				return err
			}
			defer func() {
				db.Close()
			}()

			driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
			if err != nil {
				appLog.Error("can't create driver", err.Error())
				return err
			}

			currentDir, _ := os.Getwd()
			m, err := migrate.NewWithDatabaseInstance(
				"file:///"+path.Join(currentDir, "db", "migrations"),
				"postgres", driver)
			if err != nil {
				appLog.Error("can't create new migrate: ", err.Error())
				return err
			}

			err = m.Up()
			if err != nil {
				appLog.Info("can't migrate up: ", err.Error())
			}

			return nil
		}),
	).Run()
}
