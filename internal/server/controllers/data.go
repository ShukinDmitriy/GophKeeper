package controllers

import (
	"errors"
	commonRequests "github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
	"github.com/ShukinDmitriy/GophKeeper/internal/helpers"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/auth"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type DataController struct {
	authService    auth.AuthServiceInterface
	dataRepository repositories.DataRepositoryInterface
}

func NewDataController(
	authService auth.AuthServiceInterface,
	dataRepository repositories.DataRepositoryInterface,
) *DataController {
	return &DataController{
		authService:    authService,
		dataRepository: dataRepository,
	}
}

// DataIndex
// @Title DataIndex
// @Description Получение списка данных
// @Tags Data
// @Accept json
// @Produce json
// @Param data query requests.DataList true "data"
// @Success 200 {array} []responses.DataInfo
// @Failure 400 "BadRequest"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /data [get]
func (controller *DataController) DataIndex() echo.HandlerFunc {
	return func(c echo.Context) error {
		var dataListRequest commonRequests.DataList
		err := c.Bind(&dataListRequest)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		dataListRequest.UserID = controller.authService.GetUserID(c)

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(dataListRequest)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ExtractErrors(err))
		}

		dataInfos, err := controller.dataRepository.List(dataListRequest)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}

		return c.JSON(http.StatusOK, dataInfos)
	}
}

// DataCreate
// @Title DataCreate
// @Description Создать данные
// @Tags Data
// @Accept json
// @Produce json
// @Param data body requests.DataModel true "data"
// @Success 201 {object} responses.DataInfo
// @Failure 400 "BadRequest"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /data [post]
func (controller *DataController) DataCreate() echo.HandlerFunc {
	return func(c echo.Context) error {
		var dataModel commonRequests.DataModel
		err := c.Bind(&dataModel)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		dataModel.UserID = controller.authService.GetUserID(c)

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(dataModel)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helpers.ExtractErrors(err))
		}

		dataInfo, err := controller.dataRepository.Create(dataModel)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}

		return c.JSON(http.StatusCreated, dataInfo)
	}
}

// DataRead
// @Title DataRead
// @Description Получить данные
// @Tags Data
// @Accept json
// @Produce json
// @Param id path number true "id"
// @Success 200 {object} responses.DataInfo
// @Failure 400 "BadRequest"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /data/{id} [get]
func (controller *DataController) DataRead() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		userID := controller.authService.GetUserID(c)

		dataInfo, err := controller.dataRepository.Find(uint(id), userID)
		if err != nil {
			errNotFound := &repositories.NotFoundError{}
			if errors.As(err, &errNotFound) {
				return c.JSON(http.StatusNotFound, "not found")
			}

			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}

		return c.JSON(http.StatusOK, dataInfo)
	}
}

// DataUpdate
// @Title DataUpdate
// @Description Обновить данные
// @Tags Data
// @Accept json
// @Produce json
// @Param id path number true "id"
// @Param data body commonRequests.DataModel true "data"
// @Success 200 {object} responses.DataInfo
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /data/{id} [put]
func (controller *DataController) DataUpdate() echo.HandlerFunc {
	return func(c echo.Context) error {
		var dataModel commonRequests.DataModel
		err := c.Bind(&dataModel)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		dataModel.UserID = controller.authService.GetUserID(c)

		dataInfo, err := controller.dataRepository.Update(uint(id), dataModel)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}

		return c.JSON(http.StatusOK, dataInfo)
	}
}

// DataDelete
// @Title DataDelete
// @Description Удалить данные
// @Tags Data
// @Accept json
// @Produce json
// @Param id path number true "id"
// @Success 202
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /data/{id} [delete]
func (controller *DataController) DataDelete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		userID := controller.authService.GetUserID(c)

		err = controller.dataRepository.Delete(uint(id), userID)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, "internal GophKeeper error")
		}

		return c.JSON(http.StatusAccepted, http.NoBody)
	}
}
