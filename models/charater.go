package models

import (
	"encoding/json"
	"net/http"
)

type Characters []Characters

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

func (u CharacterDetailedResponse) ConsumeError(err error, w http.ResponseWriter, code int) {
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

func (u CharacterDetailedResponse) OK(user Character, w http.ResponseWriter) {
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

func (u CharactersDetailedResponse) OK(user Users, w http.ResponseWriter) {
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
