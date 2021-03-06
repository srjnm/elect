package models

import (
	"log"
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

func (user *User) AfterDelete(db *gorm.DB) error {
	err := db.Model(&Participant{}).Where("user_id = ?", user.UserID.String()).Delete(&Participant{}).Error
	if err != nil {
		log.Println("gorm:")
		log.Println(err)
		return err
	}

	err = db.Model(&Candidate{}).Where("user_id = ?", user.UserID.String()).Delete(&Candidate{}).Error
	if err != nil {
		log.Println("gorm:")
		log.Println(err)
		return err
	}

	err = db.Model(&Blacklist{}).Where("user_id = ?", user.UserID.String()).Delete(&Blacklist{}).Error
	if err != nil {
		log.Println("gorm:")
		log.Println(err)
		return err
	}

	return nil
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

func (election *Election) AfterDelete(db *gorm.DB) error {
	err := db.Model(&Participant{}).Where("election_id = ?", election.ElectionID.String()).Delete(&Participant{}).Error
	if err != nil {
		log.Println("gorm:")
		log.Println(err)
		return err
	}

	err = db.Model(&Candidate{}).Where("election_id = ?", election.ElectionID.String()).Delete(&Candidate{}).Error
	if err != nil {
		log.Println("gorm:")
		log.Println(err)
		return err
	}

	return nil
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
	gorm.Model
	User   User      `gorm:"foreignKey: UserID; constraint:OnDelete:CASCADE;"`
	UserID uuid.UUID `gorm:"unique"`
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
	Base
}

type ResetToken struct {
	gorm.Model
	Email     string    `validate:"email,optional" gorm:"not null; type: varchar(384)"`
	Token     string    `gorm:"not null; unique"`
	ExpiresAt time.Time `gorm:"not null"`
}
