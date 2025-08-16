package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	domain "saunalog/domain/user"
	"saunalog/usecase"
	"time"
)

type UserHandler struct {
	Usecase *usecase.UserUsecase
}

type UserCreateRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Gender     string `json:"gender"`
	Age        int    `json:"age"`
	Prefecture string `json:"prefecture"`
}

type UserCreateResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "400 Bad Request")
	}

	user, err := domain.NewUser(domain.User{
		ID:         "",
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		Gender:     req.Gender,
		Age:        req.Age,
		Prefecture: req.Prefecture,
		CreatedAt:  time.Now(),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	out, err := h.Usecase.Create(c.Request().Context(), user)
	if err != nil {
		// TODO: エラー型に応じて 400/409/500 などへマッピング
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, UserCreateResponse{
		ID:      string(out.ID),
		Message: "ユーザーの新規作成に成功しました",
	})
}
