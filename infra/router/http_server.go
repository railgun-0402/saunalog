package router

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	domain "saunalog/domain/user"
	"saunalog/handler"
	"saunalog/infra/db"
	"saunalog/usecase/repository"

	"saunalog/usecase"
)

type Server struct {
	e     *echo.Echo
	logUC usecase.ExperienceLogUseCase
}

// TODO: APIにする
func tryUserCreate() {
	userRepo := db.NewUserRepo(db.Conn)

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

func NewServer(logUC usecase.ExperienceLogUseCase) *Server {
	return &Server{
		e:     echo.New(),
		logUC: logUC,
	}
}

func (s *Server) router() {
	logHandler := handler.ExperienceLogHandler{
		Usecase: s.logUC,
	}
	s.e.POST("/logs", logHandler.CreateExperienceLog)
}

func (s *Server) Listen(addr string) error {
	s.router()
	if err := s.e.Start(addr); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server start: %w", err)
	}
	return nil
}

func Start() {
	mysql, err := db.NewMySQLFromEnv()
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer mysql.Close()

	// TODO: APIにする
	var userRepo repository.UserRepository = db.NewUserRepo(mysql)
	_ = userRepo

	logUC := usecase.NewExperienceLogUseCase()

	srv := NewServer(logUC)
	if err := srv.Listen(":8080"); err != nil {
		log.Fatalf("listen error: %v", err)
	}

}
