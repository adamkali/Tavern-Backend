package controllers

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Make storage for your data #2
type userHandler struct {
	db     *gorm.DB
	config *lib.Configuration
}

func NewUserHandler(database gorm.DB, config lib.Configuration) *userHandler {
	return &userHandler{
		db:     &database,
		config: &config,
	}
}

// make a post enforcement function it should take
// http.ResponseWriter, *http.Request, and a Generic T that can either be
// a UserDetailedResponse or a UsersDetailedResponse or an AuthTokenDetailedResponse
func (h *userHandler) methodEnforce(r *http.Request, m string) bool {
	return r.Method == strings.ToUpper(m)
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

	var logger lib.LogEntryObject
	startTime := time.Now()

	var prep []models.User
	var users models.Users

	result := h.db.Preload("Characters").Preload("Plots").Find(&prep)

	// var response models.UsersDetailedResponse
	var response models.DetailedResponse[models.Users]

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		logger.TimeTaken = time.Since(startTime).Milliseconds()
		logger.Log(r, http.StatusInternalServerError, 0, result.Error.Error())
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

// PUT /api/users/:id
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

// GET /api/auth/users/:id
// Get the User for the Profile screen
// 
func (h *userHandler) GetAuthUserByID(w http.ResponseWriter, r *http.Request) {
	var response models.UserDetailedResponse
	var token models.AuthToken
	var data models.User

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 5 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent path...",
			false,
		)
		return
	}

	// check authentification
	auth := r.Header.Get("Authorization")
	if auth == "" {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Authorization header not found",
			false,
		)
		return
	}

	tokenString := strings.Split(auth, "Bearer ")
	if len(tokenString) != 2 {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Authorization header not found",
			false,
		)
		return
	}

	token.AuthHash = tokenString[1]
	res := h.db.Where("auth_hash = ?", token.AuthHash).First(&token)
	if res.Error != nil {
		response.ConsumeError(
			res.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	userId := string(path[4])
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

// PUT /api/auth/users/:id
// 
func (h *userHandler) AuthUpdateUserByID(w http.ResponseWriter, r *http.Request) {
	var response models.UserDetailedResponse
	var token models.AuthToken

	// check authentification from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Authorization header not found.",
			false,
		)
		return
	}

	// check if authentification is valid
	tokenString := strings.Split(authHeader, "Bearer ")
	if len(tokenString) != 2 {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Authorization header not valid.",
			false,
		)
		return
	}

	// check if token is valid
	result := h.db.Where("auth = ?", tokenString[1]).First(&token)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	} else if result.RowsAffected == 0 {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Authorization token not found.",
			false,
		)
		return
	}

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

	// Verify that the user is the same as the one in the token
	if userId != token.UserID {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"The User you are trying to update is not you, Please reconnect and try again.",
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
			err,
			w,
			http.StatusBadRequest,
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

	// update the user in the database with the new data user
	old := &models.User{ID: userId}
	result = h.db.First(old)
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

	// return ok
	response.OK(user, w)
	return
}

// === === === === === === === === === === === ===
//	>=> AUTHTOKEN/VERIFICATION CONTROLLER <=<
//   	>=> `/api/auth/<path>` 		      <=<
// === === === === === === === === === === === ===

// >=> POST /api/login
// login with username and password
// 
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

// >=> DELETE /api/auth/login
// this will be deleting the credentials from the AuthToken table
//  TODO: make accesible by the user and hav it delete the user as well.
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
// 
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

		// save the new user to the database
		data = models.AuthToken{
			Username:  username,
			AuthHash:  hash,
			UserID:    userId,
			UserEmail: useremail,
			ID:        generateUUID(),
		}

		// now creat the AuthRegister
		authRegister := models.AuthTokenActivation{
			ID:        generateUUID(),
			AuthID:    data.ID,
			AuthPin:   generatePin(),
			AuthEmail: useremail,
		}
		// save the AuthRegister to the database
		err = authRegister.SendRegistrationEmail(*h.config)
		if err != nil {
			response.ConsumeError(
				w,
				err,
			)
			return
		}

		userResult := h.db.Create(&user)
		if userResult.Error != nil {
			response.ConsumeError(
				w,
				userResult.Error,
			)
			return
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

		authRegisterResult := h.db.Create(&authRegister)
		if authRegisterResult.Error != nil {
			response.ConsumeError(
				w,
				authRegisterResult.Error,
			)
			return
		}
		fmt.Printf("AuthRegister: %+v", authRegister)

		// since there is a UserID associated with the hash in the AuthToken table, we can get the
		// userId form result add it to data and write a TokenDetailedResonse to the client
		response.OK(w, data)
		return
	}
}

