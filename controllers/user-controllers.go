package controllers

import (
	"elect/dto"
	"elect/email"
	"elect/services"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/hgfischer/go-otp"
	"golang.org/x/crypto/bcrypt"
)

type UserController interface {
	Login(*gin.Context) error
	Refresh(*gin.Context) error
	Verify(*gin.Context) error
	CheckToken(*gin.Context) error
	OTPVerication(*gin.Context) error
	GetOTP(*gin.Context) (string, error)
}

type userController struct {
	userService services.UserService
	jwtService  services.JWTService
}

func NewUserController(userService services.UserService, jwtService services.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (controller *userController) Login(cxt *gin.Context) error {
	var authUser dto.AuthUserDTO

	err := cxt.ShouldBindJSON(&authUser)
	if err != nil {
		return err
	}

	dbUser, err := controller.userService.GetUserForAuth(authUser.Email)
	if err != nil {
		return err
	}

	auth := CheckPasswordHash(authUser.Password, dbUser.Password)
	if !auth {
		return errors.New("Invalid Email or Password!")
	}

	totp := &otp.TOTP{Secret: os.Getenv("OTP_SECRET") + dbUser.Email, Period: 240}

	err = email.SendOTPEmail(dbUser.Email, totp.Get(), "otptemplate.html")
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

	value := map[string]string{
		"otp_token": controller.jwtService.GenerateOTPToken(dbUser.Email),
	}

	if encoded, err := s.Encode("token", value); err == nil {
		http.SetCookie(
			cxt.Writer,
			&http.Cookie{
				Name:     "otp",
				Value:    encoded,
				MaxAge:   240,
				Path:     "/",
				Domain:   "",
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteDefaultMode,
			},
		)
	}

	return nil
}

func (controller *userController) Refresh(cxt *gin.Context) error {
	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

	cookie, err := cxt.Cookie("token")
	if err != nil {
		cxt.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Message: "Unauthorized User",
		})
		return err
	}

	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}
	refreshToken, err := controller.jwtService.ValidateRefreshToken(value["refresh_token"])
	if err != nil && err.Error() != "Token is expired" {
		return err
	}

	if !refreshToken.Valid {
		return err
	}

	newAccessToken, newRefreshToken, err := controller.jwtService.GenerateNewTokens(value["refresh_token"])
	if err != nil {
		cxt.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Message: "Session Invalid",
		})
		return err
	}

	newValue := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	if encoded, err := s.Encode("tokens", newValue); err == nil {
		http.SetCookie(
			cxt.Writer,
			&http.Cookie{
				Name:     "token",
				Value:    encoded,
				MaxAge:   3600 * 24 * 7,
				Path:     "/",
				Domain:   "",
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteDefaultMode,
			},
		)
	}

	if err != nil {
		return err
	}

	return nil
}

func (controller *userController) Verify(cxt *gin.Context) error {
	var setPasswordDTO dto.SetPasswordDTO
	err := cxt.ShouldBindJSON(&setPasswordDTO)
	if err != nil {
		return err
	}

	return controller.userService.VerifyUserAndSetPassword(setPasswordDTO)
}

func (controller *userController) CheckToken(cxt *gin.Context) error {
	if cxt.Param("token") == "" {
		return errors.New("Bad Request")
	}

	err := controller.userService.TokenValidity(cxt.Param("token"))
	if err != nil {
		return err
	}

	return nil
}

func (controller *userController) OTPVerication(cxt *gin.Context) error {
	var otpDTO dto.OTPDTO
	err := cxt.ShouldBindJSON(&otpDTO)
	if err != nil {
		return err
	}

	totp := &otp.TOTP{Secret: os.Getenv("OTP_SECRET") + otpDTO.Email, Period: 240}
	valid := totp.Verify(otpDTO.OTP)
	if !valid {
		return errors.New("Invalid OTP!")
	}

	dbUser, err := controller.userService.GetUserForAuth(otpDTO.Email)
	if err != nil {
		return err
	}

	role, err := controller.userService.GetUserRole(otpDTO.Email)
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

	value := map[string]string{
		"access_token":  controller.jwtService.GenerateToken(dbUser, role),
		"refresh_token": controller.jwtService.GenerateRefreshToken(dbUser),
	}

	if encoded, err := s.Encode("tokens", value); err == nil {
		http.SetCookie(
			cxt.Writer,
			&http.Cookie{
				Name:     "token",
				Value:    encoded,
				MaxAge:   3600 * 24 * 7,
				Path:     "/",
				Domain:   "",
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteDefaultMode,
			},
		)
	}

	return nil
}

func (controller *userController) GetOTP(cxt *gin.Context) (string, error) {
	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

	cookie, err := cxt.Cookie("otp")
	if err != nil {
		return "", err
	}

	value := make(map[string]string)
	err = s.Decode("token", cookie, &value)
	if err != nil {
		return "", err
	}

	email, err := controller.jwtService.GetEmailFromOTPToken(value["otp_token"])
	if err != nil {
		return "", err
	}

	return email, nil
}

//Bcrypt Functions
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
