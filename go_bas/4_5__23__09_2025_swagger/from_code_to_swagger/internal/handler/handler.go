package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gritsulyak/go_otus_additionals/go_bas/4_5__23__09_2025_swagger/from_code_to_swagger/internal/middleware"
	"github.com/gritsulyak/go_otus_additionals/go_bas/4_5__23__09_2025_swagger/from_code_to_swagger/internal/models"
	"github.com/gritsulyak/go_otus_additionals/go_bas/4_5__23__09_2025_swagger/from_code_to_swagger/internal/utils/jwtgen"
)

// Структуры для пользователя
var users = make(map[string]string) // Картотека пользователей: username -> password

// Handler структура, которая имеет методы для работы с каждым эндпоинтом.
// По желанию можно добавить поля (например различные валидаторы)
type Handler struct {
}

// New конструктор
func New() *Handler {
	return &Handler{}
}

// Register
// @Summary Регистрация
// @Tags auths
// @Accept	json
// @Produce	json
// @Param input body models.RegisterRequest true "Модель которую принимает метод (Param)"
// @Success 200 {string}  string "Registration successful"
// @Failure 400 {string} string "Invalid request"
// @Router /register [post]
func (h *Handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		// Проверка, что имя пользователя не существует
		if _, exists := users[req.Username]; exists {
			c.JSON(http.StatusConflict, gin.H{"message": "Username already exists"})
			return
		}

		// Регистрируем пользователя
		users[req.Username] = req.Password
		c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
	}
}

func (h *Handler) Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		// Проверка существования пользователя
		storedPassword, exists := users[req.Username]
		if !exists || storedPassword != req.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
			return
		}

		// Генерация токена
		token, err := jwtgen.GenerateToken(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
		})
	}
}

func (h *Handler) Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello, %s! You have access to this protected route.", username),
		})
	}
}

// InitHandler Инициализация маршрутов и обработчиков запросов.
// Здесь регистрируются эндпоинты (на всю программу может быть несколько подобных инициализаторов)
func InitHandler(api *gin.Engine, h *Handler) {
	api.POST("/register", h.Register())
	api.POST("/login", h.Login())
	api.GET("/protected", middleware.TokenAuthMiddleware(), h.Protected())
}
