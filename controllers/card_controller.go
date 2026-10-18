package controllers

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/services"
	"github.com/MaulanaBarzaqi/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type CardController struct {
	cardService services.CardService
	attachmentService services.AttachmentService
	uploadDir string
}

func NewCardController(
	cardService services.CardService, 
	attachmentService services.AttachmentService, 
	uploadDir string,
	) *CardController {
		return &CardController{
			cardService: cardService, 
			attachmentService: attachmentService,
			uploadDir: uploadDir,
		}
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
	if err := c.cardService.Create(card, req.ListPublicID); err != nil {
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
	if err := c.cardService.Update(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "failed to pubdate data", err.Error())
	}
	// // masalah created at
	// updatedCard, err := c.cardService.GetByPublicID(publicID)
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
	card, err := c.cardService.GetByPublicID(PublicID)
	if err != nil {
		return utils.NotFound(ctx, "card not found", err.Error())
	}
	if err := c.cardService.Delete(uint(card.InternalID)); err != nil {
		return utils.BadRequest(ctx, "failed to delete card", err.Error())
	}
	return utils.Success(ctx, "success to delete card", PublicID)
}

func (c *CardController) GetCardDetail(ctx *fiber.Ctx) error {
	cardPublicID := ctx.Params("id")
	
	card, err := c.cardService.GetByPublicID(cardPublicID)
	if err != nil {
		return utils.InternalServerError(ctx, "error fetch detail card", err.Error())
	}
	if card == nil {
		return utils.NotFound(ctx, "card not found", err.Error())
	}
	return utils.Success(ctx, "success to fetch card", card)
}

func (c *CardController) GetListCard(ctx *fiber.Ctx) error {
    listPublicID := ctx.Params("id")
    if _, err := uuid.Parse(listPublicID); err != nil {
        return utils.BadRequest(ctx, "id not valid", err.Error())
    }
    
    // Pastikan service card Anda memiliki method untuk mencari berdasarkan list ID (misal: GetByListID)
    cards, err := c.cardService.GetByListPublicID(listPublicID)
    if err != nil {
        return utils.NotFound(ctx, "card not found", err.Error())
    }
    
    return utils.Success(ctx, "success to fetch cards", cards)
}

func (c *CardController) AddCardLabel(ctx *fiber.Ctx) error {
	cardPublicID := ctx.Params("id")

	var body struct {
		LabelID string `json:"label_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "invalid id", err.Error())
	}
	if err := c.cardService.AddLabel(cardPublicID, body.LabelID); err != nil {
		return utils.BadRequest(ctx, "failed to add label", err.Error())
	}
	return utils.Success(ctx, "Success to add label", nil)
}

func (c *CardController) RemoveCardLabel(ctx *fiber.Ctx) error {
	cardPublicID := ctx.Params("id")

	var body struct {
		LabelID string `json:"label_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "invalid id", err.Error())
	}
	if err := c.cardService.RemoveLabel(cardPublicID, body.LabelID); err != nil {
		return utils.BadRequest(ctx, "failed to remove label", err.Error())
	}
	return utils.Success(ctx, "Success to remove label", nil)
}

