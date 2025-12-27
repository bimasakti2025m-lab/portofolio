package delivery

import (
	"database/sql"
	"fmt"

	"enigmacamp.com/livecode-catatan-keuangan/config"
	"enigmacamp.com/livecode-catatan-keuangan/delivery/controller"
	"enigmacamp.com/livecode-catatan-keuangan/delivery/middleware"
	"enigmacamp.com/livecode-catatan-keuangan/repository"
	"enigmacamp.com/livecode-catatan-keuangan/shared/service"
	"enigmacamp.com/livecode-catatan-keuangan/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	expenseUc  usecase.ExpenseUseCase
	userUc     usecase.UserUseCase
	authUsc    usecase.AuthUseCase
	jwtService service.JwtService
	engine     *gin.Engine
	host       string
}

func (s *Server) initRoute() {
	rg := s.engine.Group(config.ApiGroup)
	authMid := middleware.NewAuthMiddleware(s.jwtService)
	controller.NewAuthController(s.authUsc, rg).Route()
	controller.NewUserController(s.userUc, rg, authMid).Route()
	controller.NewExpenseController(s.expenseUc, rg, authMid).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, becauce error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	cfg, _ := config.NewConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		panic("connection error")
	}
	jwtService := service.NewJwtService(cfg.TokenConfig)
	expenseRepo := repository.NewExpenseRepository(db)
	userRepo := repository.NewUserRepository(db)
	taskUC := usecase.NewExpenseUseCase(expenseRepo)
	userUc := usecase.NewUserUseCase(userRepo)
	authUc := usecase.NewAuthUseCase(userUc, jwtService)
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)
	return &Server{
		expenseUc:  taskUC,
		userUc:     userUc,
		authUsc:    authUc,
		jwtService: jwtService,
		engine:     engine,
		host:       host,
	}
}
