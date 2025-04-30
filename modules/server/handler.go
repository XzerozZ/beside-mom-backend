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

	setupAuthRoutes(app, db, jwt, mail)
	setupQuestRoutes(app, db, jwt)
	setupHistoryRoutes(app, db, jwt)
	setupLikeRoutes(app, db, jwt)
	setupAppointRoutes(app, db, jwt)
	setupEvaluateRoutes(app, db, jwt)
	setupGrowthRoutes(app, db, jwt)
	setupVideoRoutes(app, db, jwt, supa)
	setupCareRoutes(app, db, jwt, supa)
	setupKidRoutes(app, db, jwt, supa)
	setupQuizRoutes(app, db, jwt, supa)
	setupUserRoutes(app, db, jwt, supa, mail)
}

func setupAuthRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, mail configs.Mail) {
	repository := repositories.NewGormUserRepository(db)
	usecase := usecases.NewAuthUseCase(repository, jwt, mail)
	controller := controllers.NewAuthController(usecase)

	authGroup := app.Group("/auth")
	authGroup.Post("/register", controller.RegisterHandler)
	authGroup.Post("/login", controller.LoginHandler)
	authGroup.Post("/forgotpassword", controller.ForgotPasswordHandler)
	authGroup.Post("/forgotpassword/otp", controller.VerifyOTPHandler)
	authGroup.Put("/forgotpassword/changepassword", controller.ChangedPasswordHandler)
}

func setupQuestRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	repository := repositories.NewGormQuestionRepository(db)
	usecase := usecases.NewQuestionUseCase(repository)
	controller := controllers.NewQuestionController(usecase)

	questionGroup := app.Group("/question", middlewares.JWTMiddleware(jwt), middlewares.AdminMiddleware)
	questionGroup.Post("/", controller.CreateQuestionHandler)
	questionGroup.Get("/", controller.GetAllQuestionHandler)
	questionGroup.Get("/:id", controller.GetQuestionByIDHandler)
	questionGroup.Put("/:id", controller.UpdateQuestionByIDHandler)
	questionGroup.Delete("/:id", controller.DeleteQuestionByIDHandler)
}

func setupVideoRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase) {
	repository := repositories.NewGormVideoRepository(db)
	likerepository := repositories.NewGormLikesRepository(db)
	usecase := usecases.NewVideoUseCase(repository, likerepository, supa)
	controller := controllers.NewVideoController(usecase)

	videoGroup := app.Group("/video", middlewares.JWTMiddleware(jwt))
	videoGroup.Post("/", middlewares.AdminMiddleware, controller.CreateVideoHandler)
	videoGroup.Get("/", controller.GetAllVideoHandler)
	videoGroup.Get("/:id", controller.GetVideoByIDHandler)
	videoGroup.Put("/:id", middlewares.AdminMiddleware, controller.UpdateVideoHandler)
	videoGroup.Delete("/:id", middlewares.AdminMiddleware, controller.DeleteVideoByIDHandler)
}

func setupUserRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase, mail configs.Mail) {
	repository := repositories.NewGormUserRepository(db)
	kidrepository := repositories.NewGormKidsRepository(db)
	usecase := usecases.NewUserUseCase(repository, supa, mail)
	kidusecase := usecases.NewKidUseCase(kidrepository, supa)
	controller := controllers.NewUserController(usecase, kidusecase)

	userGroup := app.Group("/user", middlewares.JWTMiddleware(jwt))
	userGroup.Post("/", middlewares.AdminMiddleware, controller.CreateUserandKidsHandler)
	userGroup.Get("/", middlewares.AdminMiddleware, controller.GetAllMomHandler)
	userGroup.Get("/info/:id", middlewares.AdminMiddleware, controller.GetMomByIDHandler)
	userGroup.Put("/", controller.UpdateUserByIDForUserHandler)
	userGroup.Put("/:id", controller.UpdateUserByIDForAdminHandler)
	userGroup.Delete("/:id", middlewares.AdminMiddleware, controller.DeleteUserHandler)
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

func setupQuizRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase) {
	repository := repositories.NewGormQuizRepository(db)
	usecase := usecases.NewQuizUseCase(repository, supa)
	controller := controllers.NewQuizController(usecase)

	quizGroup := app.Group("/quiz", middlewares.JWTMiddleware(jwt))
	quizGroup.Post("/", middlewares.AdminMiddleware, controller.CreateQuizHandler)
	quizGroup.Get("/", controller.GetAllQuizHandler)
	quizGroup.Get("/period/:period/category/:category/question/:id", controller.GetQuizByIDandPeriodHandler)
	quizGroup.Get("/period/:period/category/:category", controller.GetQuizByCategoryandPeriodHandler)
	quizGroup.Get("/:id", controller.GetQuizByIDHandler)
	quizGroup.Put("/:id", middlewares.AdminMiddleware, controller.UpdateQuizByIDHandler)
	quizGroup.Delete("/:id", middlewares.AdminMiddleware, controller.DeleteQuizByIDHandler)
}

func setupKidRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase) {
	repository := repositories.NewGormKidsRepository(db)
	usecase := usecases.NewKidUseCase(repository, supa)
	controller := controllers.NewKidController(usecase)

	kidGroup := app.Group("/kid", middlewares.JWTMiddleware(jwt))
	kidGroup.Post("/:id", middlewares.AdminMiddleware, controller.CreateKidHandler)
	kidGroup.Get("/:id", controller.GetKidByIDHandler)
	kidGroup.Put("/:id", middlewares.AdminMiddleware, controller.UpdateKidByIDHandler)
}

func setupAppointRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	repository := repositories.NewGormAppRepository(db)
	usecase := usecases.NewAppUseCase(repository)
	controller := controllers.NewAppController(usecase)

	appointGroup := app.Group("/appoint", middlewares.JWTMiddleware(jwt))
	appointGroup.Post("/:userID", middlewares.AdminMiddleware, controller.CreateAppointmentHandler)
	appointGroup.Get("/", controller.GetAppHandler)
	appointGroup.Get("/:id", controller.GetAppByIDHandler)
	appointGroup.Get("/history/progress", controller.GetAppInProgressUserIDHandler)
	appointGroup.Get("/history/mom/:id", middlewares.AdminMiddleware, controller.GetAllAppUserIDHandler)
	appointGroup.Put("/:id", middlewares.AdminMiddleware, controller.UpdateAppByIDHandler)
	appointGroup.Delete("/:id", middlewares.AdminMiddleware, controller.DeleteAppByIDHandler)
}

func setupCareRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT, supa configs.Supabase) {
	repository := repositories.NewGormCareRepository(db)
	usecase := usecases.NewCareUseCase(repository, supa)
	controller := controllers.NewCareController(usecase)

	careGroup := app.Group("/care", middlewares.JWTMiddleware(jwt))
	careGroup.Post("/", middlewares.AdminMiddleware, controller.CreateCareHandler)
	careGroup.Get("/", controller.GetAllCareHandler)
	careGroup.Get("/:id", controller.GetCareByID)
	careGroup.Put("/:id", middlewares.AdminMiddleware, controller.UpdateCareHandler)
	careGroup.Delete("/:id", middlewares.AdminMiddleware, controller.DeleteCareCareHandler)
}

func setupHistoryRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	repository := repositories.NewGormHistoryRepository(db)
	evaluate := repositories.NewGormEvaluateRepository(db)
	usecase := usecases.NewHistoryUseCase(repository, evaluate)
	controller := controllers.NewHistoryController(usecase)

	historyGroup := app.Group("/history", middlewares.JWTMiddleware(jwt))
	historyGroup.Post("/evaluate/:times/category/:category/kid/:id", controller.CreateHistoryHandler)
	historyGroup.Get("/evaluate/:times/kid/:id", controller.GetHistoryHandler)
	historyGroup.Get("/latest/:times/category/:category/kid/:id", controller.GetLatestHistoryHandler)
}

func setupEvaluateRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	repository := repositories.NewGormEvaluateRepository(db)
	usecase := usecases.NewEvaluateUseCase(repository)
	controller := controllers.NewEvaluateController(usecase)

	evaluateGroup := app.Group("/evaluate", middlewares.JWTMiddleware(jwt))
	evaluateGroup.Get("/all/:id", controller.GetAllEvaluateHandler)
}

func setupGrowthRoutes(app *fiber.App, db *gorm.DB, jwt configs.JWT) {
	repository := repositories.NewGormGrowthRepository(db)
	kidrepository := repositories.NewGormKidsRepository(db)
	usecase := usecases.NewGrowthUseCase(repository, kidrepository)
	controller := controllers.NewGrowthController(usecase)

	growthGroup := app.Group("/growth", middlewares.JWTMiddleware(jwt))
	growthGroup.Post("/kid/:id", controller.CreateGrowthHandler)
	growthGroup.Get("/kid/:id/summary", controller.GetSummary)
	growthGroup.Get("/kid/:id/all", controller.GetAllGrowth)
	growthGroup.Put("/:id", controller.UpdateGrowthByID)
}
