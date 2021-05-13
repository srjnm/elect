package apis

import (
	"elect/controllers"
	"elect/dto"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ElectionAPI struct {
	electionController controllers.ElectionController
}

func NewElectionAPI(electionController controllers.ElectionController) *ElectionAPI {
	return &ElectionAPI{
		electionController: electionController,
	}
}

var update = make(chan []byte, 1)
var wg sync.WaitGroup

// CreateElection godoc
// @Summary Create Election if you are an Admin
// @Tags election
// @Produce json
// @Param election body dto.CreateElectionDTO true "Election Details"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/election [post]
func (election *ElectionAPI) CreateElectionHandler(cxt *gin.Context) {
	err := election.electionController.CreateElection(cxt)

	if err != nil {
		log.Println(err.Error())
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	update <- []byte("update")
	wg.Done()
	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Created Election.",
	})
	return
}

// EditElection godoc
// @Summary Edit the election you created
// @Tags election
// @Produce json
// @Param election body dto.EditElectionDTO true "Edit Election"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/election [put]
func (election *ElectionAPI) EditElectionHandler(cxt *gin.Context) {
	err := election.electionController.EditElection(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	update <- []byte("update")
	wg.Done()
	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Election edited.",
	})
	return
}

// DeleteElection godoc
// @Summary Delete the election you created
// @Tags election
// @Produce json
// @Param id path string true "Election ID"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/election/{id} [delete]
func (election *ElectionAPI) DeleteElectionHandler(cxt *gin.Context) {
	err := election.electionController.DeleteElection(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	update <- []byte("update")
	wg.Done()
	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Election deleted.",
	})
	return
}

// AddParticipants godoc
// @Summary Add participants to the election you created
// @Tags participant
// @Consume multipart/form-data
// @Produce json
// @Param id path string true "Election ID"
// @Param participants formData file true "Participants List"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/participants/{id} [post]
func (election *ElectionAPI) AddParticipantsHandler(cxt *gin.Context) {
	pCount, err := election.electionController.AddParticipants(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: strconv.Itoa(pCount) + " participants added.",
	})
	return
}

// DeleteParticipant godoc
// @Summary Delete the participant of the election you created
// @Tags participant
// @Produce json
// @Param participant body dto.DeleteParticipantDTO true "Delete Participant"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/participant [delete]
func (election *ElectionAPI) DeleteParticipantHandler(cxt *gin.Context) {
	err := election.electionController.DeleteParticipant(cxt)

	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Participant deleted.",
	})
	return
}

// Elections godoc
// @Summary Get a list of election you are part of OR you have created
// @Tags election
// @Produce json
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Param orderby query string false "Order By - starting_at"
// @Param order query string false "Order - asc or desc"
// @Success 200 {object} []dto.Elections
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/elections [get]
func (election *ElectionAPI) GetElectionsHandler(cxt *gin.Context) {
	elections, err := election.electionController.GetElections(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, elections)
	return
}

// EnrollCandidate godoc
// @Summary Enroll as a candidate for the election you are part of
// @Tags candidate
// @Consume multipart/form-data
// @Produce json
// @Param candidateDetails formData dto.CandidateInputs true "Candidate Details"
// @Param display_picture formData file true "Display Picture"
// @Param poster formData file true "poster"
// @Param id_proof formData file true "ID Proof"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/candidate [post]
func (election *ElectionAPI) EnrollCandidateHandler(cxt *gin.Context) {
	err := election.electionController.EnrollCandidate(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Enrolled as candidate.",
	})
	return
}

// ApproveCandidate godoc
// @Summary Approve enrolled candidates to the election you created
// @Tags candidate
// @Produce json
// @Param id path string true "Candidate ID"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/candidate/approve/{id} [post]
func (election *ElectionAPI) ApproveCandidateHandler(cxt *gin.Context) {
	err := election.electionController.ApproveCandidate(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Candidate approved.",
	})
	return
}

// UnapproveCandidate godoc
// @Summary Unapprove enrolled candidates to the election you created
// @Tags candidate
// @Produce json
// @Param id path string true "Candidate ID"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/candidate/unapprove/{id} [post]
func (election *ElectionAPI) UnapproveCandidateHandler(cxt *gin.Context) {
	err := election.electionController.UnapproveCandidate(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Candidate unapproved.",
	})
	return
}

// GetElection godoc
// @Summary Get details of the election you created or you are part of
// @Tags election
// @Produce json
// @Param id path string true "Election ID"
// @Success 200 {object} dto.GeneralElectionDTO
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/election/{id} [get]
func (election *ElectionAPI) GetElectionHandler(cxt *gin.Context) {
	generalElectionDTO, err := election.electionController.GetElection(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, generalElectionDTO)
	return
}

// CastVote godoc
// @Summary Cast vote to the candidate of the election you are part of
// @Tags participant
// @Produce json
// @Param vote body dto.CastVoteDTO true "Cast Vote"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/vote [post]
func (election *ElectionAPI) CastVoteHandler(cxt *gin.Context) {
	err := election.electionController.CastVote(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Vote casted.",
	})
	return
}

// GetElectionResults godoc
// @Summary Get the results of the election you were part of or you created
// @Tags election
// @Produce json
// @Param id path string true "Election ID"
// @Success 200 {object} dto.GeneralElectionResultsDTO
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/results/{id} [get]
func (election *ElectionAPI) GetElectionResultsHandler(cxt *gin.Context) {
	generalElectionResultsDTO, err := election.electionController.GetElectionResults(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, generalElectionResultsDTO)
	return
}

func (election *ElectionAPI) ElectionUpdatesHandler(cxt *gin.Context) {
	electionWS(cxt.Writer, cxt.Request)
}

var wsUpgrader = websocket.Upgrader{}

func electionWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade the connection: " + err.Error())
		return
	}

	for {
		wg.Add(1)
		u := <-update
		wg.Wait()
		conn.WriteMessage(1, u)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
