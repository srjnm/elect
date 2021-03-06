package apis

import (
	"elect/controllers"
	"elect/dto"
	"errors"
	"log"
	"net/http"
	"strconv"

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

// CreateElection godoc
// @Summary Create Election if you are an Admin
// @ID election
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

	select {
	case update <- []byte("update"):

	default:
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Created Election.",
	})
	return
}

// EditElection godoc
// @Summary Edit the election you created
// @ID election
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

	select {
	case update <- []byte("update"):

	default:
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Election edited.",
	})
	return
}

// DeleteElection godoc
// @Summary Delete the election you created
// @ID electionP
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

	select {
	case update <- []byte("update"):

	default:
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Election deleted.",
	})
	return
}

// AddParticipants godoc
// @Summary Add participants to the election you created
// @ID addParticipants
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
// @ID participant
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
// @ID elections
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
// @ID enrollCandidate
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
// @ID approveCandidate
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
// @ID unapproveCandidate
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
// @ID electionP
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
// @ID castVote
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
// @ID getElectionResults
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

var connections = make(map[*websocket.Conn]bool)

func electionWS(w http.ResponseWriter, r *http.Request) {
	var wsUpgrader = websocket.Upgrader{}
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade the connection: " + err.Error())
		return
	}

	connections[conn] = true

	go func() {
		for {
			err := readFromConnections()
			if err != nil {
				break
			}
		}
		return
	}()

	go func() {
		for {
			u := <-update
			err := writeToConnections(u)
			if err != nil {
				break
			}
		}
		return
	}()

}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func writeToConnections(msg []byte) error {
	if len(connections) == 0 || connections == nil {
		return errors.New("No connections!")
	} else {
		for index := range connections {
			err := index.WriteMessage(1, msg)
			if err != nil {
				index.Close()
				delete(connections, index)
				return nil
			}
		}
		return nil
	}
}

func readFromConnections() error {
	if len(connections) == 0 || connections == nil {
		return errors.New("No connections!")
	} else {
		for index := range connections {
			_, _, err := index.ReadMessage()
			if err != nil {
				index.Close()
				delete(connections, index)
				return nil
			}
		}
		return nil
	}
}
