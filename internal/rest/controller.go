package rest

import (
	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/Alieksieiev0/goshop/internal/providers"
	"github.com/Alieksieiev0/goshop/internal/services"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Activate(app *fiber.App)
}

type BaseController struct {
	path     string
	pathSlug string
}

type ProductController interface {
	Controller
}

type ProductRestController struct {
	service services.ProductService
	base    BaseController
}

func NewProductRestController(service services.ProductService) ProductController {
	return &ProductRestController{
		service: service,
		base: BaseController{
			path:     "/products",
			pathSlug: "/products/:id",
		},
	}
}

func (prc *ProductRestController) Activate(app *fiber.App) {
	app.Get(prc.base.pathSlug, get[models.Product](prc.service))
	app.Get(prc.base.path, getAll[models.Product](prc.service))
	app.Delete(
		prc.base.pathSlug,
		AuthMiddleware,
		AdminMiddleware,
		delete[models.Product](prc.service),
	)
	app.Post(prc.base.path, AuthMiddleware, AdminMiddleware, save[models.Product](prc.service))
}

type CategoryController interface {
	Controller
}

type CategoryRestController struct {
	service services.CategoryService
	base    BaseController
}

func NewCategoryRestController(service services.CategoryService) CategoryController {
	return &CategoryRestController{
		service: service,
		base: BaseController{
			path:     "/categories",
			pathSlug: "/categories/:id",
		},
	}
}

func (crs *CategoryRestController) Activate(app *fiber.App) {
	app.Get(crs.base.pathSlug, get[models.Category](crs.service))
	app.Get(crs.base.path, getAll[models.Category](crs.service))
	app.Delete(
		crs.base.pathSlug,
		AuthMiddleware,
		AdminMiddleware,
		delete[models.Category](crs.service),
	)
	app.Post(crs.base.path, AuthMiddleware, AdminMiddleware, save[models.Category](crs.service))
}

type AuthController interface {
	Controller
}

type AuthRestController struct {
	service       services.AuthService
	tokenProvider providers.TokenProvider
}

func NewAuthRestController(
	service services.AuthService,
	tokenProvider providers.TokenProvider,
) AuthController {
	return &AuthRestController{
		service:       service,
		tokenProvider: tokenProvider,
	}
}

func (ars *AuthRestController) Activate(app *fiber.App) {
	app.Post("/register", register(ars.service))
	app.Post("/login", login(ars.service, ars.tokenProvider))
}
