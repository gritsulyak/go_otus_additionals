package jwtgen

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gritsulyak/go_otus_additionals/go_bas/4_5__23__09_2025_swagger/from_code_to_swagger/internal/models"
)

// Структуры для пользователя и токенов
var (
	secretKey               = []byte("mySecretKey")
	ErrInvalidSigningMethod = errors.New("Invalid signing method")
	ErrInvalidToken         = errors.New("Invalid token")
)

// GenerateToken Функция для создания JWT токена
func GenerateToken(username string) (string, error) {
	// Устанавливаем срок действия токена (1 час)
	expirationTime := time.Now().Add(1 * time.Hour)

	// Создаем JWT токен
	claims := &models.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "go-gin-jwt-example",
		},
	}

	// Создаем JWT токен и шифруем его методом HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (*models.Claims, error) {
	// Проверяем и парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return secretKey, nil
	})

	// Если у нас произошла ошибка или токен не валиден, то вернем ошибку
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
