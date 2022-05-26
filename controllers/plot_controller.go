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

type plotHandler struct {
	db *gorm.DB
}

func NewPlotHandler(database gorm.DB) *plotHandler {
	return &plotHandler{
		db: &database,
	}
}

/*
=== === === === === === === === === === === === === === === === === === ===
		>=> PLOT CONTROLLER PAGES <=<
=== === === === === === === === === === === === === === === === === === === */

func (h *plotHandler) Plot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getPlotByID(w, r)
		return
	case "PUT":
		h.updatePlotByID(w, r)
		return
	case "DELETE":
		h.deletePlotByID(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}
}

func (h *plotHandler) Plots(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getPlotsByUserID(w, r)
		return
	case "POST":
		h.postPlotByUserID(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
		return
	}
}

func (h *plotHandler) getPlotByID(w http.ResponseWriter, r *http.Request) {
	var response models.PlotDetailedResponse
	var data models.Plot

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

	plotId := string(path[3])
	if len(plotId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid length not long enough",
			false,
		)
		return
	}

	data.ID = plotId
	result := h.db.First(&data)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	if data.ID == "" {
		response.Data = models.Plot{ID: plotId}
		response.UDRWrite(
			w,
			http.StatusNotModified,
			"Plot not found.",
			false,
		)
		return
	}

	w.Header().Add("content-type", "application/json")
	response.OK(data, w)
	return
}

func (h *plotHandler) updatePlotByID(w http.ResponseWriter, r *http.Request) {

	var response models.PlotDetailedResponse

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

	plotId := string(path[3])
	if len(plotId) != 32 {
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

	var plot models.Plot
	err = json.Unmarshal(bodyBytes, &plot)
	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
		return
	}
	
	old := &models.Plot{ ID: plotId }
	result := h.db.First(&old)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}
	plot.ID = plotId
	plot.UserID = old.UserID
	result = h.db.Save(&plot)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	response.OK(plot, w)
	return
}

func (h *plotHandler) deletePlotByID(w http.ResponseWriter, r *http.Request) {

	var response models.PlotDetailedResponse
	var data models.Plot

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

	plotId := string(path[3])
	if len(plotId) != 32 {
		response.UDRWrite(
			w,
			http.StatusBadRequest,
			"Guid not long enough",
			false,
		)
		return
	}

	result := h.db.Delete(&models.Plot{}, plotId)

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

// GET /api/plot/userId/{userId}
func (h *plotHandler) getPlotsByUserID(w http.ResponseWriter, r *http.Request) {

	var response models.PlotsDetailedResponse

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

	var prep []models.Plot
	var plots models.Plots
	var plot models.Plot
	plot.UserID = userId
	result := h.db.Model(&plot).Find(&prep)

	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	for _, chara := range prep {
		plots = append(plots, chara)
	}

	_, err := json.Marshal(plots)
	if err != nil {
		response.ConsumeError(err, w, http.StatusInternalServerError)
		return
	}

	// On success we can now write the response!
	w.Header().Add("content-type", "application/json")
	response.OK(plots, w)
	return
}

// 	>=> POST /api/plots/userId/{userId}
func (h *plotHandler) postPlotByUserID(w http.ResponseWriter, r *http.Request) {

	var response models.PlotDetailedResponse

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

	var plot models.Plot
	err = json.Unmarshal(bodyBytes, &plot)
	if err != nil {
		response.ConsumeError(
			err,
			w,
			http.StatusBadRequest,
		)
		return
	}

	plot.ID = string(generateUUID())
	plot.UserID = userId

	result := h.db.Create(&plot)
	if result.Error != nil {
		response.ConsumeError(
			result.Error,
			w,
			http.StatusInternalServerError,
		)
		return
	}

	response.OK(plot, w)
	return
}
