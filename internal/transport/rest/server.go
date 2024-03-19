package rest

import (
	"github.com/Alieksieiev0/goshop/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Server struct {
	app *fiber.App
	ps  services.ProductService
	cs  services.CategoryService
}

func NewServer(app *fiber.App, ps services.ProductService, cs services.CategoryService) *Server {
	return &Server{
		app: app,
		ps:  ps,
		cs:  cs,
	}
}

func (s *Server) Start(addr string) error {
	s.app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path} | ${error}\nResponse Body: ${resBody}\n",
	}))
	s.app.Use(cors.New())
	s.app.Get("/categories/:id", getHandler(s.cs))
	s.app.Get("/categories", getAllHandler(s.ps))
	s.app.Post("/categories", createHandler(s.cs))
	s.app.Get("/products/:id", getHandler(s.ps))
	s.app.Get("/products", getAllHandler(s.ps))
	s.app.Post("/products", createHandler(s.ps))
	return s.app.Listen(addr)
}
