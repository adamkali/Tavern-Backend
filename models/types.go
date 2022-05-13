package models

import (
	"encoding/json"
)

type User struct {
	ID              string   `json:"id"`
	Username        string   `json:"username"`
	Bio             string   `json:"bio"`
	Tags            []string `json:"tags"`
	PlayerPrefrence string   `json:"player_prefrence"`

	Plots struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Plot   string `json:"plot"`
		UserFk string `json:"user_fk"`
	} `json:"plots,omitempty"`

	Character struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Backstory    string `json:"backstory"`
		Bio          string `json:"bio"`
		Strength     int    `json:"strength"`
		Dexterity    int    `json:"dexterity"`
		Constitution int    `json:"constitution"`
		Intellegence int    `json:"intellegence"`
		Wisdom       int    `json:"wisdom"`
		Charisma     int    `json:"charisma"`
		UserFk       string `json:"user_fk"`
	} `json:"character,omitempty"`
	GroupFk string `json:"group_fk,omitempty"`
}

type Group struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	GroupMembers Users  `json:"group_members"`
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

func (u UserDetailedResponse) UDRWrite(w http.ResponseWritter) {
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	return jsonBytes
}

func (u UsersDetailedResponse) UDRWrite(w http.ResponseWritter) {
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	return jsonBytes
}

func (u GroupDetailedResponse) UDRWrite(w http.ResponseWritter) {
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		panic(err.Error())
	}
	return jsonBytes
}

func (u UserDetailedResponse) ConsumeError(err error, w http.ResponseWritter) {
	u.Message = err.Error()
	u.Successful = false
	u.Data = User{}
	return u.UDRWrite()
}

func (u UserDetailedResponse) OK(user User, w http.ResponseWritter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = user
	return u.UDRWrite()
}

func (u UsersDetailedResponse) ConsumeError(err error, w http.ResponseWritter) {
	u.Message = err.Error()
	u.Successful = false
	u.Data = nil
	return u.UDRWrite()
}

func (u UsersDetailedResponse) OK(user Users, w http.ResponseWritter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = user
	return u.UDRWrite()
}

func (u GroupDetailedResponse) ConsumeError(err error, w http.ResponseWritter) {
	u.Message = err.Error()
	u.Successful = false
	u.Data = Group{}
	return u.UDRWrite()
}

func (u GroupDetailedResponse) OK(user Group, w http.ResponseWritter) {
	u.Message = "OK"
	u.Successful = true
	u.Data = user
	return u.UDRWrite()
}
