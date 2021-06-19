package services

import (
	"elect/database"
	"elect/dto"
	"elect/mappers"
	"elect/models"
	"errors"
	"log"

	uuid "github.com/satori/go.uuid"
)

type ElectionService interface {
	CreateElection(userId string, createElectionDTO dto.CreateElectionDTO) error
	EditElection(userId string, editElectionDTO dto.EditElectionDTO) error
	DeleteElection(userId string, electionId string) error
	AddParticipants(userId string, electionId string, participants []dto.CreateParticipantDTO) (int, error)
	GetElectionsForAdmins(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralElectionDTO, error)
	GetElectionsForStudents(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralElectionDTO, error)
	DeleteParticipant(userId string, electionId string, participantId string) error
	EnrollCandidate(userId string, createCandidateDTO dto.CreateCandidateDTO) error
	CheckCandidateEligibility(userId string, electionId string) error
	ApproveCandidate(userId string, candidateId string) error
	UnapproveCandidate(userId string, candidateId string) error
	GetElectionForAdmins(userId string, electionId string) (dto.GeneralElectionDTO, error)
	GetElectionForStudents(userId string, electionId string) (dto.GeneralElectionDTO, error)
	CastVote(userId string, castVoteDTO dto.CastVoteDTO) error
	GetElectionResults(userId string, role int, electionId string) (dto.GeneralElectionResultsDTO, error)
}

type electionService struct {
	database database.Database
}

func NewElectionService(database database.Database) ElectionService {
	return &electionService{
		database: database,
	}
}

func (service *electionService) CreateElection(userId string, createElectionDTO dto.CreateElectionDTO) error {
	election := mappers.ToElectionFromCreateElectionDTO(createElectionDTO)
	election.CreatedBy = userId

	if election.LockingAt.After(election.StartingAt) {
		return errors.New("Locking At is after Starting At!")
	}
	if election.StartingAt.After(election.EndingAt) {
		return errors.New("Starting At is after Ending At!")
	}

	err := service.database.CreateElection(election)
	if err != nil {
		return err
	}

	return nil
}

func (service *electionService) EditElection(userId string, editElectionDTO dto.EditElectionDTO) error {
	election := mappers.ToElectionFromEditElectionDTO(editElectionDTO)

	if election.LockingAt.After(election.StartingAt) {
		return errors.New("Locking At is after Starting At!")
	}
	if election.StartingAt.After(election.EndingAt) {
		return errors.New("Starting At is after Ending At!")
	}

	return service.database.EditElection(userId, election)
}

func (service *electionService) DeleteElection(userId string, electionId string) error {
	return service.database.DeleteElection(userId, electionId)
}

func (service *electionService) AddParticipants(userId string, electionId string, participants []dto.CreateParticipantDTO) (int, error) {
	count := 0
	for index, participant := range participants {
		if index != 0 {
			err := service.database.AddParticipant(userId, electionId, participant.RegisterNumber)
			if err == nil {
				count++
			}
		}
	}

	return count, nil
}

func (service *electionService) GetElectionsForAdmins(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralElectionDTO, error) {
	var elections []models.Election
	var err error
	if paginatorParams.Page == "" {
		elections, err = service.database.GetElectionsForAdmins(userId, dto.PaginatorParams{})
		if err != nil {
			return nil, err
		}
	} else {
		elections, err = service.database.GetElectionsForAdmins(userId, paginatorParams)
		if err != nil {
			return nil, err
		}
	}

	var generalElectionsDTO []dto.GeneralElectionDTO
	for _, election := range elections {
		generalElectionsDTO = append(generalElectionsDTO, mappers.ToGeneralElectionDTOFromElection(election))
	}

	return generalElectionsDTO, nil
}

func (service *electionService) GetElectionsForStudents(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralElectionDTO, error) {
	var elections []models.Election
	var err error
	if paginatorParams.Page == "" {
		elections, err = service.database.GetElectionsForStudents(userId, dto.PaginatorParams{})
		if err != nil {
			return nil, err
		}
	} else {
		elections, err = service.database.GetElectionsForStudents(userId, paginatorParams)
		if err != nil {
			return nil, err
		}
	}

	var generalElectionsDTO []dto.GeneralElectionDTO
	for _, election := range elections {
		generalElectionsDTO = append(generalElectionsDTO, mappers.ToGeneralElectionDTOFromElection(election))
	}

	return generalElectionsDTO, nil
}

