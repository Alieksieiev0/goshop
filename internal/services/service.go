package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Alieksieiev0/goshop/internal/database"
	"github.com/Alieksieiev0/goshop/internal/models"
	"github.com/Alieksieiev0/goshop/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service[T any] interface {
	GetById(ctx context.Context, id string) (*T, error)
	GetAll(ctx context.Context, filters map[string]string) ([]T, error)
	SaveEntity(ctx context.Context, entity *T) error
	DeleteById(ctx context.Context, id string) error
}

type DatabaseService[T any] struct {
	repository repositories.DatabaseRepository[T]
}

func NewDatabaseService[T any](repository repositories.DatabaseRepository[T]) *DatabaseService[T] {
	return &DatabaseService[T]{
		repository: repository,
	}
}

func (cs *DatabaseService[T]) GetById(ctx context.Context, id string) (*T, error) {
	return cs.repository.Get(ctx, id)
}

func (cs *DatabaseService[T]) GetAll(ctx context.Context, filters map[string]string) ([]T, error) {
	params := []database.Param{}

	limit := 10
	if filterLimit, err := strconv.Atoi(filters["limit"]); err == nil {
		limit = filterLimit
	}
	params = append(params, database.Limit(limit))

	offset := 0
	if offsetFilter, err := strconv.Atoi(filters["offset"]); err == nil {
		offset = offsetFilter
	}
	params = append(params, database.Offset(offset))

	sortBy, ok := filters["sort_by"]
	if ok {
		params = append(params, database.Order(sortBy, filters["order_by"]))
	}

	params = database.AppendFilters[T](params, filters)
	return cs.repository.GetAllWithFilters(ctx, params...)
}

func (cs *DatabaseService[T]) SaveEntity(ctx context.Context, entity *T) error {
	return cs.repository.Save(ctx, entity)
}

func (cs *DatabaseService[T]) DeleteById(ctx context.Context, id string) error {
	return cs.repository.Delete(ctx, id)
}

type CategoryService interface {
	Service[models.Category]
}

type CategoryDatabaseService struct {
	*DatabaseService[models.Category]
}

func NewCategoryDatabaseService(db *gorm.DB) CategoryService {
	return &CategoryDatabaseService{
		DatabaseService: NewDatabaseService(repositories.NewGormRepository[models.Category](db)),
	}
}

type ProductService interface {
	Service[models.Product]
}

type ProductDatabaseService struct {
	*DatabaseService[models.Product]
}

func NewProductDatabaseService(db *gorm.DB) ProductService {
	return &ProductDatabaseService{
		DatabaseService: NewDatabaseService(repositories.NewGormRepository[models.Product](db)),
	}
}

type UserService interface {
	Service[models.User]
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type UserDatabaseService struct {
	*DatabaseService[models.User]
}

func NewUserDatabaseService(db *gorm.DB) UserService {
	return &UserDatabaseService{
		DatabaseService: NewDatabaseService(repositories.NewGormRepository[models.User](db)),
	}
}

func (rs *UserDatabaseService) GetByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {
	return rs.repository.GetWithFilters(ctx, database.Filter("username", username, true))
}

type RoleService interface {
	Service[models.Role]
	GetByName(ctx context.Context, name string) (*models.Role, error)
}

type RoleDatabaseService struct {
	*DatabaseService[models.Role]
}

func NewRoleDatabaseService(db *gorm.DB) RoleService {
	return &RoleDatabaseService{
		DatabaseService: NewDatabaseService(repositories.NewGormRepository[models.Role](db)),
	}
}

func (rs *RoleDatabaseService) GetByName(ctx context.Context, name string) (*models.Role, error) {
	return rs.repository.GetWithFilters(ctx, database.Filter("name", name, true))
}

type AuthService interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, user *models.User) (*models.User, error)
}

type DefaultAuthService struct {
	userService UserService
	roleService RoleService
}

func NewDefaultAuthService(userService UserService, roleService RoleService) AuthService {
	return &DefaultAuthService{
		userService: userService,
		roleService: roleService,
	}
}

func (das *DefaultAuthService) Register(ctx context.Context, user *models.User) error {
	if user.Username == "" && user.Email == "" && user.Password == "" {
		return fmt.Errorf("insufficient user data")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	adminRole, err := das.roleService.GetByName(ctx, "USER")
	if err != nil {
		return nil
	}

	user.Password = hashedPassword
	user.Roles = []models.Role{*adminRole}
	err = das.userService.SaveEntity(ctx, user)
	if err != nil {
		return err
	}

	user.Password = ""
	return nil
}

func (das *DefaultAuthService) Login(ctx context.Context, user *models.User) (*models.User, error) {
	if user.Username == "" && user.Password == "" {
		return nil, fmt.Errorf("insufficient user data")
	}

	dbUser, err := das.userService.GetByUsername(ctx, user.Username)

	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user with such username does not exist")
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, fmt.Errorf("provided password is incorrect")
	}

	return dbUser, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
