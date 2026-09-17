package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/MaulanaBarzaqi/project-management/config"
	"github.com/MaulanaBarzaqi/project-management/models"
	"github.com/MaulanaBarzaqi/project-management/models/types"
	"github.com/MaulanaBarzaqi/project-management/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardService interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error
	GetByListPublicID(listPublicID string) ([]models.Card, error)
	GetByID(id uint) (*models.Card, error)
	GetByPublicID(publicID string) (*models.Card, error)
	AddLabel(cardPublicID, labelPublicID string) error
	RemoveLabel(cardPublicID, labelPublicID string) error
	AddAssignees(cardPublicID string, userPublicIDs []string) error
	RemoveAssignee(cardPublicID, userPublicID string) error
	UpdatePositions(listPublicID string, positions []string) error
	UpdateCardList(card *models.Card, listPublicID string) error
	GetCardDetail(cardPublicID string) (*models.Card, error)
}

type cardService struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepository
	userRepo repositories.UserRepository
	labelRepo repositories.LabelRepository
}

func NewCardService(
	cardRepo repositories.CardRepository,
	listRepo repositories.ListRepository,
	userRepo repositories.UserRepository,
	labelRepo repositories.LabelRepository,
) CardService {
	return &cardService{
		cardRepo: cardRepo,
		listRepo: listRepo,
		userRepo: userRepo,
		labelRepo: labelRepo,
	}
}

func (s *cardService) Create(card *models.Card, listPublicID string) error {
	// ambil list nya dari listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}
	// set list_internal_id
	card.ListID = list.InternalID
	// generate public_id jika belum ada
	if card.PublicID == uuid.Nil {
		card.PublicID = uuid.New()
	}
	card.CreatedAt = time.Now()
	// transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	// simpan card
	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create card: %w", err)
	}
	// update atau buat card_position
	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).
			  Where("list_internal_id = ?",list.InternalID).
			  First(&position).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// buat baru jika belum ada
					position = models.CardPosition{
						PublicID: uuid.New(),
						ListID: list.InternalID,
						CardOrder: types.UUIDArray{card.PublicID},
					}
					if err := tx.Create(&position).Error; err != nil {
						tx.Rollback()
						return fmt.Errorf("failed to create card position: %w", err)
					}
				} else {
					tx.Rollback()
					return fmt.Errorf("failed to get card position: %w", err)
				}
			  } else {
				// tambahkan card baru ke urutan
				position.CardOrder = append(position.CardOrder, card.PublicID)
				if err := tx.Model(&models.CardPosition{}).
						  Where("internal_id = ?", position.InternalID).
						  Update("card_order", position.CardOrder).Error; err != nil {
							tx.Rollback()
							return fmt.Errorf("failed to update card position: %w", err)
						  }
			  }
			//   commit transaksi
			if err := tx.Commit().Error; err != nil {
				return fmt.Errorf("transaction commit failed: %w", err)
			}
			return nil
}

func (s *cardService) Update(card *models.Card, listPublicID string) error {
    // Ambil card dulu
    existingCard, err := s.cardRepo.FindByPublicID(card.PublicID.String())
    if err != nil {
        return fmt.Errorf("card not found: %w", err)
    }
    // Ambil list baru
    newList, err := s.listRepo.FindByPublicID(listPublicID)
    if err != nil {
        return fmt.Errorf("list not found: %w", err)
    }
    // Mulai transaksi
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()
    // Jika pindah list -> hapus dari posisi lamanya & tambah ke list baru
    if existingCard.ListID != newList.InternalID {
        // Hapus dari list lama
        var oldPos models.CardPosition
        errOld := tx.Where("list_internal_id = ?", existingCard.ListID).First(&oldPos).Error 
        if errOld == nil {
            filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
            for _, id := range oldPos.CardOrder {
                if id != existingCard.PublicID {
                    filtered = append(filtered, id)
                }
            }
            // Update posisi list lama
            if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?", oldPos.InternalID).
                Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
                tx.Rollback()
                return fmt.Errorf("failed to update old card position: %w", err)
            }
        } else if !errors.Is(errOld, gorm.ErrRecordNotFound) {
            tx.Rollback()
            return fmt.Errorf("failed to get old card position: %w", errOld)
        }
        // Tambah ke list baru
        var newPos models.CardPosition
        res := tx.Where("list_internal_id = ?", newList.InternalID).First(&newPos)
        if errors.Is(res.Error, gorm.ErrRecordNotFound) {
            newPos = models.CardPosition{
                PublicID:  uuid.New(),
                ListID:    newList.InternalID,
                CardOrder: types.UUIDArray{existingCard.PublicID},
            }
            if err := tx.Create(&newPos).Error; err != nil {
                tx.Rollback()
                return fmt.Errorf("failed to create card position for new list: %w", err)
            }
        } else if res.Error == nil {
            updateOrder := append(newPos.CardOrder, existingCard.PublicID)
            if err := tx.Model(&models.CardPosition{}).
                Where("internal_id = ?", newPos.InternalID).
                Update("card_order", types.UUIDArray(updateOrder)).Error; err != nil {
                tx.Rollback()
                return fmt.Errorf("failed to update new card position: %w", err)
            }
        } else {
            tx.Rollback()
            return fmt.Errorf("failed to get new card position: %w", res.Error)
        }
    }

    // Update data card (Pastikan ListID ikut diperbarui ke newList.InternalID)
    card.InternalID = existingCard.InternalID
    card.PublicID = existingCard.PublicID
    card.ListID = newList.InternalID

    if err := tx.Save(card).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update card: %w", err)
    }
    // Commit transaksi
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("transaction commit failed: %w", err)
    }
    return nil
}

