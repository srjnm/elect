package apis

import (
	"elect/controllers"
	"elect/dto"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ElectionAPI struct {
	electionController controllers.ElectionController
}

func NewElectionAPI(electionController controllers.ElectionController) *ElectionAPI {
	return &ElectionAPI{
		electionController: electionController,
	}
}

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
		log.Fatalln(err.Error())
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

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

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Election deleted.",
	})
	return
}

// AddParticipants godoc
// @Summary Add participants to the election you created
// @Tags election
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
// @Tags election
// @Produce json
// @Param election body dto.DeleteParticipantDTO true "Delete Participant"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/participant [delete]
func (election *ElectionAPI) DeleteParticipantHandler(cxt *gin.Context) {
	err := election.electionController.DeleteElection(cxt)

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
