package dto

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Verify struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTP struct {
	Email string `json:"email" binding:"email,required"`
	OTP   string `json:"otp" binding:"required"`
}

type Elections struct {
	ElectionID     string `json:"election_id"`
	Title          string `json:"title"`
	StartingAt     string `json:"starting_at"`
	EndingAt       string `json:"ending_at"`
	LockingAt      string `json:"locking_at"`
	GenderSpecific bool   `json:"gender_specific,omitempty"`
}
