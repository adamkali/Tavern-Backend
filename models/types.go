package models

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
)

// #region TYPES
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
	ID               string            `json:"id" gorm:"column:id;type:varchar(32)"`
	Username         string            `json:"username" gorm:"column:username;type:varchar(128) not null"`
	Bio              string            `json:"bio" grom:"column:bio;type:text not null"`
	Plots            []Plot            `json:"user_plots,omitempty"`
	Characters       []Character       `json:"user_characters,omitempty"`
	Tags             []Tag             `json:"user_tags,omitempty"`
	PlayerPrefrences []PlayerPrefrence `json:"user_player_prefrences,omitempty"`

	// GroupID string `json:"group_fk,omitempty" gorm:"foreignKey:GroupID;refernces:ID"`
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
	Role      Role   `json:"role" gorm:"foreignKey:RoleID;refernces:ID"`
}

type AuthTokenActivation struct {
	ID        string `json:"id" gorm:"column:id;type:varchar(32);primaryKey"`
	AuthID    string `json:"auth_fk,omitempty" gorm:"foreignKey:AuthID;refernces:ID"`
	AuthPin   string `json:"auth_pin" gorm:"column:auth_pin;type:varchar(8) not null"`
	AuthEmail string `json:"auth_email" gorm:"column:auth_email;type:varchar(128) not null"`
}

type UserRelationship struct {
	ID        string       `json:"id" gorm:"column:id;type:varchar(32);primaryKey"`
	Self      string       `json:"self" gorm:"column:self;type:varchar(32)"`
	SelfUser  User         `json:"self_user, omitempty"`
	Other     string       `json:"other" gorm:"column:other;type:varchar(32)"`
	OtherUser User         `json:"other_user, omitempty"`
	Type      Relationship `json:"type" gorm:"column:type;type:varchar(32)"`
}

type UserRelationships struct {
	UserRelationships []UserRelationship `json:"user_relationships"`
	TempID            string             `json:"temp_id"`
}
type Characters struct {
	Characters []Character `json:"characters"`
	TempID     string      `json:"temp_id"`
}
type Users struct {
	Users  []User `json:"users"`
	TempID string `json:"temp_id"`
}
type Plots struct {
	Plots  []Plot `json:"plots"`
	TempID string `json:"temp_id"`
}

type AuthEmailConfiglette struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// #endregion

// #region AUTHENTICATION

func (a AuthTokenActivation) SendRegistrationEmail(config AuthEmailConfiglette) error {
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

	msg.SetHeader("From", config.Username)
	msg.SetHeader("To", a.AuthEmail)
	msg.SetHeader("Subject", "Tavern Registration")
	msg.SetBody("text/html", fl)

	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return err
}

// Make a function to take a username, password, and userID
// and then return a token
// This function will use a hash function to create a hash
// by using the username, password, and userID
func (t *AuthToken) GenerateToken(username string, password string, user_email string) {
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
	t.UserEmail = user_email
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

func (t *AuthToken) IsAdmin() bool {
	return t.Role.RoleName == "Admin"
}

func (t *AuthToken) IsPremium() bool {
	return t.Role.RoleName == "Premium" || t.Role.RoleName == "Admin" || t.Role.RoleName == "Lifetime Premium"
}

func (t *AuthToken) IsLifetimePremium() bool {
	return t.Role.RoleName == "Lifetime Premium" || t.Role.RoleName == "Admin"
}

// #endregion

// #region COMMON FUNCTIONS
func (u User) SetID(id string)                { u.ID = id }
func (u Plot) SetID(id string)                { u.ID = id }
func (u Character) SetID(id string)           { u.ID = id }
func (u UserRelationship) SetID(id string)    { u.ID = id }
func (u AuthToken) SetID(id string)           { u.ID = id }
func (u AuthTokenActivation) SetID(id string) { u.ID = id }
func (u User) GetID() string                  { return u.ID }
func (u Plot) GetID() string                  { return u.ID }
func (u Character) GetID() string             { return u.ID }
func (u UserRelationship) GetID() string      { return u.ID }
func (u AuthToken) GetID() string             { return u.ID }
func (u AuthTokenActivation) GetID() string   { return u.ID }

// implement the GetID and SetID functions for
// Users, Plots, Characters, UserRelationships, AuthTokens, and AuthTokenActivations
func (u Users) GetID() string {
	var ret []string
	for _, v := range u.Users {
		ret = append(ret, v.GetID())
	}
	completeRet := strings.Join(ret, ",")
	return completeRet
}
func (u Users) SetID(id string) {
	splitIds := strings.Split(id, ",")
	for i, v := range splitIds {
		u.Users[i].SetID(v)
	}
}

// Follow the above pattern for the rest of the
// IData Types
func (u Plots) GetID() string {
	var ret []string
	for _, v := range u.Plots {
		ret = append(ret, v.GetID())
	}
	completeRet := strings.Join(ret, ",")
	return completeRet
}
func (u Plots) SetID(id string) {
	splitIds := strings.Split(id, ",")
	for i, v := range splitIds {
		u.Plots[i].SetID(v)
	}
}
func (u Characters) GetID() string {
	var ret []string
	for _, v := range u.Characters {
		ret = append(ret, v.GetID())
	}
	completeRet := strings.Join(ret, ",")
	return completeRet
}
func (u Characters) SetID(id string) {
	splitIds := strings.Split(id, ",")
	for i, v := range splitIds {
		u.Characters[i].SetID(v)
	}
}

//Tags
func (u Tags) GetID() string {
	var ret []string
	for _, v := range u.Tags {
		ret = append(ret, v.GetID())
	}
	completeRet := strings.Join(ret, ",")
	return completeRet
}
func (u Tags) SetID(id string) {
	splitIds := strings.Split(id, ",")
	for i, v := range splitIds {
		u.Tags[i].SetID(v)
	}
}

//Player Preference
func (u PlayerPrefrences) GetID() string {
	var ret []string
	for _, v := range u.PlayerPrefrences {
		ret = append(ret, v.GetID())
	}
	completeRet := strings.Join(ret, ",")
	return completeRet
}
func (u PlayerPrefrences) SetID(id string) {
	splitIds := strings.Split(id, ",")
	for i, v := range splitIds {
		u.PlayerPrefrences[i].SetID(v)
	}
}

//Relationships
func (u Relationships) GetID() string {
	var ret []string
	for _, v := range u.Relationships {
		ret = append(ret, v.GetID())
	}
	completeRet := strings.Join(ret, ",")
	return completeRet
}
func (u Relationships) SetID(id string) {
	splitIds := strings.Split(id, ",")
	for i, v := range splitIds {
		u.Relationships[i].SetID(v)
	}
}

//UserRelationships
func (u UserRelationships) GetID() string {
	var ret []string
	for _, v := range u.UserRelationships {
		ret = append(ret, v.GetID())
	}
	completeRet := strings.Join(ret, ",")
	return completeRet
}
func (u UserRelationships) SetID(id string) {
	splitIds := strings.Split(id, ",")
	for i, v := range splitIds {
		u.UserRelationships[i].SetID(v)
	}
}

// #endregion
