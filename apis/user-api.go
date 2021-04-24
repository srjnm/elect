package apis

import (
	"elect/controllers"
	"elect/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	userController controllers.UserController
}

func NewUserAPI(userController controllers.UserController) *UserAPI {
	return &UserAPI{
		userController: userController,
	}
}

// RegisterStudents godoc
// @Summary Register Students if you are an Admin
// @Tags user
// @Consume multipart/form-data
// @Produce json
// @Param register formData file true "Student List"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/registerstudents [post]
func (user *UserAPI) RegisterStudentsHandler(cxt *gin.Context) {
	success, err := user.userController.RegisterStudents(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, dto.Response{
		Message: "Registered " + strconv.Itoa(success) + " students.",
	})
	return
}

// RegisteredStudents godoc
// @Summary Get a list of Students you've registered
// @Tags user
// @Produce json
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Param orderby query string false "Order By - reg_number"
// @Param order query string false "Order - asc or desc"
// @Success 200 {object} []dto.GeneralStudentDTO
// @Failure 401 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/registeredstudents [get]
func (user *UserAPI) RegisteredStudentsHandler(cxt *gin.Context) {
	registeredStudents, err := user.userController.RegisteredStudents(cxt)
	if err != nil {
		cxt.JSON(http.StatusBadRequest, dto.Response{
			Message: err.Error(),
		})
		return
	}

	cxt.JSON(http.StatusOK, registeredStudents)
	return
}
