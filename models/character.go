package models

import (
	"encoding/json"
	"net/http"
)

type Characters []Character

type CharacterDetailedResponse struct {
	Data       Character `json:"data"`
	Successful bool      `json:"successful"`
	Message    string    `json:"message"`
}

type CharactersDetailedResponse struct {
	Data       Characters `json:"data"`
	Successful bool       `json:"successful"`
	Message    string     `json:"message"`
}

// UDRWrite(w http.ResponseWriter, code int, message string)
//
// This is a method to diplay any text in detailed responses
// To the http.ResponseWriter. It will dump the message into
// The json displayed the client

func (u CharacterDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
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
func (u CharactersDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
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
func (u CharacterDetailedResponse) ConsumeError(err error, w http.ResponseWriter, code int) {
	u.Message = err.Error()
	u.Successful = false
	u.Data = Character{}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

// OK(character Character, w http.ResponseWriter)
//
// A friendly status Ok writter to the web console to write the data
// set that the request was successful, and set the that the function
// is ready to be returned.
//
// params:
// 	character Character -> This is the main data that we want
//		to send back to who ever is requesting the data.
// 	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
func (u CharacterDetailedResponse) OK(character Character, w http.ResponseWriter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = character
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
func (u CharactersDetailedResponse) ConsumeError(err error, w http.ResponseWriter, code int) {
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

// OK(characters Characters, w http.ResponseWriter)
//
// A friendly status Ok writter to the web console to write the data
// set that the request was successful, and set the that the function
// is ready to be returned.
//
// params:
// 	characters Characters -> This is the main data that we want
//		to send back to who ever is requesting the data.
// 	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
func (u CharactersDetailedResponse) OK(characters Characters, w http.ResponseWriter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = characters
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
