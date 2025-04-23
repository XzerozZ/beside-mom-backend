package main

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/server"
	"Beside-Mom-BE/pkg/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := configs.LoadConfigs()
	database.InitDB(config.PostgreSQL)
	app := fiber.New(fiber.Config{
		BodyLimit:         2 * 1024 * 1024 * 1024,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      30 * time.Minute,
		IdleTimeout:       10 * time.Minute,
		StreamRequestBody: true,
	})

	server.SetupRoutes(app, config.JWT, config.Supabase, config.Mail)
	serverAddress := config.App.Host + ":" + config.App.Port
	log.Printf("Server is running on %s", serverAddress)
	log.Fatal(app.Listen(serverAddress))
}
