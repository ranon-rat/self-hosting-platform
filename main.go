package main

import (
	"github.com/joho/godotenv"
	"github.com/ranon-rat/self-hosting-manager/src/database"
	executionlogsDB "github.com/ranon-rat/self-hosting-manager/src/database/executionlogs"
	projectsDB "github.com/ranon-rat/self-hosting-manager/src/database/projects"
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
	executionerServices.Setup(&repos)
	router.Setup(&repos)
}
