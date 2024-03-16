package rest

import (
	"log"
	"os"

	"github.com/Alieksieiev0/goshop/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r  *gin.Engine
	ps services.ProductService
}

func NewServer(r *gin.Engine, ps services.ProductService) *Server {
	return &Server{
		r:  r,
		ps: ps,
	}
}

func (s *Server) Start(addr ...string) error {
	logger := log.New(os.Stdout, "[APP-debug] [INFO] ", log.Lshortfile)
	s.r.Use(Logger(logger))
	s.r.Use(cors.Default())
	s.r.GET("/products", s.handleGetAllProducts)
	s.r.POST("/products", s.handleCreateProduct)
	s.r.DELETE("/products", s.handleDeleteAllProducts)
	return s.r.Run(addr...)
}
