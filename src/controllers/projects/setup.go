package projectsC

import (
	"github.com/gofiber/fiber/v2"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
	"github.com/ranon-rat/self-hosting-manager/src/domain/repositories"
)

var pRepo projectsD.ProjectsRepoDB

func Setup(app *fiber.App, repos *repositories.Repositories) {
	pRepo = repos.ProjectRepo
	group := app.Group("/projects")
	group.Post("/", POSTCreateProject)
	group.Put("/", PUTProjects)
	group.Put("/pause", PauseProject)
	group.Get("/by-id", GetByID)
	group.Get("/", SearchProject)
}
