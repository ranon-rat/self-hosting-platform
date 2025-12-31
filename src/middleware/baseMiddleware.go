package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func BaseMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Password")
		if auth == "" {
			auth = c.Query("password")
		}
		if auth != os.Getenv("PASSWORD") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "wrong password",
			})
		}
		return c.Next()

	}
}
