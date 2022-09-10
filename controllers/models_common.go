package controllers

import (
	"Tavern-Backend/models"
	"crypto/rand"
	"fmt"
	"log"

	"gorm.io/gorm"
)

var bank = "ABCDEF0123456789"

func EmptyGuid() string {
	return "00000000000000000000000000000000"
}

func generateUUID() string {
	//00000000-0000-0000-0000-000000000000
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x%x%x%x%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func generatePin() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	pin := fmt.Sprintf("%x%x%x%x",
		b[0:1], b[1:2], b[2:3], b[3:4])
	return pin
}

func verifyAuthorizationToken(db gorm.DB, authToken string) (models.AuthToken, error) {
	var data models.AuthToken
	result := db.Where("token = ?", authToken).First(&data)
	if result.Error != nil {
		return models.AuthToken{}, result.Error
	}
	return data, nil
}
