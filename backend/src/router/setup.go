package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func Setup() {
	app := fiber.New()
	port := os.Getenv("PORT")
	if port == "" {

		port = "9239"
	}
	app.Listen(":" + port)

}
