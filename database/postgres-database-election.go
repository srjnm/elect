package database

import (
	"elect/dto"
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
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) EditElection(userId string, election models.Election) error {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", election.ElectionID, userId).Count(&count)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Fatalln("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", election.ElectionID).First(&findElection)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Fatalln("Election Locked!")
		return errors.New("Election Locked!")
	}

	res = db.connection.Update(&election)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) DeleteElection(userId string, electionId string) error {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", electionId, userId).Count(&count)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Fatalln("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).First(&findElection)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Fatalln("Election Locked!")
		return errors.New("Election Locked!")
	}

	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Delete(&models.Election{ElectionID: uuid.FromStringOrNil(electionId)})
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) AddParticipant(userId string, electId string, regno string) error {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", electId, userId).Count(&count)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Fatalln("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electId).First(&findElection)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Fatalln("Election Locked!")
		return errors.New("Election Locked!")
	}

	res = db.connection.Model(&models.User{}).Where("reg_number = ?", regno).Count(&count)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}
	if count == 0 {
		log.Fatalln(regno + ": Student not registered!")
		return errors.New("Student not registered!")
	}

	var student models.User
	res = db.connection.Model(&models.User{}).Where("reg_number = ?", regno).First(&student)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	res = db.connection.Model(&models.Participant{}).Where("election_id = ? AND user_id = ?", electId, student.UserID.String()).Count(&count)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}
	if count > 0 {
		log.Fatalln(regno + ": Student already registered!")
		return errors.New("Student already registered!")
	}

	participant := models.Participant{
		UserID:     student.UserID,
		ElectionID: uuid.FromStringOrNil(electId),
	}

	res = db.connection.Model(&models.Participant{}).Create(&participant)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) GetElectionForAdmins(userId string, paginatorParams dto.PaginatorParams) ([]models.Election, error) {
	var elections []models.Election
	if paginatorParams.Page == "" {
		res := db.connection.Model(&models.Election{}).Where("created_by = ?", userId).Find(&elections)
		if res.Error != nil {
			log.Fatalln(res.Error.Error())
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

func (db *postgresDatabase) GetElectionForStudents(userId string, paginatorParams dto.PaginatorParams) ([]models.Election, error) {
	var electionIds []models.Participant
	res := db.connection.Model(&models.Participant{}).Where("user_id = ?", userId).Find(&electionIds)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return nil, res.Error
	}

	var elections []models.Election
	if paginatorParams.Page == "" {
		for _, eId := range electionIds {
			var election models.Election
			res := db.connection.Model(&models.Election{}).Where("election_id = ?", eId.ElectionID.String()).First(&election)
			if res.Error != nil {
				log.Fatalln(res.Error.Error())
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
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Fatalln("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).First(&findElection)
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Fatalln("Election Locked!")
		return errors.New("Election Locked!")
	}

	res = db.connection.Model(&models.Participant{}).Where("participant_id = ?", participantId).Delete(&models.Participant{ParticipantID: uuid.FromStringOrNil(participantId)})
	if res.Error != nil {
		log.Fatalln(res.Error.Error())
		return res.Error
	}

	return nil
}
