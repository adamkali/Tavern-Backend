package models

import (
	"Tavern-Backend/lib"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
)

type Plot struct {
	ID     string `json:"id" gorm:"column:id;type:varchar(32)"`
	Name   string `json:"plot_name" gorm:"column:name;type:varchar(128) not null"`
	Plot   string `json:"plot" grom:"column:plot;type:text not null"`
	UserID string `json:"-"`
	//	Parent User   `json:"user;omitempty" gorm:"foreignKey:UserFk;refernces:ID"`
}

type Character struct {
	ID                 string `json:"id" gorm:"column:id;type:varchar(32)"`
	Name               string `json:"character_name" gorm:"column:name;type:varchar(128) not null"`
	Backstory          string `json:"back_story" grom:"column:back_story;type:text not null"`
	Bio                string `json:"bio" grom:"column:bio;type:text not null"`
	Strength           int    `json:"strength" gorm:"column:strength;type: tinyint not null"`
	Dexterity          int    `json:"dexterity" gorm:"column:dexterity;type: tinyint not null"`
	Constitution       int    `json:"constitution" gorm:"column:constitution;type: tinyint not null"`
	Intelligence       int    `json:"intelligence" gorm:"column:intelligence;type: tinyint not null"`
	Wisdom             int    `json:"wisdom" gorm:"column:wisdom;type: tinyint not null"`
	Charisma           int    `json:"charisma" gorm:"column:charisma;type: tinyint not null"`
	CharacterClass     string `json:"character_class" grom:"column:character_class;type:varchar(64) not null"`
	CharacterLevel     int    `json:"character_level" gorm:"column:character_level;type: tinyint not null"`
	CharacterTraits    string `json:"character_traits" grom:"column:character_traits;type:text not null"`
	CharacterHitPoints string `json:"character_hit_points" gorm:"column:character_hit_points;type:varchar(32) not null"`
	UserID             string `json:"-"`
	//	Parent          User   `json:"user;omitempty" gorm:"foreignKey:UserFk;refernces:ID"`
}

type User struct {
	ID              string `json:"id" gorm:"column:id;type:varchar(32)"`
	Username        string `json:"username" gorm:"column:username;type:varchar(128) not null"`
	Bio             string `json:"bio" grom:"column:bio;type:text not null"`
	Tags            string `json:"tags" grom:"column:tags;type:text not null"`
	PlayerPrefrence string `json:"player_prefrence" gorm:"column:player_preference;type:varchar(32)"`
	//	Plots           []Plot      `json:"user_plots,omitempty" gorm:"foreignKey:UserID;refernces:ID"`
	//	Characters      []Character `json:"user_characters,omitempty" gorm:"foreignKey:UserID;refernces:ID"`
	Plots      []Plot      `json:"user_plots,omitempty"`
	Characters []Character `json:"user_characters,omitempty"`

	GroupID string `json:"group_fk,omitempty" gorm:"foreignKey:GroupID;refernces:ID"`
}

type Group struct {
	ID   string `json:"id" gorm:"column:id;type:varchar(32);primaryKey"`
	Name string `json:"group_name" gorm:"column:name;type:varchar(128) not null"`
}

type AuthToken struct {
	ID        string `json:"id" gorm:"column:id;type:varchar(32);primaryKey"`
	UserID    string `json:"user_fk,omitempty" gorm:"foreignKey:UserID;refernces:ID"`
	Username  string `json:"username" gorm:"column:username;type:varchar(32) not null"`
	UserEmail string `json:"user_email" gorm:"column:email;type:varchar(128) not null"`
	AuthHash  string `json:"auth_hash" gorm:"column:auth_hash;type:varchar(128) not null"`
	Active    bool   `json:"active" gorm:"column:active;type:tinyint(1) not null"`
}

type AuthTokenActivation struct {
	ID        string `json:"id" gorm:"column:id;type:varchar(32);primaryKey"`
	AuthID    string `json:"auth_fk,omitempty" gorm:"foreignKey:AuthID;refernces:ID"`
	AuthPin   string `json:"auth_pin" gorm:"column:auth_pin;type:varchar(8) not null"`
	AuthEmail string `json:"auth_email" gorm:"column:auth_email;type:varchar(128) not null"`
}

