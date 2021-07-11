package controllers

import (
	"bytes"
	"elect/dto"
	"elect/services"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

type ElectionController interface {
	CreateElection(cxt *gin.Context) error
	EditElection(cxt *gin.Context) error
	DeleteElection(cxt *gin.Context) error
	AddParticipants(cxt *gin.Context) (int, error)
	DeleteParticipant(cxt *gin.Context) error
	GetElections(cxt *gin.Context) ([]dto.GeneralElectionDTO, error)
	EnrollCandidate(cxt *gin.Context) error
	ApproveCandidate(cxt *gin.Context) error
	UnapproveCandidate(cxt *gin.Context) error
	GetElection(cxt *gin.Context) (dto.GeneralElectionDTO, error)
	CastVote(cxt *gin.Context) error
	GetElectionResults(cxt *gin.Context) (dto.GeneralElectionResultsDTO, error)
}

type electionController struct {
	electionService services.ElectionService
	jwtService      services.JWTService
}

func NewElectionController(electionService services.ElectionService, jwtService services.JWTService) ElectionController {
	return &electionController{
		electionService: electionService,
		jwtService:      jwtService,
	}
}

func (controller *electionController) CreateElection(cxt *gin.Context) error {
	var createElectionDTO dto.CreateElectionDTO
	err := cxt.ShouldBindJSON(&createElectionDTO)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return controller.electionService.CreateElection(userId, createElectionDTO)
}

func (controller *electionController) EditElection(cxt *gin.Context) error {
	var editElectionDTO dto.EditElectionDTO
	err := cxt.ShouldBindJSON(&editElectionDTO)
	if err != nil {
		return err
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return err
	}

	return controller.electionService.EditElection(userId, editElectionDTO)
}

func (controller *electionController) DeleteElection(cxt *gin.Context) error {
	electionId := cxt.Param("id")
	if electionId == "" {
		log.Println("Invalid Election ID!")
		return errors.New("Invalid Election ID!")
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return err
	}

	return controller.electionService.DeleteElection(userId, electionId)
}

func (controller *electionController) AddParticipants(cxt *gin.Context) (int, error) {
	electionId := cxt.Param("id")
	if electionId == "" {
		log.Println("Invalid Election ID!")
		return 0, errors.New("Invalid Election ID!")
	}

	file, _, err := cxt.Request.FormFile("participants")
	defer file.Close()
	if err != nil {
		return 0, err
	}

	excFile, err := excelize.OpenReader(file)
	if err != nil {
		return 0, err
	}

	cols, err := excFile.GetCols("Sheet1")
	if err != nil {
		return 0, err
	}

	if strings.ToUpper(cols[0][0]) != "REGNO" {
		return 0, errors.New("Invalid Column name!")
	}

	rows, err := excFile.GetRows("Sheet1")
	if err != nil {
		return 0, err
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return 0, err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return 0, err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return 0, err
	}

	var participants []dto.CreateParticipantDTO
	for _, participant := range rows {
		if len(participant) != 0 {
			participants = append(participants, dto.CreateParticipantDTO{RegisterNumber: strings.ReplaceAll(participant[0], ".0", "")})
		}
	}

	return controller.electionService.AddParticipants(userId, electionId, participants)
}

func (controller *electionController) GetElections(cxt *gin.Context) ([]dto.GeneralElectionDTO, error) {
	cookie, err := cxt.Cookie("token")
	if err != nil {
		return nil, err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return nil, err
	}

	userId, role, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return nil, err
	}

	paginatorParams := dto.PaginatorParams{
		Page:    cxt.Query("page"),
		Limit:   cxt.Query("limit"),
		OrderBy: cxt.Query("orderby") + " " + cxt.Query("order"),
	}

	if paginatorParams.Page != "" || paginatorParams.Limit != "" || paginatorParams.OrderBy != " " {
		if paginatorParams.Page == "" || paginatorParams.Limit == "" || paginatorParams.OrderBy == " " {
			return nil, errors.New("Invalid Query!")
		}
	}

	if role == 1 || role == 2 {
		return controller.electionService.GetElectionsForAdmins(userId, paginatorParams)
	} else if role == 0 {
		return controller.electionService.GetElectionsForStudents(userId, paginatorParams)
	}

	return nil, errors.New("Invalid Role!")
}

