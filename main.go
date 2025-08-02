package main

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"saunalog/handler"
	"saunalog/usecase"
)

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
