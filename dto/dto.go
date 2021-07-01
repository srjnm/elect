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
	Verified       bool   `json:"verified"`
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

type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type CreateResetTokenDTO struct {
	Email string `json:"email" binding:"email,required"`
}

type ResetPasswordDTO struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
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

// Election DTOs
type CreateElectionDTO struct {
	Title          string `json:"title" binding:"required"`
	StartingAt     string `json:"starting_at" binding:"required"`
	EndingAt       string `json:"ending_at" binding:"required"`
	LockingAt      string `json:"locking_at" binding:"required"`
	GenderSpecific bool   `json:"gender_specific"`
}

type EditElectionDTO struct {
	ElectionId     string `json:"election_id" binding:"required"`
	Title          string `json:"title,omitempty"`
	StartingAt     string `json:"starting_at,omitempty"`
	EndingAt       string `json:"ending_at,omitempty"`
	LockingAt      string `json:"locking_at,omitempty"`
	GenderSpecific bool   `json:"gender_specific,omitempty"`
}

type CreateParticipantDTO struct {
	RegisterNumber string `json:"register_number"`
}

type DeleteParticipantDTO struct {
	ElectionId    string `json:"election_id" binding:"required"`
	ParticipantId string `json:"participant_id" binding:"required"`
}

type CreateCandidateDTO struct {
	ElectionId     string `json:"election_id"`
	Sex            int    `json:"sex"`
	DisplayPicture string `json:"display_picture"`
	Poster         string `json:"poster"`
	IdProof        string `json:"id_proof"`
}

type CastVoteDTO struct {
	ElectionId  string `json:"election_id" binding:"required"`
	CandidateId string `json:"candidate_id" binding:"required"`
}

type CandidateResultsDTO struct {
	CandidateID    string `json:"candidate_id"`
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	Sex            int    `json:"sex"`
	ElectionID     string `json:"election_id"`
	DisplayPicture string `json:"display_picture"`
	Votes          int    `json:"votes"`
}

type GeneralCandidateDTO struct {
	CandidateID    string `json:"candidate_id"`
	UserID         string `json:"user_id"`
	ElectionID     string `json:"election_id"`
	RegisterNo     string `json:"register_no"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Sex            int    `json:"sex"`
	DisplayPicture string `json:"display_picture"`
	Poster         string `json:"poster"`
	IDProof        string `json:"id_proof,omitempty"`
	Approved       bool   `json:"approved"`
}

type GeneralParticipantDTO struct {
	ParticipantID string `json:"participant_id"`
	UserID        string `json:"user_id"`
	RegNumber     string `json:"register_number"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Voted         bool   `json:"voted"`
}

type GeneralElectionDTO struct {
	ElectionID     string                  `json:"election_id"`
	Title          string                  `json:"title"`
	StartingAt     string                  `json:"starting_at"`
	EndingAt       string                  `json:"ending_at"`
	LockingAt      string                  `json:"locking_at"`
	Voted          bool                    `json:"voted,omitempty"`
	GenderSpecific bool                    `json:"gender_specific,omitempty"`
	Participants   []GeneralParticipantDTO `json:"participants,omitempty"`
	Candidates     []GeneralCandidateDTO   `json:"candidates,omitempty"`
	Candidate      *GeneralCandidateDTO    `json:"candidate,omitempty"`
}

type GeneralElectionResultsDTO struct {
	ElectionID        string                `json:"election_id"`
	Title             string                `json:"title"`
	StartingAt        string                `json:"starting_at"`
	EndingAt          string                `json:"ending_at"`
	LockingAt         string                `json:"locking_at"`
	GenderSpecific    bool                  `json:"gender_specific"`
	TotalVotes        int                   `json:"total_votes"`
	TotalParticipants int                   `json:"total_participants"`
	CandidateResults  []CandidateResultsDTO `json:"candidate_results,omitempty"`
	MCandidateResults []CandidateResultsDTO `json:"mcandidate_results,omitempty"`
	FCandidateResults []CandidateResultsDTO `json:"fcandidate_results,omitempty"`
	OCandidateResults []CandidateResultsDTO `json:"ocandidate_results,omitempty"`
}
