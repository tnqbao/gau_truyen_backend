//go:build dev

package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnqbao/gau_truyen_backend/config"
	"github.com/tnqbao/gau_truyen_backend/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//config.InitRedis()
	db := config.InitDB()

	router := routes.SetupRouter(db)
	router.Run(":8084")
}
