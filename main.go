package main

import (
	"Tavern-Backend/controllers"
	"Tavern-Backend/models"
	"net/http"

	"github.com/rs/cors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
	var token models.AuthToken
	db.AutoMigrate(&user, &plot, &character, &token)
	userH := controllers.NewUserHandler(*db) //#2
	characterH := controllers.NewCharacterHandler(*db)
	plotH := controllers.NewPlotHandler(*db)

	// Create a cors middleware to allow cross-origin requests.
	// have it return the handler function.
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{
			"*",
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions},
	})
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Write([]byte("Hello, World!"))
		return
	}))
	http.Handle("/api/users", cors.Handler(http.HandlerFunc(userH.Users)))
	http.Handle("/api/characters/userId/", cors.Handler(http.HandlerFunc(characterH.Characters)))
	http.Handle("/api/plots/userId/", cors.Handler(http.HandlerFunc(plotH.Plots)))
	http.Handle("/api/characters/", cors.Handler(http.HandlerFunc(characterH.Character)))
	http.Handle("/api/plots/", cors.Handler(http.HandlerFunc(plotH.Plot)))
	http.Handle("/api/users/", cors.Handler(http.HandlerFunc(userH.User)))
	http.ListenAndServe(":8000", nil)
}
