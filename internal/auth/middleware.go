package auth

import (
	"Gerlix/internal/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Middleware для проверки JWT-токена
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		// Если заголовка нет — возвращаем ошибку
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		var tokenString string

		// Проверяем, если токен в заголовке без "Bearer ", значит это сам токен
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ") // Убираем "Bearer "
		} else {
			tokenString = authHeader // Просто сам токен
		}

		// Разбираем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Извлекаем user_id из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Проверяем, что "user_id" существует в claims
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
		}

		// Преобразуем строковый user_id в UUID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID format"})
		}

		// Проверяем, существует ли пользователь
		var exists bool
		err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
		if err != nil || !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		// Сохраняем user_id в Locals для использования в других хендлерах
		c.Locals("user_id", userID)

		// Переходим к следующему обработчику
		return c.Next()
	}
}
