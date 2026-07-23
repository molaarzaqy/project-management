package services

import (
	"errors"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/repositories"
	"github.com/google/uuid"
)

type BoardService interface {
	Create(board *models.Board) error
	Update(board *models.Board) error 
	GetByPublicID(publicID string) (*models.Board, error)
	AddMembers(boardPublicID string, userPublicIDS []string) error
	RemoveMembers(boardPublicID string, userPublicIDs []string) error
	GetAllByUserPaginate(userID, filter, sort string, limit, offset int) ([]models.Board, int64, error)
	GetMembers(boardPublicID string) ([]models.User, error)
}

type boardService struct {
	boardRepo repositories.BoardRepository
	userRepo repositories.UserRepository
	boardMemberRepo repositories.BoardMemberRepository
}

func NewBoardService(
	boardRepo repositories.BoardRepository, 
	userRepo repositories.UserRepository,
	boardMemberRepo repositories.BoardMemberRepository,
	) BoardService {
	return &boardService{boardRepo,userRepo, boardMemberRepo}
}

func (s *boardService) Create(board *models.Board) error {
	user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("owner not found")
	}
	board.PublicID = uuid.New()
	board.OwnerID = user.InternalID
	return s.boardRepo.Create(board)
}

func (s *boardService) Update(board *models.Board) error {
	return s.boardRepo.Update(board)
}

func (s *boardService) GetByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepo.FindByPublicID(publicID)
}

func (s *boardService) AddMembers(boardPublicID string, userPublicIDS []string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}
	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDS {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user not found"+ userPublicID)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}
	// cek anggota sudah ada di board
	existingMember, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return  err
	}
	// cek menggunakan map untuk efisiensi pencarian
	memberMap := make(map[uint]bool)
	for _, member := range existingMember {
		memberMap[uint(member.InternalID)] = true
	}
	var newMemberIDs []uint
	for _, userID := range userInternalIDs {
		if !memberMap[userID] {
			newMemberIDs = append(newMemberIDs, userID)
		}
	}
	if len(newMemberIDs) == 0 {
		return  nil
	}
	return s.boardRepo.AddMember(uint(board.InternalID), newMemberIDs)
}

func (s *boardService) RemoveMembers(boardPublicID string, userPublicIDs []string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}
	// validasi user
	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user not found"+ userPublicID)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}
	existingMember, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return  err
	}
	memberMap := make(map[uint]bool)
	for _, member := range existingMember {
		memberMap[uint(member.InternalID)] = true
	}

	var membersToRemove []uint
	for _,userID := range userInternalIDs {
		if memberMap[userID] {
			membersToRemove = append(membersToRemove, userID)
		}
	}
	return s.boardRepo.RemoveMember(uint(board.InternalID), membersToRemove)
}

func (s *boardService) GetAllByUserPaginate(userID, filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	return s.boardRepo.FindAllUserPaginate(userID, filter, sort, limit, offset)
}

func (s *boardService) GetMembers(boardPublicID string) ([]models.User, error) {
    // Validasi apakah board ada
    _, err := s.boardRepo.FindByPublicID(boardPublicID)
    if err != nil {
        return nil, errors.New("board not found")
    }

    // Ambil data user yang menjadi member melalui repository
    return s.boardMemberRepo.GetMembers(boardPublicID)
}