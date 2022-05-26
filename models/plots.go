package models

import (
	"encoding/json"
	"net/http"
)

type Plots []Plot

type PlotDetailedResponse struct {
	Data       Plot   `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type PlotsDetailedResponse struct {
	Data       Plots  `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

// UDRWrite(w http.ResponseWriter, code int, message string)
//
// This is a method to diplay any text in detailed responses
// To the http.ResponseWriter. It will dump the message into
// The json displayed the client

func (u PlotDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
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
func (u PlotsDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
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
func (u PlotDetailedResponse) ConsumeError(err error, w http.ResponseWriter, code int) {
	u.Message = err.Error()
	u.Successful = false
	u.Data = Plot{}
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

// OK(plot Plot, w http.ResponseWriter)
//
// A friendly status Ok writter to the web console to write the data
// set that the request was successful, and set the that the function
// is ready to be returned.
//
// params:
// 	plot Plot -> This is the main data that we want
//		to send back to who ever is requesting the data.
// 	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
func (u PlotDetailedResponse) OK(plot Plot, w http.ResponseWriter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = plot
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
func (u PlotsDetailedResponse) ConsumeError(err error, w http.ResponseWriter, code int) {
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

// OK(plots Plots, w http.ResponseWriter)
//
// A friendly status Ok writter to the web console to write the data
// set that the request was successful, and set the that the function
// is ready to be returned.
//
// params:
// 	plots Plots -> This is the main data that we want
//		to send back to who ever is requesting the data.
// 	w http.ResponseWriter -> The writer incharge of outputting
//		to the web console that gets responses.
func (u PlotsDetailedResponse) OK(plots Plots, w http.ResponseWriter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = plots
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
