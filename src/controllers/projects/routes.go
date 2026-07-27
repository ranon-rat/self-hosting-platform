package projectsC

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
	executionerServices "github.com/ranon-rat/self-hosting-manager/src/services/executioner"
)

// POST /projects/
func POSTCreateProject(c *fiber.Ctx) error {
	newProject := new(projectsD.NewProject)
	if err := c.BodyParser(newProject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	id, err := pRepo.Create(newProject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	project, err := pRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	executionerServices.Executioner(project)
	return c.Status(200).JSON(fiber.Map{
		"message": "everything seems good :D",
	})
}

// PUT /projects/
func PUTProjects(c *fiber.Ctx) error {
	putProject := new(projectsD.UpdateProject)
	if err := c.BodyParser(putProject); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if err := pRepo.UpdateProject(putProject); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "everything seems good :D",
	})

}

// PUT /projects/pause?id=1234&pause=false-true
func PauseProject(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	pause := c.QueryBool("pause")
	if pause {
		if err := executionerServices.StopProject(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"message": "everything seems good :D",
		})
	}
	if err := executionerServices.StartProject(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "everything seems good :D",
	})
}

// GET /projects/by-id?id=1234
func GetByID(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	project, err := pRepo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(project)
}

// GET /projects?search=asdf
func SearchProject(c *fiber.Ctx) error {
	search := c.Query("search")
	projects, err := pRepo.Search(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(projects)
}

// DELETE /projects?id=id
func DeleteProject(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	fmt.Println("deleting project", id)

	// first stop the running process and close its log channel,
	// then remove it from the database
	if err := executionerServices.DeleteProject(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if err := pRepo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "project deleted",
	})
}
