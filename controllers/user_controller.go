package controllers

import (
	"Tavern-Backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// Make storage for your data #2
type userHandler struct {
	sync.Mutex
	store map[string]models.User
}

func NewUserHandler() *userHandler {
	return &userHandler{
		store: map[string]models.User{},
	}
}

/*
=== === === === === === === === === === === === === === === === === === ===
		>=> USERS CONTROLLER PAGES <=<
=== === === === === === === === === === === === === === === === === === ===
*/

func (h *userHandler) User(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getUserByID(w, r)
		return
	/*case "PUT"
		// h.updateUserByID(w, r)
	case "DELETE":
		// h.deleteUserByID(w, r)*/
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}
}

func (h *userHandler) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
		return
	}
}

/*
=== === === === === === === === === === === === === === === === === === ===
		>=> USERS CONTROLLER ENDPOINTS `/api/users` <=<
=== === === === === === === === === === === === === === === === === === ===
*/

// Create a getter from the userHandler #2
// 	>=> GET /api/users
func (h *userHandler) get(w http.ResponseWriter, r *http.Request) {
	users := make(models.Users, len(h.store))

	var response models.UsersDetailedResponse

	h.Lock()
	i := 0
	for _, user := range h.store {
		users[i] = user
		i++
	}
	h.Unlock()

	_, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response.ConsumeError(err))
		return
	}

	// On success we can now write the response!
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response.OK(users))
	return
}

// 	>=> POST /api/users
func (h *userHandler) post(w http.ResponseWriter, r *http.Request) {
	var response models.UserDetailedResponse

	bodyBytes, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response.ConsumeError(err))
		return
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		response.Data = models.User{}
		response.Successful = false
		response.Message = fmt.Sprintf("Application data is not application/json, got: {%s}", contentType)
		return
	}

	var user models.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.ConsumeError(err))
		return
	}

	user.ID = string(generateUUID())
	h.Lock()
	h.store[user.ID] = user
	defer h.Unlock()

	w.Write(response.OK(user))
	return
}

/*
=== === === === === === === === === === === === === === === === === === ===
	>=> USERS CONTROLLER ENDPOINTS `/api/users/:id` <=<
=== === === === === === === === === === === === === === === === === === ===
*/

// GET /api/users/:id
func (h *userHandler) getUserByID(w http.ResponseWriter, r *http.Request) {
	var response models.UserDetailedResponse
	var data models.User

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 3 {
		response.Message = "Insufficent path..."
		response.Successful = false
		response.Data = models.User{}
		w.WriteHeader(http.StatusNotFound)
		w.Write(response.UDRWrite())
	}

	userId := string(path[2])
	if len(userId) != 32 {
		response.Message = "Guid length not long enough"
		response.Successful = false
		response.Data = models.User{}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.UDRWrite())
	}

	h.Lock()
	for _, user := range h.store {
		temp := user
		if temp.ID == userId {
			data = user
		}
	}

	if data.ID == "" {
		response.Message = "Could not find user with that Guid..."
		response.Successful = false
		response.Data = models.User{ID: userId}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.UDRWrite())
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response.OK(data))
	return
}
