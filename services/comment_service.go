package services

import (
	"errors"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/repositories"
)

type CommentService interface {
	GetByCardPublicID(cardPublicID string) ([]models.Comment, error)
	GetByID(id uint) (*models.Comment, error)
	GetByPublicID(publicID string) (*models.Comment, error)
	Create(comment *models.Comment) error
	Update(comment *models.Comment) error
	Delete(id uint) error
}

type commentService struct {
	commentRepo repositories.CommentRepository
	cardRepo	repositories.CardRepository
	userRepo	repositories.UserRepository
}

func NewCommentService(
		commentRepo repositories.CommentRepository,
		cardRepo repositories.CardRepository,
		userRepo repositories.UserRepository,
	) CommentService {
		return &commentService{
			commentRepo: commentRepo,
			cardRepo: cardRepo,
			userRepo: userRepo,
		}
}

func (s *commentService) GetByCardPublicID(cardPublicID string) ([]models.Comment, error) {
	// Verifikasi card ada
	_, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return nil, errors.New("card not found")
	}

	return s.commentRepo.FindByCardPublicID(cardPublicID)
}

func (s *commentService) GetByID(id uint) (*models.Comment, error) {
	return s.commentRepo.FindByID(id)
}

func (s *commentService) GetByPublicID(id string) (*models.Comment, error) {
	return s.commentRepo.FindByPublicID(id)
}

func (s *commentService) Create(comment *models.Comment) error {
	// Pastikan card dan user ada
	_, err := s.cardRepo.FindByPublicID(comment.CardPubID.String())
	if err != nil {
		return errors.New("card not found")
	}

	_, err = s.userRepo.FindByPublicID(comment.UserPubID.String())
	if err != nil {
		return errors.New("user not found")
	}

	return s.commentRepo.Create(comment)
}

func (s *commentService) Update(comment *models.Comment) error {
	return s.commentRepo.Update(comment)
}

func (s *commentService) Delete(id uint) error {
	return s.commentRepo.Delete(id)
}