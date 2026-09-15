package repositories

import (
	"time"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardRepository interface {
	FindAll() ([]models.Board, error)
	FindAllByUser(userPublicID string) ([]models.Board, error)
	FindByID(id uint) (*models.Board, error)
	Create(board *models.Board) error
	Update(board *models.Board) error
	Delete(id uint) error
	FindByPublicID(publicID string) (*models.Board, error)
	AddMembers(boardID uint, userIDs []uint) error
	AddMember(boardID, userID uint) error
	RemoveMember(boardID, userIDs uint) error
	RemoveMembers(boardID uint, userIDs []uint) error
	GetListPositions(boardPublicID string) ([]uuid.UUID, error)
	FindAllPaginate(filter, sort string, limit, offset int) ([]models.Board, int64, error)
	FindAllByUserPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error)
}

type boardRepository struct {
	DB *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{DB: db}
}

func (r *boardRepository) FindAll() ([]models.Board, error) {
	var boards []models.Board
	err := r.DB.Find(&boards).Error
	return boards, err
}

func (r *boardRepository) FindAllByUser(userPublicID string) ([]models.Board, error) {
	var boards []models.Board
	err := r.DB.Where("owner_public_id = ? OR internal_id IN (SELECT board_members.board_internal_id FROM board_members JOIN users ON users.internal_id = board_members.user_internal_id Where users.public_id = ?)", userPublicID, userPublicID).
		Where(&boards).Error
	return boards, err
}

func (r *boardRepository) FindByID(id uint) (*models.Board, error) {
	var board models.Board
	err := r.DB.First(&board, id).Error
	return &board, err
}

func (r *boardRepository) Create(board *models.Board) error {
	return r.DB.Create(board).Error
}

func (r *boardRepository) Update(board *models.Board) error {
	return r.DB.Model(&models.Board{}).Where("public_id = ?",board.PublicID).Updates(map[string]interface{}{
		"title": board.Title,
		"description": board.Description,
		"due_date": board.DueDate,
	}).Error
}

func (r *boardRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Board{}, id).Error
}

func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
	var board models.Board
	err := r.DB.Where("public_id = ?", publicID).First(&board).Error
	return &board, err
}

func (r *boardRepository) AddMembers(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}
	now := time.Now()
	var members []models.BoardMember
	for _, userID := range userIDs {
		members = append(members, models.BoardMember{
			BoardID: int64(boardID),
			UserID: int64(userID),
			JoinedAt: now,
		})
	}
	return r.DB.Create(&members).Error
}

func (r *boardRepository) AddMember(boardID, userID uint) error {
	member := models.BoardMember{
		BoardID: int64(boardID),
		UserID: int64(userID),
		JoinedAt: time.Now(),
	}
	return r.DB.Create(&member).Error
}

func (r *boardRepository) RemoveMember(boardID, userID uint) error {
	return r.DB.Where("board_internal_id = ? AND user_internal_id = ?", boardID, userID).
		Delete(&models.BoardMember{}).Error
}

func (r *boardRepository) RemoveMembers(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}
	return r.DB.Where("board_internal_id = ? AND user_internal_id IN (?)", boardID, userIDs).
	Delete(models.BoardMember{}).Error
}

func (r *boardRepository) GetListPositions(boardPublicID string) ([]uuid.UUID, error) {
	var position models.ListPosition
	err := r.DB.Joins("JOIN boards ON boards.internal_id = list_positions.board_internal_id").
		Where("boards.public_id ?", boardPublicID).First(&position).Error
	return position.ListOrder, err
}

func (r *boardRepository) FindAllPaginate(filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	var boards []models.Board
	var total int64

	query := r.DB.Model(&models.Board{})

	// filter
	if filter != "" {
		query = query.Where("title ILIKE ?", "%"+filter+"%")
	}
	// hitung total sebelum limit/offset
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// sorting
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at desc")
	}
	// paginate
	if err := query.Limit(limit).Offset(offset).Find(&boards).Error; err != nil {
		return nil, 0, err
	}
	return boards, total, nil
}

func (r *boardRepository) FindAllByUserPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	var board []models.Board
	var total int64

	query := r.DB.Model(&models.Board{}).
	Where("owner_public_id = ? OR internal_id IN ("+
		  "SELECT board_members.board_internal_id FROM board_members "+
		  "JOIN users ON users.internal_id = board_members.user_internal_id "+
		  "WHERE users.public_id = ?)", userPublicID, userPublicID)
	if filter != "" {
		query = query.Where("title ILIKE ?", "%"+filter+"%")
	}
	// counting
	if err := query.Count(&total).Error;err != nil {
		return nil, 0, err
	}
	// sorting created at
	if sort != "" {
		query = query.Order(sort)
	}else {
		query = query.Order("created_at desc")
	}
	if err := query.Limit(limit).Offset(offset).Find(&board).Error; err != nil {
		return nil, 0, err
	}
	return board, total, nil
}