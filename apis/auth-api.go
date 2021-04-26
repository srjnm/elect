package apis

import (
	"elect/controllers"
	"elect/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthAPI struct {
	userController controllers.UserController
}

func NewAuthAPI(userController controllers.UserController) *AuthAPI {
	return &AuthAPI{
		userController: userController,
	}
}

// Login godoc
// @Summary User Login
// @Tags auth
// @Consume json
// @Produce json
// @Param login body dto.Login true "User Login"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /login [post]
func (auth *AuthAPI) LoginHandler(cxt *gin.Context) {
	email, err := auth.userController.Login(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, dto.LoginResponse{
		Email:   email,
		Message: "OTP Sent",
	})
	return
}

// Logout godoc
// @Summary User Logout
// @Tags auth
// @Description A user has to be logged in currently to access this endpoint.
// @Produce json
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /logout [post]
func (auth *AuthAPI) LogoutHandler(cxt *gin.Context) {
	err := auth.userController.Logout(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	http.SetCookie(
		cxt.Writer,
		&http.Cookie{
			Name:     "token",
			Value:    "",
			MaxAge:   -1,
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	)

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Logged out successfully",
	})
	return
}

// Refresh godoc
// @Summary Refresh Token
// @Tags auth
// @Description A user needs a valid refresh token to access this endpoint.
// @Produce json
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /refresh [post]
func (auth *AuthAPI) RefreshHandler(cxt *gin.Context) {
	err := auth.userController.Refresh(cxt)

	if err.Error() == "Logged in other device!" {
		cxt.JSON(http.StatusNetworkAuthenticationRequired, dto.Response{
			Message: err.Error(),
		})
		return
	}

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	} else {
		cxt.JSON(http.StatusOK, dto.Response{
			Message: "Token Refreshed",
		})
		return
	}
}

// Verify godoc
// @Summary Verify Email and Set Password
// @Tags auth
// @Produce json
// @Param password body dto.Verify true "Verify"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /setpassword [post]
func (auth *AuthAPI) VerifyHandler(cxt *gin.Context) {
	err := auth.userController.Verify(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	} else {
		cxt.JSON(http.StatusOK, dto.Response{
			Message: "Account Verified",
		})
		return
	}
}

// SetPassword godoc
// @Summary Set Password
// @Tags auth
// @Produce text/html
// @Success 200
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /verify/{token} [get]
func (auth *AuthAPI) VerifyGETHandler(cxt *gin.Context) {
	err := auth.userController.CheckToken(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.HTML(http.StatusOK, "index.html", gin.H{"token": cxt.Param("token")})
}

// EnterOTP godoc
// @Summary Enter OTP
// @Tags auth
// @Produce text/html
// @Success 200
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /otp [get]
func (auth *AuthAPI) OTPGETHandler(cxt *gin.Context) {
	email, err := auth.userController.GetOTP(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.HTML(http.StatusOK, "otp.html", gin.H{"email": email})
}

// EnterOTP godoc
// @Summary Enter OTP
// @Tags auth
// @Produce json
// @Param otp body dto.OTP true "Verify OTP"
// @Success 200 {object} dto.OTPResponse
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /otp [post]
func (auth *AuthAPI) OTPHandler(cxt *gin.Context) {
	userId, email, role, err := auth.userController.OTPVerication(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	} else {
		http.SetCookie(
			cxt.Writer,
			&http.Cookie{
				Name:     "otp",
				Value:    "",
				MaxAge:   -1,
				Path:     "/",
				Domain:   "",
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			},
		)
		cxt.JSON(http.StatusOK, dto.OTPResponse{
			UserId:  userId,
			Email:   email,
			Role:    role,
			Message: "Login Sucessful!",
		})
		return
	}
}