type Users []User

type UserDetailedResponse struct {
	Data       User   `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type UsersDetailedResponse struct {
	Data       Users  `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type GroupDetailedResponse struct {
	Data       Group  `json:"group"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type TokenDetailedResponse struct {
	Data       AuthToken `json:"data"`
	Successful bool      `json:"successful"`
	Message    string    `json:"message"`
}

type AuthRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserEmail string `json:"user_email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a AuthTokenActivation) SendRegistrationEmail(config lib.Configuration) error {
	msg := gomail.NewMessage()

	// read a file from Tavern-Backend/lib/html/Registration.html
	// and use it as the body of the email
	// replace the placeholder text with the AuthPin

	var err error = nil

	f, err := filepath.Abs("lib\\html\\Register.html")
	if err != nil {
		return err
	}
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	file_string, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	defer file.Close()
	fl := string(file_string)
	// replace the placeholder text with the AuthPin
	// send the email
	fl = strings.Replace(fl, "<<<code>>>", a.AuthPin, -1)

	msg.SetHeader("From", config.Email.Username)
	msg.SetHeader("To", a.AuthEmail)
	msg.SetHeader("Subject", "Tavern Registration")
	msg.SetBody("text/html", fl)

	d := gomail.NewDialer(config.Email.Host, config.Email.Port, config.Email.Username, config.Email.Password)
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return err
}

func (t TokenDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
	t.Successful = successful
	t.Message = message
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

func (t *TokenDetailedResponse) OK(w http.ResponseWriter, auth AuthToken) {
	t.Data = auth
	t.UDRWrite(w, http.StatusOK, "OK", true)
}

func (t *TokenDetailedResponse) ConsumeError(w http.ResponseWriter, err error) {
	t.UDRWrite(w, http.StatusInternalServerError, err.Error(), false)
}

// Make a function to take a username, password, and userID
// and then return a token
// This function will use a hash function to create a hash
// by using the username, password, and userID
func (t *AuthToken) GenerateToken(username string, password string) {
	// Create a hash function
	hash := sha256.New()
	// Write the username, password, and userID to the hash function
	hash.Write([]byte(username))
	hash.Write([]byte(password))
	// Get the hash value
	hashValue := hash.Sum(nil)
	// Convert the hash value to a string
	hashString := hex.EncodeToString(hashValue)
	// Set the AuthToken struct values
	t.Username = username
	t.UserEmail = password
	t.AuthHash = hashString
}

func (t *AuthToken) VerifyToken(username string, password string) bool {
	// Create a hash function
	hash := sha256.New()
	// Write the username, password, and userID to the hash function
	hash.Write([]byte(username))
	hash.Write([]byte(password))
	// Get the hash value
	hashValue := hash.Sum(nil)
	// Convert the hash value to a string
	hashString := hex.EncodeToString(hashValue)
	// Compare the hashString to the AuthToken struct's AuthHash
	return hashString == t.AuthHash
}

func (t *AuthToken) GenAuth(username string, password string) string {
	// Create a hash function
	hash := sha256.New()
	// Write the username, password, and to the hash function
	hash.Write([]byte(username))
	hash.Write([]byte(password))
	// Get the hash value
	hashValue := hash.Sum(nil)
	// Convert the hash value to a string
	hashString := hex.EncodeToString(hashValue)
	return hashString
}

func (t *AuthToken) ValidatePassword(password string) bool {
	// check that the password is valid by checking the following rules
	// 1. Password must be at least 8 characters long
	// 2. Password must contain at least one number
	// 3. Password must contain at least one uppercase letter
	// 4. Password must contain at least one lowercase letter
	// 5. Password must contain at least one special character (!@#$%^&*()_+)
	// 6. Password must not contain any spaces
	// 7. Password must not contain any of the following characters: /'^(){}|:"<>?`~;[]\=-,
	//     ( this is to sanitize the password for the database )

	// Check that the password is at least 8 characters long
	if len(password) < 8 {
		return false
	}
	// Check that the password contains at least one number
	if !strings.ContainsAny(password, "0123456789") {
		return false
	}
	// Check that the password contains at least one uppercase letter
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}
	// Check that the password contains at least one lowercase letter
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}
	// Check that the password contains at least one special character (!@#$%^&*()_+)
	if !strings.ContainsAny(password, "!@#$%^&*()_+") {
		return false
	}
	// Check that the password does not contain any spaces
	if strings.ContainsAny(password, " ") {
		return false
	}
	// Check that the password does not contain any of the following characters: /'^(){}|:"<>?`~;[]\=-,
	if strings.ContainsAny(password, "/'^(){}|:\"<>?`~;[]\\=-,") {
		return false
	}
	return true
}