func (service *electionService) DeleteParticipant(userId string, electionId string, participantId string) error {
	return service.database.DeleteParticipant(userId, electionId, participantId)
}

func (service *electionService) EnrollCandidate(userId string, createCandidateDTO dto.CreateCandidateDTO) error {
	candidate := mappers.ToCandidateFromCreateCandidateDTO(createCandidateDTO)
	candidate.UserID = uuid.FromStringOrNil(userId)

	return service.database.EnrollCandidate(candidate)
}

func (service *electionService) CheckCandidateEligibility(userId string, electionId string) error {
	return service.database.CheckCandidateEligibility(userId, electionId)
}

func (service *electionService) ApproveCandidate(userId string, candidateId string) error {
	return service.database.ApproveCandidate(userId, candidateId)
}

func (service *electionService) UnapproveCandidate(userId string, candidateId string) error {
	return service.database.UnapproveCandidate(userId, candidateId)
}

func (service *electionService) GetElectionForAdmins(userId string, electionId string) (dto.GeneralElectionDTO, error) {
	election, generalParticipantDTOs, candidates, err := service.database.GetElectionForAdmins(userId, electionId)
	if err != nil {
		return dto.GeneralElectionDTO{}, err
	}

	var generalCandidateDTOs []dto.GeneralCandidateDTO
	for _, candidate := range candidates {
		user, err := service.database.GetUser(candidate.UserID.String())
		if err != nil {
			return dto.GeneralElectionDTO{}, err
		}

		generalCandidateDTOs = append(generalCandidateDTOs, mappers.ToGeneralCandidateDTOFromCandidate(candidate, user))
	}

	return mappers.ToGeneralElectionDTOForAdmins(election, generalParticipantDTOs, generalCandidateDTOs), nil
}

func (service *electionService) GetElectionForStudents(userId string, electionId string) (dto.GeneralElectionDTO, error) {
	election, candidates, err := service.database.GetElectionForStudents(userId, electionId)
	if err != nil {
		return dto.GeneralElectionDTO{}, err
	}

	var generalCandidateDTOs []dto.GeneralCandidateDTO
	for _, candidate := range candidates {
		user, err := service.database.GetUser(candidate.UserID.String())
		if err != nil {
			return dto.GeneralElectionDTO{}, err
		}

		generalCandidateDTOs = append(generalCandidateDTOs, mappers.ToGeneralCandidateDTOFromCandidate(candidate, user))
	}

	return mappers.ToGeneralElectionDTOForStudents(election, generalCandidateDTOs), nil
}

func (service *electionService) CastVote(userId string, castVoteDTO dto.CastVoteDTO) error {
	return service.database.CastVote(userId, castVoteDTO.ElectionId, castVoteDTO.CandidateId)
}

func (service *electionService) GetElectionResults(userId string, role int, electionId string) (dto.GeneralElectionResultsDTO, error) {
	election, candidates, err := service.database.GetResults(userId, role, electionId)
	if err != nil {
		return dto.GeneralElectionResultsDTO{}, err
	}

	var candidateResultsDTOs []dto.CandidateResultsDTO
	for _, candidate := range candidates {
		user, err := service.database.GetUser(candidate.UserID.String())
		if err != nil {
			return dto.GeneralElectionResultsDTO{}, err
		}
		candidateResultsDTOs = append(candidateResultsDTOs, mappers.ToCandidateResultsDTOFromCandidate(candidate, user.FirstName+" "+user.LastName))
	}

	if role == 1 || role == 2 {
		generalParticipantDTOs, err := service.database.GetElectionParticipants(userId, electionId)
		if err != nil {
			return dto.GeneralElectionResultsDTO{}, err
		}

		return mappers.ToGeneralElectionResultsDTOForAdmins(election, generalParticipantDTOs, candidateResultsDTOs), nil
	} else if role == 0 {
		return mappers.ToGeneralElectionResultsDTOForStudents(election, candidateResultsDTOs), nil
	}

	log.Println("Invalid role!")
	return dto.GeneralElectionResultsDTO{}, errors.New("Invalid role!")
}
