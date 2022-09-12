package models

import (
	"encoding/json"
	"net/http"
)

type Tags struct {
	ID      string `json:"id" gorm:"primaryKey; not null; type:varchar(32);"`
	TagID   int    `json:"tag_id" gorm:"column:tag_id;type:smallint(255);not null"`
	TagName string `json:"tag_name" gorm:"column:tag_name;varchar(32) not null"`
}

type TagsDetailedResponse struct {
	Data       []Tags `json:"data"`
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}

func (t *TagsDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
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

func (t *TagsDetailedResponse) OK(w http.ResponseWriter, tags []Tags) {
	t.Data = tags
	t.UDRWrite(w, http.StatusOK, "OK", true)
}

func (t *TagsDetailedResponse) ConsumeError(w http.ResponseWriter, err error) {
	t.UDRWrite(w, http.StatusInternalServerError, err.Error(), false)
}

type PlayerPrefrence struct {
	ID           string `json:"id" gorm:"primaryKey; not null; type:varchar(32);"`
	PreferenceID int    `json:"pref_id" gorm:"column:pref_id;type:smallint(255);not null"`
	Preference   string `json:"pref_name" gorm:"column:pref_name;varchar(32) not null"`
}

type PlayerPrefrenceDetailedResponse struct {
	Data       []PlayerPrefrence `json:"data"`
	Successful bool              `json:"successful"`
	Message    string            `json:"message"`
}

func (t *PlayerPrefrenceDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
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

func (t *PlayerPrefrenceDetailedResponse) OK(w http.ResponseWriter, tags []PlayerPrefrence) {
	t.Data = tags
	t.UDRWrite(w, http.StatusOK, "OK", true)
}

func (t *PlayerPrefrenceDetailedResponse) ConsumeError(w http.ResponseWriter, err error) {
	t.UDRWrite(w, http.StatusInternalServerError, err.Error(), false)
}
