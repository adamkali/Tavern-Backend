package models

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"
)

type IData interface {
	User |
		Users |
		Character |
		Characters |
		Plot |
		Plots |
		Tag |
		Tags |
		PlayerPrefrence |
		PlayerPrefrences |
		AuthToken |
		Relationship |
		Relationships |
		UserRelationship |
		UserRelationships |
		Role
	SetID(string)
	GetID() string
	NewData() interface{}
}

type DetailedResponse[T IData] struct {
	Data       T      `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

type DetailedResponseList[T IData] struct {
	Data       []T    `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

//#region DetailedResponse Methods

func (dr DetailedResponse[T]) UDRWrite(
	w http.ResponseWriter,
	code int,
	message string,
	successful bool,
) {
	dr.Successful = successful
	dr.Message = message
	jsonBytes, err := json.Marshal(dr)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

func (dr DetailedResponse[T]) ConsumeError(
	w http.ResponseWriter,
	err error,
	code int,
) {
	dr.UDRWrite(w, code, err.Error(), false)
}

func (dr DetailedResponse[T]) OK(
	w http.ResponseWriter,
	data T,
) {
	dr.Data = data
	dr.UDRWrite(w, http.StatusOK, "OK", true)
}

func (dr DetailedResponse[T]) SizeOf() float32 {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(dr)
	return float32(network.Len() / 1024)
}

func (dr DetailedResponseList[T]) UDRWrite(
	w http.ResponseWriter,
	code int,
	message string,
	successful bool,
) {
	dr.Successful = successful
	dr.Message = message
	jsonBytes, err := json.Marshal(dr)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

func (dr DetailedResponseList[T]) ConsumeError(
	w http.ResponseWriter,
	err error,
	code int,
) {
	dr.UDRWrite(w, code, err.Error(), false)
}

func (dr DetailedResponseList[T]) OK(
	w http.ResponseWriter,
	data []T,
) {
	dr.Data = data
	dr.UDRWrite(w, http.StatusOK, "OK", true)
}

func (dr DetailedResponseList[T]) SizeOf() float32 {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(dr)
	return float32(network.Len() / 1024)
}

//#endregion
