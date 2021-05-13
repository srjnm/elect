package database

import (
	"crypto/sha256"
	"elect/dto"
	"elect/email"
	"elect/models"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/validations"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Database interface {
	// Auth
	FindUserForAuth(email string) (dto.AuthUserDTO, error)
	FindUserByID(userID string) (dto.GeneralUserDTO, error)
	GetUserRole(email string) (int, error)
	VerifyAndSetPassword(setPasswordDTO dto.SetPasswordDTO) error
	TokenValidity(token string) error
	StoreActiveRefreshToken(token string, email string) error
	GetActiveRefreshToken(email string) (string, error)
	ClearActiveRefreshToken(email string) error
	ChangePassword(userId string, changePasswordDTO dto.ChangePasswordDTO) error
	GenerateResetToken(email string) (string, string, error)
	CheckResetTokenValidity(token string) error
	ResetPassword(resetPasswordDTO dto.ResetPasswordDTO) error

	// Users
	RegisterStudent(user models.User) error
	RegisteredStudents(userId string, paginatorParams dto.PaginatorParams) ([]models.User, error)
	DeleteRegisteredStudent(userId string, studentUserId string) error

	// Election
	CreateElection(election models.Election) error
	EditElection(userId string, election models.Election) error
	DeleteElection(userId string, electionId string) error
	AddParticipant(userId string, electId string, regno string) error
	DeleteParticipant(userId string, electionId string, participantId string) error
	GetElectionParticipants(userId string, electionId string) ([]dto.GeneralParticipantDTO, error)
	GetElectionsForAdmins(userId string, paginatorParams dto.PaginatorParams) ([]models.Election, error)
	GetElectionsForStudents(userId string, paginatorParams dto.PaginatorParams) ([]models.Election, error)
	EnrollCandidate(candidate models.Candidate) error
	CheckCandidateEligibility(userId string, electionId string) error
	ApproveCandidate(userId string, candidateId string) error
	UnapproveCandidate(userId string, candidateId string) error
	GetElectionForAdmins(userId string, electionId string) (models.Election, []dto.GeneralParticipantDTO, []models.Candidate, error)
	GetElectionForStudents(userId string, electionId string) (models.Election, []models.Candidate, error)
	CastVote(userId string, electionId string, candidateId string) error
	GetResults(userId string, role int, electionId string) (models.Election, []models.Candidate, error)
}

