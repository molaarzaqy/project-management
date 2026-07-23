package controllers

import (
	"time"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/services"
	"github.com/MaulanaBarzaqi/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardController struct {
	services services.CardService
}

func NewCardController(s services.CardService) *CardController {
	return &CardController{services: s}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPublicID 	string 		`json:"list_id"`
		Title			string		`json:"title"`
		Description		string		`json:"description"`
		DueDate			time.Time	`json:"due_date"`
		Position		int			`json:"position"`
	}
	var req CreateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "failed to get data", err.Error())
	}
	card := &models.Card{
		Title: req.Title,
		Description: req.Description,
		DueDate: &req.DueDate,
		Position: req.Position,
	}
	if err := c.services.Create(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "failed to create card", err.Error())
	}
	return utils.Success(ctx, "success to created card", card)
}

func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	type UpdateCardRequest struct {
		ListPublicID	string		`json:"list_id"`
		Title			string 		`json:"title"`
		Description		string		`json:"description"`
		DueDate			*time.Time	`json:"due_date"`
		Position		int			`json:"position"`
	}
	var req UpdateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "failed to parsing data", err.Error())
	}
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "id not valid", err.Error())
	}
	card := &models.Card{
		Title: req.Title,
		Description: req.Description,
		DueDate: req.DueDate,
		Position: req.Position,
		PublicID: uuid.MustParse(publicID),
	}
	if err := c.services.Update(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "failed to pubdate data", err.Error())
	}
	// // masalah created at
	// updatedCard, err := c.services.GetByPublicID(publicID)
	// if err != nil {
	// 	return utils.InternalServerError(ctx, "failed to fetch updated data", err.Error())
	// }
	return utils.Success(ctx, "success to update card", card)
}

func (c *CardController) DeleteCard(ctx *fiber.Ctx) error {
	PublicID := ctx.Params("id")

	if _, err := uuid.Parse(PublicID); err != nil {
		return utils.BadRequest(ctx, "ID not valid", err.Error())
	}
	card, err := c.services.GetByPublicID(PublicID)
	if err != nil {
		return utils.NotFound(ctx, "card not found", err.Error())
	}
	if err := c.services.Delete(uint(card.InternalID)); err != nil {
		return utils.BadRequest(ctx, "failed to delete card", err.Error())
	}
	return utils.Success(ctx, "success to delete card", PublicID)
}

func (c *CardController) GetCardDetail(ctx *fiber.Ctx) error {
	cardPublicID := ctx.Params("id")
	
	card, err := c.services.GetByPublicID(cardPublicID)
	if err != nil {
		return utils.InternalServerError(ctx, "error fetch detail card", err.Error())
	}
	if card == nil {
		return utils.NotFound(ctx, "card not found", err.Error())
	}
	return utils.Success(ctx, "success to fetch card", card)
}

func (c *CardController) GetCardsByList(ctx *fiber.Ctx) error {
    listPublicID := ctx.Params("id")
    if _, err := uuid.Parse(listPublicID); err != nil {
        return utils.BadRequest(ctx, "id not valid", err.Error())
    }
    
    // Pastikan service card Anda memiliki method untuk mencari berdasarkan list ID (misal: GetByListID)
    cards, err := c.services.GetByListID(listPublicID)
    if err != nil {
        return utils.InternalServerError(ctx, "failed to fetch cards", err.Error())
    }
    
    return utils.Success(ctx, "success to fetch cards", cards)
}