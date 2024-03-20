package rest

import (
	"fmt"
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

	token, err := jwt.Parse(parsedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, c.Status(http.StatusUnauthorized).
				JSON(fiber.Map{"error": "You're Unauthorized due to error parsing the JWT"})
		}
		return []byte("secret-key"), nil

	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "You're Unauthorized due to invalid token"})
	}

	userRole, err := getUserRoleFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "You're Unauthorized due to invalid token claims"})
	}
	c.Locals("userRole", userRole)
	return c.Next()
}

func AdminMiddleware(c *fiber.Ctx) error {
	userRole := c.Locals("userRole")
	fmt.Println(userRole)
	if userRole != models.Admin {
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"error": "You don`t have a persmissions to access a page"})
	}
	return nil
}

func getUserRoleFromToken(token *jwt.Token) (models.UserRole, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error processing token claims")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("error parsing user role")
	}

	var userRole models.UserRole
	switch role {
	case string(models.Admin):
		userRole = models.Admin
	case string(models.Usr):
		userRole = models.Usr
	case string(models.Guest):
		userRole = models.Guest
	}

	return userRole, nil
}
