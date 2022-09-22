package controllers

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"net/http"

	"gorm.io/gorm"
)

type CharacterController struct{ H BaseHandler[models.Character] }

func NewCharacterController(DB *gorm.DB) *CharacterController {
	return &CharacterController{
		H: *NewHandler(DB, models.Character{}, "character"),
	}
}

func (c *CharacterController) AuthGetAllCByUserID(w http.ResponseWriter, r *http.Request) {
	c.H.Model = models.Character{}

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
	var characters []models.Character
	res := c.H.DB.Where("user_fk = ?", userID).Find(&characters)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	c.H.ResponseList.OK(w, characters)
	size := c.H.ResponseList.SizeOf()
	logger.Log(size, http.StatusOK, "OK")

}
