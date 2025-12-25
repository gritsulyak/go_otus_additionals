package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gritsulyak/go_otus_additionals/go_bas/4_5__23__09_2025_swagger/from_code_to_swagger/internal/utils/jwtgen"
)

func ExMidlew() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Работаем с данными

		// Если ошибка вызываем c.Abort()

		// Иначе c.Next()
	}
}

// TokenAuthMiddleware Middleware для проверки JWT токена
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()

			return
		}

		// Удаляем Bearer из токена
		tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

		// Проверяем и парсим токен
		claims, err := jwtgen.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()

			return
		}

		// Если токен валиден, добавляем пользователя в контекст запроса
		c.Set("username", claims.Username)

		c.Next()
	}
}
