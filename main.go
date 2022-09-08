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

// Make the server proper.
func main() {

	// Get the arguments from the command line.
	// if the first argument is "dev" then pass true to the LoadConfiguration function.
	// else pass false.

	var config lib.Configuration
	println(os.Args[1])
	if os.Args[1] == "dev" {
		config = lib.LoadConfiguration(true)
	} else if len(os.Args) == 1 {
		config = lib.LoadConfiguration(false)
	} else {
		panic("Invalid Argument")
	}

	dsn := config.GetDatabaseConnectionString()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var user models.User
	var plot models.Plot
	var character models.Character
	var token models.AuthToken
	db.AutoMigrate(&user, &plot, &character, &token)
	userH := controllers.NewUserHandler(*db) //#2
	characterH := controllers.NewCharacterHandler(*db)
	plotH := controllers.NewPlotHandler(*db)

	// Create a cors middleware to allow cross-origin requests.
	// have it return the handler function.
	fmt.Print("Starting Cors Config\n")
	cors := cors.New(cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowCredentials: config.Cors.Credentials,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		AllowedMethods:   config.Cors.AllowedMethods,
	})
	fmt.Print("Starting Server\n")
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Write([]byte("Hello, World!"))
	}))
	http.Handle("/api/signup", cors.Handler(http.HandlerFunc(userH.Signup)))
	http.Handle("/api/login", cors.Handler(http.HandlerFunc(userH.Login)))
	http.Handle("/api/users", cors.Handler(http.HandlerFunc(userH.Users)))
	http.Handle("/api/characters/userId/", cors.Handler(http.HandlerFunc(characterH.Characters)))
	http.Handle("/api/plots/userId/", cors.Handler(http.HandlerFunc(plotH.Plots)))
	http.Handle("/api/characters/", cors.Handler(http.HandlerFunc(characterH.Character)))
	http.Handle("/api/plots/", cors.Handler(http.HandlerFunc(plotH.Plot)))
	http.Handle("/api/users/", cors.Handler(http.HandlerFunc(userH.User)))
	http.ListenAndServe(fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort), nil)
}
