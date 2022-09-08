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

// make a post enforcement function it should take
// http.ResponseWriter, *http.Request, and a Generic T that can either be
// a UserDetailedResponse or a UsersDetailedResponse or an AuthTokenDetailedResponse
func (h *userHandler) methodEnforce(r *http.Request, m string) bool {
	return r.Method == strings.ToUpper(m)
}

func (h *userHandler) hashEnfoce(r *http.Request) bool {
	hash := r.Header.Get("AuthorizationHash")
	userID := r.Header.Get("UserID")

	var token models.AuthToken

	hashResult := h.db.Where("auth_hash = ?", hash).First(&token)
	if hashResult.Error != nil {
		return false
	}
	return token.UserID == userID
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

// 	AuthToken Handler /api/auth
func (h *userHandler) AuthToken(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		h.Login(w, r)
		return
	case "DELETE":
		h.deleteAuthToken(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
		return
	}
}

//    Verification Handler /api/signup
func (h *userHandler) Verify(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		h.Signup(w, r)
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

// === === === === === === === === === === === ===
//	>=> AUTHTOKEN/VERIFICATION CONTROLLER <=<
//   	>=> `/api/auth/<path>` 		      <=<
// === === === === === === === === === === === ===

// >=> POST /api/login
// login with username and password
func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	// get the username and password from the request
	var response models.TokenDetailedResponse
	var data models.AuthToken
	var req models.LoginRequest

	good := h.methodEnforce(r, "POST")
	if !good {
		response.UDRWrite(
			w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
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
	if err != nil {
		response.ConsumeError(
			w,
			err,
		)
		return
	}

	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		response.ConsumeError(
			w,
			err,
		)
		return
	}

	// get the username and password from the request
	hash := data.GenAuth(req.Username, req.Password)

	if hash == "" {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Invalid username or password",
			false,
		)
		return
	}

	// get check if the user exists in the AuthToken table to see
	// that there is a userID associated with the hash
	result := h.db.Where("auth_hash = ?", hash).Find(&data)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		return
	}

	// since there is a UserID associated with the hash in the AuthToken table, we can get the
	// userId form result add it to data and write a TokenDetailedResonse to the client
	response.OK(w, data)
	return
}

// >=> DELETE /api/login
// this will be deleting the credentials from the AuthToken table
func (h *userHandler) deleteAuthToken(w http.ResponseWriter, r *http.Request) {
	var response models.TokenDetailedResponse
	var data models.AuthToken

	// get the hash from the request
	hash := r.Header.Get("Authorization")
	if hash == "" {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Invalid hash",
			false,
		)
		return
	}

	// delete the hash from the AuthToken table
	result := h.db.Delete(&models.AuthToken{}, hash)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		return
	}

	response.OK(w, data)
	return
}

// >=> POST /api/signup
// create a new user
func (h *userHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var response models.TokenDetailedResponse
	var data models.AuthToken
	var req models.AuthRequest

	good := h.methodEnforce(r, "POST")
	if !good {
		response.UDRWrite(
			w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
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
	if err != nil {
		response.ConsumeError(
			w,
			err,
		)
		return
	}

	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		response.ConsumeError(
			w,
			err,
		)
		return
	}

	if err != nil {
		response.ConsumeError(
			w,
			err,
		)
		return
	}

	// get the username and password from the request
	username := req.Username
	password := req.Password
	useremail := req.UserEmail

	// check if the username is already taken
	result := h.db.Where("username = ?", username).Find(&data)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		return
	}

	//check that the username is not longer than 32 characters
	if len(username) > 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Username is too long",
			false,
		)
		return
	}

	// check that the result is empty (no user found) if it is not
	// then the username is taken and throw and return a 409 Conflict
	if result.RowsAffected != 0 {
		response.UDRWrite(
			w,
			http.StatusConflict,
			"Username already taken",
			false,
		)
		return
	}

	// check that the password length is greater than 8,
	// has a lowercase letter, an uppercase letter, a number and a special character
	if !data.ValidatePassword(password) {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Password must be at least 8 characters long and contain at least one lowercase letter, one uppercase letter, one number and one special character (!@#$%^&*()_+-=)",
			false,
		)
		return
	}

	hash := data.GenAuth(username, password)

	if hash == "" {
		response.UDRWrite(
			w,
			http.StatusInternalServerError,
			"Invalid username or password",
			false,
		)
		return
	}

	// get check if the user exists in the AuthToken table to see
	// that there is a userID associated with the hash
	result = h.db.Where("auth_hash = ?", hash).Find(&data)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		return
	}

	if result.RowsAffected != 0 && data.Username == username {
		// since there is a UserID associated with the hash in the AuthToken table
		// return a 409 Conflict
		response.UDRWrite(
			w,
			http.StatusConflict,
			"Username already taken",
			false,
		)
		return

	} else {
		// create a new user and make a new AuthToken for the user
		// Make a new user
		userId := generateUUID()
		user := models.User{
			Username:        username,
			ID:              userId,
			Bio:             "",
			Tags:            "",
			PlayerPrefrence: "",
			Plots:           nil,
			Characters:      nil,
			GroupID:         generateUUID(),
		}
		userResult := h.db.Create(&user)
		if userResult.Error != nil {
			response.ConsumeError(
				w,
				userResult.Error,
			)
			return
		}

		// save the new user to the database
		data = models.AuthToken{
			Username:  username,
			AuthHash:  hash,
			UserID:    userId,
			UserEmail: useremail,
			ID:        generateUUID(),
		}

		// print data to the console
		result = h.db.Create(&data) // this will create a new user and a new hash
		if result.Error != nil {
			response.ConsumeError(
				w,
				result.Error,
			)
			return
		}

		// since there is a UserID associated with the hash in the AuthToken table, we can get the
		// userId form result add it to data and write a TokenDetailedResonse to the client
		response.OK(w, data)
		return
	}
}
