package main

import (
	"log"

	"github.com/MaulanaBarzaqi/project-management/config"
	"github.com/MaulanaBarzaqi/project-management/controllers"
	"github.com/MaulanaBarzaqi/project-management/database/seed"
	_ "github.com/MaulanaBarzaqi/project-management/docs"
	"github.com/MaulanaBarzaqi/project-management/repositories"
	"github.com/MaulanaBarzaqi/project-management/routes"
	"github.com/MaulanaBarzaqi/project-management/services"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Project Management API
// @version 1.0
// @description Ini adalah dokumentasi API untuk aplikasi Project Management.
// @termsOfService http://swagger.io/terms/

// @contact.name Support Team
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3030
// @BasePath /
func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()
	app := fiber.New()
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	// user setup
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)
	// board setup
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)
    // list
	listPosRepo := repositories.NewListPositionRepository()
	listRepo := repositories.NewListRepository()
	listService := services.NewListService(listRepo, listPosRepo, boardRepo)
	listController := controllers.NewListController(listService)
	// card
	cardRepo := repositories.NewCardRepository()
	cardService := services.NewCardService(cardRepo, listRepo, userRepo)
	cardController := controllers.NewCardController(cardService) 

	routes.Setup(app,userController, boardController, listController, cardController)

	port := config.AppConfig.AppPort
	log.Println("server is running on port :", port)
	log.Fatal(app.Listen(":" + port))
}