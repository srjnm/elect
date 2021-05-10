package services

import (
	"elect/database"
	"elect/dto"
	"elect/mappers"
	"elect/models"
)

type ElectionService interface {
	CreateElection(userId string, createElectionDTO dto.CreateElectionDTO) error
	EditElection(userId string, editElectionDTO dto.EditElectionDTO) error
	DeleteElection(userId string, electionId string) error
	AddParticipants(userId string, electionId string, participants []dto.CreateParticipantDTO) (int, error)
	GetElectionForAdmins(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralElectionDTO, error)
	GetElectionForStudents(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralElectionDTO, error)
	DeleteParticipant(userId string, electionId string, participantId string) error
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

func (service *electionService) GetElectionForAdmins(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralElectionDTO, error) {
	var elections []models.Election
	var err error
	if paginatorParams.Page == "" {
		elections, err = service.database.GetElectionForAdmins(userId, dto.PaginatorParams{})
		if err != nil {
			return nil, err
		}
	} else {
		elections, err = service.database.GetElectionForAdmins(userId, paginatorParams)
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

func (service *electionService) GetElectionForStudents(userId string, paginatorParams dto.PaginatorParams) ([]dto.GeneralElectionDTO, error) {
	var elections []models.Election
	var err error
	if paginatorParams.Page == "" {
		elections, err = service.database.GetElectionForStudents(userId, dto.PaginatorParams{})
		if err != nil {
			return nil, err
		}
	} else {
		elections, err = service.database.GetElectionForStudents(userId, paginatorParams)
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
