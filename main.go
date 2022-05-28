package main

import (
	"Tavern-Backend/controllers"
	"Tavern-Backend/models"
	"fmt"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/rs/cors"
)

// Make a handler for response #1
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

// Make the server proper.
func main() {
	dsn := "root:Sierra&Adam4ever@tcp(127.0.0.1:3306)/taverndatabase?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var user models.User
	var plot models.Plot
	var character models.Character
	db.AutoMigrate(&user, &plot, &character)
	userH := controllers.NewUserHandler(*db) //#2
	characterH := controllers.NewCharacterHandler(*db)
	plotH := controllers.NewPlotHandler(*db)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:19000",
			"http://localhost:19006",
		},
		AllowCredentials: true,
		Debug: true,
	})

	http.HandleFunc("/api", handler)

	// USER PAGES
	http.HandleFunc("/api/users", userH.Users) // #2
	http.HandleFunc("/api/users/", userH.User)

	// CHARACTER PAGES
	http.HandleFunc("/api/character/", characterH.Character)
	http.HandleFunc("/api/characters/userId/", characterH.Characters)

	// PLOTS PAGES
	http.HandleFunc("/api/plot/", plotH.Plot)
	http.HandleFunc("/api/plots/userId/", plotH.Plots)

	// Handle errors // #2
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
