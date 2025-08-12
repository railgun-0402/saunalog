package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	domain "saunalog/domain/user"
	"saunalog/handler"
	infra "saunalog/infra/db"
	"saunalog/usecase"
)

// TODO: APIにする
func tryUserCreate() {
	userRepo := infra.NewUserRepo(infra.Conn)

	u := &domain.User{
		Name:       "テスト太郎",
		Email:      "test@example.com",
		Gender:     "male",
		Age:        28,
		Password:   "hashed-password",
		Prefecture: "Tokyo",
	}

	newUser, err := userRepo.CreateUser(context.Background(), u)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created User: &+\n", newUser)
}

func main() {
	e := echo.New()

	logUsecase := usecase.NewExperienceLogUseCase()
	logHandler := handler.ExperienceLogHandler{
		Usecase: logUsecase,
	}
	e.POST("/logs", logHandler.CreateExperienceLog)

	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed Server Start: %v", err)
	}
}
