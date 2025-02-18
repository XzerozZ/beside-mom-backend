package main

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/pkg/database"
	"log"
	"math"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := configs.LoadConfigs()
	database.InitDB(config.PostgreSQL)
	app := fiber.New(fiber.Config{
		BodyLimit: math.MaxInt64,
	})

	serverAddress := config.App.Host + ":" + config.App.Port
	log.Printf("Server is running on %s", serverAddress)
	log.Fatal(app.Listen(serverAddress))
}
