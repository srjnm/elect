package database

import (
	"crypto/sha256"
	"elect/dto"
	"elect/mappers"
	"elect/models"
	"elect/roles"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type postgresDatabase struct {
	connection *gorm.DB
}

func NewPostgresDatabase() (Database, *http.ServeMux) {
	source := os.Getenv("DATABASE_URL")
	db, err := gorm.Open("postgres", source)
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.User{})

	count := 0
	if db.Model(models.User{}).Where("email = ?", os.Getenv("ADMIN_EMAIL")).Count(&count); count == 0 {
		hashedPassword, err := HashPassword(os.Getenv("ADMIN_PASSWORD"))
		if err != nil {
			panic("Failed to initialize database!")
		}

		t := sha256.Sum256([]byte(strconv.Itoa(rand.Int()) + os.Getenv("ADMIN_EMAIL") + strconv.Itoa(rand.Int())))
		token := t[:]

		ret := db.Create(&models.User{
			FirstName:   "Suraj",
			LastName:    "N M",
			Email:       os.Getenv("ADMIN_EMAIL"),
			Password:    hashedPassword,
			Role:        roles.SuperAdmin,
			VerifyToken: hex.EncodeToString(token),
			Verified:    true})
		if ret.Error != nil {
			panic(ret.Error.Error())
		}
	}

	mux := SetUpQORAdmin(db)

	return &postgresDatabase{
		connection: db,
	}, mux
}

func (db *postgresDatabase) FindUserForAuth(email string) (dto.AuthUserDTO, error) {
	var count int
	user := models.User{}
	if db.connection.Model(&user).Where("email = ?", email).Count(&count); count == 0 {
		return dto.AuthUserDTO{}, errors.New("Invalid user!")
	}

	if db.connection.Where("email = ?", email).Find(&user).RowsAffected == 0 {
		return dto.AuthUserDTO{}, errors.New("Invalid user!")
	}

	if user.Verified == false {
		return dto.AuthUserDTO{}, errors.New("Account not verified!")
	}

	if user.Password == "" {
		return dto.AuthUserDTO{}, errors.New("Account not verified!")
	}

	return mappers.ToAuthUserDTO(user), nil
}

func (db *postgresDatabase) FindUserByID(userID string) (dto.GeneralUserDTO, error) {
	var count int
	user := models.User{}
	if db.connection.Model(&user).Where("user_id = ?", userID).Count(&count); count == 0 {
		return dto.GeneralUserDTO{}, errors.New("Invalid user!")
	}

	if db.connection.Where("user_id = ?", userID).Find(&user).RowsAffected == 0 {
		return dto.GeneralUserDTO{}, errors.New("Invalid user!")
	}

	return mappers.ToGeneralUserDTO(user), nil
}

func (db *postgresDatabase) GetUserRole(email string) (int, error) {
	var count int
	user := models.User{}
	if db.connection.Model(&user).Where("email = ?", email).Count(&count); count == 0 {
		return -1, errors.New("Invalid user!")
	}

	if db.connection.Where("email = ?", email).Find(&user).RowsAffected == 0 {
		return -1, errors.New("Invalid user!")
	}

	return user.Role, nil
}

func (db *postgresDatabase) VerifyAndSetPassword(setPasswordDTO dto.SetPasswordDTO) error {
	var user models.User
	res := db.connection.Model(&models.User{}).Where("verify_token = ?", setPasswordDTO.Token).Find(&user)

	fmt.Println(".")

	if res.RowsAffected == 0 {
		return errors.New("Invalid!")
	}

	if user.Verified == true {
		return errors.New("Already Verified!")
	}

	hashedPassword, err := HashPassword(setPasswordDTO.Password)
	if err != nil {
		return err
	}

	err = db.connection.Model(&models.User{}).Where("verify_token = ?", setPasswordDTO.Token).Update("password", hashedPassword).Error
	if err != nil {
		return err
	}

	fmt.Println("...")

	err = db.connection.Model(&models.User{}).Where("verify_token = ?", setPasswordDTO.Token).Update("verified", true).Error
	if err != nil {
		return err
	}

	fmt.Println(".....")

	return nil
}

func (db *postgresDatabase) TokenValidity(token string) error {
	var user models.User
	res := db.connection.Model(&models.User{}).Where("verify_token = ?", token).Find(&user)

	if res.RowsAffected == 0 {
		return errors.New("Invalid!")
	}
	if user.Verified {
		return errors.New("Already verified!")
	}

	return nil
}

//Bcrypt Functions
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
