package repositories

import (
	"errors"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	FindByCardID(cardPublicID string) ([]models.CardAttachment, error)
	Create(attachment *models.CardAttachment) error
	DeleteByPublicID(pubID uuid.UUID) error
	GetByPublicID(pubID uuid.UUID) (*models.CardAttachment, error)
}

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) FindByCardID(cardPublicID string) ([]models.CardAttachment, error) {
	var card models.Card
	if err := r.db.Where("publid_id = ?", cardPublicID).Find(&card).Error; err != nil {
		return nil, err
	}
	var attachments []models.CardAttachment
	if err := r.db.Where("card_internal_id = ?", card.InternalID).Find((&attachments)).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *attachmentRepository) Create(attachment *models.CardAttachment) error {
	return r.db.Create(attachment).Error
}

func (r *attachmentRepository) DeleteByPublicID(pubID uuid.UUID) error {
	return r.db.Where("public_id = ?", pubID).Delete(&models.CardAttachment{}).Error
}

func (r *attachmentRepository) GetByPublicID(pubID uuid.UUID) (*models.CardAttachment, error) {
	var att models.CardAttachment
	if err := r.db.Where("public_id = ?", pubID).First(&att).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &att, nil
}