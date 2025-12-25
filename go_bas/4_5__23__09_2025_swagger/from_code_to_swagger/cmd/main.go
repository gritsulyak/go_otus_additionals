package main

import (
	_ "github.com/SaveliiYam/simple-http-server/docs"
	"github.com/SaveliiYam/simple-http-server/internal/handler"
	"github.com/gin-gonic/gin"
	_ "github.com/gritsulyak/go_otus_additionals/go_bas/4_5__23__09_2025_swagger/from_code_to_swagger/internal/docs"
	"github.com/gritsulyak/go_otus_additionals/go_bas/4_5__23__09_2025_swagger/from_code_to_swagger/internal/handler"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Базовый сервис авторизации
// @version 1
// @description API Server

// @host localhost:8080/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Инициализация роутера (сервера) Gin
	r := gin.Default()

	// Инициализация обработчика запросов
	// Именно он будет отвечать за хендлеры и обрабатывать каждый запрос.
	handle := handler.New()

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Инициализация маршрутов и обработчиков запросов
	// Функция нужна чтобы добавить каждый хендлер в сервер

	handler.InitHandler(r, handle)

	// Запуск сервера, если при запуске произойдет ошибка, то программа остановится
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
