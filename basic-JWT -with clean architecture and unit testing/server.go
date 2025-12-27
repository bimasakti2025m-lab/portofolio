package main

import (
	"basic-JWT/config"
	"basic-JWT/controller"
	"basic-JWT/middleware"
	"basic-JWT/repository"
	"basic-JWT/usecase"
	"basic-JWT/utils/service"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	userUc usecase.UserUsecase
	authUc usecase.AuthenticationUsecase
	jwtSvc service.JWTservice
	engine *gin.Engine
	host   string
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtSvc)

	controller.NewUserController(rg, s.userUc, authMiddleware).Route()
	controller.NewAuthController(rg, s.authUc).Route()

}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, becauce error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		panic("failed to load config")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.Database)
	db, err := sql.Open(cfg.DB.Driver, dsn)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect to database")
	}

	err = db.Ping()
	if err != nil {
		panic("failed to ping database")
	}

	jwtService := service.NewJWTService(cfg.Token)
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	authUsecase := usecase.NewAuthenticationUsecase(userUsecase, jwtService)
	JwtService := service.NewJWTService(cfg.Token)

	return &Server{
		userUc: userUsecase,
		authUc: authUsecase,
		jwtSvc: JwtService,
		engine: gin.Default(),
		host:   ":" + cfg.API.Port,
	}

}
