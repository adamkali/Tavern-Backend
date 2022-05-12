package main

import (
	"Tavern-Backend/controllers"
	"fmt"
	"net/http"
)

// Make a handler for response #1
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

// Make the server proper.
func main() {
	h := controllers.NewUserHandler() //#2

	http.HandleFunc("/api", handler)

	// ALLOWED METHODS GET, POST
	http.HandleFunc("/api/users", h.Users) // #2

	// ALLOWED METHODS GET, PUT, DELETE
	http.HandleFunc("/api/users/", h.User)

	// Handle errors // #2
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
