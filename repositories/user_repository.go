package repositories

import (
	"strings"

	"github.com/MaulanaBarzaqi/project-management/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
	FindAll() ([]models.User, error)
	FindAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error){
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error){
	var user models.User
	err := r.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByPublicID(publicID string) (*models.User, error){
	var user models.User
	err := r.DB.Where("public_id = ?", publicID).First(&user).Error
	return &user, err
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *userRepository) FindAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	db := r.DB.Model(&models.User{})
	// filtering
	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("name Ilike ? OR email Ilike ?", filterPattern, filterPattern)
	}
	// count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// sorting
	if sort != "" {
		if sort == "-id"{
			sort = "-internal_id"
		}else if sort == "id"{
			sort = "internal_id"
		}
		if strings.HasPrefix(sort, "-") {
			sort = strings.TrimPrefix(sort, "-") + " DESC"
		}else {
			sort += " ASC"
		}
		db = db.Order(sort)
	}
	// pagination
	err := db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.DB.Model(&models.User{}).Where("public_id = ?", user.PublicID).Updates(map[string]interface{}{
		"name": user.Name,
	}).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}