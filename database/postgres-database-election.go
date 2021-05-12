package database

import (
	"elect/dto"
	"elect/mappers"
	"elect/models"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/biezhi/gorm-paginator/pagination"
	uuid "github.com/satori/go.uuid"
)

func (db *postgresDatabase) CreateElection(election models.Election) error {
	res := db.connection.Create(&election)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) EditElection(userId string, election models.Election) error {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", election.ElectionID, userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Println("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", election.ElectionID).First(&findElection)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Println("Election Locked!")
		return errors.New("Election Locked!")
	}

	res = db.connection.Update(&election)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) DeleteElection(userId string, electionId string) error {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", electionId, userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Println("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).First(&findElection)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Println("Election Locked!")
		return errors.New("Election Locked!")
	}

	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Delete(&models.Election{ElectionID: uuid.FromStringOrNil(electionId)})
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) AddParticipant(userId string, electId string, regno string) error {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", electId, userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Println("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electId).First(&findElection)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Println("Election Locked!")
		return errors.New("Election Locked!")
	}

	res = db.connection.Model(&models.User{}).Where("reg_number = ?", regno).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}
	if count == 0 {
		log.Println(regno + ": Student not registered!")
		return errors.New("Student not registered!")
	}

	var student models.User
	res = db.connection.Model(&models.User{}).Where("reg_number = ?", regno).First(&student)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	res = db.connection.Model(&models.Participant{}).Where("election_id = ? AND user_id = ?", electId, student.UserID.String()).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}
	if count > 0 {
		log.Println(regno + ": Student already registered!")
		return errors.New("Student already registered!")
	}

	participant := models.Participant{
		UserID:     student.UserID,
		ElectionID: uuid.FromStringOrNil(electId),
	}

	res = db.connection.Model(&models.Participant{}).Create(&participant)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) GetElectionsForAdmins(userId string, paginatorParams dto.PaginatorParams) ([]models.Election, error) {
	var elections []models.Election
	if paginatorParams.Page == "" {
		res := db.connection.Model(&models.Election{}).Where("created_by = ?", userId).Find(&elections)
		if res.Error != nil {
			log.Println(res.Error.Error())
			return nil, res.Error
		}

		return elections, nil
	}

	page, _ := strconv.Atoi(paginatorParams.Page)
	limit, _ := strconv.Atoi(paginatorParams.Limit)

	_ = pagination.Paging(
		&pagination.Param{
			DB:      db.connection.Model(&models.Election{}).Where("created_by = ?", userId),
			Page:    page,
			Limit:   limit,
			OrderBy: []string{paginatorParams.OrderBy},
			ShowSQL: false,
		},
		&elections,
	)

	return elections, nil
}

func (db *postgresDatabase) GetElectionsForStudents(userId string, paginatorParams dto.PaginatorParams) ([]models.Election, error) {
	var electionIds []models.Participant
	res := db.connection.Model(&models.Participant{}).Where("user_id = ?", userId).Find(&electionIds)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return nil, res.Error
	}

	var elections []models.Election
	if paginatorParams.Page == "" {
		for _, eId := range electionIds {
			var election models.Election
			res := db.connection.Model(&models.Election{}).Where("election_id = ?", eId.ElectionID.String()).First(&election)
			if res.Error != nil {
				log.Println(res.Error.Error())
				return nil, res.Error
			}

			elections = append(elections, election)
		}

		return elections, nil
	}

	var electIdS []string

	for _, eId := range electionIds {
		electIdS = append(electIdS, eId.ElectionID.String())
	}

	page, _ := strconv.Atoi(paginatorParams.Page)
	limit, _ := strconv.Atoi(paginatorParams.Limit)

	_ = pagination.Paging(
		&pagination.Param{
			DB:      db.connection.Model(&models.Election{}).Where(map[string]interface{}{"election_id": electIdS}),
			Page:    page,
			Limit:   limit,
			OrderBy: []string{paginatorParams.OrderBy},
			ShowSQL: false,
		},
		&elections,
	)

	return elections, nil
}

