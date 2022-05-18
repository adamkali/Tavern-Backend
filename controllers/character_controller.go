package controllers

import (
	"Tavern-Backend/models"
	"encoding/json"
	"net/http"
	"sync"
)

type characterHandler struct {
	sync.Mutex
	store userHandler.store
}

func NewCharacterHandler() characterHandler {
	return &characterHandler{
		store: map[string]models.Character{},
	}
}

/*
=== === === === === === === === === === === === === === === === === === ===
		>=> CHARACTERS CONTROLLER PAGES <=<
=== === === === === === === === === === === === === === === === === === ===
*/

func (h *characterHandler) Character(w http.ResponseWriter, r *http.Request) {
	/*switch r.Method {
	case "GET":

	}
	*/
}

/*
=== === === === === === === === === === === === === === === === === === ===
	>=> CHARACTERS CONTROLLER ENDPOINTS `/api/users` <=<
=== === === === === === === === === === === === === === === === === === ===
*/

// 	>=> GET /api/characters/{userId}
func (h characterHandler) get(w http.ResponseWriter, r *http.Request) {
	characters := make(models.Characters, len(h.store))

	var response models.CharactersDetailedResponse

	h.Lock()
	i := 0
	for _, character := range h.store {
		characters[i] = character
		i++
	}
	h.Unlock()

	_, err := json.Marshal(character)
	if err != nil {
		response.ConsumeError(err, w, http.StatusInternalServerError)
		return
	}

	// On success we can now write the response!
	w.Header().Add("content-type", "application/json")
	response.OK(characters, w)
	return
}
