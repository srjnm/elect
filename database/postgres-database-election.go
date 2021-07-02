package database

import (
	"elect/dto"
	"elect/mappers"
	"elect/models"
	"errors"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/jinzhu/gorm"
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
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", election.ElectionID.String(), userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if count == 0 {
		log.Println("Unauthorized!")
		return errors.New("Unauthorized!")
	}

	var findElection models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", election.ElectionID.String()).First(&findElection)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if findElection.LockingAt.UTC().Before(time.Now().UTC()) {
		log.Println("Election Locked!")
		return errors.New("Election Locked!")
	}

	election.ElectionID = uuid.Nil
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", findElection.ElectionID.String()).Update(&election)
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
		res := db.connection.Model(&models.Election{}).Where("created_by = ?", userId).Order("created_at DESC").Find(&elections)
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

func (db *postgresDatabase) GetElectionsForStudents(userId string, paginatorParams dto.PaginatorParams) ([]models.Election, []bool, error) {
	var electionIds []models.Participant
	res := db.connection.Model(&models.Participant{}).Where("user_id = ?", userId).Order("created_at DESC").Find(&electionIds)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return nil, nil, res.Error
	}

	var elections []models.Election
	var voted []bool
	if paginatorParams.Page == "" {
		for _, eId := range electionIds {
			var election models.Election
			res := db.connection.Model(&models.Election{}).Where("election_id = ?", eId.ElectionID.String()).First(&election)
			if res.Error != nil {
				log.Println(res.Error.Error())
				return nil, nil, res.Error
			}

			elections = append(elections, election)
			voted = append(voted, eId.Voted)
		}

		return elections, voted, nil
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

	for _, election := range elections {
		for _, eId := range electionIds {
			if election.ElectionID.String() == eId.ElectionID.String() {
				voted = append(voted, eId.Voted)
				break
			}
		}
	}

	return elections, voted, nil
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

func (db *postgresDatabase) GetElectionParticipants(userId string, electionId string) ([]dto.GeneralParticipantDTO, error) {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", electionId, userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return nil, res.Error
	}
	if count == 0 {
		log.Println("Unauthorized!")
		return nil, errors.New("Unauthorized!")
	}

	var participants []models.Participant
	res = db.connection.Model(&models.Participant{}).Where("election_id = ?", electionId).Find(&participants)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return nil, res.Error
	}

	var generalParticipantDTOs []dto.GeneralParticipantDTO
	for _, participant := range participants {
		var user models.User
		res = db.connection.Model(&models.User{}).Where("user_id = ?", participant.UserID.String()).Find(&user)
		if res.Error != nil {
			log.Println(res.Error.Error())
			return nil, res.Error
		}

		generalParticipantDTOs = append(generalParticipantDTOs, mappers.ToGeneralParticipantDTOFromUser(participant.ParticipantID.String(), participant.Voted, user))
	}

	return generalParticipantDTOs, nil
}

func (db *postgresDatabase) GetTotalElectionParticipants(electionId string, userId string) (int, error) {
	var pcount int
	res := db.connection.Model(&models.Participant{}).Where("user_id = ? AND election_id = ?", userId, electionId).Count(&pcount)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return 0, res.Error
	}
	var ccount int
	res = db.connection.Model(&models.Election{}).Where("election_id = ? AND created_by = ?", electionId, userId).Count(&ccount)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return 0, res.Error
	}
	if ccount == 0 && pcount == 0 {
		log.Println("Unauthorized!")
		return 0, errors.New("Unauthorized!")
	}

	var count int
	res = db.connection.Model(&models.Participant{}).Where("election_id = ?", electionId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return 0, res.Error
	}

	return count, nil
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

	res = db.connection.Model(&models.Candidate{}).Where("candidate_id = ? AND election_id = ?", candidateId, candidate.ElectionID.String()).Updates(map[string]interface{}{"approved": true})
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

	res = db.connection.Model(&models.Candidate{}).Where("candidate_id = ? AND election_id = ?", candidateId, candidate.ElectionID.String()).Updates(map[string]interface{}{"approved": false})
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

		generalParticipantDTOs = append(generalParticipantDTOs, mappers.ToGeneralParticipantDTOFromUser(participant.ParticipantID.String(), participant.Voted, user))
	}

	var candidates []models.Candidate
	res = db.connection.Model(&models.Candidate{}).Where("election_id = ?", electionId).Find(&candidates)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, nil, res.Error
	}

	return election, generalParticipantDTOs, candidates, nil
}

func (db *postgresDatabase) GetElectionForStudents(userId string, electionId string) (models.Election, []models.Candidate, models.Candidate, bool, error) {
	var count int
	res := db.connection.Model(&models.Participant{}).Where("election_id = ? AND user_id = ?", electionId, userId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, models.Candidate{}, false, res.Error
	}
	if count == 0 {
		log.Println("Unauthorized!")
		return models.Election{}, nil, models.Candidate{}, false, errors.New("Unauthorized!")
	}

	var election models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Find(&election)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, models.Candidate{}, false, res.Error
	}

	var candidates []models.Candidate
	res = db.connection.Model(&models.Candidate{}).Where("election_id = ? AND approved = ?", electionId, true).Find(&candidates)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, models.Candidate{}, false, res.Error
	}

	var candidate models.Candidate
	var ccount int
	res = db.connection.Model(&models.Candidate{}).Where("election_id = ? AND user_id = ?", electionId, userId).Count(&ccount)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Println(res.Error.Error())
		return models.Election{}, nil, models.Candidate{}, false, res.Error
	} else {
		if ccount > 0 {
			res = db.connection.Model(&models.Candidate{}).Where("election_id = ? AND user_id = ?", electionId, userId).Find(&candidate)
			if res.Error != nil {
				log.Println(res.Error.Error())
				return models.Election{}, nil, models.Candidate{}, false, res.Error
			}
		}
	}

	var participant models.Participant
	res = db.connection.Model(&models.Participant{}).Where("election_id = ? AND user_id = ?", electionId, userId).Find(&participant)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, models.Candidate{}, false, res.Error
	}

	return election, candidates, candidate, participant.Voted, nil
}

