package controllers

import (
	"Tavern-Backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	enable(&w)
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
	enable(&w)
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
		return
	}

	result := h.db.First(&data, characterId)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
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
		return
	}

	w.Header().Add("content-type", "application/json")
	response.OK(data, w)
	return
}

func (h *characterHandler) updateCharacterByID(w http.ResponseWriter, r *http.Request) {

	var response models.CharacterDetailedResponse

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 4 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
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
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)

	var character models.Character
	err = json.Unmarshal(bodyBytes, &character)
	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
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
		return
	}

	response.OK(character, w)
	return
}

func (h *characterHandler) deleteCharacterByID(w http.ResponseWriter, r *http.Request) {

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
		return
	}

	result := h.db.Delete(&models.Character{}, characterId)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}
	response.OK(data, w)
	return
}

// GET /api/character/userId/{userId}
func (h *characterHandler) getCharactersByUserID(w http.ResponseWriter, r *http.Request) {

	var response models.CharactersDetailedResponse

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 5 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
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
		return
	}

	for _, chara := range prep {
		characters = append(characters, chara)
	}

	_, err := json.Marshal(characters)
	if err != nil {
		response.ConsumeError(err, w, http.StatusInternalServerError)
		return
	}

	// On success we can now write the response!
	w.Header().Add("content-type", "application/json")
	response.OK(characters, w)
	return
}

// 	>=> POST /api/characters/userId/{userId}
func (h *characterHandler) postCharacterByUserID(w http.ResponseWriter, r *http.Request) {

	var response models.CharacterDetailedResponse

	path := strings.Split(r.URL.String(), "/")
	if len(path) != 5 {
		response.UDRWrite(
			w,
			http.StatusNotFound,
			"Insufficent Path",
			false,
		)
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
		return
	}

	response.OK(character, w)
	return
}
