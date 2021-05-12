package services

import (
	"elect/database"
	"elect/dto"
	"elect/mappers"
	"elect/models"

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

	err := service.database.CreateElection(election)
	if err != nil {
		return err
	}

	return nil
}

func (service *electionService) EditElection(userId string, editElectionDTO dto.EditElectionDTO) error {
	return service.database.EditElection(userId, mappers.ToElectionFromEditElectionDTO(editElectionDTO))
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
		generalCandidateDTOs = append(generalCandidateDTOs, mappers.ToGeneralCandidateDTOFromCandidate(candidate))
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
		generalCandidateDTOs = append(generalCandidateDTOs, mappers.ToGeneralCandidateDTOFromCandidate(candidate))
	}

	return mappers.ToGeneralElectionDTOForStudents(election, generalCandidateDTOs), nil
}
