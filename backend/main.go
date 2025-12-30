package main

import (
	"github.com/joho/godotenv"
	"github.com/ranon-rat/self-hosting-manager/src/database"
	"github.com/ranon-rat/self-hosting-manager/src/router"
)

func main() {
	godotenv.Load(".env.local")
	database.Setup()
	router.Setup()
}
