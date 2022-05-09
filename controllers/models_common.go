package controllers

import (
	"crypto/rand"
	"fmt"
	"log"
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
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
