package controllers

import (
	"elect/dto"
	"elect/services"
	"errors"
	"log"
	"os"
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
	GetElections(cxt *gin.Context) ([]dto.GeneralElectionDTO, error)
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
		log.Fatalln(err.Error())
		return err
	}

	cookie, err := cxt.Cookie("token")
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	var s = securecookie.New([]byte(os.Getenv("COOKIE_HASH_SECRET")), nil)
	value := make(map[string]string)
	err = s.Decode("tokens", cookie, &value)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	userId, _, err := controller.jwtService.GetUserIDAndRole(value["access_token"])
	if err != nil {
		log.Fatalln(err.Error())
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
		log.Fatalln("Invalid Election ID!")
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
		log.Fatalln("Invalid Election ID!")
		return 0, errors.New("Invalid Election ID!")
	}

	file, _, err := cxt.Request.FormFile("participants")
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
		participants = append(participants, dto.CreateParticipantDTO{RegisterNumber: strings.ReplaceAll(participant[0], ".0", "")})
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
		return controller.electionService.GetElectionForAdmins(userId, paginatorParams)
	} else if role == 0 {
		return controller.electionService.GetElectionForStudents(userId, paginatorParams)
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
