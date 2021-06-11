package services

import (
	"elect/database"
	"elect/dto"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(dto.AuthUserDTO, int) string
	GenerateRefreshToken(dto.AuthUserDTO) string
	GenerateNewTokens(tokenString string) (string, string, error)
	ValidateAccessToken(tokenString string) (*jwt.Token, error)
	ValidateRefreshToken(tokenString string) (*jwt.Token, error)
	GetUserIDAndRole(tokenString string) (string, int, error)
	GetRole(tokenString string) (int, error)
	GetEmail(tokenString string) (string, error)
	GenerateOTPToken(email string) string
	ValidateOTPToken(tokenString string) (*jwt.Token, error)
	GetEmailFromOTPToken(tokenString string) (string, error)
}

type jwtService struct {
	issuer   string
	database database.Database
}

func NewJWTService(issuer string, database database.Database) JWTService {
	return &jwtService{
		issuer:   issuer,
		database: database,
	}
}

func (service *jwtService) GenerateToken(authUserDTO dto.AuthUserDTO, role int) string {
	claims := &jwt.MapClaims{
		"userid": authUserDTO.UserID,
		"email":  authUserDTO.Email,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
		"iss":    service.issuer,
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtService) GenerateNewTokens(tokenString string) (string, string, error) {
	token, err := service.ValidateRefreshToken(tokenString)

	if err != nil {
		return "", "", err
	}

	if claims := token.Claims.(jwt.MapClaims); token.Valid {
		//Sprintf to convert interface{} to string
		email := fmt.Sprintf("%v", claims["email"])

		authDTO, err := service.database.FindUserForAuth(email)
		if err != nil {
			return "", "", err
		}

		role, err := service.database.GetUserRole(email)
		if err != nil {
			return "", "", err
		}

		return service.GenerateToken(authDTO, role), service.GenerateRefreshToken(authDTO), nil
	} else {
		return "", "", errors.New("Failed to extract JWT claims.")
	}
}

func (service *jwtService) GenerateRefreshToken(authUserDTO dto.AuthUserDTO) string {
	claims := &jwt.MapClaims{
		"userid": authUserDTO.UserID,
		"email":  authUserDTO.Email,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iss":    service.issuer,
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET") + authUserDTO.Password))
	if err != nil {
		panic(err)
	}

	return t
}

func (service *jwtService) ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method.")
		}

		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
}

func (service *jwtService) ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	token, _ := jwt.Parse(tokenString, nil)
	if token == nil {
		return nil, errors.New("Failed to parse refresh token.")
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	//Sprintf to convert interface{} to string
	authUserDTO, err := service.database.FindUserForAuth(fmt.Sprintf("%v", claims["email"]))
	if err != nil {
		return nil, err
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method.")
		}

		return []byte(os.Getenv("REFRESH_TOKEN_SECRET") + authUserDTO.Password), nil
	})
}

func (service *jwtService) GetUserIDAndRole(tokenString string) (string, int, error) {
	token, err := service.ValidateAccessToken(tokenString)

	if err != nil && err.Error() != "Token is expired" {
		return "", -1, err
	}
	err = nil

	if claims := token.Claims.(jwt.MapClaims); token.Valid {
		//Sprintf to convert interface{} to string
		userid := fmt.Sprintf("%v", claims["userid"])

		////Sprintf to convert interface{} to string and then convert it to int
		role, err := strconv.Atoi(fmt.Sprintf("%v", claims["role"]))
		if err != nil {
			return "", -1, err
		}

		return userid, role, nil
	} else {
		return "", -1, errors.New("Failed to extract JWT claims.")
	}
}

func (service *jwtService) GetRole(tokenString string) (int, error) {
	token, err := service.ValidateAccessToken(tokenString)

	if err != nil && err.Error() != "Token is expired" {
		return -1, errors.New("Failed to extract JWT claims.")
	}
	err = nil

	claims := token.Claims.(jwt.MapClaims)

	////Sprintf to convert interface{} to string and then convert it to int
	role, err := strconv.Atoi(fmt.Sprintf("%v", claims["role"]))
	if err != nil {
		return -1, err
	}

	return role, nil
}

func (service *jwtService) GetEmail(tokenString string) (string, error) {
	token, err := service.ValidateAccessToken(tokenString)

	if err != nil && err.Error() != "Token is expired" {
		return "", errors.New("Failed to extract JWT claims.")
	}
	err = nil

	claims := token.Claims.(jwt.MapClaims)

	////Sprintf to convert interface{} to string and then convert it to int
	email := fmt.Sprintf("%v", claims["email"])
	if err != nil {
		return "", err
	}

	return email, nil
}

func (service *jwtService) GenerateOTPToken(email string) string {
	claims := &jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Minute * 4).Unix(),
		"iss":   service.issuer,
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("OTP_TOKEN_SECRET")))
	if err != nil {
		panic(err)
	}

	return t
}

func (service *jwtService) ValidateOTPToken(tokenString string) (*jwt.Token, error) {
	token, _ := jwt.Parse(tokenString, nil)
	if token == nil {
		return nil, errors.New("Failed to parse otp token.")
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method.")
		}

		return []byte(os.Getenv("OTP_TOKEN_SECRET")), nil
	})
}

func (service *jwtService) GetEmailFromOTPToken(tokenString string) (string, error) {
	token, err := service.ValidateOTPToken(tokenString)

	if err != nil && err.Error() != "Token is expired" {
		return "", errors.New("Failed to extract JWT claims.")
	}

	claims := token.Claims.(jwt.MapClaims)

	////Sprintf to convert interface{} to string and then convert it to int
	email := fmt.Sprintf("%v", claims["email"])

	return email, nil
}
