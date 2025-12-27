package main

import (
	"E-commerce-Sederhana/config"
	"E-commerce-Sederhana/controller"
	"E-commerce-Sederhana/middleware"
	"E-commerce-Sederhana/repository"
	"E-commerce-Sederhana/usecase"
	"E-commerce-Sederhana/utils/service"
	"E-commerce-Sederhana/utils/service/midtrans"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	userUc       usecase.UserUsecase
	authUc       usecase.AuthenticationUsecase
	prodUc       usecase.ProductUsecase
	cartUc       usecase.CartUseCase
	cartItemUc   usecase.CartItemUseCase
	orderUc      usecase.OrderUsecase
	orderItemUc  usecase.OrderItemUsecase
	jwtSvc       service.JWTservice
	midtransHdlr *midtrans.MidtransHandler
	engine       *gin.Engine
	host         string
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtSvc)

	controller.NewUserController(rg, s.userUc, authMiddleware).Route()
	controller.NewAuthController(rg, s.authUc).Route()
	controller.NewProductController(rg, s.prodUc, authMiddleware).Route()
	controller.NewCartController(rg, s.cartUc, authMiddleware).Route()
	controller.NewCartItemController(s.cartItemUc, rg, authMiddleware).Route()
	controller.NewOrderController(rg, s.orderUc, authMiddleware).Route()
	controller.NewOrderItemController(rg, s.orderItemUc, authMiddleware).Route()

	// Daftarkan route untuk notifikasi Midtrans
	rg.POST("/midtrans/notification", s.midtransHdlr.HandleNotification)

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
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	cartItemRepo := repository.NewCartItemRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderItemRepo := repository.NewOrderItemRepository(db)
	
	userUsecase := usecase.NewUserUsecase(userRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
	cartUsecase := usecase.NewCartUseCase(cartRepo)
	cartItemUsecase := usecase.NewCartItemUseCase(cartItemRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, midtrans.NewMidtransService(cfg.Midtrans, orderRepo))
	orderItemUsecase := usecase.NewOrderItemUsecase(orderItemRepo)
	authUsecase := usecase.NewAuthenticationUsecase(userUsecase, jwtService)
	
	midtransHandler := midtrans.NewMidtransHandler(midtrans.NewMidtransService(cfg.Midtrans, orderRepo))

	return &Server{
		userUc:       userUsecase,
		authUc:       authUsecase,
		prodUc:       productUsecase,
		cartUc:       cartUsecase,
		cartItemUc:   cartItemUsecase,
		orderUc:      orderUsecase,
		orderItemUc:  orderItemUsecase,
		jwtSvc:       jwtService,
		midtransHdlr: midtransHandler,
		engine:       gin.Default(),
		host:         ":" + cfg.API.Port,
	}

}
