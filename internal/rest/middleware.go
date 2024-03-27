package rest

import (
	"net/http"
	"slices"
	"strings"

	"github.com/Alieksieiev0/goshop/internal/providers"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx, tokenProvider providers.TokenProvider) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No authorization code"})
	}
	claims, err := tokenProvider.Read(strings.Split(tokenString, " ")[1])
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	c.Locals("userRoles", claims.Roles)
	return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {
	userRoles := c.Locals("userRoles").([]string)
	if !slices.Contains(userRoles, "ADMIN") {
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"error": "You don`t have persmissions to perfom an action or acess a resource"})
	}
	return nil
}
