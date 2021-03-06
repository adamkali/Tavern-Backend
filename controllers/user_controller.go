package controllers

import (
	"Tavern-Backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Make storage for your data #2
type userHandler struct {
	db *gorm.DB
}

func NewUserHandler(database gorm.DB) *userHandler {
	return &userHandler{
		db: &database,
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

func (h *userHandler) Relationships(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		h.getRelationships(w, r)
		return
	case "POST":
		h.postRelationships(w, r)
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
	var prep []models.User
	var users models.Users
	
	result := h.db.Preload("Characters").Preload("Plots").Find(&prep)

	var response models.UsersDetailedResponse

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	for _, user := range prep {
		users = append(users, user)
	}

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

	result := h.db.Create(&user)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

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
			"Guid length not long enough",
			false,
		)
		return
	}

	data.ID = userId
	result := h.db.Preload(clause.Associations).First(&data)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

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

	old := &models.User{ID: userId}
	result := h.db.First(old)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}
	user.ID = userId
	user.ID = old.GroupID
	result = h.db.Save(&user)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	response.OK(user, w)
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

	result := h.db.Delete(&models.User{}, userId)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}
	response.OK(data, w)
	return
}


//=== === === === === === === === === === === === === === === === === ===
//	>=> USERRELATION CONTROLLER ENDPOINTS 
//      `/api/users/user/:id` <=<
//=== === === === === === === === === === === === === === === === === ===

// GET /api/users/user/:id
// get all relations for :id
func (h *userHandler) getUserRelations(w http.ResponseWriter, r *http.Request) {
	var response models.UserRelationshipsDetailedResponse
	var data models.UserRelationships
	var result models.UserRelationships

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 5 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
		return
	}

	userId := string(path[4])
	if len(userId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough",
			false,
		)
		return
	}

	other := string(path[4])
	if len(other) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough",
			false,
		)
		return
	}

	// get all relations for the userId
	result := h.db.Where("self = ?", userId).Find(&data)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	if data.len() !>= 0 {
		response.Data = models.UserRelationships{}
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

// POST /api/users/user/:id/relate/:other
// create a new relation between :id and :other
// this will create a new relation if one does not exist
// there will be a post request to create a new relation
// There will be two users in the body of the request.
// The first user will be the user that is creating the relation
// Its ID will be the :id and self in the UserRelationship
// The second user will be the User that is being related to
// Its ID will be the :other and Other in the UserRelationship
func (h *userHandler) createUserRelations(w http.ResponseWriter, r *http.Request) {
	var response models.UserRelationshipDetailedResponse
	var data models.UserRelationship
	jsonBytes, err := json.Unmarshal(r.Body)

	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
		return
	}

	var self models.User
	var other models.User
	

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 7 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
		return
	}

	userId := string(path[6])
	otherId := string(path[4])
	if len(userId) != 32 || len(otherId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough for either one of the users.",
			false,
		)
		return
	}

	// get all relations for the userId
	result := h.db.Where("self = ?", userId).Find(&data)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}


}