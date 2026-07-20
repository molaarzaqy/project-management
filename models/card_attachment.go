package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardAttachment struct {
	InternalID 	int64		`json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement;column:internal_id"`
	PublicID   	uuid.UUID	`json:"public_id" db:"public_id" gorm:"type:uuid;column:public_id"`
	CardID 		int64		`json:"card_internal_id" db:"card_internal_id" gorm:"column:card_internal_id"`
	UserID 		int64		`json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id"`
	File 		string		`json:"file" db:"file" gorm:"column:file"`
	CreatedAt 	time.Time	`json:"created_at" db:"created_at" gorm:"column:created_at"`

	FileURL 	string 		`json:"file_url" gorm:"-"`
}

func(CardAttachment) TableName() string {
	return "card_attachment"
}

func (a *CardAttachment) BeforeCreate(tx *gorm.DB) (err error) {
	if a.PublicID == uuid.Nil {
		a.PublicID = uuid.New()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	return nil;
}