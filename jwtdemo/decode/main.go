package main

import (
	"fmt"
	"log"

	"github.com/chaosannals/jwtdemo/common"
	"github.com/golang-jwt/jwt/v5"
)

const tokonString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsIlVzZXJOYW1lIjoiYWFhYSIsImlzcyI6InRlc3QiLCJzdWIiOiJzb21lYm9keSIsImF1ZCI6WyJzb21lYm9keV9lbHNlIl0sImV4cCI6MTcxNDgxMTQ1MCwibmJmIjoxNzE0NzI1MDUwLCJpYXQiOjE3MTQ3MjUwNTAsImp0aSI6IjEifQ.lZTraJKxAqaqtOdtUQnzMqs7ihw38hXhtx_zodTNT9w"

func parseDefault() {
	token, err := jwt.Parse(tokonString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("加密算法不是 HMAC %v", token.Header["alg"])
		}
		fmt.Printf("HMAC %s\n", token.Header["alg"])
		return common.DEMO_KEY, nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	// 可以直接转换成 jwt 提供的 通用的字典类型
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Printf("MapClaims: %v\n", claims)
	} else {
		fmt.Printf("MapClaims 转换失败: %v\n", err)
	}

	// 不能直接转换成自定义类型
	if claims, ok := token.Claims.(common.MyCustomClaims); ok {
		fmt.Printf("MyCustomClaims: %v\n", claims)
	} else {
		fmt.Printf("MyCustomClaims 转换失败: %v\n", err)
	}
}

func parseCustom() {
	claims := &common.MyCustomClaims{} // 直接获取填充的值。
	token, err := jwt.ParseWithClaims(tokonString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("加密算法不是 HMAC %v", token.Header["alg"])
		}
		fmt.Printf("HMAC %s\n", token.Header["alg"])
		return common.DEMO_KEY, nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	// 直接转换成自定义类型
	fmt.Printf("MyCustomClaims: %v\n", claims)
	// token.Claims 会变成 nil

	// 现在类型变成指定的类型了，所以 默认的类型 不可用。
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Printf("MapClaims: %v\n", claims)
	} else {
		fmt.Printf("MapClaims 转换失败: %v\n", err)
	}

	// 而且 token.Claims 也不能直接转换成自定义类型.
	if claims, ok := token.Claims.(common.MyCustomClaims); ok {
		fmt.Printf("MyCustomClaims: %v\n", claims)
	} else {
		fmt.Printf("MyCustomClaims 转换失败: %v\n", err)
	}
}

func main() {
	parseDefault()
	parseCustom()
}
