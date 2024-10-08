package server

import (
	"strings"

	"github.com/ShukinDmitriy/GophKeeper/internal/common/router"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/auth"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/controllers"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewHTTPServer(
	conf *config.Config,
	authService *auth.AuthService,
	userController *controllers.UserController,
	dataController *controllers.DataController,
) *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(conf.LogLevel)

	// middleware
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}

			skipByAcceptEncodingHeader := true
			skipByContentTypeHeader := true

			acceptEncodingRaw := c.Request().Header.Get("Accept-Encoding")
			acceptEncodingValues := strings.Split(acceptEncodingRaw, ",")

			for _, value := range acceptEncodingValues {
				parts := strings.Split(value, ";")
				format := strings.TrimSpace(parts[0])

				if format == "gzip" {
					skipByAcceptEncodingHeader = false
					break
				}
			}

			contentTypeRaw := c.Request().Header.Get("Content-Type")
			contentTypeValues := strings.Split(contentTypeRaw, ",")

			for _, value := range contentTypeValues {
				if value == "application/json" || value == "text/html" {
					skipByContentTypeHeader = false
					break
				}
			}

			return skipByAcceptEncodingHeader && skipByContentTypeHeader
		},
	}))

	// decompress
	e.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{
		Skipper: func(c echo.Context) bool {
      return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		BeforeFunc: authService.BeforeFunc,
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
			return &auth.Claims{}
		},
		SigningKey:    []byte(auth.GetJWTSecret()),
		SigningMethod: jwt.SigningMethodHS256.Alg(),
		TokenLookup:   "cookie:access-token", // "<source>:<name>"
		ErrorHandler:  authService.JWTErrorChecker,
	})

	// routes
	// GET /swagger — swagger;
	// POST /api/user/register — регистрация пользователя;
	// POST /api/user/login — аутентификация пользователя;
	// GET /api/data — список данных;
	// POST /api/data — создать данные;
	// GET /api/data/:id — получить данные;
	// PUT /api/data/:id — изменить данные;
	// DELETE /api/data/:id — удалить данные;

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST(router.ApiRegisterPath, userController.UserRegister())
	e.POST(router.ApiLoginPath, userController.UserLogin())
	e.GET(router.ApiDataListPath, dataController.DataIndex(), jwtMiddleware)
	e.POST(router.ApiDataCreatePath, dataController.DataCreate(), jwtMiddleware)
	e.GET(router.ApiDataReadPath, dataController.DataRead(), jwtMiddleware)
	e.PUT(router.ApiDataUpdatePath, dataController.DataUpdate(), jwtMiddleware)
	e.DELETE(router.ApiDataDeletePath, dataController.DataDelete(), jwtMiddleware)

	return e
}
