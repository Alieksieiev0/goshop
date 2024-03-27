package types

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	Id        string   `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Roles     []string `json:"role"`
	Algorithm string   `json:"alg"`
	jwt.RegisteredClaims
}
