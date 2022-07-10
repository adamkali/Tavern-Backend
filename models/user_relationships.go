package models

import import (
	"encoding/json"
	"net/http"
)

// UserRelationships is a struct for the user_relationships table
// It will relate different users to each other
// For example, a user can like, dislike, or join another user
// This is a many-to-many relationship
type UserRelationship struct {
	ID 			string 	`json:"id" gorm:"column:id;type:varchar(32);primaryKey"`
	Self 		string 	`json:"self" gorm:"column:self;type:varchar(32)"`
	SelfUser 	User 	`json:"self_user, omitempty"`
	Other 		string 	`json:"other" gorm:"column:other;type:varchar(32)"`
	OtherUser 	User 	`json:"other_user, omitempty"`
	Type 		string 	`json:"type" gorm:"column:type;type:varchar(32)"`
}

type UserRelationships []UserRelationship

type UserRelationshipDetailedResponse struct {
	Data       UserRelationship `json:"data"`
	Successful bool             `json:"successful"`
	Message    string           `json:"message"`
}

type UserRelationshipsDetailedResponse struct {
	Data       UserRelationships `json:"data"`
	Successful bool              `json:"successful"`
	Message    string            `json:"message"`
}

func (u UserRelationshipDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
	u.Successful = successful
	u.Message = message
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

func (u UserRelationshipsDetailedResponse) UDRWrite(w http.ResponseWriter, code int, message string, successful bool) {
	u.Successful = successful
	u.Message = message
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

func (u UserRelationshipDetailedResponse) ConsumeError(
	err error, w http.ResponseWriter, code int
) {
	if err != nil {
		u.UDRWrite(w, code, err.Error(), false)
	}
}

func (u UserRelationshipsDetailedResponse) ConsumeError(
	err error, w http.ResponseWriter, code int
) {
	if err != nil {
		u.UDRWrite(w, code, err.Error(), false)
	}
}

func (u UserRelationshipDetialedResponse) OK(
	rel UserRelationship, w http.ResponseWriter
) {
	u.Data = rel
	u.UDRWrite(w, http.StatusOK, "OK", true)
}

func (u UserRelationshipsDetailedResponse) OK(
	rel UserRelationships, w http.ResponseWriter
) {
	u.Data = rel
	u.UDRWrite(w, http.StatusOK, "OK", true)
}