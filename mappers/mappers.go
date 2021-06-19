package mappers

import (
	"elect/dto"
	"elect/models"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
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
		UserID:         user.UserID.String(),
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		RegisterNumber: user.RegNumber,
		Verified:       user.Verified,
	}
}

func ToElectionFromCreateElectionDTO(electionDTO dto.CreateElectionDTO) models.Election {
	var sTime, eTime, lTime time.Time
	if strings.Contains(electionDTO.StartingAt, "(") {
		sT := strings.SplitAfter(electionDTO.StartingAt, "(")[0]
		sTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", sT[:len(sT)-2])
	} else {
		sTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", electionDTO.StartingAt)
	}

	if strings.Contains(electionDTO.EndingAt, "(") {
		eT := strings.SplitAfter(electionDTO.EndingAt, "(")[0]
		eTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", eT[:len(eT)-2])
	} else {
		eTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", electionDTO.EndingAt)
	}

	if strings.Contains(electionDTO.LockingAt, "(") {
		lT := strings.SplitAfter(electionDTO.LockingAt, "(")[0]
		lTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", lT[:len(lT)-2])
	} else {
		lTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", electionDTO.LockingAt)
	}

	return models.Election{
		Title:          electionDTO.Title,
		StartingAt:     sTime,
		EndingAt:       eTime,
		LockingAt:      lTime,
		GenderSpecific: electionDTO.GenderSpecific,
	}
}

func ToElectionFromEditElectionDTO(editElectionDTO dto.EditElectionDTO) models.Election {
	var sTime, eTime, lTime time.Time
	if editElectionDTO.StartingAt != "" {
		if strings.Contains(editElectionDTO.StartingAt, "(") {
			sT := strings.SplitAfter(editElectionDTO.StartingAt, "(")[0]
			sTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", sT[:len(sT)-2])
		} else {
			sTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", editElectionDTO.StartingAt)
		}
		sTime = sTime.UTC()
	}

	if editElectionDTO.EndingAt != "" {
		if strings.Contains(editElectionDTO.EndingAt, "(") {
			eT := strings.SplitAfter(editElectionDTO.EndingAt, "(")[0]
			eTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", eT[:len(eT)-2])
		} else {
			eTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", editElectionDTO.EndingAt)
		}
		eTime = eTime.UTC()
	}

	if editElectionDTO.LockingAt != "" {
		if strings.Contains(editElectionDTO.LockingAt, "(") {
			lT := strings.SplitAfter(editElectionDTO.LockingAt, "(")[0]
			lTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", lT[:len(lT)-2])
		} else {
			lTime, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700", editElectionDTO.LockingAt)
		}
		lTime = lTime.UTC()
	}

	return models.Election{
		ElectionID:     uuid.FromStringOrNil(editElectionDTO.ElectionId),
		Title:          editElectionDTO.Title,
		StartingAt:     sTime,
		EndingAt:       eTime,
		LockingAt:      lTime,
		GenderSpecific: editElectionDTO.GenderSpecific,
	}
}

func ToGeneralElectionDTOFromElection(election models.Election) dto.GeneralElectionDTO {
	return dto.GeneralElectionDTO{
		ElectionID:     election.ElectionID.String(),
		Title:          election.Title,
		StartingAt:     election.StartingAt.String(),
		EndingAt:       election.EndingAt.String(),
		LockingAt:      election.LockingAt.String(),
		GenderSpecific: election.GenderSpecific,
	}
}

func ToCandidateFromCreateCandidateDTO(createCandidateDTO dto.CreateCandidateDTO) models.Candidate {
	return models.Candidate{
		ElectionID:     uuid.FromStringOrNil(createCandidateDTO.ElectionId),
		Sex:            createCandidateDTO.Sex,
		DisplayPicture: createCandidateDTO.DisplayPicture,
		Poster:         createCandidateDTO.Poster,
		IDProof:        createCandidateDTO.IdProof,
	}
}

func ToGeneralParticipantDTOFromUser(participantId string, user models.User) dto.GeneralParticipantDTO {
	return dto.GeneralParticipantDTO{
		ParticipantID: participantId,
		UserID:        user.UserID.String(),
		RegNumber:     user.RegNumber,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
	}
}

