package rest

import (
	"net/http"
	"strings"

	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": "No authorization code"})
	}
	parsedToken := strings.Split(tokenString, " ")[1]

	token, err := jwt.ParseWithClaims(
		parsedToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, c.Status(http.StatusUnauthorized).
					JSON(fiber.Map{"error": "You're Unauthorized due to error parsing the JWT"})
			}
			return []byte("secret-key"), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "You're Unauthorized due to invalid token"})
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "You're Unauthorized due to claims parsing error"})
	}

	c.Locals("userRole", claims.Role)
	return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {
	userRole := c.Locals("userRole")
	if userRole != models.Admin {
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"error": "You don`t have persmissions to perfom an action or acess a resource"})
	}
	return nil
}
