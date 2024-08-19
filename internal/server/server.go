package server

import (
	"github.com/ShukinDmitriy/GophKeeper/internal/server/auth"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewHTTPServer(
	conf *config.Config,
	authService *auth.AuthService,
	userController *controllers.UserController,
) *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(conf.LogLevel)

	// middleware
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// decompress
	e.Use(middleware.Decompress())

	//jwtMiddleware := echojwt.WithConfig(echojwt.Config{
	//	BeforeFunc: authService.BeforeFunc,
	//	NewClaimsFunc: func(_ echo.Context) jwt.Claims {
	//		return &auth.Claims{}
	//	},
	//	SigningKey:    []byte(auth.GetJWTSecret()),
	//	SigningMethod: jwt.SigningMethodHS256.Alg(),
	//	TokenLookup:   "cookie:access-token", // "<source>:<name>"
	//	ErrorHandler:  authService.JWTErrorChecker,
	//})

	// routes
	// POST /api/user/register — регистрация пользователя;
	// POST /api/user/login — аутентификация пользователя;

	e.POST("/api/user/register", userController.UserRegister())
	e.POST("/api/user/login", userController.UserLogin())
	// e.POST("/api/user/orders", orderController.CreateOrder(), jwtMiddleware)

	return e
}
