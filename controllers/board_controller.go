package controllers

import (
	"math"
	"strconv"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/services"
	"github.com/MaulanaBarzaqi/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	var err error

	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "failed to read request", err.Error())
	}

	userID, err = uuid.Parse(claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "failed to read request", err.Error())
	}
	board.OwnerPublicID = userID

	if err:= c.service.Create(board); err !=nil {
		return utils.BadRequest(ctx, "failed to create board", err.Error())
	}
	return utils.Success(ctx, "board created successfully", board)
}

func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(models.Board)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "failed to parsing data", err.Error())
	}
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID not Valid", err.Error())
	}
	existingBoard, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board not found", err.Error())
	}
	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerID = existingBoard.OwnerID
	board.OwnerPublicID = existingBoard.OwnerPublicID
	board.CreatedAt = existingBoard.CreatedAt

	if err := c.service.Update(board); err != nil {
		return utils.BadRequest(ctx, "failed to update board", err.Error())
	}
	return utils.Success(ctx, "Success to Update Board", board)
}

func (c *BoardController) AddBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed to parsing data", err.Error())
	}
	if err := c.service.AddMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed to Added member", err.Error())
	}
	return utils.Success(ctx, "Success to Added member", nil)
}

func (c *BoardController) RemoveBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed to parsing data", err.Error())
	}
	if err := c.service.RemoveMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed to remove member", err.Error())
	}
	return utils.Success(ctx, "Success to remove member", nil)
}

func (c *BoardController) GetMyBoardPaginate(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["pub_id"].(string)

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	board, total, err := c.service.GetAllByUserPaginate(userID,filter,sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "failed to get board data", err.Error())
	}

	meta := utils.PaginationMeta{
		Page: page,
		Limit: limit,
		Total: int(total),
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		Filter: filter,
		Sort: sort,
	}
	return utils.SuccessPagination(ctx, "Success to Get Board data", board, meta)
}

func (c *BoardController) GetDetail(ctx *fiber.Ctx) error {
    publicID := ctx.Params("id")
    
    board, err := c.service.GetByPublicID(publicID)
    if err != nil {
        return utils.NotFound(ctx, "Board not found", err.Error())
    }
    
    return utils.Success(ctx, "Success get board detail", board)
}

func (c *BoardController) GetBoardMembers(ctx *fiber.Ctx) error {
    publicID := ctx.Params("id")
    
    members, err := c.service.GetMembers(publicID)
    if err != nil {
        return utils.NotFound(ctx, "Failed to get members", err.Error())
    }
    
    return utils.Success(ctx, "Success get members", members)
}