func ToGeneralCandidateDTOFromCandidate(candidate models.Candidate, user models.User) dto.GeneralCandidateDTO {
	return dto.GeneralCandidateDTO{
		CandidateID:    candidate.CandidateID.String(),
		UserID:         candidate.UserID.String(),
		ElectionID:     candidate.ElectionID.String(),
		RegisterNo:     user.RegNumber,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Sex:            candidate.Sex,
		DisplayPicture: candidate.DisplayPicture,
		Poster:         candidate.Poster,
		IDProof:        candidate.IDProof,
		Approved:       candidate.Approved,
	}
}

func ToGeneralElectionDTOForAdmins(election models.Election, generalParticipantDTOs []dto.GeneralParticipantDTO, generalCandidateDTOs []dto.GeneralCandidateDTO) dto.GeneralElectionDTO {
	return dto.GeneralElectionDTO{
		ElectionID:     election.ElectionID.String(),
		Title:          election.Title,
		StartingAt:     election.StartingAt.String(),
		EndingAt:       election.EndingAt.String(),
		LockingAt:      election.LockingAt.String(),
		GenderSpecific: election.GenderSpecific,
		Participants:   generalParticipantDTOs,
		Candidates:     generalCandidateDTOs,
	}
}

func ToGeneralElectionDTOForStudents(election models.Election, generalCandidateDTOs []dto.GeneralCandidateDTO) dto.GeneralElectionDTO {
	return dto.GeneralElectionDTO{
		ElectionID:     election.ElectionID.String(),
		Title:          election.Title,
		StartingAt:     election.StartingAt.String(),
		EndingAt:       election.EndingAt.String(),
		LockingAt:      election.LockingAt.String(),
		GenderSpecific: election.GenderSpecific,
		Candidates:     generalCandidateDTOs,
	}
}

func ToGeneralElectionResultsDTOForAdmins(election models.Election, generalParticipantDTOs []dto.GeneralParticipantDTO, candidateResultsDTOs []dto.CandidateResultsDTO, mCandidateResultsDTOs []dto.CandidateResultsDTO, fCandidateResultsDTOs []dto.CandidateResultsDTO, oCandidateResultsDTOs []dto.CandidateResultsDTO, total int) dto.GeneralElectionResultsDTO {
	return dto.GeneralElectionResultsDTO{
		ElectionID:        election.ElectionID.String(),
		Title:             election.Title,
		StartingAt:        election.StartingAt.String(),
		EndingAt:          election.EndingAt.String(),
		LockingAt:         election.LockingAt.String(),
		GenderSpecific:    election.GenderSpecific,
		TotalVotes:        total,
		Participants:      generalParticipantDTOs,
		CandidateResults:  candidateResultsDTOs,
		MCandidateResults: mCandidateResultsDTOs,
		FCandidateResults: fCandidateResultsDTOs,
		OCandidateResults: oCandidateResultsDTOs,
	}
}

func ToGeneralElectionResultsDTOForStudents(election models.Election, candidateResultsDTOs []dto.CandidateResultsDTO, mCandidateResultsDTOs []dto.CandidateResultsDTO, fCandidateResultsDTOs []dto.CandidateResultsDTO, oCandidateResultsDTOs []dto.CandidateResultsDTO, total int) dto.GeneralElectionResultsDTO {
	return dto.GeneralElectionResultsDTO{
		ElectionID:        election.ElectionID.String(),
		Title:             election.Title,
		StartingAt:        election.StartingAt.String(),
		EndingAt:          election.EndingAt.String(),
		LockingAt:         election.LockingAt.String(),
		GenderSpecific:    election.GenderSpecific,
		TotalVotes:        total,
		CandidateResults:  candidateResultsDTOs,
		MCandidateResults: mCandidateResultsDTOs,
		FCandidateResults: fCandidateResultsDTOs,
		OCandidateResults: oCandidateResultsDTOs,
	}
}

func ToCandidateResultsDTOFromCandidate(candidate models.Candidate, name string) dto.CandidateResultsDTO {
	return dto.CandidateResultsDTO{
		CandidateID:    candidate.CandidateID.String(),
		UserID:         candidate.UserID.String(),
		Name:           name,
		Sex:            candidate.Sex,
		ElectionID:     candidate.ElectionID.String(),
		DisplayPicture: candidate.DisplayPicture,
		Votes:          candidate.Votes,
	}
}
