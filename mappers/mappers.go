package mappers

import (
	"elect/dto"
	"elect/models"
)

func ToAuthUserDTO(user models.User) dto.AuthUserDTO {
	return dto.AuthUserDTO{
		UserID:   user.UserID.String(),
		Email:    user.Email,
		Password: user.Password,
	}
}

func ToGeneralUserDTO(user models.User) dto.GeneralUserDTO {
	return dto.GeneralUserDTO{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}
}

func ToUserFromGeneralUserDTO(genUserDTO dto.GeneralUserDTO) models.User {
	return models.User{
		Email:     genUserDTO.Email,
		FirstName: genUserDTO.FirstName,
		LastName:  genUserDTO.LastName,
		Role:      genUserDTO.Role,
	}
}

func ToUserFromRegisterStudentDTO(registerStudentDTO dto.RegisterStudentDTO) models.User {
	return models.User{
		Email:        registerStudentDTO.Email,
		FirstName:    registerStudentDTO.FirstName,
		LastName:     registerStudentDTO.LastName,
		RegNumber:    registerStudentDTO.RegNumber,
		RegisteredBy: registerStudentDTO.RegisteredBy,
	}
}

func ToGeneralStudentDTOFromUser(user models.User) dto.GeneralStudentDTO {
	return dto.GeneralStudentDTO{
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		RegisterNumber: user.RegNumber,
	}
}
