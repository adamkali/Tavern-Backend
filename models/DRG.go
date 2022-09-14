package models

import (
	"encoding/json"
	"net/http"
)

type DetailedResponse[T any] struct {
	Data       T      `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

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
	err error,
	w http.ResponseWriter,
	code int,
) {
	dr.Data = nil
	dr.UDRWrite(w, code, err.Error(), false)
}

func (dr DetailedResponse[T]) OK(
	data T,
	w http.ResponseWriter,
) {
	dr.Data = data
	dr.UDRWrite(w, http.StatusOK, "OK", true)
}