// UDRWrite(w http.ResponseWriter, code int, message string)
//
// This is a method to diplay any text in detailed responses
// To the http.ResponseWriter. It will dump the message into
// The json displayed the client

func (u UserDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
	u.Successful = successful
	u.Message = message
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

// UDRWrite(w http.ResponseWriter, code int, message string)
//
// This is a method to diplay any text in detailed responses
// To the http.ResponseWriter. It will dump the message into
// The json displayed the client
func (u UsersDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
	u.Successful = successful
	u.Message = message
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

// UDRWrite(w http.ResponseWriter, code int, message string)
//
// This is a method to diplay any text in detailed responses
// To the http.ResponseWriter. It will dump the message into
// The json displayed the client
func (u GroupDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
	u.Successful = successful
	u.Message = message
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

// ConumeError(err error, w http.ResponseWriter, code int)
//
// An error cosumer made to make the server requests as client
// friendly as possible.
//
// params:
// 	err Error -> This is any error that can be produced by go
//	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
//	code int -> This is designed to contian http.StatusOK or
//		any of the http statuses.
func (u UserDetailedResponse) ConsumeError(err error, w http.ResponseWriter, code int) {
	u.Message = err.Error()
	u.Successful = false
	u.Data = User{}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

// OK(user User, w http.ResponseWriter)
//
// A friendly status Ok writter to the web console to write the data
// set that the request was successful, and set the that the function
// is ready to be returned.
//
// params:
// 	user User -> This is the main data that we want to send back
//		to who ever is requesting the data.
// 	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
func (u UserDetailedResponse) OK(user User, w http.ResponseWriter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = user
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// ConumeError(err error, w http.ResponseWriter, code int)
//
// An error cosumer made to make the server requests as client
// friendly as possible.
//
// params:
// 	err Error -> This is any error that can be produced by go
//	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
//	code int -> This is designed to contian http.StatusOK or
//		any of the http statuses.
func (u UsersDetailedResponse) ConsumeError(err error, w http.ResponseWriter, code int) {
	u.Message = err.Error()
	u.Successful = false
	u.Data = nil
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

// OK(user User, w http.ResponseWriter)
//
// A friendly status Ok writter to the web console to write the data
// set that the request was successful, and set the that the function
// is ready to be returned.
//
// params:
// 	user Users -> This is the main data that we want to send back
//		to who ever is requesting the data.
// 	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
func (u UsersDetailedResponse) OK(user Users, w http.ResponseWriter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = user
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// ConumeError(err error, w http.ResponseWriter, code int)
//
// An error cosumer made to make the server requests as client
// friendly as possible.
//
// params:
// 	err Error -> This is any error that can be produced by go
//	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
//	code int -> This is designed to contian http.StatusOK or
//		any of the http statuses.
func (u GroupDetailedResponse) ConsumeError(err error, w http.ResponseWriter, code int) {
	u.Message = err.Error()
	u.Successful = false
	u.Data = Group{}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

// OK(user Group, w http.ResponseWriter)
//
// A friendly status Ok writter to the web console to write the data
// set that the request was successful, and set the that the function
// is ready to be returned.
//
// params:
// 	user Group -> This is the main data that we want to send back
//		to who ever is requesting the data.
// 	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
func (u GroupDetailedResponse) OK(user Group, w http.ResponseWriter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = user
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
