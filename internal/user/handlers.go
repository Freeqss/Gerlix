package user

import (
	"Gerlix/internal/database"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Получение информации о пользователе (только с авторизацией)
func UserInfo(c *fiber.Ctx) error {
	// Получаем user_id из middleware
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user session"})
	}

	// Достаем информацию о пользователе из БД
	var email string
	err := database.DB.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&email)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	return c.JSON(fiber.Map{
		"user_id": userID.String(),
		"email":   email,
	})
}

func GetAllUsers(c *fiber.Ctx) error {
	// Структура для хранения данных о пользователе
	type User struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	// Запрос к БД
	rows, err := database.DB.Query("SELECT id, email FROM users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()

	// Массив для хранения пользователей
	var users []User

	// Обход строк результата запроса
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to scan user"})
		}
		users = append(users, user)
	}

	// Проверка на ошибку после итерации
	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}

	// Возвращаем JSON
	return c.JSON(users)
}

func UserInfoById(c *fiber.Ctx) error {
	// Получаем user_id из query-параметров
	userId := c.Query("user_id")

	// Логируем полученное значение
	fmt.Println("Received user_id:", userId)

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id query parameter is required"})
	}

	var email string
	// Запрос к базе данных с использованием полученного user_id
	err := database.DB.QueryRow("SELECT email FROM users WHERE id = $1", userId).Scan(&email)

	// Логируем ошибку
	if err != nil {
		// Логирование ошибки
		fmt.Println("Error querying database:", err)

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	fmt.Println("Received email:", email)
	return c.JSON(fiber.Map{
		"email": email,
	})
}
