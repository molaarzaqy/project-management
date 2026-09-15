package repositories

import (
	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardPositionRepository interface {
	GetByList(listPublicID string) (*models.CardPosition, error)
	CreateOrUpdate(listPublicID string, cardOrder []uuid.UUID) error
	GetCardOrder(listPublicID string) ([]uuid.UUID, error)
	RebuildCardOrder(listPublicID string) error
}

type cardPositionRepository struct {
	DB *gorm.DB
}

func NewCardPositionRepository(db *gorm.DB) CardPositionRepository {
	return &cardPositionRepository{DB: db}
}

func (r *cardPositionRepository) GetByList(listPublicID string) (*models.CardPosition, error) {
	var position models.CardPosition
	err := r.DB.Joins("JOIN lists ON lists.internal_id = card_positions.list_internal_id").
		Where("lists.public_id = ?", listPublicID).
		First(&position).Error
	return &position, err
}

func (r *cardPositionRepository) CreateOrUpdate(listPublicID string, cardOrder []uuid.UUID) error {
	return r.DB.Exec(`
		INSERT INTO card_positions (list_internal_id, card_order)
		SELECT internal_id, ? FROM lists WHERE public_id = ?
		ON CONFLICT (list_internal_id) 
		DO UPDATE SET card_order = EXCLUDED.card_order
	`, cardOrder, listPublicID).Error
}

func (r *cardPositionRepository) GetCardOrder(listPublicID string) ([]uuid.UUID, error) {
	position, err := r.GetByList(listPublicID)
	if err != nil {
		return nil, err
	}
	return position.CardOrder, nil
}

func (r *cardPositionRepository) RebuildCardOrder(listPublicID string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var cards []models.Card
		if err := tx.Joins("JOIN lists ON lists.internal_id = cards.list_internal_id").
			Where("lists.public_id = ?", listPublicID).
			Order("position ASC"). // Gunakan kolom position sebagai fallback
			Find(&cards).Error; err != nil {
			return err
		}

		var cardOrder []uuid.UUID
		for _, card := range cards {
			cardOrder = append(cardOrder, card.PublicID)
		}

		return r.CreateOrUpdate(listPublicID, cardOrder)
	})
}
