package repositories

import (
	"github.com/MaulanaBarzaqi/project-management/models"
	"gorm.io/gorm"
)

type LabelRepository interface {
	Create(label *models.Label) error
	Update(label *models.Label) error
	Delete(id uint) error
	FindAll() ([]models.Label, error)
	FindByID(id uint) (*models.Label, error)
	FindByPublicID(publicID string) (*models.Label, error)
}

type labelRepository struct {
	DB *gorm.DB
}

func NewLabelRepository(db *gorm.DB) LabelRepository {
	return &labelRepository{DB: db}
}

func (r *labelRepository) Create(label *models.Label) error {
	return r.DB.Create(label).Error
}

func (r *labelRepository) Update(label *models.Label) error {
	return r.DB.Save(label).Error
}

func (r *labelRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Label{}, id).Error
}

func (r *labelRepository) FindAll() ([]models.Label, error) {
	var labels []models.Label
	err := r.DB.Find(&labels).Error
	return labels, err
}

func (r *labelRepository) FindByID(id uint) (*models.Label, error) {
	var label models.Label
	err := r.DB.First(&label, id).Error
	return &label, err
}

func (r *labelRepository) FindByPublicID(publicID string) (*models.Label, error) {
	var label models.Label
	err := r.DB.Where("public_id = ?", publicID).First(&label).Error
	return &label, err
}