func (controller *electionController) DeleteParticipant(cxt *gin.Context) error {
	var deleteParticipantDTO dto.DeleteParticipantDTO

	err := cxt.ShouldBindJSON(&deleteParticipantDTO)
	if err != nil {
		return err
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return err
	}

	return controller.electionService.DeleteParticipant(userId, deleteParticipantDTO.ElectionId, deleteParticipantDTO.ParticipantId)
}

func (controller *electionController) EnrollCandidate(cxt *gin.Context) error {
	electionId := cxt.PostForm("election_id")
	sex, err := strconv.Atoi(cxt.PostForm("sex"))
	if err != nil {
		return err
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = controller.electionService.CheckCandidateEligibility(userId, electionId)
	if err != nil {
		return err
	}

	dpFile, _, err := cxt.Request.FormFile("display_picture")
	defer dpFile.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	dpURL, err := uploadImage(dpFile)
	if err != nil {
		return err
	}

	posterFile, _, err := cxt.Request.FormFile("poster")
	defer posterFile.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	posterURL, err := uploadImage(posterFile)
	if err != nil {
		return err
	}

	idFile, _, err := cxt.Request.FormFile("id_proof")
	defer idFile.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	idURL, err := uploadImage(idFile)
	if err != nil {
		return err
	}

	createCandidateDTO := dto.CreateCandidateDTO{
		ElectionId:     electionId,
		Sex:            sex,
		DisplayPicture: dpURL,
		Poster:         posterURL,
		IdProof:        idURL,
	}

	return controller.electionService.EnrollCandidate(userId, createCandidateDTO)
}

func (controller *electionController) ApproveCandidate(cxt *gin.Context) error {
	candidateId := cxt.Param("id")
	if candidateId == "" {
		log.Println("Invalid ID!")
		return errors.New("Invalid ID!")
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return controller.electionService.ApproveCandidate(userId, candidateId)
}

func (controller *electionController) UnapproveCandidate(cxt *gin.Context) error {
	candidateId := cxt.Param("id")
	if candidateId == "" {
		log.Println("Invalid ID!")
		return errors.New("Invalid ID!")
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return controller.electionService.UnapproveCandidate(userId, candidateId)
}

func (controller *electionController) GetElection(cxt *gin.Context) (dto.GeneralElectionDTO, error) {
	electionId := cxt.Param("id")
	if electionId == "" {
		log.Println("Invalid ID!")
		return dto.GeneralElectionDTO{}, errors.New("Invalid ID!")
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return dto.GeneralElectionDTO{}, err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return dto.GeneralElectionDTO{}, err
	}

	userId, role, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return dto.GeneralElectionDTO{}, err
	}

	if role == 1 || role == 2 {
		return controller.electionService.GetElectionForAdmins(userId, electionId)
	} else if role == 0 {
		return controller.electionService.GetElectionForStudents(userId, electionId)
	}

	return dto.GeneralElectionDTO{}, errors.New("Invalid Role!")
}

func (controller *electionController) CastVote(cxt *gin.Context) error {
	var castVoteDTO dto.CastVoteDTO
	cxt.ShouldBindJSON(&castVoteDTO)

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return err
	}

	return controller.electionService.CastVote(userId, castVoteDTO)
}

func (controller *electionController) GetElectionResults(cxt *gin.Context) (dto.GeneralElectionResultsDTO, error) {
	electionId := cxt.Param("id")
	if electionId == "" {
		log.Println("Invalid ID!")
		return dto.GeneralElectionResultsDTO{}, errors.New("Invalid ID!")
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		return dto.GeneralElectionResultsDTO{}, err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		return dto.GeneralElectionResultsDTO{}, err
	}

	userId, role, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		return dto.GeneralElectionResultsDTO{}, err
	}

	return controller.electionService.GetElectionResults(userId, role, electionId)
}

//Private functions
func uploadImage(file multipart.File) (string, error) {
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Println(err.Error())
		return "", err
	}
	url, err := services.UploadImageToCandidateAzureStorage(buf.Bytes())
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return url, nil
}
