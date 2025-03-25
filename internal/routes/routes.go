package routes

import (
	"Gerlix/internal/auth"
	"Gerlix/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// RegisterRoutes подключает маршруты к Fiber-приложению
func RegisterRoutes(app *fiber.App) {

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3501", // Укажи точный домен фронтенда
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // Разрешение на использование cookies
	}))
	// Роуты авторизации
	app.Post("/register", auth.Register)
	app.Post("/login", auth.Login)
	// app.Post("/refresh", auth.RefreshTokenHandler) // Добавили
	app.Get("/user/info", auth.JWTMiddleware(), user.UserInfo)
	app.Get("/user", auth.JWTMiddleware(), user.UserInfoById)
	app.Get("/user/all", auth.JWTMiddleware(), user.GetAllUsers)

	// Защищённый маршрут (пример)
	app.Get("/protected", auth.JWTMiddleware(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to protected route!"})
	})
}
