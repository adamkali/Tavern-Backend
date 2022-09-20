package main

import (
	"Tavern-Backend/controllers"
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/cors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title Tavern Profile API
// @version 1.0
// @description This is the API for the Tavern Profile Application.

// @License MIT

// @host localhost:8000
// @BasePath /api
func main() {

	// Get the arguments from the command line.
	// if the first argument is "dev" then pass true to the LoadConfiguration function.
	// else pass false.

	var config lib.Configuration
	if os.Args[1] == "dev" {
		config = lib.LoadConfiguration(true)
	} else if os.Args[1] == "prod" {
		config = lib.LoadConfiguration(false)
	} else {
		panic("Invalid Argument")
	}

	dsn := config.GetDatabaseConnectionString()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// TODO: Refactor For Scalability
	var user models.User
	var plot models.Plot
	var character models.Character
	var token models.AuthToken
	var authTA models.AuthTokenActivation
	var tags models.Tags
	var pref models.PlayerPrefrence
	db.AutoMigrate(
		&user,
		&plot,
		&character,
		&token,
		&authTA,
		&tags,
		&pref,
	)

	// Instantiate the controllers
	userController := controllers.NewUserController(db)
	authController := controllers.NewAuthController(db)
	plotController := controllers.NewPlotController(db)
	characterController := controllers.NewCharacterController(db)
	relationshipController := controllers.NewRelationshipController(db)
	// :ENDTODO

	// Create a cors middleware to allow cross-origin requests.
	// have it return the handler function.
	cors := cors.New(cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowCredentials: config.Cors.Credentials,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		AllowedMethods:   config.Cors.AllowedMethods,
	})

	// TODO: Refactor For Scalability
	entries := lib.LogEntries{}
	entries.StartLogging()
	http.Handle("/api/admin/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entries.RenderHtml(w)
	}))
	http.Handle(userController.H.AuthPath,
		cors.Handler(http.HandlerFunc(userController.H.Controller)))
	http.Handle(userController.H.AdmnAllPath,
		cors.Handler(http.HandlerFunc(userController.AdminGetAll)))
	http.Handle(userController.H.AuthAllPath,
		cors.Handler(http.HandlerFunc(userController.UserQueue)))
	http.Handle(userController.H.AuthPath+"full",
		cors.Handler(http.HandlerFunc(userController.GetAuthenticatedUser)))
	http.Handle(userController.H.AuthPath+"full/",
		cors.Handler(http.HandlerFunc(userController.AuthGetByIDFull)))
	http.Handle("/api/login",
		cors.Handler(http.HandlerFunc(authController.Login)))
	http.Handle("/api/signup",
		cors.Handler(http.HandlerFunc(authController.SignUp)))
	http.Handle("/api/verify",
		cors.Handler(http.HandlerFunc(authController.Verify)))
	http.Handle(characterController.H.AuthPath,
		cors.Handler(http.HandlerFunc(characterController.H.Controller)))
	http.Handle(characterController.H.AuthAllPath,
		cors.Handler(http.HandlerFunc(characterController.AuthGetAllCByUserID)))
	http.Handle(plotController.H.AuthPath,
		cors.Handler(http.HandlerFunc(plotController.H.Controller)))
	http.Handle(plotController.H.AuthAllPath,
		cors.Handler(http.HandlerFunc(plotController.AuthGetAllPByUserID)))
	http.Handle(relationshipController.H.AuthPath,
		cors.Handler(http.HandlerFunc(relationshipController.H.Controller)))

	http.ListenAndServe(fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort), nil)
	// :ENDTODO
}
