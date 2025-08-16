package router

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"saunalog/handler"
	"saunalog/infra/db"
	"saunalog/usecase"
)

type AppDB struct {
	DB *sql.DB
}

func Router(e *echo.Echo, apps AppDB) {
	api := e.Group("/api")
	v1 := api.Group("/v1")

	{
		userRepo := db.NewUserRepo(apps.DB)
		userUC := usecase.NewUserCreate(userRepo)
		userHandler := handler.UserHandler{
			Usecase: userUC,
		}
		v1.POST("/users", userHandler.CreateUser)
	}

	{
		logUC := usecase.NewExperienceLogUseCase()
		logHandler := handler.ExperienceLogHandler{
			Usecase: logUC,
		}
		v1.POST("/logs", logHandler.CreateExperienceLog)
	}

}

func Listen(addr string) error {
	mysql, err := db.NewMySQLFromEnv()
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer mysql.Close()

	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	Router(e, AppDB{DB: mysql})

	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server start: %w", err)
	}
	return nil
}

func Start() {
	if err := Listen(":8080"); err != nil {
		log.Fatalf("listen error: %v", err)
	}
}
