package repositories

import (
	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListRepository interface {
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePositions(boardPublicID string, positions []string) error
	GetCardPositions(listPublicID string) ([]uuid.UUID, error)
	FindByBoardID(boardPublicID string) ([]models.List, error)
	FindByPublicID(publicID string) (*models.List, error)
	FindByID(id uint) (*models.List, error)
}

type listRepository struct {
	DB *gorm.DB
}

func NewListRepository(db *gorm.DB) ListRepository {
	return &listRepository{DB: db}
}

func (r *listRepository) Create(list *models.List) error {
	return r.DB.Create(list).Error
}

func (r *listRepository) Update( list *models.List) error {
	return r.DB.Model(&models.List{}).
	Where("public_id = ?", list.PublicID).Updates(map[string]interface{}{
		"title" : list.Title,
	}).Error
}

func (r *listRepository) Delete(id uint) error {
	return r.DB.Delete(&models.List{},id).Error
}

func (r *listRepository) UpdatePositions(boardPublicID string, positions []string) error {
	return r.DB.Model(&models.ListPosition{}).
	Where("board_internal_id = (Select internal_id FROM boards Where public_id = ?)", boardPublicID).
	Update("list_order", positions).Error
}

func (r *listRepository) GetCardPositions(listPublicID string) ([]uuid.UUID, error) {
	var position models.CardPosition
	err := r.DB.Joins("JOIN lists ON list.internal_id = card_positions.list_internal_id").
	Where("list.public_id = ?", listPublicID).Error
	return position.CardOrder, err
}

func (r *listRepository) FindByBoardID(boardPublicID string) ([]models.List, error) {
	var list []models.List
	err := r.DB.Where("board_public_id = ?", boardPublicID).Order("internal_id ASC").Find(&list).Error

	return list, err
}

func (r *listRepository) FindByPublicID(publicID string) (*models.List, error) {
	var list models.List
	err := r.DB.Where("public_id = ?", publicID).First(&list).Error

	return &list, err
}

func (r *listRepository) FindByID(id uint) (*models.List, error) {
	var list models.List

	err := r.DB.First(&list, id).Error
	return &list, err
}