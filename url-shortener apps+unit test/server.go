package main

import (
	"database/sql"
	"fmt"
	"log"

	"enigmacamp.com/url-shortener/config"
	"enigmacamp.com/url-shortener/controller"
	"enigmacamp.com/url-shortener/middleware"
	"enigmacamp.com/url-shortener/repository"
	"enigmacamp.com/url-shortener/usecase"
	"enigmacamp.com/url-shortener/utils/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Mendeklarasikan struct
type Server struct {
	userUC usecase.UserUseCase
	authUC usecase.AuthenticateUsecase
	urlUC  usecase.UrlUsecase
	jwtSvc service.JwtService
	engine *gin.Engine
	host   string
}

// Mendeklarasikan method initRoute
func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")

	authMiddleware := middleware.NewAuthMiddleware(s.jwtSvc)
	controller.NewUserController(s.userUC, rg, authMiddleware).Route()
	controller.NewAuthController(s.authUC, rg).Route()
	controller.NewUrlController(s.urlUC, s.engine, authMiddleware).Route()
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
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// Dependencies
	jwtService := service.NewJwtService(cfg.TokenConfig)
	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewUrlRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	authUseCase := usecase.NewAuthenticateUsecase(userUseCase, jwtService)
	urlUseCase := usecase.NewUrlUsecase(urlRepo)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)
	return &Server{
		userUC: userUseCase,
		authUC: authUseCase,
		urlUC:  urlUseCase,
		jwtSvc: jwtService,
		engine: engine,
		host:   host,
	}
}