func (s *cardService) Delete(id uint) error {
	return s.cardRepo.Delete(id)
}

func (s *cardService) GetByListPublicID(listPublicID string) ([]models.Card, error) {
    // Verifikasi kalau list-nya ada 
    list, err := s.listRepo.FindByPublicID(listPublicID)
    if err != nil {
        return nil, fmt.Errorf("list not found: %w", err)
    }
    // Ambil semua card di list tersebut terlebih dahulu
    cards, err := s.cardRepo.FindByListPublicID(listPublicID)
    if err != nil {
        return nil, fmt.Errorf("failed to get cards: %w", err)
    }
    // Ambil card position-nya 
    position, err := s.cardRepo.FindByListID(list.InternalID)
    if err != nil {
        if !errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("failed to get card position: %w", err)
        }
    }
    // Sorting jika position-nya ada dan tidak kosong
    if position != nil && len(position.CardOrder) > 0 {
        cards = sortCardByPosition(cards, position.CardOrder)
    }
    return cards, nil
}

func (s *cardService) GetByID(id uint) (*models.Card, error) {
	return s.cardRepo.FindByID(id)
}

func (s *cardService) GetByPublicID(publicID string) (*models.Card, error) {
	return s.cardRepo.FindByPublicID(publicID)
}

func (s *cardService) AddLabel(cardPublicID, labelPublicID string) error {
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return errors.New("card not found")
	}
	label, err := s.labelRepo.FindByPublicID(labelPublicID)
	if err != nil {
		return errors.New("label not found")
	}

	return s.cardRepo.AddLabel(uint(card.InternalID), uint(label.InternalID))
}

func (s *cardService) RemoveLabel(cardPublicID, labelPublicID string) error {
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return errors.New("card not found")
	}
	label, err := s.labelRepo.FindByPublicID(labelPublicID)
	if err != nil {
		return errors.New("label not found")
	}

	return s.cardRepo.RemoveLabel(uint(card.InternalID), uint(label.InternalID))
}

func (s *cardService) AddAssignees(cardPublicID string, userPublicIDs []string) error {
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return errors.New("card not found")
	}

	var userIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return fmt.Errorf("user %s not found", userPublicID)
		}
		userIDs = append(userIDs, uint(user.InternalID))
	}
	return s.cardRepo.AddAssignees(uint(card.InternalID), userIDs)
}

func (s *cardService) RemoveAssignee(cardPublicID, userPublicID string) error {
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return errors.New("card not found")
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.cardRepo.RemoveAssignee(uint(card.InternalID), uint(user.InternalID))
}

func (s *cardService) UpdatePositions(listPublicID string, positions []string) error {
	// Pastikan list ada
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	// Convert []string → UUIDArray
	var cardOrder types.UUIDArray
	for _, p := range positions {
		id, err := uuid.Parse(p)
		if err != nil {
			return fmt.Errorf("invalid card id: %w", err)
		}
		cardOrder = append(cardOrder, id)
	}

	// Mulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Update card_positions
	var position models.CardPosition
	res := tx.Where("list_internal_id = ?", list.InternalID).First(&position)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// Kalau belum ada, bikin baru
		position = models.CardPosition{
			PublicID:  uuid.New(),
			ListID:    list.InternalID,
			CardOrder: cardOrder,
		}
		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create card position: %w", err)
		}
	} else if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get card position: %w", res.Error)
	} else {
		// Kalau sudah ada, update array-nya
		if err := tx.Model(&position).
			Where("internal_id = ?", position.InternalID).
			Update("card_order", cardOrder).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update card positions: %w", err)
		}
	}

	// Commit
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (s *cardService) UpdateCardList(card *models.Card, listPublicID string) error {
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}
	card.ListID = list.InternalID
	return nil
}

func (s *cardService) GetCardDetail(cardPublicID string) (*models.Card, error) {
	// Ambil card dari repository, preload assignees dan attachments
	card, err := s.cardRepo.FindByPublicID(cardPublicID)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func sortCardByPosition(cards []models.Card, order []uuid.UUID) []models.Card {
	// buat map untuk pencarian cepat
	orderMap := make(map[uuid.UUID]int)
	for i, id := range order {
		orderMap[id] = i
	}
	defaultIndex := len(order)
	// sorting slice
	sort.SliceStable(cards, func (i,j int) bool {
		idxI, okI := orderMap[cards[i].PublicID]
		if !okI {
			idxI = defaultIndex
		}
		idxJ, okJ :=orderMap[cards[j].PublicID]
		if !okJ {
			idxJ = defaultIndex
		}
		if idxI == idxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}
		return idxI < idxJ
	})
	return cards
}