// GET /api/activate/{pin}
// 
func (h *userHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var response models.TokenDetailedResponse
	var active models.AuthTokenActivation
	var data models.AuthToken

	// enforce get
	good := h.methodEnforce(r, "GET")
	if !good {
		response.UDRWrite(
			w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
			false,
		)
		return
	}
	// get the pin from the request from http://localhost:8080/api/activate/{pin}
	pin := strings.Split(r.URL.Path, "/")[3]
	// check if the pin is in the AuthTokenActivation table
	result := h.db.Where("auth_pin = ?", pin).Find(&active)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		return
	}
	// check if that the authID is in the AuthToken table
	result = h.db.Where("id = ?", active.AuthID).Find(&data)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		return
	}
	// check if the user is already activated
	if data.Active {
		response.UDRWrite(
			w,
			http.StatusConflict,
			"User already activated",
			false,
		)
		return
	}
	// update the auth token to be active
	result = h.db.Model(&data).Where("id = ?", data.ID).Update("active", true)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		return
	}
	// delete the AuthTokenActivation
	result = h.db.Where("id = ?", active.ID).Delete(&active)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		return
	}
	// write a TokenDetailedResponse to the client
	response.OK(w, data)
	return
}

// === === === === === === === === === === === ===
// >=> 			ENUM CONTROLLER 			   <=<
// >=> 			`/api/auth/enum` 		       <=<
// === === === === === === === === === === === ===

// GET /api/auth/enum/tags
// 
func (h *userHandler) GetAuthTags(w http.ResponseWriter, r *http.Request) {

	var logger lib.LogEntryObject
	startTime := time.Now()

	var response models.TagsDetailedResponse
	var data []models.Tags
	var token models.AuthToken

	// enforce get
	good := h.methodEnforce(r, "GET")
	if !good {
		response.UDRWrite(
			w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
			false,
		)
		logger.TimeTaken = time.Since(startTime).Milliseconds()
		logger.Log(r, http.StatusMethodNotAllowed, 0, "Method not allowed")
		return
	}
	// get the token from the request
	auth_hash := r.Header.Get("Authorization")
	tokenString := strings.Split(auth_hash, "Bearer ")
	if len(tokenString) != 2 {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Authorization header not found",
			false,
		)
		logger.TimeTaken = time.Since(startTime).Milliseconds()
		logger.Log(r, http.StatusUnauthorized, 0, "Authorization header not found")
		return
	}
	result := h.db.Where("auth_hash = ?", auth_hash).Find(&token)
	if result.Error != nil || result.RowsAffected == 0 {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Unauthorized",
			false,
		)
		logger.TimeTaken = time.Since(startTime).Milliseconds()
		logger.Log(r, http.StatusUnauthorized, 0, "Unauthorized", result.Error)
		return
	}

	// get the tags from the database
	result = h.db.Find(&data)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		logger.TimeTaken = time.Since(startTime).Milliseconds()
		logger.Log(r, http.StatusInternalServerError, 0, "Internal Server Error", result.Error)
		return
	}

	// write the tags to the client
	response.OK(w, data)
	logger.TimeTaken = time.Since(startTime).Milliseconds()
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(data)
	logger.Log(r, http.StatusOK, float64(network.Len()/1000), "OK")
	return
}

// GET /api/auth/enum/preferences
// 
func (h *userHandler) GetAuthPreferences(w http.ResponseWriter, r *http.Request) {

	var logger lib.LogEntryObject
	time := time.Now()

	var response models.PlayerPrefrenceDetailedResponse
	var data []models.PlayerPrefrence
	var token models.AuthToken

	// enforce get
	good := h.methodEnforce(r, "GET")
	if !good {
		response.UDRWrite(
			w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
			false,
		)
		logger.TimeTaken = int(time.Since(time).Milliseconds())
		logger.Log(r, http.StatusMethodNotAllowed, 0, "Method not allowed")
		return
	}
	// get the token from the request
	auth_hash := r.Header.Get("Authorization")
	tokenString := strings.Split(auth_hash, "Bearer ")
	if len(tokenString) != 2 {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Authorization header not found",
			false,
		)
		logger.TimeTaken = int(time.Since(time).Milliseconds())
		logger.Log(r, http.StatusUnauthorized, 0, "Authorization header not found")
		return
	}
	// check if the token is in the AuthToken table
	result := h.db.Where("auth_hash = ?", auth_hash).Find(&token)
	if result.Error != nil || result.RowsAffected == 0 {
		response.UDRWrite(
			w,
			http.StatusUnauthorized,
			"Unauthorized",
			false,
		)
		logger.TimeTaken = int(time.Since(time).Milliseconds())
		logger.Log(r, http.StatusUnauthorized, 0, "Unauthorized", result.Error)
		return
	}

	// get the preferences from the database
	result = h.db.Find(&data)
	if result.Error != nil {
		response.ConsumeError(
			w,
			result.Error,
		)
		logger.TimeTaken = int(time.Since(time).Milliseconds())
		logger.Log(r, http.StatusInternalServerError, 0, "Internal Server Error", result.Error)
		return
	}

	// write the preferences to the client
	response.OK(w, data)
	logger.TimeTaken = int(time.Since(time).Milliseconds())
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(data)
	logger.Log(r, http.StatusOK, float64(network.Len()/1000), "OK")
	return
}
