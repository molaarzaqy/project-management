package services

import (
	"errors"
	"time"

	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/repositories"
	"github.com/google/uuid"
)

type AttachmentService interface {
	GetByPublicID(pubID uuid.UUID) (*models.CardAttachment, error)
	Create(cardPublidID, userPublicID, fileName string) (*models.CardAttachment, error)
	DeleteByPublicID(pubID uuid.UUID) error
}

type attachmentService struct {
	attactmentRepo repositories.AttachmentRepository
	cardRepo repositories.CardRepository
	userRepo repositories.UserRepository
}

func NewAttachmentService(
	attactmentRepo repositories.AttachmentRepository,
	cardRepo repositories.CardRepository,
	userRepo repositories.UserRepository,
	) AttachmentService {
		return &attachmentService{
			attactmentRepo: attactmentRepo,
			cardRepo: cardRepo,
			userRepo: userRepo,
		}
	}

func (s *attachmentService) GetByPublicID(pubID uuid.UUID) (*models.CardAttachment, error) {
	return s.attactmentRepo.GetByPublicID(pubID)
}

func (s *attachmentService) Create(cardPublidID, userPublicID, fileName string) (*models.CardAttachment, error) {
	card, err := s.cardRepo.FindByPublicID(cardPublidID)
	if err != nil {
		return nil, errors.New("card not found")
	} 
	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	attach := &models.CardAttachment{
		PublicID: uuid.New(),
		CardID: card.InternalID,
		UserID: user.InternalID,
		File: fileName,
		CreatedAt: time.Now(),
	}
	if err := s.attactmentRepo.Create(attach); err != nil {
		return nil, err
	}
	return attach, nil
}

func (s *attachmentService) DeleteByPublicID(pubID uuid.UUID) error {
	return s.attactmentRepo.DeleteByPublicID(pubID)
}
