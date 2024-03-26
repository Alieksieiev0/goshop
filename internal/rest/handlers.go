package rest

import (
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/Alieksieiev0/goshop/internal/database"
	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/Alieksieiev0/goshop/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func registerHandler(service services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if user.Username == "" && user.Email == "" && user.Password == "" {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "Insufficient user data"})
		}

		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "Error processing data"})
		}
		user.Password = hashedPassword
		user.Role = models.Usr

		err = service.Create(c.Context(), user)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		user.Password = ""
		return c.Status(http.StatusOK).JSON(fiber.Map{"data": user})
	}
}

func loginHandler(service services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := &models.User{}
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if user.Email == "" && user.Password == "" {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "Insufficient user data"})
		}

		dbUser, err := service.GetWithFilters(
			c.Context(),
			database.Filter("email", user.Email, true),
		)

		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnprocessableEntity).
				JSON(fiber.Map{"error": "User with such email does not exist"})
		} else if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "Provided password is incorrect"})
		}

		token, err := generateJWT(dbUser, (time.Hour * 6))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
	}
}

func get[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entity, err := service.Get(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{"data": entity})
	}
}

func getAll[T any](service services.Service[T]) fiber.Handler {
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

func save[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		entity := new(T)
		if err := c.BodyParser(entity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := service.Create(c.Context(), entity); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(entity)
	}
}

func delete[T any](service services.Service[T]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := service.Delete(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		c.Status(http.StatusOK)
		return nil
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

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func generateJWT(user *models.User, expiresAfter time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.String(),
		Algorithm: jwt.SigningMethodHS256.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAfter)),
		},
	})

	signedToken, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

type UserClaims struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Algorithm string `json:"alg"`
	jwt.RegisteredClaims
}
