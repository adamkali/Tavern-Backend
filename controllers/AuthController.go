package controllers

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type AuthController struct {
	H BaseHandler[models.AuthToken]
	C models.AuthEmailConfiglette
}

func NewAuthController(DB *gorm.DB, C models.AuthEmailConfiglette) *AuthController {
	return &AuthController{
		H: *NewHandler(DB, models.AuthToken{}, "auth"),
		C: C,
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	c.H.AuthToken = models.AuthToken{}

	logger := lib.New(r)
	c.H.ForcePOST(w, r)

	// get the email and password from the request body
	var req models.LoginRequest
	body := r.Body
	if r.Header.Get("Content-Type") != "application/json" {
		c.H.Response.UDRWrite(w, http.StatusBadRequest, "Invalid Content-Type", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Invalid Content-Type")
		return
	}
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		c.H.Response.ConsumeError(w, err, http.StatusBadRequest)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	defer body.Close()

	// generate the hash
	hash := c.H.AuthToken.GenAuth(req.Username, req.Password)
	res := c.H.DB.Where("auth_hash = ?", hash).First(&c.H.AuthToken)
	if res.Error != nil {
		c.H.Response.UDRWrite(w, http.StatusUnauthorized, "Invalid credentials", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Invalid credentials", res.Error)
		return
	}

	// OK Response
	c.H.Response.OK(w, c.H.AuthToken)
	size := c.H.Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")

}

func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	c.H.AuthToken = models.AuthToken{}

	logger := lib.New(r)
	c.H.ForcePOST(w, r)

	// Get the request body
	var req models.AuthRequest
	body := r.Body
	if r.Header.Get("Content-Type") != "application/json" {
		c.H.Response.UDRWrite(w, http.StatusBadRequest, "Invalid Content-Type", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Invalid Content-Type")
		return
	}
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		c.H.Response.ConsumeError(w, err, http.StatusBadRequest)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	defer body.Close()

	// Check if the username is taken
	res := c.H.DB.Where("username = ?", req.Username).First(&c.H.AuthToken)
	if res.Error == nil {
		c.H.Response.UDRWrite(w, http.StatusConflict, "Username is taken", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusConflict, "Username is taken")
		return
	}

	// check if the email is taken
	res = c.H.DB.Where("user_email = ?", req.UserEmail).First(&c.H.AuthToken)
	if res.Error == nil {
		c.H.Response.UDRWrite(w, http.StatusConflict, "Email is taken", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusConflict, "Email is taken")
		return
	}

	// check if the password is valid
	if !c.H.AuthToken.ValidatePassword(req.Password) {
		c.H.Response.UDRWrite(w, http.StatusBadRequest, "Invalid password", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Invalid password")
		return
	}

	// generate the hash
	c.H.AuthToken.ID = (generateUUID())
	c.H.AuthToken.GenerateToken(req.Username, req.Password, req.UserEmail)
	c.H.AuthToken.RoleFK = "4915B1FE0F7643F692FC25B3A60CC762"

	c.H.Response.Data = c.H.AuthToken
	res = c.H.DB.Preload("Role").Create(&c.H.AuthToken)
	if res.Error != nil {
		c.H.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}
	// create a new user for the token
	user := models.User{
		ID:       generateUUID(),
		Username: req.Username,
		Bio:      "Write a bio, to tell everyone about you!",
	}
	res = c.H.DB.Create(&user)
	if res.Error != nil {
		c.H.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	// add the user to the token
	c.H.AuthToken.UserID = user.ID

	// Give the Token a role of "Not Verified" by default

	// update the auth token with the role
	res = c.H.DB.Save(&c.H.AuthToken)
	if res.Error != nil {
		c.H.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	// make a new auth activation token
	aat := models.AuthTokenActivation{
		ID:        generateUUID(),
		AuthPin:   generatePin(),
		AuthEmail: req.UserEmail,
		AuthID:    c.H.AuthToken.ID,
	}
	res = c.H.DB.Create(&aat)
	if res.Error != nil {
		c.H.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}
	aat.SendRegistrationEmail(c.C)

	// OK Response
	c.H.Response.OK(w, c.H.AuthToken)
	size := c.H.Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

func (c *AuthController) Verify(w http.ResponseWriter, r *http.Request) {
	c.H.AuthToken = models.AuthToken{}
	logger := lib.New(r)
	c.H.ForcePOST(w, r)
	var authActivation models.AuthTokenActivation

	// ensure application/json
	if r.Header.Get("Content-Type") != "application/json" {
		c.H.Response.UDRWrite(w, http.StatusBadRequest, "Invalid Content-Type", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Invalid Content-Type")
		return
	}
	// get the request body
	var req models.VerifyRequest
	body := r.Body
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		c.H.Response.ConsumeError(w, err, http.StatusBadRequest)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	defer body.Close()
	fmt.Println(req)

	// get the token from the database
	res := c.H.DB.Where("auth_pin = ? AND auth_email = ?", req.Pin, req.UserEmail).First(&authActivation)
	if res.Error != nil {
		c.H.Response.UDRWrite(w, http.StatusUnauthorized, "Invalid credentials", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Invalid credentials", res.Error)
		return
	}
	// get the auth token from the database using the auth_hash from the authActivation
	res = c.H.DB.Where("id = ?", authActivation.AuthID).Preload("Role").First(&c.H.AuthToken)
	if res.Error != nil {
		c.H.Response.UDRWrite(w, http.StatusUnauthorized, "Invalid credentials", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, "Invalid credentials", res.Error)
		return
	}

	// check if the token is already verified
	if c.H.AuthToken.Active {
		c.H.Response.UDRWrite(w, http.StatusConflict, "Account already verified", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusConflict, "Account already verified")
		return
	}

	c.H.AuthToken.Active = true
	c.H.AuthToken.RoleFK = "747A97752DA547348E21E93DAF207A43"
	// get the new role from the database
	c.H.AuthToken.Role = models.Role{}
	fmt.Printf("%s", c.H.AuthToken.RoleFK)
	res = c.H.DB.Where("id = ?", c.H.AuthToken.RoleFK).First(&c.H.AuthToken.Role)
	if res.Error != nil {
		c.H.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	// update the auth token with the role
	res = c.H.DB.Preload("Role").Save(&c.H.AuthToken)
	if res.Error != nil {
		c.H.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	// delete the activation token
	res = c.H.DB.Delete(&authActivation)
	if res.Error != nil {
		c.H.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, "Internal Server Error", res.Error)
		return
	}

	// OK Response
	c.H.Response.OK(w, c.H.AuthToken)
	size := c.H.Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}
