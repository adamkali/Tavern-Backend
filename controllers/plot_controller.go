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

	var logger lib.LogEntryObject
	startTime := time.Now()

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
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusNotFound, 0, "Insufficent path...")
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
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, "Guid length not long enough")
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
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, result.Error.Error(), result.Error)
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
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusNotModified, 0, "Plot not found.")
		return
	}

	w.Header().Add("content-type", "application/json")
	response.OK(data, w)
	logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(data)
	logger.Log(r, http.StatusOK, float64(network.Len()/1000), "Plot found.")
	return
}

func (h *plotHandler) updatePlotByID(w http.ResponseWriter, r *http.Request) {

	var logger lib.LogEntryObject
	startTime := time.Now()

	var response models.PlotDetailedResponse

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

	plotId := string(path[3])
	if len(plotId) != 32 {
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

	var plot models.Plot
	err = json.Unmarshal(bodyBytes, &plot)
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

	old := &models.Plot{ID: plotId}
	result := h.db.First(&old)
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
	plot.ID = plotId
	plot.UserID = old.UserID
	result = h.db.Save(&plot)
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

	response.OK(plot, w)
	logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(plot)
	return
}

func (h *plotHandler) deletePlotByID(w http.ResponseWriter, r *http.Request) {

	var logger lib.LogEntryObject
	startTime := time.Now()

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
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusNotFound, 0, "Insufficent Path")
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
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusBadRequest, 0, "Guid not long enough")
		return
	}

	result := h.db.Delete(&models.Plot{}, plotId)

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
	return
}

// GET /api/plot/userId/{userId}
func (h *plotHandler) getPlotsByUserID(w http.ResponseWriter, r *http.Request) {

	var logger lib.LogEntryObject
	startTime := time.Now()

	var response models.PlotsDetailedResponse

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
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, result.Error.Error(), result.Error)
		return
	}

	for _, chara := range prep {
		plots = append(plots, chara)
	}

	_, err := json.Marshal(plots)
	if err != nil {
		response.ConsumeError(err, w, http.StatusInternalServerError)
		logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
		logger.Log(r, http.StatusInternalServerError, 0, err.Error(), err)
		return
	}

	// On success we can now write the response!
	w.Header().Add("content-type", "application/json")
	response.OK(plots, w)
	logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(plots)
	logger.Log(r, http.StatusOK, float64(network.Len()/1000), "Success")
	return
}

// 	>=> POST /api/plots/userId/{userId}
func (h *plotHandler) postPlotByUserID(w http.ResponseWriter, r *http.Request) {

	var logger lib.LogEntryObject
	startTime := time.Now()

	var response models.PlotDetailedResponse

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

	var plot models.Plot
	err = json.Unmarshal(bodyBytes, &plot)
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

	plot.ID = string(generateUUID())
	plot.UserID = userId

	result := h.db.Create(&plot)
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

	response.OK(plot, w)
	logger.TimeTaken = int64(time.Since(startTime) * time.Millisecond)
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(plot)
	logger.Log(r, http.StatusOK, float64(network.Len()/1000), "Success")
	return
}
