package dto

type AuthUserDTO struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GeneralUserDTO struct {
	Email     string `json:"email" binding:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      int    `json:"role"`
}

type SetPasswordDTO struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPDTO struct {
	Email string `json:"email" binding:"email,required"`
	OTP   string `json:"otp" binding:"required"`
}

type Response struct {
	Message string `json:"message"`
}
