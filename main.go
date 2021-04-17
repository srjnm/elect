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

	port = os.Getenv("PORT")

	config := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000", "http://localhost:3000", "http://localhost:8080", "https://e1ect.herokuapp.com/"},
		AllowCredentials: true,
	})

	server := gin.Default()
	server.Use(config)

	//Setting up the templates
	server.Static("/assets", "./templates/assets")
	server.LoadHTMLGlob("templates/html/*.html")

	server.Use()

	//Login
	server.POST("/login", authAPI.LoginHandler, middlewares.Authorizer(jwtService, authEnforcer))
	//Logout
	server.POST("/logout", authAPI.LogoutHandler, middlewares.Authorizer(jwtService, authEnforcer))
	//Refresh
	server.POST("/refresh", middlewares.EnsureValidity(jwtService), middlewares.Authorizer(jwtService, authEnforcer), authAPI.RefreshHandler)
	//VerifyFrontEnd
	server.GET("/verify/:token", authAPI.VerifyGETHandler, middlewares.Authorizer(jwtService, authEnforcer))
	//Verify Account
	server.POST("/setpassword", authAPI.VerifyHandler, middlewares.Authorizer(jwtService, authEnforcer))
	//OTP Verification FrontEnd
	server.GET("/otp", authAPI.OTPGETHandler, middlewares.Authorizer(jwtService, authEnforcer))
	//OTP Verification
	server.POST("/otp", authAPI.OTPHandler, middlewares.Authorizer(jwtService, authEnforcer))

	//Swagger Endpoint Integration
	server.GET("/docs", func(cxt *gin.Context) {
		cxt.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	},
		middlewares.Authorizer(jwtService, authEnforcer),
	)
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url), middlewares.Authorizer(jwtService, authEnforcer))

	//QOR Admin Endpoint Integration
	server.Any("/admin/*resources", middlewares.Authorizer(jwtService, authEnforcer), middlewares.Authorization(jwtService), gin.WrapH(mux))

	server.Run(":" + port)
}