func (db *postgresDatabase) CastVote(userId string, electionId string, candidateId string) error {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}
	if count == 0 {
		log.Println("Invalid Election!")
		return errors.New("Invalid Election!")
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

	var election models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Find(&election)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if !(time.Now().UTC().After(election.StartingAt.UTC()) && time.Now().UTC().Before(election.EndingAt.UTC())) {
		log.Println("Election Locked!")
		return errors.New("Election Locked!")
	}

	var participant models.Participant
	res = db.connection.Model(&models.Participant{}).Where("election_id = ? AND user_id = ?", electionId, userId).Find(&participant)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if participant.Voted {
		log.Println("Already voted!")
		return errors.New("Already voted!")
	}

	var candidate models.Candidate
	res = db.connection.Model(&models.Candidate{}).Where("candidate_id = ?", candidateId).Find(&candidate)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	if candidate.Approved == false {
		return errors.New("Unapproved candidate!")
	}

	participant.Voted = true
	candidate.Votes++

	res = db.connection.Model(&models.Participant{}).Update(&participant)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	res = db.connection.Model(&models.Candidate{}).Update(&candidate)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) GetResults(userId string, role int, electionId string) (models.Election, []models.Candidate, []models.Candidate, []models.Candidate, []models.Candidate, int, error) {
	var count int
	res := db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Count(&count)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, nil, nil, nil, 0, res.Error
	}
	if count == 0 {
		log.Println("Invalid Election!")
		return models.Election{}, nil, nil, nil, nil, 0, errors.New("Invalid Election!")
	}

	if role == 0 {
		res = db.connection.Model(&models.Participant{}).Where("user_id = ? AND election_id = ?", userId, electionId).Count(&count)
		if res.Error != nil {
			log.Println(res.Error.Error())
			return models.Election{}, nil, nil, nil, nil, 0, res.Error
		}
		if count == 0 {
			log.Println("You are not the part of the election!")
			return models.Election{}, nil, nil, nil, nil, 0, errors.New("You are not the part of the election!")
		}
	} else if role == 1 || role == 2 {
		res = db.connection.Model(&models.Election{}).Where("created_by = ? AND election_id = ?", userId, electionId).Count(&count)
		if res.Error != nil {
			log.Println(res.Error.Error())
			return models.Election{}, nil, nil, nil, nil, 0, res.Error
		}
		if count == 0 {
			log.Println("Unauthorized!")
			return models.Election{}, nil, nil, nil, nil, 0, errors.New("Unauthorized!")
		}
	}

	var election models.Election
	res = db.connection.Model(&models.Election{}).Where("election_id = ?", electionId).Find(&election)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return models.Election{}, nil, nil, nil, nil, 0, res.Error
	}

	if !(time.Now().UTC().After(election.EndingAt.UTC())) {
		log.Println("Election has not completed!")
		return models.Election{}, nil, nil, nil, nil, 0, errors.New("Election has not completed!")
	}

	if !election.GenderSpecific {
		var candidates []models.Candidate
		res = db.connection.Model(&models.Candidate{}).Where("election_id = ? AND approved = ?", electionId, true).Find(&candidates)
		if res.Error != nil {
			log.Println(res.Error.Error())
			return models.Election{}, nil, nil, nil, nil, 0, res.Error
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Votes > candidates[j].Votes
		})

		total := 0
		for _, candidate := range candidates {
			total += candidate.Votes
		}

		return election, candidates, nil, nil, nil, total, nil
	} else {
		total := 0

		//Get male candidates
		var mCandidates []models.Candidate
		res = db.connection.Model(&models.Candidate{}).Where("election_id = ? AND approved = ? AND sex = ?", electionId, true, 0).Order("votes desc").Find(&mCandidates)
		if res.Error != nil {
			log.Println(res.Error.Error())
			return models.Election{}, nil, nil, nil, nil, 0, res.Error
		}

		for _, candidate := range mCandidates {
			total += candidate.Votes
		}

		//Get female candidates
		var fCandidates []models.Candidate
		res = db.connection.Model(&models.Candidate{}).Where("election_id = ? AND approved = ? AND sex = ?", electionId, true, 1).Order("votes desc").Find(&fCandidates)
		if res.Error != nil {
			log.Println(res.Error.Error())
			return models.Election{}, nil, nil, nil, nil, 0, res.Error
		}

		for _, candidate := range fCandidates {
			total += candidate.Votes
		}

		//Get other candidates
		var oCandidates []models.Candidate
		res = db.connection.Model(&models.Candidate{}).Where("election_id = ? AND approved = ? AND sex = ?", electionId, true, 2).Order("votes desc").Find(&oCandidates)
		if res.Error != nil {
			log.Println(res.Error.Error())
			return models.Election{}, nil, nil, nil, nil, 0, res.Error
		}

		for _, candidate := range oCandidates {
			total += candidate.Votes
		}

		return election, nil, mCandidates, fCandidates, oCandidates, total, nil
	}
}
