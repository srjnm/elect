package database

import (
	"crypto/sha256"
	"elect/dto"
	"elect/email"
	"elect/models"
	"encoding/hex"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/validations"
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

	// Users
	RegisterStudent(user models.User) error
	RegisteredStudents(userId string, paginatorParams dto.PaginatorParams) ([]models.User, error)
	DeleteRegisteredStudent(userId string, studentUserId string) error
}

func SetUpQORAdmin(db *gorm.DB) *http.ServeMux {
	adm := admin.New(&admin.AdminConfig{SiteName: "ELECT", DB: db})
	mux := http.NewServeMux()
	adm.MountTo("/admin", mux)

	usr := adm.AddResource(models.User{}, &admin.Config{Menu: []string{"User Management"}})
	usr.IndexAttrs("-Password", "-VerifyToken", "-ActiveRefreshToken")
	usr.NewAttrs("-Password", "-ActiveRefreshToken", "-RegisteredBy")
	usr.EditAttrs("-VerifyToken", "-Email", "-ActiveRefreshToken", "-RegisteredBy")
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

	validations.RegisterCallbacks(db)

	return mux
}
