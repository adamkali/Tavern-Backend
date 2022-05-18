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
	case "PUT":
		h.updateUserByID(w, r)
		return
	case "DELETE":
		h.deleteUserByID(w, r)
		return
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
		response.ConsumeError(err, w, http.StatusInternalServerError)
		return
	}

	// On success we can now write the response!
	w.Header().Add("content-type", "application/json")
	response.OK(users, w)
	return
}

// 	>=> POST /api/users
func (h *userHandler) post(w http.ResponseWriter, r *http.Request) {
	var response models.UserDetailedResponse

	bodyBytes, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
		return
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		response.UDRWrite(
			w,
			http.StatusUnsupportedMediaType,
			fmt.Sprintf("Application data is not application/json, got: {%s}", contentType),
			false,
		)
		return
	}

	var user models.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
		return
	}

	user.ID = string(generateUUID())
	h.Lock()
	h.store[user.ID] = user
	defer h.Unlock()

	response.OK(user, w)
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
	if len(path) != 4 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent path...",
			false,
		)
		return
	}

	userId := string(path[3])
	if len(userId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid length not lonf enough",
			false,
		)
		return
	}

	h.Lock()
	i := 0
	for _, user := range h.store {
		temp := user
		if temp.ID == userId {
			data = user
		}
		i++
	}
	h.Unlock()

	if data.ID == "" {
		response.Data = models.User{ID: userId}
		response.UDRWrite(
			w,
			http.StatusNotModified,
			"User not found.",
			false,
		)
		return
	}

	w.Header().Add("content-type", "application/json")
	response.OK(data, w)
	return
}

// PUT api/users/:id
func (h *userHandler) updateUserByID(w http.ResponseWriter, r *http.Request) {
	var response models.UserDetailedResponse

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 4 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
		return
	}

	userId := string(path[3])
	if len(userId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough",
			false,
		)
		return
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		response.UDRWrite(
			w,
			http.StatusUnsupportedMediaType,
			"Content Type needs to be application/json.",
			false,
		)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)

	var user models.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
		return
	}

	var temp models.User
	h.Lock()
	for _, usr := range h.store {
		if usr.ID == userId {
			temp = usr
		}
	}

	if temp.Username == "" {
		response.UDRWrite(
			w,
			http.StatusNotModified,
			"User was not found.",
			false,
		)
		h.Unlock()
		return
	}
	user.ID = userId
	h.store[temp.ID] = user
	defer h.Unlock()

	response.Ok(user, w)
	return
}

// DELETE /api/users/{id}
func (h *userHandler) deleteUserByID(w http.ResponseWriter, r *http.Request) {
	var response models.UserDetailedResponse
	var data models.User

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 4 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
		return
	}

	userId := string(path[3])
	if len(userId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough",
			false,
		)
		return
	}

	h.Lock()
	for _, user := range h.store {
		temp := user
		if temp.ID == userId {
			data = user
		}
	}
	if data.ID == "" {
		response.UDRWrite(
			w,
			http.StatusNotModified,
			"User was not found.",
			false,
		)
		h.Unlock()
		return
	}
	delete(h.store, userId)
	h.Unlock()

	response.OK(data, w)
	return
}
