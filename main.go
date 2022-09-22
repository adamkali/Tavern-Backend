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
	user := models.User{}
	plot := models.Plot{}
	character := models.Character{}
	token := models.AuthToken{}
	authTA := models.AuthTokenActivation{}
	tags := models.Tag{}
	pref := models.PlayerPrefrence{}
	role := models.Role{}
	relationship := models.UserRelationship{}
	rels := models.Relationship{}

	err = db.AutoMigrate(
		&user, &plot, &character,
		&token, &authTA, &tags,
		&pref, &role, &relationship,
		&rels)
	if err != nil {
		panic(err)
	}

	// Instantiate the controllers
	userController := controllers.NewUserController(db)
	authController := controllers.NewAuthController(db, models.AuthEmailConfiglette(config.Email))
	plotController := controllers.NewPlotController(db)
	characterController := controllers.NewCharacterController(db)
	relationshipController := controllers.NewRelationshipController(db)
	tagController := controllers.NewTagController(db)
	prefController := controllers.NewPlayerPrefrenceController(db)
	roleController := controllers.NewRoleController(db)
	relsController := controllers.NewRelsController(db)
	// :ENDTODO

	// Create a cors middleware to allow cross-origin requests.
	// have it return the handler function.
	cors := cors.New(cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowCredentials: config.Cors.Credentials,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		AllowedMethods:   config.Cors.AllowedMethods,
	})

	// #region API Routes
	entries := lib.LogEntries{}
	entries.StartLogging()
	http.Handle("/api/admin/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entries.RenderHtml(w)
	}))
	http.Handle(userController.H.AuthPath,
		cors.Handler(http.HandlerFunc(
			userController.H.Sanitize(userController.H.Controller))))
	http.Handle(userController.H.AdmnAllPath,
		cors.Handler(http.HandlerFunc(
			userController.H.Sanitize(userController.AdminGetAll))))
	http.Handle(userController.H.AuthAllPath,
		cors.Handler(http.HandlerFunc(
			userController.H.Sanitize(userController.UserQueue))))
	http.Handle(userController.H.AuthPath+"full",
		cors.Handler(http.HandlerFunc(
			userController.H.Sanitize(userController.GetAuthenticatedUser))))
	http.Handle(userController.H.AuthPath+"full/",
		cors.Handler(http.HandlerFunc(
			userController.H.Sanitize(userController.AuthGetByIDFull))))
	http.Handle("/api/login",
		cors.Handler(http.HandlerFunc(
			authController.H.Sanitize(authController.Login))))
	http.Handle("/api/signup",
		cors.Handler(http.HandlerFunc(
			authController.H.Sanitize(authController.SignUp))))
	http.Handle("/api/verify",
		cors.Handler(http.HandlerFunc(
			authController.H.Sanitize(authController.Verify))))
	http.Handle(characterController.H.AuthPath,
		cors.Handler(http.HandlerFunc(
			characterController.H.Sanitize(characterController.H.Controller))))
	http.Handle(characterController.H.AuthAllPath,
		cors.Handler(http.HandlerFunc(
			characterController.H.Sanitize(characterController.AuthGetAllCByUserID))))
	http.Handle(plotController.H.AuthPath,
		cors.Handler(http.HandlerFunc(
			plotController.H.Sanitize(plotController.H.Controller))))
	http.Handle(plotController.H.AuthAllPath,
		cors.Handler(http.HandlerFunc(
			plotController.H.Sanitize(plotController.AuthGetAllPByUserID))))
	http.Handle(relationshipController.H.AuthPath,
		cors.Handler(http.HandlerFunc(
			relationshipController.H.Sanitize(relationshipController.H.Controller))))
	http.Handle(tagController.H.AuthPath,
		cors.Handler(http.HandlerFunc(
			tagController.H.Sanitize(tagController.H.Controller))))
	http.Handle(tagController.H.AuthPath+"add",
		cors.Handler(http.HandlerFunc(
			tagController.H.Sanitize(tagController.AuthPostTagToUser))))
	http.Handle(tagController.H.AuthAllPath+"add",
		cors.Handler(http.HandlerFunc(
			tagController.H.Sanitize(tagController.AuthPostTagsToUser))))
	http.Handle(tagController.H.AuthAllPath,
		cors.Handler(http.HandlerFunc(
			tagController.H.Sanitize(tagController.AuthGetTags))))
	http.Handle(prefController.H.AuthPath,
		cors.Handler(http.HandlerFunc(
			prefController.H.Sanitize(prefController.H.Controller))))
	http.Handle(prefController.H.AuthPath+"add",
		cors.Handler(http.HandlerFunc(
			prefController.H.Sanitize(prefController.AuthPostPlayerPrefrenceToUser))))
	http.Handle(prefController.H.AuthAllPath,
		cors.Handler(http.HandlerFunc(
			prefController.H.Sanitize(prefController.AuthGetPrefrences))))
	http.Handle(roleController.H.AdmnPath,
		cors.Handler(http.HandlerFunc(
			roleController.H.Sanitize(roleController.AdminChangeRole))))
	http.Handle(relsController.H.AuthPath,
		cors.Handler(http.HandlerFunc(
			relsController.H.Sanitize(relsController.H.Controller))))
	http.Handle(relsController.H.AuthAllPath,
		cors.Handler(http.HandlerFunc(
			relsController.H.Sanitize(relsController.AuthGetRelationships))))
	http.ListenAndServe(fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort), nil)
	// #endregion
}
