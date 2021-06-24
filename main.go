package main

import (
	"elect/apis"
	"elect/controllers"
	"elect/database"
	"elect/middlewares"
	"elect/services"
	"log"
	"net/http"
	"os"

	_ "elect/docs"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	cors "github.com/rs/cors/wrapper/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title ELECT REST API
// @version 1.0
// @description This is the backend server of ELECT web application.

// @contact.name ELECT API Support
// @contact.email surajnm15@gmail.com

// @host localhost:8080
// @BasePath /
func main() {
	var port string
	err := godotenv.Load()
	if err != nil {
		port = "8080"
	}

	authEnforcer, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	gin.SetMode(gin.ReleaseMode)

	//Declaring all layers
	postgresDatabase, mux := database.NewPostgresDatabase()
	userService := services.NewUserService(postgresDatabase)
	electionService := services.NewElectionService(postgresDatabase)
	jwtService := services.NewJWTService("e1ect.herokuapp.com", postgresDatabase)
	userController := controllers.NewUserController(userService, jwtService)
	electionController := controllers.NewElectionController(electionService, jwtService)
	authAPI := apis.NewAuthAPI(userController)
	userAPI := apis.NewUserAPI(userController)
	electionAPI := apis.NewElectionAPI(electionController)

	port = os.Getenv("PORT")

	config := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000", "http://localhost:3000", "http://localhost:8080", "https://e1ect.herokuapp.com/", "http://192.168.1.248:3000"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowedMethods:   []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
	})

	server := gin.Default()
	server.Use(config)

	//Setting up the templates
	server.Static("/assets", "./templates/assets")
	server.LoadHTMLGlob("templates/html/*.html")

	//Login
	server.POST("/login", middlewares.Authorizer(jwtService, authEnforcer), authAPI.LoginHandler)
	//Logout
	server.POST("/logout", middlewares.Authorizer(jwtService, authEnforcer), authAPI.LogoutHandler)
	//Refresh
	server.POST("/refresh", middlewares.EnsureValidity(jwtService), authAPI.RefreshHandler)
	//VerifyFrontEnd
	server.GET("/verify/:token", middlewares.Authorizer(jwtService, authEnforcer), authAPI.VerifyGETHandler)
	//Check Verify Token Validity
	server.POST("/verifytoken/:token", middlewares.Authorizer(jwtService, authEnforcer), authAPI.CheckVerifyTokenValidityHandler)
	//Verify Account
	server.POST("/setpassword", middlewares.Authorizer(jwtService, authEnforcer), authAPI.VerifyHandler)
	//OTP Verification FrontEnd
	server.GET("/otp", middlewares.Authorizer(jwtService, authEnforcer), authAPI.OTPGETHandler)
	//OTP Verification
	server.POST("/otp", middlewares.Authorizer(jwtService, authEnforcer), authAPI.OTPHandler)
	//Change Password
	server.POST("/changepassword", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), authAPI.ChangePasswordHandler)
	//Reset Password
	server.POST("/resetpassword", middlewares.Authorizer(jwtService, authEnforcer), authAPI.ResetPasswordHandler)
	//Create Reset Token
	server.POST("/createresettoken", middlewares.Authorizer(jwtService, authEnforcer), authAPI.CreateResetTokenHandler)
	//Check Reset Token Validity
	server.POST("/resettoken/:token", middlewares.Authorizer(jwtService, authEnforcer), authAPI.CheckResetTokenValidityHandler)

	apiRoutes := server.Group("/api")
	//Register Students
	apiRoutes.POST("/registerstudents", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), userAPI.RegisterStudentsHandler)
	//Registered Students
	apiRoutes.GET("/registeredstudents", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), userAPI.RegisteredStudentsHandler)
	//Delete Registered Student
	apiRoutes.DELETE("/registeredstudent/:id", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), userAPI.DeleteRegisteredStudentHandler)

	//Get Elections
	apiRoutes.GET("/elections", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.GetElectionsHandler)
	//Get Election
	apiRoutes.GET("/election/:id", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.GetElectionHandler)
	//Create Election
	apiRoutes.POST("/election", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.CreateElectionHandler)
	//Edit Election
	apiRoutes.PUT("/election", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.EditElectionHandler)
	//Delete Election
	apiRoutes.DELETE("/election/:id", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.DeleteElectionHandler)
	//Add Participants
	apiRoutes.POST("/participants/:id", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.AddParticipantsHandler)
	//Delete Participant
	apiRoutes.DELETE("/participant", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.DeleteParticipantHandler)
	//Enroll Candidate
	apiRoutes.POST("/candidate", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.EnrollCandidateHandler)
	//Approve Candidate
	apiRoutes.POST("/candidate/approve/:id", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.ApproveCandidateHandler)
	//Unapprove Candidate
	apiRoutes.POST("/candidate/unapprove/:id", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.UnapproveCandidateHandler)
	//Cast Vote
	apiRoutes.POST("/vote", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.CastVoteHandler)
	//Get Election Results
	apiRoutes.GET("/results/:id", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), electionAPI.GetElectionResultsHandler)

	//Elections Update WebSocket
	apiRoutes.GET("/ws/election" /*middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService),*/, electionAPI.ElectionUpdatesHandler)

	//Swagger Endpoint Integration
	server.GET("/docs", middlewares.Authorizer(jwtService, authEnforcer), func(cxt *gin.Context) {
		cxt.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	server.GET("/swagger/*any", middlewares.Authorizer(jwtService, authEnforcer), ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//QOR Admin Endpoint Integration
	server.Any("/admin/*resources", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), gin.WrapH(mux))

	server.Run(":" + port)
}
