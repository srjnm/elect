package services

import (
	"elect/database"
	"elect/dto"
)

type UserService interface {
	GetUserForAuth(email string) (dto.AuthUserDTO, error)
	GetUserByID(userID string) (dto.GeneralUserDTO, error)
	GetUserRole(email string) (int, error)
	VerifyUserAndSetPassword(setPasswordDTO dto.SetPasswordDTO) error
	TokenValidity(token string) error
}

type userService struct {
	database database.Database
}

func NewUserService(database database.Database) UserService {
	return &userService{
		database: database,
	}
}

func (service *userService) GetUserForAuth(email string) (dto.AuthUserDTO, error) {
	return service.database.FindUserForAuth(email)
}

func (service *userService) GetUserByID(userID string) (dto.GeneralUserDTO, error) {
	return service.GetUserByID(userID)
}

func (service *userService) GetUserRole(email string) (int, error) {
	return service.database.GetUserRole(email)
}

func (service *userService) VerifyUserAndSetPassword(setPasswordDTO dto.SetPasswordDTO) error {
	return service.database.VerifyAndSetPassword(setPasswordDTO)
}

func (service *userService) TokenValidity(token string) error {
	return service.database.TokenValidity(token)
}
