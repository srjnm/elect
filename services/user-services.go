package services

import (
	"crypto/sha256"
	"elect/database"
	"elect/dto"
	"elect/email"
	"elect/mappers"
	"elect/models"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
)

type UserService interface {
	GetUserForAuth(email string) (dto.AuthUserDTO, error)
	GetUserByID(userID string) (dto.GeneralUserDTO, error)
	GetUserRole(email string) (int, error)
	VerifyUserAndSetPassword(setPasswordDTO dto.SetPasswordDTO) error
	TokenValidity(token string) error
	RegisterStudent(registerStudentDTO dto.RegisterStudentDTO) error
	RegisteredStudents(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralStudentDTO, error)
	DeleteRegisteredStudent(userId string, studentUserId string) error
	SetActiveRefreshToken(token string, email string) error
	CheckIfActiveRefreshToken(token string, email string) error
	ClearActiveRefreshToken(email string) error
	ChangePassword(userId string, changePasswordDTO dto.ChangePasswordDTO) error
	CheckVerifyTokenValidity(token string) error
	CheckResetTokenValidity(token string) error
	GenerateResetToken(createResetTokenDTO dto.CreateResetTokenDTO) error
	ResetPassword(resetPasswordDTO dto.ResetPasswordDTO) error
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
	user, err := service.database.GetUser(userID)
	if err != nil {
		return dto.GeneralUserDTO{}, err
	}

	return mappers.ToGeneralUserDTO(user), nil
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

func (service *userService) SetActiveRefreshToken(token string, email string) error {
	err := service.database.StoreActiveRefreshToken(token, email)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) CheckIfActiveRefreshToken(token string, email string) error {
	activeToken, err := service.database.GetActiveRefreshToken(email)
	if err != nil {
		return err
	}

	if activeToken != token {
		return errors.New("Logged in other device!")
	}

	return nil
}

func (service *userService) ClearActiveRefreshToken(email string) error {
	err := service.database.ClearActiveRefreshToken(email)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) RegisterStudent(registerStudentDTO dto.RegisterStudentDTO) error {
	user := mappers.ToUserFromRegisterStudentDTO(registerStudentDTO)

	t := sha256.Sum256([]byte(strconv.Itoa(rand.Int()) + user.Email + strconv.Itoa(rand.Int())))
	token := t[:]
	user.VerifyToken = hex.EncodeToString(token)

	err := service.database.RegisterStudent(user)
	if err != nil {
		return err
	}

	err = email.SendVerificationEmail(user.FirstName, user.Email, user.VerifyToken, "template.html")
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) RegisteredStudents(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralStudentDTO, error) {
	var regStudents []models.User
	var err error
	if paginatorParams.Page == "" {
		regStudents, err = service.database.RegisteredStudents(userId, dto.PaginatorParams{})
		if err != nil {
			return nil, err
		}
	} else {
		regStudents, err = service.database.RegisteredStudents(userId, paginatorParams)
		if err != nil {
			return nil, err
		}
	}

	var generalStudentsDTO []dto.GeneralStudentDTO
	for _, student := range regStudents {
		generalStudentsDTO = append(generalStudentsDTO, mappers.ToGeneralStudentDTOFromUser(student))
	}

	return generalStudentsDTO, nil
}

func (service *userService) DeleteRegisteredStudent(userId string, studentUserId string) error {
	return service.database.DeleteRegisteredStudent(userId, studentUserId)
}

func (service *userService) ChangePassword(userId string, changePasswordDTO dto.ChangePasswordDTO) error {
	return service.database.ChangePassword(userId, changePasswordDTO)
}

func (service *userService) CheckVerifyTokenValidity(token string) error {
	return service.database.TokenValidity(token)
}

func (service *userService) CheckResetTokenValidity(token string) error {
	return service.database.CheckResetTokenValidity(token)
}

func (service *userService) GenerateResetToken(createResetTokenDTO dto.CreateResetTokenDTO) error {
	name, token, err := service.database.GenerateResetToken(createResetTokenDTO.Email)
	if err != nil {
		return err
	}

	err = email.SendResetPasswordEmail(name, createResetTokenDTO.Email, token, "reset.html")
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) ResetPassword(resetPasswordDTO dto.ResetPasswordDTO) error {
	return service.database.ResetPassword(resetPasswordDTO)
}
