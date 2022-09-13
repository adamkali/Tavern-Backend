package models

import (
	"Tavern-Backend/lib"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
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

type AuthToken struct {
	ID        string `json:"id" gorm:"column:id;type:varchar(32);primaryKey"`
	UserID    string `json:"user_fk,omitempty" gorm:"foreignKey:UserID;refernces:ID"`
	Username  string `json:"username" gorm:"column:username;type:varchar(32) not null"`
	UserEmail string `json:"user_email" gorm:"column:email;type:varchar(128) not null"`
	AuthHash  string `json:"auth_hash" gorm:"column:auth_hash;type:varchar(128) not null"`
	Active    bool   `json:"active" gorm:"column:active;type:tinyint(1) not null"`
}

type AuthTokenActivation struct {
	ID        string `json:"id" gorm:"column:id;type:varchar(32);primaryKey"`
	AuthID    string `json:"auth_fk,omitempty" gorm:"foreignKey:AuthID;refernces:ID"`
	AuthPin   string `json:"auth_pin" gorm:"column:auth_pin;type:varchar(8) not null"`
	AuthEmail string `json:"auth_email" gorm:"column:auth_email;type:varchar(128) not null"`
}

type Users []User

type AuthRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserEmail string `json:"user_email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a AuthTokenActivation) SendRegistrationEmail(config lib.Configuration) error {
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

	msg.SetHeader("From", config.Email.Username)
	msg.SetHeader("To", a.AuthEmail)
	msg.SetHeader("Subject", "Tavern Registration")
	msg.SetBody("text/html", fl)

	d := gomail.NewDialer(config.Email.Host, config.Email.Port, config.Email.Username, config.Email.Password)
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return err
}

// Make a function to take a username, password, and userID
// and then return a token
// This function will use a hash function to create a hash
// by using the username, password, and userID
func (t *AuthToken) GenerateToken(username string, password string) {
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
	t.UserEmail = password
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
