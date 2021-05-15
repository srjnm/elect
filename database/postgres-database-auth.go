package database

import (
	"elect/dto"
	"elect/mappers"
	"elect/models"
	"errors"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

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

	err = db.connection.Model(&models.User{}).Where("verify_token = ?", setPasswordDTO.Token).Update("verified", true).Error
	if err != nil {
		return err
	}

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

func (db *postgresDatabase) StoreActiveRefreshToken(token string, email string) error {
	var user models.User
	res := db.connection.Model(&user).Where("email = ?", email).Update("active_refresh_token", token)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) GetActiveRefreshToken(email string) (string, error) {
	var user models.User
	res := db.connection.Model(&models.User{}).Where("email = ?", email).Find(&user)
	if res.Error != nil {
		return "", res.Error
	}

	return user.ActiveRefreshToken, nil
}

func (db *postgresDatabase) ClearActiveRefreshToken(email string) error {
	var user models.User

	res := db.connection.Model(&user).Where("email = ?", email).Update("active_refresh_token", "")
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) ChangePassword(userId string, changePasswordDTO dto.ChangePasswordDTO) error {
	var user models.User
	res := db.connection.Model(&models.User{}).Where("user_id = ?", userId).Find(&user)
	if res.Error != nil {
		return res.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePasswordDTO.CurrentPassword))
	if err != nil {
		return errors.New("Invalid Current Password!")
	}

	newPassword, err := HashPassword(changePasswordDTO.NewPassword)
	if err != nil {
		return errors.New("Failed to hash password!")
	}

	user.Password = newPassword

	res = db.connection.Model(&models.User{}).Where("user_id = ?", userId).Update(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) GenerateResetToken(email string) (string, string, error) {
	var count int
	res := db.connection.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if res.Error != nil {
		return "", "", res.Error
	}
	if count == 0 {
		return "", "", errors.New("Invalid email!")
	}

	var user models.User
	res = db.connection.Model(&models.User{}).Where("email = ?", email).Find(&user)
	if res.Error != nil {
		return "", "", res.Error
	}
	if !user.Verified {
		return "", "", errors.New("Account not verified yet!")
	}

	token := uuid.NewV4().String()
	expiresAt := time.Now().Add(time.Minute * 30).UTC()

	res = db.connection.Model(&models.ResetToken{}).Create(&models.ResetToken{Email: email, Token: token, ExpiresAt: expiresAt})
	if res.Error != nil {
		return "", "", res.Error
	}

	return user.FirstName, token, nil
}

func (db *postgresDatabase) CheckResetTokenValidity(token string) error {
	var count int
	res := db.connection.Model(&models.ResetToken{}).Where("token = ?", token).Count(&count)
	if res.Error != nil {
		return res.Error
	}
	if count == 0 {
		return errors.New("Invalid Reset Token!")
	}

	var resetToken models.ResetToken
	res = db.connection.Model(&models.ResetToken{}).Where("token = ?", token).Find(&resetToken)
	if res.Error != nil {
		return res.Error
	}

	if time.Now().UTC().After(resetToken.ExpiresAt) {
		return errors.New("Reset Token Expired!")
	}

	return nil
}

func (db *postgresDatabase) ResetPassword(resetPasswordDTO dto.ResetPasswordDTO) error {
	var count int
	res := db.connection.Model(&models.ResetToken{}).Where("token = ?", resetPasswordDTO.Token).Count(&count)
	if res.Error != nil {
		return res.Error
	}
	if count == 0 {
		return errors.New("Invalid Reset Token!")
	}

	var resetToken models.ResetToken
	res = db.connection.Model(&models.ResetToken{}).Where("token = ?", resetPasswordDTO.Token).Find(&resetToken)
	if res.Error != nil {
		return res.Error
	}

	if time.Now().UTC().After(resetToken.ExpiresAt) {
		return errors.New("Reset Token Expired!")
	}

	var user models.User
	res = db.connection.Model(&models.User{}).Where("email = ?", resetToken.Email).Find(&user)
	if res.Error != nil {
		return res.Error
	}

	newPassword, err := HashPassword(resetPasswordDTO.NewPassword)
	if err != nil {
		return errors.New("Failed to hash password!")
	}

	user.Password = newPassword

	res = db.connection.Model(&models.User{}).Update(&user)
	if res.Error != nil {
		return res.Error
	}

	resetToken.ExpiresAt = time.Now().UTC()
	res = db.connection.Model(&models.ResetToken{}).Update(&resetToken)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

//Bcrypt Functions
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
