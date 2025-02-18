package server

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/controllers"
	"Beside-Mom-BE/modules/repositories"
	"Beside-Mom-BE/modules/usecases"
	"Beside-Mom-BE/pkg/database"
	"Beside-Mom-BE/pkg/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) {
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":  "Success",
			"message": "Welcome to the Nursing House System!",
		})
	})

	setupAuthRoutes(app, db, jwt)
	setupQuestRoutes(app, db, jwt)
}

func setupAuthRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	repository := repositories.NewGormUserRepository(db)
	usecase := usecases.NewAuthUseCase(repository, jwt)
	controller := controllers.NewAuthController(usecase)

	authGroup := app.Group("/auth")
	authGroup.Post("/register", controller.RegisterHandler)
	authGroup.Post("/login", controller.LoginHandler)
}

func setupQuestRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	repository := repositories.NewGormQuestionRepository(db)
	usecase := usecases.NewQuestionUseCase(repository)
	controller := controllers.NewQuestionController(usecase)

	questionGroup := app.Group("/question")
	questionGroup.Post("/", middlewares.JWTMiddleware(jwt), controller.CreateQuestionHandler)
	questionGroup.Get("/", controller.GetAllQuestionHandler)
	questionGroup.Get("/:id", middlewares.JWTMiddleware(jwt), controller.GetQuestionByIDHandler)
	questionGroup.Put("/:id", middlewares.JWTMiddleware(jwt), controller.UpdateQuestionByIDHandler)
	questionGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), controller.DeleteQuestionByIDHandler)
}
