package repositories

import (
	"errors"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	FindByCardPublicID(cardPublicID string) ([]models.CardAttachment, error)
	Create(attachment *models.CardAttachment) error
	DeleteByPublicID(publicID uuid.UUID) error
	FindByID(id uint) (*models.CardAttachment, error)
	Delete(id uint) error
	FindByPublicID(publicID uuid.UUID) (*models.CardAttachment, error)
}

type attachmentRepository struct {
	DB *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{DB: db}
}

func (r *attachmentRepository) FindByCardPublicID(cardPublicID string) ([]models.CardAttachment, error) {
	var card models.Card
	if err := r.DB.Where("publid_id = ?", cardPublicID).Find(&card).Error; err != nil {
		return nil, err
	}
	var attachments []models.CardAttachment
	if err := r.DB.Where("card_internal_id = ?", card.InternalID).Find((&attachments)).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *attachmentRepository) Create(attachment *models.CardAttachment) error {
	return r.DB.Create(attachment).Error
}

func (r *attachmentRepository) DeleteByPublicID(pubID uuid.UUID) error {
	return r.DB.Where("public_id = ?", pubID).Delete(&models.CardAttachment{}).Error
}

func (r *attachmentRepository) FindByID(id uint) (*models.CardAttachment, error) {
	var att models.CardAttachment
	if err := r.DB.Where("internal_id = ?", id).First(&att).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &att, nil
}

func (r *attachmentRepository) Delete(id uint) error {
	return r.DB.Where("internal_id = ?", id).Delete(&models.CardAttachment{}).Error
}

func (r *attachmentRepository) FindByPublicID(pubID uuid.UUID) (*models.CardAttachment, error) {
	var att models.CardAttachment
	if err := r.DB.Where("public_id = ?", pubID).First(&att).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &att, nil
}
