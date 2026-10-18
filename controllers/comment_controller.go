package controllers

import (
	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/services"
	"github.com/MaulanaBarzaqi/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CommentController struct {
	service services.CommentService
}

func NewCommentController(s services.CommentService) *CommentController {
	return &CommentController{service: s}
}

func (c *CommentController) CreateComment(ctx *fiber.Ctx) error {
	comment := new(models.Comment)
	if err := ctx.BodyParser(comment); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	// Set user ID from JWT
	userID := ctx.Locals("userID").(string)
	comment.UserPubID = uuid.MustParse(userID)

	if err := c.service.Create(comment); err != nil {
		return utils.BadRequest(ctx, "Gagal Membuat Komentar", err.Error())
	}

	return utils.Success(ctx, "Komentar berhasil dibuat", comment)
}

func (c *CommentController) UpdateComment(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	comment := new(models.Comment)

	if err := ctx.BodyParser(comment); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	//existingComment, err := c.service.GetByID(publicID)
	existingComment, err := c.service.GetByPublicID(publicID)

	if err != nil {
		return utils.NotFound(ctx, "Komentar tidak ditemukan", err.Error())
	}

	// Check if user is the owner of the comment
	userID := ctx.Locals("userID").(string)
	if userID != existingComment.UserPubID.String() {
		return utils.Unauthorized(ctx, "Anda tidak memiliki akses untuk mengubah komentar ini", err.Error())
	}

	comment.InternalID = existingComment.InternalID
	comment.PublicID = existingComment.PublicID

	if err := c.service.Update(comment); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Komentar", err.Error())
	}

	return utils.Success(ctx, "Komentar berhasil diperbarui", comment)
}


func (c *CommentController) GetCardComments(ctx *fiber.Ctx) error {
	cardID := ctx.Params("card_id")
	if _, err := uuid.Parse(cardID); err != nil {
		return utils.BadRequest(ctx, "ID Card tidak valid", err.Error())
	}

	comments, err := c.service.GetByPublicID(cardID)
	if err != nil {
		return utils.NotFound(ctx, "Komentar tidak ditemukan", err.Error())
	}

	return utils.Success(ctx, "Data komentar berhasil diambil", comments)
}

func (c *CommentController) DeleteComment(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	comment, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Komentar tidak ditemukan", err.Error())
	}

	// Check if user is the owner of the comment
	userID := ctx.Locals("userID").(string)
	if userID != comment.UserPubID.String() {
		return utils.Unauthorized(ctx, "Anda tidak memiliki akses untuk menghapus komentar ini", err.Error())
	}

	if err := c.service.Delete(uint(comment.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Komentar", err.Error())
	}

	return utils.Success(ctx, "Komentar berhasil dihapus", publicID)
}