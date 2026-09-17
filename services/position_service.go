package services

import (
	"errors"

	"github.com/MaulanaBarzaqi/project-management/repositories"
	"github.com/google/uuid"
)

type PositionService interface {
	GetListPositions(boardPublicID string) ([]uuid.UUID, error)
	UpdateListPositions(boardPublicID string, positions []uuid.UUID) error
	GetCardPositions(listPublicID string) ([]uuid.UUID, error)
	UpdateCardPositions(listPublicID string, positions []uuid.UUID) error
	RebuildCardPositions(listPublicID string) error
}

type positionService struct {
	listPositionRepo repositories.ListPositionRepository
	cardPositionRepo repositories.CardPositionRepository
	boardRepo 		 repositories.BoardRepository
	listRepo		 repositories.ListRepository
}

func NewPositionService(
		listPosRepo repositories.ListPositionRepository,
		cardPosRepo repositories.CardPositionRepository,
		boardRepo repositories.BoardRepository,
		listRepo repositories.ListRepository,
	) PositionService {
		return &positionService{
			listPositionRepo: listPosRepo,
			cardPositionRepo: cardPosRepo,
			boardRepo: boardRepo,
			listRepo: listRepo,
		}
}

func (s *positionService) GetListPositions(boardPublicID string) ([]uuid.UUID, error) {
	// Verifikasi board ada
	_, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, errors.New("board not found")
	}

	return s.listPositionRepo.GetListOrder(boardPublicID)
}

func (s *positionService) UpdateListPositions(boardPublicID string, positions []uuid.UUID) error {
	// Verifikasi board ada
	_, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	return s.listPositionRepo.CreateOrUpdate(boardPublicID, positions)
}

func (s *positionService) GetCardPositions(listPublicID string) ([]uuid.UUID, error) {
	// Verifikasi list ada
	_, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return nil, errors.New("list not found")
	}

	return s.cardPositionRepo.GetCardOrder(listPublicID)
}

func (s *positionService) UpdateCardPositions(listPublicID string, positions []uuid.UUID) error {
	// Verifikasi list ada
	_, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return errors.New("list not found")
	}

	return s.cardPositionRepo.CreateOrUpdate(listPublicID, positions)
}

func (s *positionService) RebuildCardPositions(listPublicID string) error {
	// Verifikasi list ada
	_, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return errors.New("list not found")
	}

	return s.cardPositionRepo.RebuildCardOrder(listPublicID)
}