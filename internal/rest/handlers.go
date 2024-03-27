package rest

import (
	"net/http"
	"time"

	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/Alieksieiev0/goshop/internal/providers"
	"github.com/Alieksieiev0/goshop/internal/services"
	"github.com/gofiber/fiber/v2"
)

func register(service services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := service.Register(c.Context(), user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error})
		}

		c.Status(http.StatusOK)
		return nil
	}
}

func login(service services.AuthService, tokenProvider providers.TokenProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := &models.User{}
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		user, err := service.Login(c.Context(), user)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		token, err := tokenProvider.Create(user, (time.Hour * 6))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
	}
}

func get[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entity, err := service.GetById(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"data": entity})
	}
}

func getAll[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entities, err := service.GetAll(c.Context(), c.Queries())
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"data": entities})
	}
}

func save[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entity := new(T)
		if err := c.BodyParser(entity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := service.SaveEntity(c.Context(), entity); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(entity)
	}
}

func delete[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := service.DeleteById(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		c.Status(http.StatusOK)
		return nil
	}
}
