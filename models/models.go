package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type User struct {
	UserID             uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	FirstName          string    `gorm:"not null; type: varchar(64)"`
	LastName           string    `gorm:"not null; type: varchar(64)"`
	RegNumber          string    `gorm:"type: varchar(12); default:null; unique"`
	Email              string    `validate:"email,optional" gorm:"not null; unique; type: varchar(384)"`
	Password           string    `gorm:"type: varchar(64); default:null"`
	Role               int       `gorm:"not null;"`
	VerifyToken        string    `gorm:"type: varchar(64); default:null"`
	Verified           bool      `gorm:"default:false"`
	ActiveRefreshToken string    `gorm:"default:null"`
	Base
}
