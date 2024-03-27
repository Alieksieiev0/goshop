package rest

import (
	"github.com/Alieksieiev0/goshop/internal/providers"
	"github.com/Alieksieiev0/goshop/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

type Server struct {
	app *fiber.App
	db  *gorm.DB
}

func NewServer(app *fiber.App, db *gorm.DB) *Server {
	return &Server{
		app: app,
		db:  db,
	}
}

func (s *Server) Start(addr string) error {
	s.app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path} | ${error}\nResponse Body: ${resBody}\n",
	}))
	s.app.Use(cors.New())

	NewAuthRestController(
		services.NewDefaultAuthService(
			services.NewUserDatabaseService(s.db),
			services.NewRoleDatabaseService(s.db),
		),
		providers.NewJWTProvider(),
	).Activate(s.app)
	NewProductRestController(services.NewProductDatabaseService(s.db)).Activate(s.app)
	NewCategoryRestController(services.NewCategoryDatabaseService(s.db)).Activate(s.app)

	return s.app.Listen(addr)
}
