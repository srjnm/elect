package models

import (
	"time"

	"github.com/jinzhu/gorm"
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
	RegisteredBy       string    `gorm:"default:null"`
	Verified           bool      `gorm:"default:false"`
	ActiveRefreshToken string    `gorm:"default:null"`
	Base
}

type Election struct {
	ElectionID     uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	Title          string    `gorm:"not null"`
	StartingAt     time.Time `gorm:"not null"`
	EndingAt       time.Time `gorm:"not null"`
	LockingAt      time.Time `gorm:"not null"`
	GenderSpecific bool      `gorm:"not null; default:false"`
	CreatedBy      string    `gorm:"not null"`
	Base
}

type Participant struct {
	ParticipantID uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	User          User      `gorm:"foreignKey: UserID; constraint:OnDelete:CASCADE;"`
	UserID        uuid.UUID `gorm:"uniqueIndex:idx_user_election"`
	Election      Election  `gorm:"foreignKey: ElectionID; constraint:OnDelete:CASCADE;"`
	ElectionID    uuid.UUID `gorm:"uniqueIndex:idx_user_election"`
	Voted         bool      `gorm:"not null; default: false"`
	Base
}

type Blacklist struct {
	BlacklistID uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	User        User      `gorm:"foreignKey: UserID; constraint:OnDelete:CASCADE;"`
	UserID      uuid.UUID `gorm:"unique"`
}

type Candidate struct {
	CandidateID    uuid.UUID `gorm:"primary_key; type:uuid; default:uuid_generate_v4()"`
	User           User      `gorm:"foreignKey: UserID; constraint:OnDelete:CASCADE;"`
	UserID         uuid.UUID `gorm:"uniqueIndex:idx_user_election"`
	Election       Election  `gorm:"foreignKey: ElectionID; constraint:OnDelete:CASCADE; unique"`
	ElectionID     uuid.UUID `gorm:"uniqueIndex:idx_user_election"`
	Sex            int       `gorm:"not null"`
	DisplayPicture string    `gorm:"not null"`
	Poster         string    `gorm:"not null"`
	IDProof        string    `gorm:"not null"`
	Approved       bool      `gorm:"not null; default: false"`
	Votes          int       `gorm:"not null; default: 0"`
}

type ResetToken struct {
	gorm.Model
	Email     string    `validate:"email,optional" gorm:"not null; type: varchar(384)"`
	Token     string    `gorm:"not null; unique"`
	ExpiresAt time.Time `gorm:"not null"`
}
