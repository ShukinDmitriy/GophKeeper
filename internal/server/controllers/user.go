package controllers

import (
	"net/http"

	commonRequests "github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"

	"github.com/ShukinDmitriy/GophKeeper/internal/helpers"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/auth"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/models/data"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/models/responses"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	authService    auth.AuthServiceInterface
	userRepository repositories.UserRepositoryInterface
}

func NewUserController(
	authService auth.AuthServiceInterface,
	userRepository repositories.UserRepositoryInterface,
) *UserController {
	return &UserController{
		authService:    authService,
		userRepository: userRepository,
	}
}

// UserRegister
// @Title UserRegister
// @Description Регистрация пользователя
// @Tags User
// @Accept json
// @Produce json
// @Param form body requests.UserRegister true "data"
// @Success 200 {object} responses.UserInfo
// @Failure 400 "Bad request"
// @Failure 409 "Conflict"
// @Failure 500 "Internal server error"
// @Router /user/register [post]
func (controller *UserController) UserRegister() echo.HandlerFunc {
	return func(c echo.Context) error {
		var userRegisterRequest commonRequests.UserRegister
		err := c.Bind(&userRegisterRequest)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(userRegisterRequest)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ExtractErrors(err))
		}

		existUser, err := controller.userRepository.FindBy(data.UserSearch{Login: userRegisterRequest.Login})
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}
		if existUser != nil {
			c.Logger().Error("login already exist")
			return c.JSON(http.StatusConflict, "login already exist")
		}

		user, err := controller.userRepository.Create(userRegisterRequest)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}

		err = controller.authService.GenerateTokensAndSetCookies(c, user)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}

		return c.JSON(http.StatusOK, user)
	}
}

// UserLogin
// @Title UserLogin
// @Description Авторизация пользователя
// @Tags User
// @Accept json
// @Produce json
// @Param form body requests.UserLogin true "data"
// @Success 200 {object} responses.UserInfo
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /user/login [post]
func (controller *UserController) UserLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var userLoginRequest commonRequests.UserLogin
		err := c.Bind(&userLoginRequest)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(userLoginRequest)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ExtractErrors(err))
		}

		existUser, err := controller.userRepository.FindBy(data.UserSearch{Login: userLoginRequest.Login})
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}
		if existUser == nil {
			c.Logger().Error("user not exist")
			return c.JSON(http.StatusUnauthorized, "user not exist")
		}

		if bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(userLoginRequest.Password)) != nil {
			c.Logger().Error("invalid password")
			return c.JSON(http.StatusUnauthorized, "invalid password")
		}
		existUser.Password = ""

		err = controller.authService.GenerateTokensAndSetCookies(c, &responses.UserInfo{
			ID:    existUser.ID,
			Login: existUser.Login,
		})
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}

		return c.JSON(http.StatusOK, existUser)
	}
}
