package rest

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/Alieksieiev0/goshop/internal/database"
	"github.com/Alieksieiev0/goshop/internal/services"
	"github.com/gofiber/fiber/v2"
)

func getHandler[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entity, err := service.Get(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(entity)
	}
}

func getAllHandler[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := []database.Param{}

		limit, err := strconv.Atoi(c.Query("limit", "10"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		params = append(params, database.Limit(limit))

		offset, err := strconv.Atoi(c.Query("offset", "0"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		params = append(params, database.Offset(offset))

		sortBy := c.Query("sort_by")
		orderBy := c.Query("order_by")
		if sortBy != "" {
			params = append(params, database.Order(sortBy, orderBy))
		}

		params = appendQueries[T](params, c.Queries())
		entities, err := service.GetAll(c.Context(), params...)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"data": entities})
	}
}

func createHandler[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entity := new(T)
		if err := c.BodyParser(entity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := service.Create(c.Context(), entity); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(entity)
	}
}

func appendQueries[T any](params []database.Param, queries map[string]string) []database.Param {
	for k, v := range queries {
		t := reflect.TypeOf(*new(T))
		for i := range t.NumField() {
			f := t.Field(i)
			if f.Tag.Get("json") == k {
				params = append(params, database.Filter(k, v, f.Type != reflect.TypeOf("")))
			}
		}
	}

	return params
}
