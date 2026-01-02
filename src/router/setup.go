package router

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	logsC "github.com/ranon-rat/self-hosting-manager/src/controllers/logs"
	projectsC "github.com/ranon-rat/self-hosting-manager/src/controllers/projects"
	"github.com/ranon-rat/self-hosting-manager/src/controllers/public"
	"github.com/ranon-rat/self-hosting-manager/src/domain/repositories"
	"github.com/ranon-rat/self-hosting-manager/src/middleware"
)

func Setup(repos *repositories.Repositories) {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Password",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path}?${queryParams} | ${body}\n",
	})) //
	app.Static("/", "./public")
	//
	app.Use(middleware.BaseMiddleware())
	// it just has one route that helps with the login
	public.Setup(app)
	//->/projects
	projectsC.Setup(app, repos)
	// ->/logs
	logsC.Setup(app, repos)
	//
	port := os.Getenv("PORT")
	if port == "" {
		port = "9239"
	}
	fmt.Printf("starting service on port :%s", port)
	app.Listen(":" + port)

}
