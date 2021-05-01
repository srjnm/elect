package dto

type GeneralUserDTO struct {
	Email     string `json:"email" binding:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      int    `json:"role"`
}

type GeneralStudentDTO struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email" binding:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	RegisterNumber string `json:"reg_number"`
}

// Auth DTOs
type AuthUserDTO struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SetPasswordDTO struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPDTO struct {
	Email string `json:"email" binding:"email,required"`
	OTP   string `json:"otp" binding:"required"`
}

// Register DTOs
type RegisterStudentDTO struct {
	Email        string `json:"email" binding:"email,required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name"`
	RegNumber    string `json:"reg_number" binding:"required"`
	RegisteredBy string `json:"registered_by" binding:"required"`
}

// Response DTOs
type Response struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

type OTPResponse struct {
	UserId  string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Message string `json:"message"`
}

// Paginator Params
type PaginatorParams struct {
	Page    string `json:"page"`
	Limit   string `json:"limit"`
	OrderBy string `json:"order_by"`
}
