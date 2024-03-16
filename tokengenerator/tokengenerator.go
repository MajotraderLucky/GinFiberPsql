package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("8GoPUxkoCEeKaEG381hL6p9RAfwgCaiDJhrwy+/k8Og=")

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	// Создаем утверждения
	claims := MyCustomClaims{
		Username: "exampleUser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен истекает через 24 часа
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен нашим ключом
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Ошибка при создании токена:", err)
		return
	}

	fmt.Println("Сгенерированный JWT:", tokenString)
}
