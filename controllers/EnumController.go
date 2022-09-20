package controllers

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

// #region Controller Definitions
type TagController struct{ H BaseHandler[models.Tag] }

func NewTagController(DB *gorm.DB) *TagController {
	return &TagController{
		H: *NewHandler(DB, models.Tag{}, "enum/Tag"),
	}
}

type PlayerPrefrenceController struct {
	H BaseHandler[models.PlayerPrefrence]
}

func NewPlayerPrefrenceController(DB *gorm.DB) *PlayerPrefrenceController {
	return &PlayerPrefrenceController{
		H: *NewHandler(DB, models.PlayerPrefrence{}, "enum/PlayerPrefrence"),
	}
}

type RoleController struct{ H BaseHandler[models.Role] }

func NewRoleController(DB *gorm.DB) *RoleController {
	return &RoleController{
		H: *NewHandler(DB, models.Role{}, "enum/Role"),
	}
}

type RelsController struct {
	H BaseHandler[models.Relationship]
}

func NewRelsController(DB *gorm.DB) *RelsController {
	return &RelsController{
		H: *NewHandler(DB, models.Relationship{}, "enum/Relationship"),
	}
}

// #endregion

// #region Enum Tags Definitions
func (c *TagController) AuthPostTagToUser(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	c.H.ForcePOST(w, r)

	// set the Detailed Response to be a User
	var Response models.DetailedResponse[models.User]

	// authenticate the header
	token, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		Response.ConsumeError(w, err, http.StatusUnauthorized)
		size := Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	c.H.SetAuthToken(token)
	var user models.User
	res := c.H.DB.First(&user, c.H.AuthToken.UserID)
	if res.Error != nil {
		Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}
	// get the tag from the body
	req := models.TagRequest{}
	// Decode the r.Body into the model
	body := r.Body
	err = json.NewDecoder(body).Decode(&req)
	if err != nil {
		Response.ConsumeError(w, err, http.StatusBadRequest)
		size := Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Bad Request", err)
		return
	} else if req.UserID != c.H.AuthToken.UserID {
		Response.UDRWrite(w, http.StatusBadRequest, "Validation Failed", false)
		size := Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Validation Failed")
	}
	// add the tag to the user
	user.Tags = append(user.Tags, req.Tag)
	// save the user
	res = c.H.DB.Save(&user)
	if res.Error != nil {
		Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	Response.OK(w, user)
	size := Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")

}

func (c *TagController) AuthPostTagsToUser(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	c.H.ForcePOST(w, r)

	// set the Detailed Response to be a User
	var Response models.DetailedResponse[models.User]

	// authenticate the header
	token, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		Response.ConsumeError(w, err, http.StatusUnauthorized)
		size := Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	c.H.SetAuthToken(token)
	var user models.User
	res := c.H.DB.First(&user, c.H.AuthToken.UserID)
	if res.Error != nil {
		Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}
	// get the []Tag from the body
	var req models.TagsRequest
	// Decode the r.Body into the model
	body := r.Body
	err = json.NewDecoder(body).Decode(&req)
	if err != nil {
		Response.ConsumeError(w, err, http.StatusBadRequest)
		size := Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Bad Request", err)
		return
	} else if req.UserID != c.H.AuthToken.UserID {
		Response.UDRWrite(w, http.StatusBadRequest, "Validation Failed", false)
		size := Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Validation Failed")
	}
	// add the tag to the user
	user.Tags = append(user.Tags, req.Tags...)
	// save the user
	res = c.H.DB.Save(&user)
	if res.Error != nil {
		Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	Response.OK(w, user)
	size := Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

func (c *TagController) AuthGetTags(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	c.H.ForceGET(w, r)

	// authenticate the header
	token, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		c.H.ResponseList.ConsumeError(w, err, http.StatusUnauthorized)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	c.H.SetAuthToken(token)
	var user models.User
	res := c.H.DB.First(&user, c.H.AuthToken.UserID)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}
	// Get the tags
	var tags []models.Tag
	res = c.H.DB.Find(&tags)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}
	c.H.ResponseList.OK(w, user.Tags)
	size := c.H.ResponseList.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

// #endregion

// #region Enum PPs Definitions
func (c *PlayerPrefrenceController) AuthPostPlayerPrefrenceToUser(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	c.H.ForcePOST(w, r)

	// set the Detailed Response to be a User
	var Response models.DetailedResponse[models.User]

	// authenticate the header
	token, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		Response.ConsumeError(w, err, http.StatusUnauthorized)
		size := Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	c.H.SetAuthToken(token)
	var user models.User
	res := c.H.DB.First(&user, c.H.AuthToken.UserID)
	if res.Error != nil {
		Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}
	// get the PlayerPrefrence from the body
	req := models.PrefrenceRequest{}
	// Decode the r.Body into the model
	body := r.Body
	err = json.NewDecoder(body).Decode(&req)
	if err != nil {
		Response.ConsumeError(w, err, http.StatusBadRequest)
		size := Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Bad Request", err)
		return
	} else if req.UserID != c.H.AuthToken.UserID {
		Response.UDRWrite(w, http.StatusBadRequest, "Validation Failed", false)
		size := Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Validation Failed")
	}
	// add the PlayerPrefrence to the user
	user.PlayerPrefrences = append(user.PlayerPrefrences, req.Pref)
	// save the user
	res = c.H.DB.Save(&user)
	if res.Error != nil {
		Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	Response.OK(w, user)
	size := Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

func (c *PlayerPrefrenceController) AuthPostPlayerPrefrencesToUser(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	c.H.ForcePOST(w, r)

	// set the Detailed Response to be a User
	var Response models.DetailedResponse[models.User]

	// authenticate the header
	token, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		Response.ConsumeError(w, err, http.StatusUnauthorized)
		size := Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	c.H.SetAuthToken(token)
	var user models.User
	res := c.H.DB.First(&user, c.H.AuthToken.UserID)
	if res.Error != nil {
		Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}
	// get the []PlayerPrefrence from the body
	var req models.PrefrencesRequest
	// Decode the r.Body into the model
	body := r.Body
	err = json.NewDecoder(body).Decode(&req)
	if err != nil {
		Response.ConsumeError(w, err, http.StatusBadRequest)
		size := Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Bad Request", err)
		return
	} else if req.UserID != c.H.AuthToken.UserID {
		Response.UDRWrite(w, http.StatusBadRequest, "Validation Failed", false)
		size := Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Validation Failed")
	}
	// add the PlayerPrefrence to the user
	user.PlayerPrefrences = append(user.PlayerPrefrences, req.Prefs...)
	// save the user
	res = c.H.DB.Save(&user)
	if res.Error != nil {
		Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	Response.OK(w, user)
	size := Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

func (c *PlayerPrefrenceController) AuthGetPrefrences(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	c.H.ForceGET(w, r)

	// check if the user is authenticated
	token, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		c.H.ResponseList.ConsumeError(w, err, http.StatusUnauthorized)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	c.H.SetAuthToken(token)
	// get the user
	var user models.User
	res := c.H.DB.First(&user, c.H.AuthToken.UserID)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}
	// get the prefrences
	var prefrences []models.PlayerPrefrence
	res = c.H.DB.Find(&prefrences)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	c.H.ResponseList.OK(w, prefrences)
	size := c.H.ResponseList.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

// #endregion

// #region Enum Relationships Definitions
func (c *RelsController) AuthGetRelationships(w http.ResponseWriter, r *http.Request) {
	logger := lib.New(r)
	c.H.ForceGET(w, r)

	// authenticate the header
	token, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		c.H.ResponseList.ConsumeError(w, err, http.StatusUnauthorized)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	c.H.SetAuthToken(token)

	// get the relationsips
	var relationships []models.Relationship
	res := c.H.DB.Find(&relationships)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	c.H.ResponseList.OK(w, relationships)
	size := c.H.ResponseList.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

// #endregion
