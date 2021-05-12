package controllers

import (
	"elect/dto"
	"elect/email"
	"elect/services"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/hgfischer/go-otp"
	"golang.org/x/crypto/bcrypt"
)

type UserController interface {
	Login(*gin.Context) (string, error)
	Refresh(*gin.Context) error
	Verify(*gin.Context) error
	CheckToken(*gin.Context) error
	OTPVerication(*gin.Context) (string, string, string, error)
	GetOTP(*gin.Context) (string, error)
	Logout(cxt *gin.Context) error
	RegisterStudents(*gin.Context) (int, error)
	RegisteredStudents(cxt *gin.Context) ([]dto.GeneralStudentDTO, error)
	DeleteRegisteredStudent(cxt *gin.Context) error
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

func (controller *userController) Login(cxt *gin.Context) (string, error) {
	var authUser dto.AuthUserDTO

	err := cxt.ShouldBindJSON(&authUser)
	if err != nil {
		return "", err
	}

	dbUser, err := controller.userService.GetUserForAuth(authUser.Email)
	if err != nil {
		return "", err
	}

	auth := CheckPasswordHash(authUser.Password, dbUser.Password)
	if !auth {
		return "", errors.New("Invalid Email or Password!")
	}

	totp := &otp.TOTP{Secret: os.Getenv("OTP_SECRET") + dbUser.Email, Period: 240}

	err = email.SendOTPEmail(dbUser.Email, totp.Get(), "otptemplate.html")
	if err != nil {
		return "", err
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
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
		)
	}

	return dbUser.Email, nil
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

	email, err := controller.jwtService.GetEmail(value["access_token"])
	if err != nil {
		return err
	}

	if !refreshToken.Valid {
		http.SetCookie(
			cxt.Writer,
			&http.Cookie{
				Name:     "token",
				Value:    "",
				MaxAge:   -1,
				Path:     "/",
				Domain:   "",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
		)

		cxt.AbortWithStatusJSON(http.StatusNetworkAuthenticationRequired, dto.Response{
			Message: "Not logged in!",
		})
		return err
	}

	err = controller.userService.CheckIfActiveRefreshToken(value["refresh_token"], email)
	if err != nil {
		http.SetCookie(
			cxt.Writer,
			&http.Cookie{
				Name:     "token",
				Value:    "",
				MaxAge:   -1,
				Path:     "/",
				Domain:   "",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
		)
		cxt.AbortWithStatusJSON(http.StatusNetworkAuthenticationRequired, dto.Response{
			Message: "Not logged in!",
		})
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
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
		)
	}

	if err != nil {
		return err
	}

	err = controller.userService.SetActiveRefreshToken(newRefreshToken, email)
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

func (controller *userController) OTPVerication(cxt *gin.Context) (string, string, string, error) {
	var otpDTO dto.OTPDTO
	err := cxt.ShouldBindJSON(&otpDTO)
	if err != nil {
		return "", "", "", err
	}

	totp := &otp.TOTP{Secret: os.Getenv("OTP_SECRET") + otpDTO.Email, Period: 240}
	valid := totp.Verify(otpDTO.OTP)
	if !valid {
		return "", "", "", errors.New("Invalid OTP!")
	}

	dbUser, err := controller.userService.GetUserForAuth(otpDTO.Email)
	if err != nil {
		return "", "", "", err
	}

	role, err := controller.userService.GetUserRole(otpDTO.Email)
	if err != nil {
		return "", "", "", err
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
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
		)
	}

	err = controller.userService.SetActiveRefreshToken(value["refresh_token"], dbUser.Email)
	if err != nil {
		return "", "", "", err
	}

	return dbUser.UserID, dbUser.Email, strconv.Itoa(role), nil
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

func (controller *userController) Logout(cxt *gin.Context) error {
	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return err
	}

	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}

	email, err := controller.jwtService.GetEmail(value["access_token"])
	if err != nil {
		return err
	}

	return controller.userService.ClearActiveRefreshToken(email)
}

func (controller *userController) RegisterStudents(cxt *gin.Context) (int, error) {
	successCount := 0

	file, _, err := cxt.Request.FormFile("register")
	defer file.Close()
	if err != nil {
		return 0, err
	}

	excFile, err := excelize.OpenReader(file)
	if err != nil {
		return 0, err
	}

	cols, err := excFile.GetCols("Sheet1")
	if err != nil {
		return 0, err
	}

	if strings.ToUpper(cols[0][0]) != "REGNO" && strings.ToUpper(cols[0][1]) != "FIRSTNAME" && strings.ToUpper(cols[0][2]) != "LASTNAME" && strings.ToUpper(cols[0][3]) == "EMAIL" {
		return 0, errors.New("Invalid Column names!")
	}

	rows, err := excFile.GetRows("Sheet1")
	if err != nil {
		return 0, err
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return 0, err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return 0, err
	}

	registeredBy, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return 0, err
	}

	for index, row := range rows {
		if index != 0 {
			if err == nil {
				err = controller.userService.RegisterStudent(dto.RegisterStudentDTO{
					RegNumber:    strings.ReplaceAll(row[0], ".0", ""),
					FirstName:    row[1],
					LastName:     row[2],
					Email:        row[3],
					RegisteredBy: registeredBy,
				})
				if err == nil {
					successCount++
				}
			}
		}
	}

	return successCount, nil
}

func (controller *userController) RegisteredStudents(cxt *gin.Context) ([]dto.GeneralStudentDTO, error) {
	cookie, err := cxt.Cookie("token")
	if err != nil {
		return nil, err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return nil, err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return nil, err
	}

	paginatorParams := dto.PaginatorParams{
		Page:    cxt.Query("page"),
		Limit:   cxt.Query("limit"),
		OrderBy: cxt.Query("orderby") + " " + cxt.Query("order"),
	}

	if paginatorParams.Page != "" || paginatorParams.Limit != "" || paginatorParams.OrderBy != " " {
		if paginatorParams.Page == "" || paginatorParams.Limit == "" || paginatorParams.OrderBy == " " {
			return nil, errors.New("Invalid Query!")
		}
	}

	regStudents, err := controller.userService.RegisteredStudents(userId, paginatorParams)
	if err != nil {
		return nil, err
	}

	return regStudents, nil
}

func (controller *userController) DeleteRegisteredStudent(cxt *gin.Context) error {
	cookie, err := cxt.Cookie("token")
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return err
	}

	studentUserId := cxt.Param("id")
	if studentUserId == "" {
		return errors.New("Invalid Student ID!")
	}

	return controller.userService.DeleteRegisteredStudent(userId, studentUserId)
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
