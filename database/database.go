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
	FindUserForAuth(email string) (dto.AuthUserDTO, error)
	FindUserByID(userID string) (dto.GeneralUserDTO, error)
	GetUserRole(email string) (int, error)
	VerifyAndSetPassword(setPasswordDTO dto.SetPasswordDTO) error
	TokenValidity(token string) error
}

func SetUpQORAdmin(db *gorm.DB) *http.ServeMux {
	adm := admin.New(&admin.AdminConfig{SiteName: "Blobber", DB: db})
	mux := http.NewServeMux()
	adm.MountTo("/admin", mux)

	usr := adm.AddResource(models.User{}, &admin.Config{Menu: []string{"User Management"}})
	usr.IndexAttrs("-Password", "-VerifyToken")
	usr.EditAttrs("-VerifyToken", "-Email")
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
