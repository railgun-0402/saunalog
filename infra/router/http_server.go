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
	"saunalog/usecase"
)

type router struct {
	e *echo.Echo
	d *db.UserRepo
}

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

func NewRouter(d *db.UserRepo) *router {
	return &router{
		e: echo.New(),
		d: d,
	}
}

func (s *Server) router() {
	logHandler := handler.ExperienceLogHandler{
		Usecase: s.logUC,
	}
	s.e.POST("/logs", logHandler.CreateExperienceLog)

}

func (r *router) userRouter() {
	uc := usecase.NewUserCreate(r.d)
	r.e.POST("/users", uc.Execute)
}

func (r *router) Listen(addr string) error {
	r.userRouter()
	if err := r.e.Start(addr); err != nil && err != http.ErrServerClosed {
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
	userRepo := &db.UserRepo{DB: mysql}
	_ = userRepo

	//logUC := usecase.NewExperienceLogUseCase()

	//srv := NewServer(logUC)
	srv := NewRouter(userRepo)
	if err := srv.Listen(":8080"); err != nil {
		log.Fatalf("listen error: %v", err)
	}

}
