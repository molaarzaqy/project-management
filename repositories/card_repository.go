package repositories

import (
	"fmt"
	"path/filepath"

	"github.com/MaulanaBarzaqi/project-management/config"
	"github.com/MaulanaBarzaqi/project-management/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CardRepository interface {
	Create(card *models.Card) error
	Update(card *models.Card) error
	Delete( id uint) error
	FindByID(id uint) (*models.Card, error)
	FindByPublicID(publicID string) (*models.Card, error)
	FindByListPublicID(listPublicID string) ([]models.Card, error)
	FindByListID(listID int64) (*models.CardPosition, error)
	UpdatePositions(listPublicID string, positions []string) error
	AddLabel(cardID, labelID uint) error
	RemoveLabel(cardID, labelID uint) error
	AddAssignees(cardID uint, userIDs []uint) error
	RemoveAssignee(cardID, userID uint) error
}

type cardRepository struct {
	DB *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{DB: db}
} 

func (r *cardRepository) Create(card *models.Card) error {
	return r.DB.Create(card).Error
}

func (r *cardRepository) Update(card *models.Card) error {
	return  r.DB.Save(card).Error
}

func (r *cardRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Card{}, id).Error
}


func (r *cardRepository) FindByID(id uint) (*models.Card, error) {
	var card models.Card
	err := r.DB.Preload("Labels").Preload("Assignees").First(&card, id).Error

	return &card, err
}

func (r *cardRepository) FindByPublicID(publicID string) (*models.Card, error) {
	var card models.Card
	if err := r.DB.Preload("Assignees.User", func (tx *gorm.DB) *gorm.DB {
		return tx.Select("internal_id", "public_id", "name", "email")
	}).Preload("Attachments").Where("public_id = ?", publicID).First(&card).Error; err != nil {
		return nil, err
	}

	baseUrl := config.AppConfig.APPURL

	for i := range card.Attachments {
		card.Attachments[i].FileURL= fmt.Sprintf(
			"%s/files/%s",
			baseUrl,
			filepath.Base(card.Attachments[i].File),
		)
	}
	return &card,nil
}

func (r *cardRepository) FindByListPublicID(listPublicID string) ([]models.Card, error) {
	var cards []models.Card
	err := r.DB.Joins("JOIN lists ON lists.internal_id = cards.list_internal_id").
			Where("lists.public_id = ?", listPublicID).
			Order("position ASC").
			Find(&cards).Error
	return cards, err
}

func (r *cardRepository) FindByListID(id int64) (*models.CardPosition, error) {
	var position models.CardPosition
	err := r.DB.Where("list_internal_id = ?",id).First(&position).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *cardRepository) UpdatePositions(listPublicID string, positions []string) error {
	return r.DB.Model(&models.CardPosition{}).
		   Where("list_internal_id = (SELECT internal_id FROM lists Where public_id = ?)", listPublicID).
		   Update("card_order", positions).Error
}

func (r *cardRepository) AddLabel(cardID, labelID uint) error {
	cardLabel := models.CardLabel{
		CardID: int64(cardID),
		LabelID: int64(labelID),
	}
	return r.DB.Create(cardLabel).Error
}

func (r *cardRepository) RemoveLabel(cardID, labelID uint) error {
	return r.DB.Where("card_internal_id = ? AND label_internal_id = ?", cardID, labelID).
		Delete(&models.CardLabel{}).Error
}

func (r *cardRepository) AddAssignees(cardID uint, userIDs []uint) error {
	var assignees []models.CardAssignee
	for _, userID := range userIDs {
		assignees = append(assignees, models.CardAssignee{
			CardID: int64(cardID),
			UserID: int64(userID),
		})
	}
	return r.DB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&assignees, 100).Error
}

func (r *cardRepository) RemoveAssignee(cardID, userID uint) error {
	return r.DB.Where("card_internal_id = ? AND user_internal_id = ?", cardID, userID).
		Delete(&models.CardAssignee{}).Error
}
