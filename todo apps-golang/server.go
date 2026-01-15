package main

import (
	"database/sql"
	"fmt"

	"enigmacamp.com/golang-jwt/config"
	"enigmacamp.com/golang-jwt/controller"
	"enigmacamp.com/golang-jwt/middleware"
	"enigmacamp.com/golang-jwt/repository"
	"enigmacamp.com/golang-jwt/usecase"
	"enigmacamp.com/golang-jwt/utils/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Mendeklarasikan struct
type Server struct {
	todosUC usecase.TodosUseCase
	userUC  usecase.UserUseCase
	authUC  usecase.AuthenticateUsecase
	jwtSvc  service.JwtService
	engine  *gin.Engine
	host    string
}

// Mendeklarasikan method initRoute
func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")

	authMiddleware := middleware.NewAuthMiddleware(s.jwtSvc)
	controller.NewTodosController(s.todosUC, rg, authMiddleware).Route() // This line was already present
	controller.NewUserController(s.userUC, rg, authMiddleware).Route()
	controller.NewAuthController(s.authUC, rg).Route()
}

// Mendeklarasikan method Run
func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, becauce error %v", s.host, err.Error()))
	}
}

// Mendeklarasikan konstruktor
func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		panic(fmt.Errorf("failed to open database configuration: %w", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("failed to connect to database at %s:%s. Check DB_HOST. Error: %w", cfg.Host, cfg.Port, err))
	}

	// Dependencies
	jwtService := service.NewJwtService(cfg.TokenConfig)
	todosRepo := repository.NewTodosRepository(db)
	userRepo := repository.NewUserRepository(db)
	todosUseCase := usecase.NewTodosUseCase(todosRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	authUseCase := usecase.NewAuthenticateUsecase(userUseCase, jwtService)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)
	return &Server{
		todosUC: todosUseCase,
		userUC:  userUseCase,
		authUC:  authUseCase,
		jwtSvc:  jwtService,
		engine:  engine,
		host:    host,
	}
}
