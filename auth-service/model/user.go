package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	FirebaseUID string    `gorm:"type:varchar(255);unique_index"`
	Email       string    `gorm:"type:varchar(255);unique_index"`
	Name        string    `gorm:"type:varchar(255)"`
	Role        string    `gorm:"type:varchar(50);default:'user'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}