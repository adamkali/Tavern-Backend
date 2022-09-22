package controllers

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type BaseHandler[Model models.IData] struct {
	DB           *gorm.DB
	Model        Model
	ModelName    string
	Response     models.DetailedResponse[Model]
	ResponseList models.DetailedResponseList[Model]
	AuthToken    models.AuthToken
	BasePath     string
	AuthPath     string
	AdmnPath     string
	BaseAllPath  string
	AuthAllPath  string
	AdmnAllPath  string
}

func NewHandler[Model models.IData](
	database *gorm.DB,
	model Model,
	modelName string,
) *BaseHandler[Model] {
	return &BaseHandler[Model]{
		DB:           database,
		Model:        model,
		Response:     models.DetailedResponse[Model]{},
		ResponseList: models.DetailedResponseList[Model]{},
		AuthToken:    models.AuthToken{},
		BasePath:     "/api/" + modelName,
		AuthPath:     "/api/auth/" + modelName,
		AdmnPath:     "/api/admin/" + modelName,
		BaseAllPath:  "/api/" + modelName + "s",
		AuthAllPath:  "/api/auth/" + modelName + "s",
		AdmnAllPath:  "/api/admin/" + modelName + "s",
		ModelName:    modelName,
	}
}

func (h *BaseHandler[Model]) SetAuthToken(token models.AuthToken) {
	h.AuthToken = token
}

func (h *BaseHandler[Model]) authGetByID(w http.ResponseWriter, r *http.Request) {
	h.Model = h.Model.NewData().(Model)
	logger := lib.New(r)

	// check the header for the auth token
	token := r.Header.Get("Authorization")
	if token == "" {
		h.Response.UDRWrite(w, http.StatusUnauthorized, "No token provided", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "No token provided")
		return
	}

	// check the token
	res := h.DB.Where("auth_hash = ?", token).First(&h.AuthToken)
	if res.Error != nil {
		h.Response.UDRWrite(w, http.StatusUnauthorized, "Invalid token", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Invalid token", res.Error)
	}

	// get the id from the url
	id := r.URL.Path[len(h.AuthPath+"/"):]
	if id == "" {
		h.Response.UDRWrite(w, http.StatusBadRequest, "No id provided", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "No id provided")
		return
	}

	// get the data from the database
	res = h.DB.Where("id = ?", id).First(&h.Model)
	if res.Error != nil {
		h.Response.UDRWrite(w, http.StatusNotFound, "No data found", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusNotFound, "No data found", res.Error)
		return
	}

	// write the Response
	h.Response.OK(w, h.Model)
	size := h.Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

func (h *BaseHandler[Model]) authUpdateOrInsert(w http.ResponseWriter, r *http.Request) {
	h.Model = h.Model.NewData().(Model)
	logger := lib.New(r)

	// check the header for the auth token
	token := r.Header.Get("Authorization")
	if token == "" {
		h.Response.UDRWrite(w, http.StatusUnauthorized, "No token provided", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "No token provided")
		return
	}

	// check the token
	res := h.DB.Where("auth_hash = ?", token).First(&h.AuthToken)
	if res.Error != nil {
		h.Response.UDRWrite(w, http.StatusUnauthorized, "Invalid token", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Invalid token", res.Error)
	}

	// read the body
	body := r.Body
	//Try to decode the request body into the struct. If there is an error, respond to the client with the error message and a 400 status code
	// ensure application/json
	if r.Header.Get("Content-Type") != "application/json" {
		h.Response.UDRWrite(w, http.StatusBadRequest, "Invalid content type", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Invalid content type")
		return
	}

	err := json.NewDecoder(body).Decode(&h.Model)
	if err != nil {
		h.Response.ConsumeError(w, err, http.StatusBadRequest)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	defer body.Close()

	var special string
	// Try to update the data in the database
	res = h.DB.Where("id = ?", h.Model.GetID()).Updates(&h.Model)
	if res.Error != nil {
		// The data does not exist, so insert it
		res = h.DB.Create(&h.Model)
		if res.Error != nil {
			h.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
			size := h.Response.SizeOf()
			logger.Log(size, http.StatusInternalServerError, "Failed to insert Data", res.Error)
			return
		}
		special = h.ModelName + " has been updated"
	} else {
		special = h.ModelName + " has been created"
	}

	// Make an OK Response with special as the message
	h.Response.UDRWrite(w, http.StatusOK, special, true)
	size := h.Response.SizeOf()
	logger.Log(size, http.StatusOK, special)

}

func (h *BaseHandler[Model]) authDeleteByID(w http.ResponseWriter, r *http.Request) {
	h.Model = h.Model.NewData().(Model)
	logger := lib.New(r)

	// check the header for the auth token
	token := r.Header.Get("Authorization")
	if token == "" {
		h.Response.UDRWrite(w, http.StatusUnauthorized, "No token provided", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "No token provided")
		return
	}

	// check the token
	res := h.DB.Where("auth_hash = ?", token).First(&h.AuthToken)
	if res.Error != nil {
		h.Response.UDRWrite(w, http.StatusUnauthorized, "Invalid token", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Invalid token", res.Error)
	}

	// get the id from the url
	id := r.URL.Path[len(h.AuthPath+"/"):]
	if id == "" {
		h.Response.UDRWrite(w, http.StatusBadRequest, "No id provided", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "No id provided")
		return
	}

	// delete the data from the database
	res = h.DB.Where("id = ?", id).Delete(&h.Model)
	if res.Error != nil {
		h.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Failed to delete Data", res.Error)
		return
	}

	// make an OK Response
	h.Response.UDRWrite(w, http.StatusOK, h.ModelName+" has been deleted", true)
	size := h.Response.SizeOf()
	logger.Log(size, http.StatusOK, h.ModelName+" has been deleted")
}

func (h *BaseHandler[Model]) Controller(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.authGetByID(w, r)
	case http.MethodPost:
		h.authUpdateOrInsert(w, r)
	case http.MethodDelete:
		h.authDeleteByID(w, r)
	default:
		h.Response.UDRWrite(w, http.StatusMethodNotAllowed, "Method not allowed", false)
	}
}

func (h *BaseHandler[Model]) ForceGET(w http.ResponseWriter, r *http.Request) {
	// check that the request is a GET
	logger := lib.New(r)
	if r.Method != http.MethodGet {
		h.Response.UDRWrite(w, http.StatusMethodNotAllowed, "Method not allowed! Try again...", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
}

func (h *BaseHandler[Model]) ForcePOST(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	// check that the request is a POST
	if r.Method != http.MethodPost {
		h.Response.UDRWrite(w, http.StatusMethodNotAllowed, "Method not allowed! Try again...", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
}

func (h *BaseHandler[Model]) ForcePUT(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	// check that the request is a PUT
	if r.Method != http.MethodPut {
		h.Response.UDRWrite(w, http.StatusMethodNotAllowed, "Method not allowed! Try again...", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
}

func (h *BaseHandler[Model]) ForceDELETE(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	// check that the request is a DELETE
	if r.Method != http.MethodDelete {
		h.Response.UDRWrite(w, http.StatusMethodNotAllowed, "Method not allowed! Try again...", false)
		size := h.Response.SizeOf()
		logger.Log(size, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
}
