package main

import (
	"fmt"
	"log"
	"time"

	"github.com/chaosannals/jwtdemo/common"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	claims := &common.MyCustomClaims{
		UserId:   1,
		UserName: "aaaa",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	text, err := token.SignedString(common.DEMO_KEY)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("token: %s", text)
}
