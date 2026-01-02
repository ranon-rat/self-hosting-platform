package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	sigChan := make(chan os.Signal, 1)

	// Registrar las señales que quieres capturar
	signal.Notify(sigChan,
		os.Interrupt,    // SIGINT (Ctrl+C)
		syscall.SIGTERM, // SIGTERM (kill)
		syscall.SIGQUIT, // SIGQUIT (Ctrl+\)
	)

	// Goroutine que escucha las señales
	go func() {
		sig := <-sigChan
		log.Printf("Recibida señal: %v. Deteniendo todos los proyectos...", sig)
		executionerServices.StoppingAll()
		// Espera opcional para que terminen los procesos
		time.Sleep(2 * time.Second)
		log.Println("Limpieza completa. Saliendo...")
		os.Exit(0)
	}()
	repos := repositories.Repositories{
		ProjectRepo: projectsDB.NewRepo(database.GetDB()),
		LogRepo:     executionlogsDB.NewRepo(database.GetDB()),
	}
	executionerServices.Setup(&repos)
	router.Setup(&repos)
}