func (c *CardController) AddCardAssignee(ctx *fiber.Ctx) error {
	cardID := ctx.Params("id")

	var body struct {
		UserID []string `json:"user_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.cardService.AddAssignees(cardID, body.UserID); err != nil {
		return utils.BadRequest(ctx, "Gagal Menambahkan Assignee", err.Error())
	}

	return utils.Success(ctx, "Assignee berhasil ditambahkan", nil)
}

func (c *CardController) RemoveCardAssignee(ctx *fiber.Ctx) error {
	cardID := ctx.Params("id")

	var body struct {
		UserID string `json:"user_id"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.cardService.RemoveAssignee(cardID, body.UserID); err != nil {
		return utils.BadRequest(ctx, "Gagal Menghapus Assignee", err.Error())
	}

	return utils.Success(ctx, "Assignee berhasil dihapus", nil)
}

func (c *CardController) UpdateCardPositions(ctx *fiber.Ctx) error {
	listID := ctx.Params("list_id")
	if _, err := uuid.Parse(listID); err != nil {
		return utils.BadRequest(ctx, "ID List tidak valid", err.Error())
	}

	// Buat struct sesuai request body
	var req struct {
		Positions []string `json:"positions"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	if err := c.cardService.UpdatePositions(listID, req.Positions); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Posisi", err.Error())
	}

	return utils.Success(ctx, "Posisi card berhasil diperbarui", nil)
}

// POST /api/v1/cards/:id/attachments
func (c *CardController) UploadAttachment(ctx *fiber.Ctx) error {
	cardPublicID := ctx.Params("id")
	if cardPublicID == "" {
		return utils.BadRequest(ctx, "card is require", "Error")
	}
	// ambil user public id
	u := ctx.Locals("user")
	if u == nil {
		return utils.Unauthorized(ctx, "unauthorized", "error unauthorized")
	}
	token, ok := u.(*jwt.Token)
	if !ok {
		return utils.Unauthorized(ctx, "unauthorized", "error unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.Unauthorized(ctx, "unauthorized", "error unauthorized")
	}
	
	userPublicID, _ := claims["pub_id"].(string)
	if userPublicID == "" {
		return utils.Unauthorized(ctx, "unauthorized user id claims not found", "error unauthorized")
	}
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return utils.BadRequest(ctx, "file is required", "error")
	}
	// validasi extension / size
	const maxSize = 20 << 20 // 20mb
	if fileHeader.Size > maxSize {
		return utils.BadRequest(ctx, "file too large (max 20 MB)", "Error")
	}
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowed := map[string]bool{".png":true,".jpg":true,".jpeg":true}
	if !allowed[ext] {
		return utils.BadRequest(ctx, "file type not allowed", "Error")
	}
	// buat folder jika belum ada
	if err := os.MkdirAll(c.uploadDir, os.ModePerm); err != nil {
		return utils.InternalServerError(ctx, "Cannot create directory upload", "Error")
	}

	// simpan nama unik
	dstFileName := uuid.New().String() + ext 
	dstPath := filepath.Join(c.uploadDir, dstFileName)

	src, err := fileHeader.Open()
	if err != nil {
		return utils.InternalServerError(ctx, "Cannot open uploaded file", "Error")
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return utils.InternalServerError(ctx, "Cannot create destination file", "Error")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return utils.InternalServerError(ctx, "Cannot save file", "Error")
	}
	attachment, err := c.attachmentService.Create(cardPublicID, userPublicID, dstPath)
	if err != nil {
		_ = os.Remove(dstPath)
		return utils.InternalServerError(ctx, err.Error(), "Error")
	}
	return ctx.Status(fiber.StatusCreated).JSON(attachment)
}

func (cc *CardController) GetAttachments(c *fiber.Ctx) error {
	cardPublicID := c.Params("id")
	attachments, err := cc.attachmentService.GetByCardPublicID(cardPublicID)
	if err != nil {
		return utils.InternalServerError(c, err.Error(), "Error")
	}
	return c.JSON(attachments)
}

func (cc *CardController) DeleteAttachment(c *fiber.Ctx) error {
	pubIDStr := c.Params("attachment_id")
	pubID, err := uuid.Parse(pubIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid attachment public_id", "Error")
	}

	att, err := cc.attachmentService.GetByPublicID(pubID)
	if err != nil {
		return utils.InternalServerError(c, err.Error(), "Error")
	}
	if att == nil {
		return utils.NotFound(c, "attachment not found", err.Error())
	}

	// 🔐 Cek permission: hanya owner yang boleh hapus
	u := c.Locals("user")
	token, _ := u.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return utils.Forbidden(c, "invalid user claims")
	}
	userID := int64(userIDFloat)

	if att.UserID != userID {
		return utils.Forbidden(c, "you don't have permission to delete this attachment")
	}

	// Hapus dari database
	if err := cc.attachmentService.DeleteByPublicID(pubID); err != nil {
		return utils.InternalServerError(c, err.Error(), "Error")
	}

	// Hapus file fisik
	_ = os.Remove(att.File)

	return utils.Success(c, "Berhasil menghapus attachment", nil)
}