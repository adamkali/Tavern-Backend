package controllers

import (
	"Tavern-Backend/lib"
	"Tavern-Backend/models"
	"errors"
	"fmt"
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

// TODO: Make DB Requests To Be more convienient
// like: c.H.DB.Preload("Prefrences").First(&user) ==> UserRepo.GetOne(id)
//	 	c.H.DB.Preload("Prefrences").Find(&users) ==> UserRepo.GetAll()
//   	In this way, we can make the code more readable and less error prone
//		and more readable and implementable
// :ENDTODO

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
	/* TESTME: Test this function for all cases
		* There will be the following cases:
		* 1. User and The User found in the protoqueue have a relationship -> ROLL
		* 2. The User found the User have no relationship -> ADD
		* 3. The User found the User have a relationship BUT it is positive -> ADD
		* 4. The User found the User have a relationship BUT it is negative -> ROLL
		* 5. LATER ON: The User does not match the found User's preferences -> ROLL
		* 6. LATER ON: The User matches the found User's preferences -> ADD
	:ENDTESTME */

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
	var users []models.User
	res := c.H.DB.Limit(20).Preload("PlayerPrefrence").Where(
		"id != ?", c.H.AuthToken.UserID,
	).Find(&users)
	if res.Error != nil {
		c.H.ResponseList.ConsumeError(w, res.Error, http.StatusInternalServerError)
		size := c.H.ResponseList.SizeOf()
		logger.Log(size, http.StatusInternalServerError, res.Error.Error())
		return
	}
	println("\n\nStarting go routines")
	var queue [20]chan models.User

	errChan := make(chan error)
	for i := range users {
		queue[i] = make(chan models.User, 1)
	}
	userSlice := len(users)

	// make quit channels equal to userSlice
	quit := make([]chan bool, userSlice)

	var wg sync.WaitGroup
	// Tell the 'wg' WaitGroup how many threads/goroutines
	//   that are about to run concurrently.
	wg.Add(userSlice)
	fmt.Println("Running for loop…")
	for i := 0; i < userSlice; i++ {
		userInLoop := users[i]
		quitInLoop := quit[i]
		queueInLoop := queue[i]

		go func(u models.User, q1 chan bool, q2 chan models.User) {

			defer wg.Done()
			// TODO: #5 Check for the prefrences and see if they are compatible\
			// for example if the user is looking for a veteran in experience
			// then we should check if c.H.AuthToken is a veteran or not
			// and if not then we should not return it
			// :ENDTODO
			var rel models.UserRelationship
			res := c.H.DB.Where("self_id = ? AND other_id = ?", c.H.AuthToken.UserID, u.ID).First(&rel)
			if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
				res := c.H.DB.Where("self_id = ? AND other_id = ?", u.ID, c.H.AuthToken.UserID).First(&rel)
				if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
					// add the user to the queue
					q2 <- u
					return
				} else if res.Error != nil {
					errChan <- res.Error
					return
				} else {
					if !rel.Relationship.Negative {
						q2 <- u
						return
					}
				}
			} else if res.Error != nil {
				errChan <- res.Error
				return
			}
			// Since all the checks have failed, the user is not in the queue
			// so we should get a new user from the database and try again
			var user models.User
			res = c.H.DB.Where("id != ?", c.H.AuthToken.UserID).First(&user)
			if res.Error != nil {
				errChan <- res.Error
				return
			}
			u = user
			// This should loop back around and try again
		}(userInLoop, quitInLoop, queueInLoop)
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
				fmt.Printf("Appending %s to Response\n", user.Username)
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

	// get the user.ID from the url /api/auth/User/full/:id
	id := strings.Split(r.URL.Path, "/")[5]
	print(id)
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
		"PlayerPrefrence").First(&m)
	fmt.Printf("Full User:%v\n", m)
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
