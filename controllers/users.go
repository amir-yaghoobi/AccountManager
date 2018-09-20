package controllers

import (
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/amir-yaghoobi/accountManager/db"
	"strconv"
	"github.com/amir-yaghoobi/accountManager/models"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"github.com/jinzhu/gorm"
	"time"
)

// -------------- request struct ----------------
type RegisterRequest struct {
	FirstName string    `json:"firstname" binding:"required"`
	LastName  string    `json:"lastname"  binding:"required"`
	Username  string    `json:"username"  binding:"required"`
	Password  string    `json:"password"  binding:"required"`
	Avatar    string    `json:"avatar"    binding:"required"`
}

func GetUser(c *gin.Context) {
	conn, err := db.GetPostgres()
	if err != nil {
		log.Errorf("cannot connect to postgres database, error: %s\n", err)
		c.JSON(500, gin.H{
			"status": 500,
			"message": "An internal server error happened, please try again!",
		})
		return
	}

	userParam := c.Param("userId")
	userId, err := strconv.ParseInt(userParam, 10, 64)
	if err != nil {
		log.Errorf("cannot parse userId: %s to integer, error: %s\n", userParam, err)
		c.JSON(400, gin.H{
			"status": 400,
			"message": "invalid userId",
		})
		return
	}

	var user models.User
	if res := conn.Preload("Accounts").Find(&user, userId); res.Error != nil {
		log.Errorf("cannot fetch userId: %d, error: %s", userId, res.Error)
		c.JSON(404, gin.H{
			"status": 404,
			"message": "user not found",
		})
		return
	}
	c.JSON(200, user)
}


func Register(c *gin.Context) {

	// ----- receiving POST request and parse values -----
	registerForm := RegisterRequest{}
	err := c.MustBindWith(&registerForm, binding.JSON)
	if err != nil { // missing required fields
		log.Errorf("error on binding to register form, error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	pConn, err := db.GetPostgres()
	if err != nil {
		postgresErrorHandler(c, err)
		return
	}

	queryResult := pConn.Where("Username = ?", registerForm.Username).First(&models.User{})
	if queryResult.Error != nil && queryResult.Error != gorm.ErrRecordNotFound {
		log.Errorf("error on fetching user with username:\"%s\", error: %s\n",
			registerForm.Username, queryResult.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": internalServerError,
		})
		return
	}

	if queryResult.Error != gorm.ErrRecordNotFound { // duplicate username
		c.JSON(http.StatusConflict, gin.H{
			"error": "username is already taken",
		})
		return
	}

	// creating new user
	user := models.User{
		FirstName:  registerForm.FirstName,
		LastName:   registerForm.LastName,
		Username:   registerForm.Username,
		Avatar:     registerForm.Avatar,
		LastIp:     c.ClientIP(),
		LastAccess: time.Now(),
	}

	err = user.HashPassword(registerForm.Password)
	if err != nil {
		log.Errorf("cannot hash user password: \"%s\", error: %s\n", registerForm.Password, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": internalServerError,
		})
		return
	}

	queryResult = pConn.Save(&user)
	if queryResult.Error != nil {
		log.Errorf("cannot save a new user record on postgres, error: %s\n", queryResult.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": internalServerError,
		})
		return
	}

	c.JSON(http.StatusOK,  user)
}
