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
	us  services.UserService
}

func NewServer(
	app *fiber.App,
	ps services.ProductService,
	cs services.CategoryService,
	us services.UserService,
) *Server {
	return &Server{
		app: app,
		ps:  ps,
		cs:  cs,
		us:  us,
	}
}

func (s *Server) Start(addr string) error {
	s.app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path} | ${error}\nResponse Body: ${resBody}\n",
	}))
	s.app.Use(cors.New())

	s.app.Post("/register", registerHandler(s.us))
	s.app.Post("/login", loginHandler(s.us))

	s.app.Get("/categories/:id", getHandler(s.cs))
	s.app.Delete("/categories/:id", AuthMiddleware, AdminMiddleware, deleteHandler(s.cs))
	s.app.Get("/categories", getAllHandler(s.cs))
	s.app.Post("/categories", AuthMiddleware, AdminMiddleware, createHandler(s.cs))

	s.app.Get("/products/:id", getHandler(s.ps))
	s.app.Delete("/products/:id", AuthMiddleware, AdminMiddleware, deleteHandler(s.ps))
	s.app.Get("/products", getAllHandler(s.ps))
	s.app.Post("/products", AuthMiddleware, AdminMiddleware, createHandler(s.ps))
	return s.app.Listen(addr)
}
