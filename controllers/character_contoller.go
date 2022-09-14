package controllers

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type characterHandler struct {
	db *gorm.DB
}

func NewCharacterHandler(database gorm.DB) *characterHandler {
	return &characterHandler{
		db: &database,
	}
}

/*
=== === === === === === === === === === === === === === === === === === ===
		>=> CHARACTER CONTROLLER PAGES <=<
=== === === === === === === === === === === === === === === === === === === */

func (h *characterHandler) Character(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getCharacterByID(w, r)
		return
	case "PUT":
		h.updateCharacterByID(w, r)
		return
	case "DELETE":
		h.deleteCharacterByID(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}
}

func (h *characterHandler) Characters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getCharactersByUserID(w, r)
		return
	case "POST":
		h.postCharacterByUserID(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
		return
	}
}

func (h *characterHandler) getCharacterByID(w http.ResponseWriter, r *http.Request) {
	adminLog := lib.New(r)

	var response models.CharacterDetailedResponse
	var data models.Character

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 4 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent path...",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusNotFound, 0, "Insufficent path...")
		return
	}

	characterId := string(path[3])
	if len(characterId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid length not long enough",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, "Guid length not long enough")
		return
	}

	result := h.db.First(&data, characterId)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, result.Error.Error(), result.Error)
		return
	}

	if data.ID == "" {
		response.Data = models.Character{ID: characterId}
		response.UDRWrite(
			w,
			http.StatusNotModified,
			"Character not found.",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusNotModified, 0, "Character not found.")
		return
	}

	w.Header().Add("content-type", "application/json")
	response.OK(data, w)
	logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(data)
	logger.Log(r, http.StatusOK, float64(len(network.Bytes()))/1000, "OK")
	return
}

func (h *characterHandler) updateCharacterByID(w http.ResponseWriter, r *http.Request) {
	var logger lib.LogEntryObject
	startTime := time.Now()

	var response models.CharacterDetailedResponse

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 4 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusNotFound, 0, "Insufficent Path")
		return
	}

	characterId := string(path[3])
	if len(characterId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, "Guid not long enough")
		return
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		response.UDRWrite(
			w,
			http.StatusUnsupportedMediaType,
			"Content Type needs to be application/json.",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusUnsupportedMediaType, 0, "Content Type needs to be application/json.")
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Error reading body.",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, "Error reading body.")
		return
	}

	var character models.Character
	err = json.Unmarshal(bodyBytes, &character)
	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, err.Error(), err)
		return
	}

	old := &models.Character{ID: characterId}
	result := h.db.First(old)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, result.Error.Error(), result.Error)
		return
	}
	character.ID = characterId
	character.ID = old.UserID
	result = h.db.Save(&character)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, result.Error.Error(), result.Error)
		return
	}

	response.OK(character, w)
	logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(character)
	logger.Log(r, http.StatusOK, float64(len(network.Bytes()))/1000, "OK")
}

func (h *characterHandler) deleteCharacterByID(w http.ResponseWriter, r *http.Request) {

	var logger lib.LogEntryObject
	startTime := time.Now()

	var response models.CharacterDetailedResponse
	var data models.Character

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 4 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusNotFound, 0, "Insufficent Path")
		return
	}

	characterId := string(path[3])
	if len(characterId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, "Guid not long enough")
		return
	}

	result := h.db.Delete(&models.Character{}, characterId)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, result.Error.Error(), result.Error)
		return
	}
	response.OK(data, w)
	logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(data)
	logger.Log(r, http.StatusOK, float64(len(network.Bytes()))/1000, "OK")
	return
}

// GET /api/character/userId/{userId}
func (h *characterHandler) getCharactersByUserID(w http.ResponseWriter, r *http.Request) {

	var logger lib.LogEntryObject
	startTime := time.Now()

	var response models.CharactersDetailedResponse

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 5 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusNotFound, 0, "Insufficent Path")
		return
	}

	userId := string(path[4])
	if len(userId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, "Guid not long enough")
		return
	}

	var prep []models.Character
	var characters models.Characters
	var character models.Character
	character.UserID = userId
	result := h.db.Model(&character).Find(&prep)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, result.Error.Error(), result.Error)
		return
	}

	for _, chara := range prep {
		characters = append(characters, chara)
	}

	_, err := json.Marshal(characters)
	if err != nil {
		response.ConsumeError(err, w, http.StatusInternalServerError)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, err.Error(), err)
		return
	}

	// On success we can now write the response!
	w.Header().Add("content-type", "application/json")
	response.OK(characters, w)
	logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(characters)
	logger.Log(r, http.StatusOK, float64(len(network.Bytes()))/1000, "OK")
	return
}

// 	>=> POST /api/characters/userId/{userId}
func (h *characterHandler) postCharacterByUserID(w http.ResponseWriter, r *http.Request) {

	var logger lib.LogEntryObject
	startTime := time.Now()

	var response models.CharacterDetailedResponse

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 5 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusNotFound, 0, "Insufficent Path")

		return
	}

	userId := string(path[4])
	if len(userId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough",
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, "Guid not long enough")
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, err.Error(), err)
		return
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		response.UDRWrite(
			w,
			http.StatusUnsupportedMediaType,
			fmt.Sprintf("Application data is not application/json, got: {%s}", contentType),
			false,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusUnsupportedMediaType, 0, fmt.Sprintf("Application data is not application/json, got: {%s}", contentType))
		return
	}

	var character models.Character
	err = json.Unmarshal(bodyBytes, &character)
	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, err.Error(), err)

		return
	}

	character.ID = string(generateUUID())
	character.UserID = userId

	result := h.db.Create(&character)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, result.Error.Error(), result.Error)
		return
	}

	response.OK(character, w)
	logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(character)
	logger.Log(r, http.StatusOK, float64(len(network.Bytes()))/1000, "OK")
	return
}