func SetUpQORAdmin(db *gorm.DB) *http.ServeMux {
	adm := admin.New(&admin.AdminConfig{SiteName: "ELECT", DB: db})
	mux := http.NewServeMux()
	adm.MountTo("/admin", mux)

	// User Management
	usr := adm.AddResource(models.User{}, &admin.Config{Menu: []string{"User Management"}})
	usr.IndexAttrs("-Password", "-VerifyToken", "-ActiveRefreshToken")
	usr.NewAttrs("-Password", "-ActiveRefreshToken", "-RegisteredBy")
	usr.EditAttrs("-VerifyToken", "-ActiveRefreshToken", "-RegisteredBy")
	usr.Meta(&admin.Meta{
		Name: "Password",
		Type: "password",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			values := metaValue.Value.([]string)
			if len(values) > 0 {
				if np := values[0]; np != "" {
					pwd, err := bcrypt.GenerateFromPassword([]byte(np), 14)
					if err != nil {
						context.DB.AddError(validations.NewError(usr, "Password", "Can't encrypt password"))
						return
					}
					u := resource.(*models.User)
					u.Password = string(pwd)
				}
			}
		},
	})
	usr.Meta(&admin.Meta{
		Name: "VerifyToken",
		Type: "hidden",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			u := resource.(*models.User)

			t := sha256.Sum256([]byte(strconv.Itoa(rand.Int()) + u.Email + strconv.Itoa(rand.Int())))
			token := t[:]
			u.VerifyToken = hex.EncodeToString(token)

			if !u.Verified {
				email.SendVerificationEmail(u.FirstName, u.Email, u.VerifyToken, "template.html")
			}
		},
	})

	blacklist := adm.AddResource(models.Blacklist{}, &admin.Config{Menu: []string{"User Management"}})
	blacklist.IndexAttrs("-User")
	blacklist.NewAttrs("-User")
	blacklist.EditAttrs("-User")
	blacklist.Meta(&admin.Meta{
		Name: "UserID",
		Type: "string",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			values := metaValue.Value.([]string)
			if len(values) > 0 {
				if id := values[0]; id != "" {
					b := resource.(*models.Blacklist)
					b.UserID = uuid.FromStringOrNil(id)
				}
			}
		},
	})

	// Election Management
	elect := adm.AddResource(models.Election{}, &admin.Config{Menu: []string{"Election Management"}, IconName: "Election"})
	elect.EditAttrs("-CreatedBy")
	elect.Meta(&admin.Meta{
		Name: "StartingAt",
		Type: "datetime",
		Valuer: func(record interface{}, context *qor.Context) interface{} {
			e := record.(*models.Election)
			tZone, err := time.LoadLocation("Asia/Kolkata")
			if err != nil {
				log.Println("IST TimeZone Error!")
			}
			e.StartingAt = e.StartingAt.In(tZone)
			return e.StartingAt
		},
		Setter: func(record interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			e := record.(*models.Election)
			e.StartingAt = e.StartingAt.UTC()
		},
	})
	elect.Meta(&admin.Meta{
		Name: "EndingAt",
		Type: "datetime",
		Valuer: func(record interface{}, context *qor.Context) interface{} {
			e := record.(*models.Election)
			tZone, err := time.LoadLocation("Asia/Kolkata")
			if err != nil {
				log.Println("IST TimeZone Error!")
			}
			e.EndingAt = e.EndingAt.In(tZone)
			return e.EndingAt
		},
		Setter: func(record interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			e := record.(*models.Election)
			e.EndingAt = e.EndingAt.UTC()
		},
	})
	elect.Meta(&admin.Meta{
		Name: "LockingAt",
		Type: "datetime",
		Valuer: func(record interface{}, context *qor.Context) interface{} {
			e := record.(*models.Election)
			tZone, err := time.LoadLocation("Asia/Kolkata")
			if err != nil {
				log.Println("IST TimeZone Error!")
			}
			e.LockingAt = e.LockingAt.In(tZone)
			return e.LockingAt
		},
		Setter: func(record interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			e := record.(*models.Election)
			e.LockingAt = e.LockingAt.UTC()
		},
	})

	part := adm.AddResource(models.Participant{}, &admin.Config{Menu: []string{"Election Management"}, IconName: "Election"})
	part.IndexAttrs("-User", "-Election")
	part.NewAttrs("-User", "-Election", "-Voted")
	part.Meta(&admin.Meta{
		Name: "UserID",
		Type: "string",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			values := metaValue.Value.([]string)
			if len(values) > 0 {
				if id := values[0]; id != "" {
					p := resource.(*models.Participant)
					p.UserID = uuid.FromStringOrNil(id)
				}
			}
		},
	})
	part.Meta(&admin.Meta{
		Name: "ElectionID",
		Type: "string",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			values := metaValue.Value.([]string)
			if len(values) > 0 {
				if id := values[0]; id != "" {
					p := resource.(*models.Participant)
					p.ElectionID = uuid.FromStringOrNil(id)
				}
			}
		},
	})

	cand := adm.AddResource(models.Candidate{}, &admin.Config{Menu: []string{"Election Management"}, IconName: "Election"})
	cand.IndexAttrs("-User", "-Election")
	cand.NewAttrs("-User", "-Election", "-Votes")
	cand.EditAttrs("-User", "-Election", "-Votes")
	cand.Meta(&admin.Meta{
		Name: "UserID",
		Type: "string",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			values := metaValue.Value.([]string)
			if len(values) > 0 {
				if id := values[0]; id != "" {
					p := resource.(*models.Participant)
					p.UserID = uuid.FromStringOrNil(id)
				}
			}
		},
	})
	cand.Meta(&admin.Meta{
		Name: "ElectionID",
		Type: "string",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			values := metaValue.Value.([]string)
			if len(values) > 0 {
				if id := values[0]; id != "" {
					p := resource.(*models.Participant)
					p.ElectionID = uuid.FromStringOrNil(id)
				}
			}
		},
	})

	validations.RegisterCallbacks(db)

	return mux
}
