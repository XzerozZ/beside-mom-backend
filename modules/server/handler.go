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
	setupLikeRoutes(app, db, jwt)
	setupAppointRoutes(app, db, jwt)
	setupVideoRoutes(app, db, jwt, supa)
	setupKidRoutes(app, db, jwt, supa)
	setupUserRoutes(app, db, jwt, supa, mail)
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

func setupVideoRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase) {
	repository := repositories.NewGormVideoRepository(db)
	likerepository := repositories.NewGormLikesRepository(db)
	usecase := usecases.NewVideoUseCase(repository, likerepository, supa)
	controller := controllers.NewVideoController(usecase)

	videoGroup := app.Group("/video")
	videoGroup.Post("/", middlewares.JWTMiddleware(jwt), controller.CreateVideoHandler)
	videoGroup.Get("/", middlewares.JWTMiddleware(jwt), controller.GetAllVideoHandler)
	videoGroup.Get("/:id", middlewares.JWTMiddleware(jwt), controller.GetVideoByIDHandler)
	videoGroup.Put("/:id", middlewares.JWTMiddleware(jwt), controller.UpdateVideoHandler)
	videoGroup.Delete("/:id", middlewares.JWTMiddleware(jwt), controller.DeleteVideoByIDHandler)
}

func setupUserRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) {
	repository := repositories.NewGormUserRepository(db)
	kidrepository := repositories.NewGormKidsRepository(db)
	usecase := usecases.NewUserUseCase(repository, supa, mail)
	kidusecase := usecases.NewKidUseCase(kidrepository, supa)
	controller := controllers.NewUserController(usecase, kidusecase)

	userGroup := app.Group("/user", middlewares.JWTMiddleware(jwt))
	userGroup.Post("/", controller.CreateUserandKidsHandler)
	userGroup.Get("/", middlewares.AdminMiddleware, controller.GetAllMomHandler)
	userGroup.Get("/info/:id", middlewares.AdminMiddleware, controller.GetMomByIDHandler)
}

func setupLikeRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	repository := repositories.NewGormLikesRepository(db)
	usecase := usecases.NewLikeUseCase(repository)
	controller := controllers.NewLikeController(usecase)

	likeGroup := app.Group("/like", middlewares.JWTMiddleware(jwt))
	likeGroup.Post("/", controller.CreateLikeHandler)
	likeGroup.Get("/", controller.GetLikeByUserIDHandler)
	likeGroup.Get("/:video_id", controller.CheckLikeHandler)
	likeGroup.Delete("/:video_id", controller.DeleteLikeByIDHandler)
}

func setupKidRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase) {
	repository := repositories.NewGormKidsRepository(db)
	usecase := usecases.NewKidUseCase(repository, supa)
	controller := controllers.NewKidController(usecase)

	kidGroup := app.Group("/kid")
	kidGroup.Get("/", middlewares.JWTMiddleware(jwt), controller.GetAllKidsByUserID)
	kidGroup.Get("/:id", middlewares.JWTMiddleware(jwt), controller.GetKidByIDHandler)
}

func setupAppointRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	repository := repositories.NewGormAppRepository(db)
	usecase := usecases.NewAppUseCase(repository)
	controller := controllers.NewAppController(usecase)

	appointGroup := app.Group("/appoint", middlewares.JWTMiddleware(jwt))
	appointGroup.Post("/:userID", middlewares.AdminMiddleware, controller.CreateAppointmentHandler)
	appointGroup.Get("/", controller.GetAppHandler)
	appointGroup.Get("/:id", controller.GetAppByIDHandler)
	appointGroup.Put("/:id", middlewares.AdminMiddleware, controller.UpdateAppByIDHandler)
	appointGroup.Delete("/:id", middlewares.AdminMiddleware, controller.DeleteAppByIDHandler)
}
