package models

import (
	"encoding/json"
	"net/http"
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
	CharacterRace      string `json:"character_race" gorm:"column:character_race;type:varchar(32) not null`
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
