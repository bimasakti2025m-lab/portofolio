// TODO :
// 1. Mendeklarasikan nama package main
// 2. Mendeklarasikan struct bernama Server
// 3. Mendeklarasikan funtion initRoute
// 4. Membubuat function Run
// 5. Membuat constructor bernama NewServer

package main

import (
	"simple-clean-architecture/config"
	"simple-clean-architecture/controller"
	"strconv"

	"simple-clean-architecture/repositori"
	"simple-clean-architecture/usecase"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	bookUsecase usecase.BookUsecase
	engine      *gin.Engine
	host        string
}

func (s *Server) initRoute() {
	rg := s.engine.Group("api/v1")

	controller.NewBookController(s.bookUsecase, rg).Route()
}

func (s *Server) Run() {
	s.initRoute()

	err := s.engine.Run(s.host)

	if err != nil {
		panic(err)
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	dsn := "host=" + cfg.DBConfig.Host + " port=" + strconv.Itoa(cfg.DBConfig.Port) + " user=" + cfg.DBConfig.Username + " password=" + cfg.DBConfig.Password + " dbname=" + cfg.DBConfig.Database + " sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)

	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	bookRepositori := repositori.NewBookRepositori(db)
	bookUsecase := usecase.NewBookUsecase(bookRepositori)

	engine := gin.Default()

	return &Server{
		bookUsecase: bookUsecase,
		engine:      engine,
		host:        cfg.APIConfig.Host + ":" + strconv.Itoa(cfg.APIConfig.Port),
	}

}
