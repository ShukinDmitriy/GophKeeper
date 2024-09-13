package auth

import (
	"github.com/ShukinDmitriy/GophKeeper/internal/server/models/responses"
	"github.com/labstack/echo/v4"
)

type AuthServiceInterface interface {
	GetUserID(c echo.Context) uint
	GenerateTokensAndSetCookies(c echo.Context, user *responses.UserInfo) error
}
