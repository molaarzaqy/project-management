package repositories

import (
	"github.com/MaulanaBarzaqi/project-management/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	FindByCardPublicID(cardPublicID string) ([]models.Comment, error)
	FindByID(id uint) (*models.Comment, error)
	FindByPublicID(publicID string) (*models.Comment, error)
	Create(comment *models.Comment) error
	Update(comment *models.Comment) error
	Delete(id uint) error
}

type commentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{DB: db}
}

func (r *commentRepository) FindByCardPublicID(cardPublicID string) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.DB.Joins("JOIN cards ON cards.internal_id ? = comments.card_internal_id").
		Where("cards.public_id = ?", cardPublicID).
		Order("created_at DESC").
		Find(&comments).Error

	return comments, err
}

func (r *commentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.DB.First(&comment, id).Error

	return &comment, err
}

func (r *commentRepository) FindByPublicID(publicID string) (*models.Comment, error) {
	var comment models.Comment
	err := r.DB.Where("public_id = ?", publicID).First(&comment).Error
	return &comment, err
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return r.DB.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Comment{}, id).Error
}