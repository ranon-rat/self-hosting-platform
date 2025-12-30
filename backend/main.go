package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/ranon-rat/self-hosting-manager/src/database"
	executionlogsDB "github.com/ranon-rat/self-hosting-manager/src/database/executionlogs"
	projectsDB "github.com/ranon-rat/self-hosting-manager/src/database/projects"
	projectsD "github.com/ranon-rat/self-hosting-manager/src/domain/projects"
	"github.com/ranon-rat/self-hosting-manager/src/domain/repositories"
	"github.com/ranon-rat/self-hosting-manager/src/router"
	executionerServices "github.com/ranon-rat/self-hosting-manager/src/services/executioner"
)

func main() {
	godotenv.Load(".env.local")
	database.Setup()
	repos := repositories.Repositories{
		ProjectRepo: projectsDB.NewRepo(database.GetDB()),
		LogRepo:     executionlogsDB.NewRepo(database.GetDB()),
	}
	_, err := repos.ProjectRepo.Create(&projectsD.NewProject{
		Name:         "test",
		Dir:          "/home/ranon-rat/Escritorio/proyectos/self-hosting-manager/backend",
		Command:      "ping 1.1.1.1",
		ThumbnailURL: "https://i.pinimg.com/736x/00/79/a9/0079a9a70d1f9fd429d3fbd44564afb9.jpg",
	})
	if err != nil {
		fmt.Println(err)
	}
	executionerServices.Setup(&repos)
	executionerServices.StartServices()
	router.Setup(&repos)
}
