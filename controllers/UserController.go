package controllers

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"net/http"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type UserController struct{ H BaseHandler[models.User] }

func NewUserController(DB *gorm.DB) *UserController {
	return &UserController{
		H: *NewHandler(DB, models.User{}, "User"),
	}
}

// #region Users
func (c *UserController) AdminGetAll(w http.ResponseWriter, r *http.Request) {
	c.H.Model = models.User{}

	logger := lib.New(r)
	c.H.ForceGET(w, r)

	// first check authenticated
	Auth, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		c.H.ResponseList.ConsumeError(w, err, http.StatusUnauthorized)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusUnauthorized, err.Error())
		return
	}
	c.H.SetAuthToken(Auth)

	// then check admin
	if !c.H.AuthToken.IsAdmin() {
		c.H.ResponseList.UDRWrite(w, http.StatusForbidden, "Forbidden", false)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusForbidden, "Forbidden")
		return
	}

	// then get all
	var ms []models.User
	res := c.H.DB.Find(&ms)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, res.Error.Error())
		return
	}

	// OK Response
	c.H.ResponseList.OK(w, ms)
	size := c.H.ResponseList.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

func (c *UserController) UserQueue(w http.ResponseWriter, r *http.Request) {
	c.H.Model = models.User{}

	logger := lib.New(r)
	c.H.ForceGET(w, r)

	// first check Authenticated
	Auth, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		c.H.ResponseList.ConsumeError(w, err, http.StatusUnauthorized)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusUnauthorized, err.Error())
		return
	}
	c.H.SetAuthToken(Auth)

	// then check admin
	if !c.H.AuthToken.IsAdmin() {
		c.H.ResponseList.UDRWrite(w, http.StatusForbidden, "Forbidden", false)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusForbidden, "Forbidden")
		return
	}

	// First Set up a new queue of 20 users
	// get 20 users from the database
	var users []models.User
	res := c.H.DB.Limit(20).Where(
		"id != ?", c.H.AuthToken.UserID,
	).Find(&users)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, res.Error.Error())
		return
	}

	// make a go routine for each user
	// make a wait group
	// wait for all go routines to finish
	var queue [20]chan models.User
	errChan := make(chan error)
	for i := range queue {
		queue[i] = make(chan models.User)
	}
	var wg sync.WaitGroup
	for i := 0; i < len(users); i++ {
		wg.Add(1)
		go func(user models.User) {
			defer wg.Done()
			var relat models.UserRelationship
			res := c.H.DB.Where(
				"Self = ? AND Other = ?",
				user.ID, c.H.AuthToken.UserID,
			).First(&relat)
			if res.Error != nil {
				errChan <- res.Error
				return
			}
			if !relat.Relationship.Negative {
				queue[i] <- user
			} else {
				// get a new user from the database
				// that is not the c.H.AuthToken.User
				var newUser models.User
				res := c.H.DB.Where(
					"id != ? AND id != ?",
					user.ID, c.H.AuthToken.UserID,
				).First(&newUser)
				if res.Error != nil {
					errChan <- res.Error
					return
				}
				// replace user with newUser
				users[i] = newUser
			}
		}(users[i])
	}
	wg.Wait()

	// check if errChan has any errors
	// if it does use c.H.Response.ConsumeError
	// to write the error to the Response
	// and return
	var ret []models.User
	select {
	case err := <-errChan:
		c.H.ResponseList.ConsumeError(w, err, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, err.Error())
		return
	default:
		// append the users in the queue
		// to the Response
		for _, q := range queue {
			select {
			case user := <-q:
				ret = append(ret, user)
			default:
				continue
			}
		}
	}

	c.H.ResponseList.OK(w, ret)
	size := c.H.ResponseList.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

// #endregion

// #region User
func (c *UserController) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	c.H.Model = models.User{}

	logger := lib.New(r)
	c.H.ForceGET(w, r)

	// first check Authenticated
	Auth, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		c.H.Response.ConsumeError(w, err, http.StatusUnauthorized)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, err.Error())
		return
	}
	c.H.SetAuthToken(Auth)
	// then get the user linking in all characters and plots
	// that have that user.ID as user_id
	var m models.User
	res := c.H.DB.Where(
		"id = ?", c.H.AuthToken.UserID,
	).Preload(
		"Characters").Preload(
		"Plots").Preload(
		"Tags").Preload(
		"PlayerPreferences").First(&m)
	if res.Error != nil {
		c.H.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, res.Error.Error())
		return
	} else if res.RowsAffected == 0 {
		c.H.Response.UDRWrite(w, http.StatusNotFound, "User not found", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusNotFound, "User not found")
		return
	}

	c.H.Response.OK(w, m)
	size := c.H.Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

func (c *UserController) AuthGetByIDFull(w http.ResponseWriter, r *http.Request) {
	c.H.Model = models.User{}

	logger := lib.New(r)
	c.H.ForceGET(w, r)

	// first check Authenticated
	Auth, err := verifyAuthorizationToken(*c.H.DB, r)
	if err != nil {
		c.H.Response.ConsumeError(w, err, http.StatusUnauthorized)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusUnauthorized, err.Error())
		return
	}
	c.H.SetAuthToken(Auth)

	// get the user.ID from the url /api/auth/User/:id
	id := strings.Split(r.URL.Path, "/")[4]
	if id == "" || len(id) != 32 {
		c.H.Response.UDRWrite(w, http.StatusBadRequest, "Bad Request", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusBadRequest, "Bad Request")
		return
	}

	// then get the user linking in all characters and plots
	// that have that user.ID as user_id
	var m models.User
	res := c.H.DB.Where(
		"id = ?", id,
	).Preload(
		"Characters").Preload(
		"Plots").Preload(
		"Tags").Preload(
		"PlayerPreferences").First(&m)
	if res.Error != nil {
		c.H.Response.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusInternalServerError, res.Error.Error())
		return
	} else if res.RowsAffected == 0 {
		c.H.Response.UDRWrite(w, http.StatusNotFound, "User not found", false)
		size := c.H.Response.SizeOf()
		logger.Log(size, http.StatusNotFound, "User not found")
		return
	}

	c.H.Response.OK(w, m)
	size := c.H.Response.SizeOf()
	logger.Log(size, http.StatusOK, "OK")
}

// #endregion
