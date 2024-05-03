package common

import "github.com/golang-jwt/jwt/v5"

type MyCustomClaims struct {
	UserId   int64
	UserName string
	jwt.RegisteredClaims
}
