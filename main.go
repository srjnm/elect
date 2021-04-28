package main

import (
	"elect/apis"
	"elect/controllers"
	"elect/database"
	"elect/middlewares"
	"elect/services"
	"net/http"
	_ "net/http"
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

// @host e1ect.herokuapp.com
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

	gin.SetMode(gin.ReleaseMode)

	//Declaring all layers
	postgresDatabase, mux := database.NewPostgresDatabase()
	userService := services.NewUserService(postgresDatabase)
	jwtService := services.NewJWTService("e1ect.herokuapp.com", postgresDatabase)
	userController := controllers.NewUserController(userService, jwtService)
	authAPI := apis.NewAuthAPI(userController)
	userAPI := apis.NewUserAPI(userController)

	port = os.Getenv("PORT")

	config := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000", "http://localhost:3000", "http://localhost:8080", "https://e1ect.herokuapp.com/"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Set-Cookie"},
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
	//Verify Account
	server.POST("/setpassword", middlewares.Authorizer(jwtService, authEnforcer), authAPI.VerifyHandler)
	//OTP Verification FrontEnd
	server.GET("/otp", middlewares.Authorizer(jwtService, authEnforcer), authAPI.OTPGETHandler)
	//OTP Verification
	server.POST("/otp", middlewares.Authorizer(jwtService, authEnforcer), authAPI.OTPHandler)

	apiRoutes := server.Group("/api")
	//Register Students
	apiRoutes.POST("/registerstudents", middlewares.Authorization(jwtService), middlewares.Authorizer(jwtService, authEnforcer), userAPI.RegisterStudentsHandler)
	//Registered Students
	apiRoutes.GET("/registeredstudents", middlewares.Authorization(jwtService), middlewares.Authorizer(jwtService, authEnforcer), userAPI.RegisteredStudentsHandler)

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
