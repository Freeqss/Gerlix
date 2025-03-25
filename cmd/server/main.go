package main

import (
	"Gerlix/internal/database"
	"Gerlix/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Подключаем базу данных
	database.ConnectDB()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                           // Разрешаем запросы со всех источников
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS", // Разрешаем эти методы
		AllowHeaders: "Content-Type, Authorization", // Разрешаем заголовки
	}))

	// Роуты
	routes.RegisterRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
