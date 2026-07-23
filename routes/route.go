package routes

import (
	"log"

	"github.com/MaulanaBarzaqi/project-management/config"
	"github.com/MaulanaBarzaqi/project-management/controllers"
	"github.com/MaulanaBarzaqi/project-management/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
)

func Setup(
	app *fiber.App, 
	uc *controllers.UserController,
	bc *controllers.BoardController,
	lc *controllers.ListController,
	cc *controllers.CardController,
	) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	} 
	app.Post("/v1/auth/register", uc.Register)
	app.Post("/v1/auth/login", uc.Login)

	// protected routes
	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c,"error unauthorized", err.Error())
		},
	}))
	// user
	userGroup := api.Group("/users")
	userGroup.Get("/page", uc.GetUserPagination)
	userGroup.Get("/:id", uc.GetUser)
	userGroup.Put("/:id", uc.UpdateUser)
	userGroup.Delete("/:id", uc.DeleteUser)
	// board
	boardGroup := api.Group("/boards")
	boardGroup.Get("/my", bc.GetMyBoardPaginate)
	boardGroup.Post("/", bc.CreateBoard)
	boardGroup.Post("/:id/members", bc.AddBoardMembers)
	boardGroup.Delete("/:id/members", bc.RemoveBoardMembers)
	boardGroup.Put("/:id", bc.UpdateBoard)
	boardGroup.Get("/:board_id/lists",lc.GetListOnBoard)
	boardGroup.Put("/:board_id/positions", lc.UpdateListPosition)
	boardGroup.Get("/:id", bc.GetDetail)
	boardGroup.Get("/:id/members", bc.GetBoardMembers)
	// list
	listGroup := api.Group("/lists")
	listGroup.Post("/", lc.CreateList)
	listGroup.Put("/:id", lc.UpdateList)
	listGroup.Delete("/:id", lc.DeleteList)
	listGroup.Get("/:id/cards", cc.GetCardsByList)
	listGroup.Put("/:list_id/positions", lc.UpdateCardPosition)
	// card
	cardGroup := api.Group("/cards")
	cardGroup.Post("/", cc.CreateCard)
	cardGroup.Put("/:id", cc.UpdateCard)
	cardGroup.Delete("/:id", cc.DeleteCard)
	cardGroup.Get("/:id", cc.GetCardDetail)

	// cardGroup.Post(":id/attachments", cc.UploadAttachment)
	// cardGroup.Get(":id/attachments", cc.GetAttachments)
	// cardGroup.Delete("/:card_id/attachments/:attachment_id", cc.DeleteAttachment)
}