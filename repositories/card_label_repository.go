package repositories

import (
	"github.com/MaulanaBarzaqi/project-management/models"
	"gorm.io/gorm"
)

type CardLabelRepository interface {
	GetLabels(cardPublicID string) ([]models.Label, error)
	GetHasLabel(cardPublicID, labelPublicID string) (bool, error)
}

type cardLabelRepository struct {
	DB *gorm.DB
}

func NewCardLabelRepository(db *gorm.DB) CardLabelRepository {
	return &cardLabelRepository{DB: db}
}

func (r *cardLabelRepository) GetLabels(cardPublicID string) ([]models.Label, error) {
	var labels []models.Label
	err := r.DB.Joins("JOIN card_labels ON card_labels.label_internal_id = labels.internal_id").
		Joins("JOIN cards ON cards.internal_id = card_labels.card_internal_id").
		Where("cards.public_id = ?", cardPublicID).
		Find(&labels).Error
	return labels, err
}

func (r *cardLabelRepository) GetHasLabel(cardPublicID, labelPublicID string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.CardLabel{}).
		Joins("JOIN cards ON cards.internal_id = card_labels.card_internal_id").
		Joins("JOIN labels ON labels.internal_id = card_labels.label_internal_id").
		Where("cards.public_id = ? AND labels.public_id = ?", cardPublicID, labelPublicID).
		Count(&count).Error
	return count > 0, err
} 