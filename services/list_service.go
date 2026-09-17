package services

import (
	"errors"
	"fmt"

	"github.com/MaulanaBarzaqi/project-management/config"
	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/models/types"
	"github.com/MaulanaBarzaqi/project-management/repositories"
	"github.com/MaulanaBarzaqi/project-management/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListService interface {
	GetByBoardPublicID(boardPublicID string) (*ListWithOrder, error)
	GetByID(id uint) (*models.List, error)
	GetByPublicID(publicID string) (*models.List, error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePositions(boardPublicID string, positions []uuid.UUID) error
}

type ListWithOrder struct {
	Positions []uuid.UUID
	Lists []models.List
}

type listService struct {
	listRepo repositories.ListRepository
	listPosRepo repositories.ListPositionRepository
	boardRepo repositories.BoardRepository
}

func NewListService(
	listRepo repositories.ListRepository, 
	listPosRepo repositories.ListPositionRepository, 
	boardRepo repositories.BoardRepository,
) ListService {
	return &listService{
		listRepo: listRepo,
		listPosRepo: listPosRepo,
		boardRepo: boardRepo,
	}
}

func (s *listService) GetByBoardPublicID(boardPublicID string) (*ListWithOrder, error) {
	// verivikasi board
	_, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, errors.New("board not found")
	}
	position, err := s.listPosRepo.GetListOrder(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to get list oreder : " + err.Error())
	}
	if len(position) == 0 {
		return nil, errors.New("list position not found")
	}
	list, err := s.listRepo.FindByBoardID(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to get list : " + err.Error())
	}
	// sorting by position
	orderedList := utils.SortingListsByPosition(list, position)

	return &ListWithOrder{
		Positions: position,
		Lists: orderedList,
	}, nil
}

func (s *listService) GetByID(id uint) (*models.List, error) {
	return s.listRepo.FindByID(id)
}

func (s *listService) GetByPublicID(publicID string) (*models.List, error) {
	return s.listRepo.FindByPublicID(publicID)
}

func (s *listService) Create(list *models.List) error {
	// validasi board
	board, err := s.boardRepo.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound) { // pertanyakan
			return errors.New("board not found")
		}
		return fmt.Errorf("failed to get board : %w", err)
	}
	list.BoardInternalID = board.InternalID
	
	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}
	// transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// simpan list baru
	if err := tx.Create(list).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create list : %w", err)
	}
	// update position
	var position models.ListPosition
	res := tx.Where("board_internal_id = ?", board.InternalID).First(&position)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// buat baru jika belum ada
		position = models.ListPosition{
			PublicID: uuid.New(),
			BoardID: board.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}
		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create list position : %w", err)
		}
	} else if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create list position : %w", res.Error)
	}else {
		// tambah id baru
		position.ListOrder = append(position.ListOrder, list.PublicID)
		// update db
		if err := tx.Model(&position).Update("list_order", position.ListOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update list position : %w", err)
		}
	}
	// commit trx
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

func (s *listService) Update(list *models.List) error {
	return s.listRepo.Update(list)
}

func (s *listService) Delete(id uint) error {
	return s.listRepo.Delete(id)
}

func (s *listService) UpdatePositions(boardPublicID string, positions []uuid.UUID) error {
	// verifikasi board
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}
	// get list position
	position, err := s.listPosRepo.GetByBoard(board.PublicID.String())
	if err != nil {
		return errors.New("list position not found")
	}
	// update list order
	position.ListOrder = positions
	return s.listPosRepo.UpdateListOrder(position)
}