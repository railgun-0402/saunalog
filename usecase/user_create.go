package usecase

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	domain "saunalog/domain/user"
	"saunalog/infra/db"
)

type UserCreate struct {
	r *db.UserRepo
}

type UserRequest struct {
	domain.User
}

type UserResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func NewUserCreate(r *db.UserRepo) *UserCreate {
	return &UserCreate{r: r}
}

func (u *UserCreate) Execute(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	userID := domain.UserID(uuid.New().String())

	user, err := domain.NewUser(
		userID,
		req.Name,
		req.Email,
		req.Password,
		req.Gender,
		req.Age,
		req.Prefecture,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newUser, err := u.r.CreateUser(c.Request().Context(), user)
	fmt.Printf("Created User: %+\n", newUser)

	return c.JSON(http.StatusCreated, UserResponse{
		ID:      string(newUser.ID),
		Message: "ユーザーの新規作成に成功しました",
	})
}
