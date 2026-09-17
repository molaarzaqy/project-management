package services

import (
	"errors"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/repositories"
	"github.com/google/uuid"
)

type BoardService interface {
	GetAll() ([]models.Board, error)
	GetAllByUser(userPublicID string) ([]models.Board, error)
	GetByID(id uint) (*models.Board, error)
	Create(board *models.Board) error
	Update(board *models.Board) error
	Delete(id uint) error
	GetByPublicID(publicID string) (*models.Board, error)
	AddMember(boardPublicID, userPublicID string) error
	AddMembers(boardPublicID string, userPublicIDS []string) error
	RemoveMember(boardPublicID, userPublicID string) error
	RemoveMembers(boardPublicID string, userPublicIDs []string) error
	IsMember(boardPublicID, userPublicID string) (bool, error)
	GetAllByUserPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error)
	GetAllPaginate(filter, sort string, limit, offset int) ([]models.Board, int64, error)
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
		return &boardService{
			boardRepo: boardRepo,
			userRepo: userRepo,
			boardMemberRepo: boardMemberRepo,
		}
}

func (s *boardService) GetAll() ([]models.Board, error) {
	return s.boardRepo.FindAll()
}

func (s *boardService) GetAllByUser(userPublicID string) ([]models.Board, error) {
	return s.boardRepo.FindAllByUser(userPublicID)
}

func (s *boardService) GetByID(id uint) (*models.Board, error) {
	return s.boardRepo.FindByID(id)
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

func (s *boardService) Delete(id uint) error {
	return s.boardRepo.Delete(id)
}

func (s *boardService) GetByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepo.FindByPublicID(publicID)
}

func (s *boardService) AddMember(boardPublicID, userPublicID string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.boardRepo.AddMember(uint(board.InternalID), uint(user.InternalID))
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
	return s.boardRepo.AddMembers(uint(board.InternalID), newMemberIDs)
}

func (s *boardService) RemoveMember(boardPublicID, userPublicID string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.boardRepo.RemoveMember(uint(board.InternalID), uint(user.InternalID))
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
	return s.boardRepo.RemoveMembers(uint(board.InternalID), membersToRemove)
}

func (s *boardService) IsMember(boardPublicID, userPublicID string) (bool, error) {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return false, errors.New("board not found")
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return false, errors.New("user not found")
	}

	return s.boardMemberRepo.IsMember(string(board.InternalID), string(user.InternalID))
}

func (s *boardService) GetAllByUserPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	return s.boardRepo.FindAllByUserPaginate(userPublicID, filter, sort, limit, offset)
}

func (s *boardService) GetAllPaginate(filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	return s.boardRepo.FindAllPaginate(filter, sort, limit, offset)
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