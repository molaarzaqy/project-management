package controllers

import (
	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/services"
	"github.com/MaulanaBarzaqi/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{service: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "failed to read request", err.Error())
	}
	if err := c.service.Create(list); err != nil {
		return utils.BadRequest(ctx, "failed to create list", err.Error())
	}
	return utils.Success(ctx, "success to created list", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "failed to parsing data", err.Error())
	}
	// validasi kalau benar public id format uuid
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "id is not valid", err.Error())
	}
	// verifikasi list kalu beneran ada yang mau diupdate
	existingList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "list not found", err.Error())
	}
	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID
	if err := c.service.Update(list); err != nil {
		return utils.BadRequest(ctx, "failed to update list", err.Error())
	}
	updatedList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "list not found", err.Error())
	}
	return utils.Success(ctx, "success to update list", updatedList)
}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "id is not valid", err.Error())
	}
	lists, err := c.service.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "list not found", err.Error())
	}
	return utils.Success(ctx, "success to get lists", lists)
}

func (c *ListController) DeleteList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "id is not valid", err.Error())
	}
	// cek listnya ada atau tidak
	list, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "list not found", err.Error())
	}
	if err := c.service.Delete(uint(list.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "failed to delete list", err.Error())
	}
	return utils.Success(ctx, "success to delete list", publicID)
}

func (c *ListController) UpdateListPosition(ctx *fiber.Ctx) error {
	boardID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardID); err != nil {
		return utils.BadRequest(ctx, "id is not valid", err.Error())
	}
	var positionUUID []uuid.UUID
	if err := ctx.BodyParser(&positionUUID); err != nil {
		// jika gagal, coba parse sebagai array of string
		var positionString []string
		if err := ctx.BodyParser(&positionString); err != nil {
			return utils.BadRequest(ctx, "invalid position format", err.Error())
		}
		// konfersi string ke uuid
		for _, s := range positionString {
			u, err := uuid.Parse(s)
			if err != nil {
				return utils.BadRequest(ctx, "invalid position format", err.Error())
			}
			positionUUID = append(positionUUID, u)
		}
	}
	if err := c.service.UpdatePositions(boardID, positionUUID); err != nil {
		return utils.InternalServerError(ctx, "failed to update list", err.Error())
	}
	return utils.Success(ctx, "success to update list position", nil)
}

func (c *ListController) UpdateCardPosition(ctx *fiber.Ctx) error {
	listID := ctx.Params("list_id")
	if _, err := uuid.Parse(listID); err != nil {
		return utils.BadRequest(ctx, "id is not valid", err.Error())
	}
	var payload struct {
		Positions []string `json:"positions"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return utils.BadRequest(ctx, "invalid request body", err.Error())
	}
	var positionUUID []uuid.UUID
	for _, s := range payload.Positions {
		u, err := uuid.Parse(s)
		if err != nil {
			return utils.BadRequest(ctx, "invalid position format", err.Error())
		}
		positionUUID = append(positionUUID, u)
	}
	if err := c.service.UpdateCardPositions(listID, positionUUID); err != nil {
		return utils.InternalServerError(ctx, "failed to update card positions", err.Error())
	}
	return utils.Success(ctx, "Posisi card berhasil diperbarui", nil)
}