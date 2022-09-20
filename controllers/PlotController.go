package controllers

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"net/http"

	"gorm.io/gorm"
)

type PlotController struct{ H BaseHandler[models.Plot] }

func NewPlotController(DB *gorm.DB) *PlotController {
	return &PlotController{
		H: *NewHandler(DB, models.Plot{}, "Plot"),
	}
}

func (c *PlotController) AuthGetAllPByUserID(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	c.H.ForceGET(w, r)

	authToken, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		c.H.ResponseList.UDRWrite(w, http.StatusUnauthorized, "Unauthorized", false)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Unauthorized")
		return
	}
	c.H.SetAuthToken(authToken)

	userID := c.H.AuthToken.UserID
	var plots []models.Plot
	res := c.H.DB.Where("user_fk = ?", userID).Find(&plots)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	c.H.ResponseList.OK(w, plots)
	size := c.H.ResponseList.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}
