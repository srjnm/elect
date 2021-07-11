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
// @ID login
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
// @ID logout
// @Tags auth
// @Description A user has to be logged in currently to access this endpoint.
// @Produce json
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /ulogout [post]
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
			HttpOnly: true,
		},
	)

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Logged out successfully",
	})
	return
}

func (auth *AuthAPI) LogoutGETHandler(cxt *gin.Context) {
	cxt.HTML(http.StatusOK, "logout.html", nil)
}

// Refresh godoc
// @Summary Refresh Token
// @ID refresh
// @Tags auth
// @Description A user needs a valid refresh token to access this endpoint.
// @Produce json
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /refresh [post]
func (auth *AuthAPI) RefreshHandler(cxt *gin.Context) {
	err := auth.userController.Refresh(cxt)

	if err != nil {
		if err.Error() == "Logged in other device!" {
			cxt.JSON(http.StatusNetworkAuthenticationRequired, dto.Response{
				Message: err.Error(),
			})
			return
		}
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
// @ID verify
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

// SubmitOTP godoc
// @Summary Submit OTP
// @ID submitOTP
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
				HttpOnly: true,
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

// ChangePassword godoc
// @Summary Change your password
// @ID changePassword
// @Tags auth
// @Produce json
// @Param changePassword body dto.ChangePasswordDTO true "Change Password"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /changepassword [post]
func (auth *AuthAPI) ChangePasswordHandler(cxt *gin.Context) {
	err := auth.userController.ChangePassword(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	} else {
		cxt.JSON(http.StatusOK, dto.Response{
			Message: "Password Changed",
		})
		return
	}
}

// CheckVerifyTokenValidity godoc
// @Summary Check if verify token is valid or not
// @ID checkVerifyTokenValidity
// @Tags auth
// @Produce json
// @Param token path string true "Verify Token"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /verifytoken/{token} [post]
func (auth *AuthAPI) CheckVerifyTokenValidityHandler(cxt *gin.Context) {
	err := auth.userController.CheckVerifyTokenValidity(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	} else {
		cxt.JSON(http.StatusOK, dto.Response{
			Message: "Valid Verify Token",
		})
		return
	}
}

// CheckResetTokenValidity godoc
// @Summary Check if reset token is valid or not
// @ID checkResetTokenValidity
// @Tags auth
// @Produce json
// @Param token path string true "Reset Token"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /resettoken/{token} [post]
func (auth *AuthAPI) CheckResetTokenValidityHandler(cxt *gin.Context) {
	err := auth.userController.CheckResetTokenValidity(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	} else {
		cxt.JSON(http.StatusOK, dto.Response{
			Message: "Valid Reset Token",
		})
		return
	}
}

// CreateResetToken godoc
// @Summary Create a reset token and send email to reset password
// @ID createResetToken
// @Tags auth
// @Produce json
// @Param createResetToken body dto.CreateResetTokenDTO true "Email"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /createresettoken [post]
func (auth *AuthAPI) CreateResetTokenHandler(cxt *gin.Context) {
	err := auth.userController.GenerateResetToken(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	} else {
		cxt.JSON(http.StatusOK, dto.Response{
			Message: "Reset token created and email sent",
		})
		return
	}
}

// ResetPassword godoc
// @Summary Reset password if you have a valid token
// @ID resetPassword
// @Tags auth
// @Produce json
// @Param createResetToken body dto.ResetPasswordDTO true "Reset Password"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /resetpassword [post]
func (auth *AuthAPI) ResetPasswordHandler(cxt *gin.Context) {
	err := auth.userController.ResetPassword(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	} else {
		cxt.JSON(http.StatusOK, dto.Response{
			Message: "Password Reset",
		})
		return
	}
}
