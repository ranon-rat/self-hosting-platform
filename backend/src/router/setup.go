package router

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Setup() {
	app := fiber.New()
	port := os.Getenv("PORT")
	if port == "" {

		port = "9239"
	}
	fmt.Println("starting service on port :", port)
	app.Listen(":" + port)

}
