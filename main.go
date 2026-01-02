package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,    // SIGINT (Ctrl+C)
		syscall.SIGTERM, // SIGTERM (kill)
		syscall.SIGQUIT, // SIGQUIT (Ctrl+\)
	)
	defer stop()

	repos := repositories.Repositories{
		ProjectRepo: projectsDB.NewRepo(database.GetDB()),
		LogRepo:     executionlogsDB.NewRepo(database.GetDB()),
	}

	executionerServices.Setup(&repos)
	go router.Setup(&repos)
	<-ctx.Done() // ⛔ wait for systemd signal
	log.Println("Shutdown signal received")
	executionerServices.StoppingAll(ctx)
	log.Println("Cleanup finished, exiting")
}