func (db *postgresDatabase) DeleteParticipant(userId string, electionId string, participantId string) error {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", electionId, userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Println("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).First(&findElection)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Println("Election Locked!")
		return errors.New("Election Locked!")
	}

	res = db.connection.Model(&models.Participant{}).Where("participant_id = ?", participantId).Delete(&models.Participant{ParticipantID: uuid.FromStringOrNil(participantId)})
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) EnrollCandidate(candidate models.Candidate) error {
	res := db.connection.Model(&models.Candidate{}).Create(&candidate)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) CheckCandidateEligibility(userId string, electionId string) error {
	var count int
	res := db.connection.Model(&models.Blacklist{}).Where("user_id = ?", userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}
	if count != 0 {
		log.Println("You have been blacklisted from enrolling as a candidate!")
		return errors.New("You have been blacklisted from enrolling as a candidate!")
	}

	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}
	if count == 0 {
		log.Println("Invalid election!")
		return errors.New("Invalid election!")
	}

	res = db.connection.Model(&models.Participant{}).Where("user_id = ? AND election_id = ?", userId, electionId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}
	if count == 0 {
		log.Println("You are not the part of the election!")
		return errors.New("You are not the part of the election!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).First(&findElection)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Println("Election Locked!")
		return errors.New("Election Locked!")
	}

	return nil
}

func (db *postgresDatabase) ApproveCandidate(userId string, candidateId string) error {
	var candidate models.Candidate
	res := db.connection.Model(&models.Candidate{}).Where("candidate_id = ?", candidateId).Find(&candidate)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	var count int
	res = db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", candidate.ElectionID.String(), userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Println("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", candidate.ElectionID.String()).First(&findElection)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Println("Election Locked!")
		return errors.New("Election Locked!")
	}

	candidate.Approved = true
	res = db.connection.Model(&models.Candidate{}).Update(&candidate)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) UnapproveCandidate(userId string, candidateId string) error {
	var candidate models.Candidate
	res := db.connection.Model(&models.Candidate{}).Where("candidate_id = ?", candidateId).Find(&candidate)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	var count int
	res = db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", candidate.ElectionID.String(), userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Println("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", candidate.ElectionID.String()).First(&findElection)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Println("Election Locked!")
		return errors.New("Election Locked!")
	}

	candidate.Approved = false
	res = db.connection.Model(&models.Candidate{}).Update(&candidate)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) GetElectionForAdmins(userId string, electionId string) (models.Election, []dto.GeneralParticipantDTO, []models.Candidate, error) {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", electionId, userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, nil, res.Error
	}
	if count == 0 {
		log.Println("Unauthorized!")
		return models.Election{}, nil, nil, errors.New("Unauthorized!")
	}

	var election models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Find(&election)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, nil, res.Error
	}

	var participants []models.Participant
	res = db.connection.Model(&models.Participant{}).Where("election_id = ?", electionId).Find(&participants)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, nil, res.Error
	}

	var generalParticipantDTOs []dto.GeneralParticipantDTO
	for _, participant := range participants {
		var user models.User
		res = db.connection.Model(&models.User{}).Where("user_id = ?", participant.UserID.String()).Find(&user)
		if res.Error != nil {
			log.Println(res.Error.Error())
			return models.Election{}, nil, nil, res.Error
		}

		generalParticipantDTOs = append(generalParticipantDTOs, mappers.ToGeneralParticipantDTOFromUser(participant.ParticipantID.String(), user))
	}

	var candidates []models.Candidate
	res = db.connection.Model(&models.Candidate{}).Where("election_id = ?", electionId).Find(&candidates)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, nil, res.Error
	}

	return election, generalParticipantDTOs, candidates, nil
}

func (db *postgresDatabase) GetElectionForStudents(userId string, electionId string) (models.Election, []models.Candidate, error) {
	var count int
	res := db.connection.Model(&models.Participant{}).Where("election_id = ? AND user_id = ?", electionId, userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, res.Error
	}
	if count == 0 {
		log.Println("Unauthorized!")
		return models.Election{}, nil, errors.New("Unauthorized!")
	}

	var election models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Find(&election)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, res.Error
	}

	var candidates []models.Candidate
	res = db.connection.Model(&models.Candidate{}).Where("election_id = ?", electionId).Find(&candidates)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, res.Error
	}

	return election, candidates, nil
}
