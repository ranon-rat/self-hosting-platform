package public

import "github.com/gofiber/fiber/v2"

func Setup(app *fiber.App) {
	group := app.Group("/public")
	group.Get("/login", func(c *fiber.Ctx) error {
		// this is empty, but you need to remember that in the router
		// the auth header only compares it with an already stored password
		return c.Status(200).JSON(fiber.Map{
			"message": "finally you are alive :)",
		})
	})
}
