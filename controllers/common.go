package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/amir-yaghoobi/accountManager/models"
)

const internalServerError = "An internal server error happened, please try again!"


func getUserFromContext(c *gin.Context) (user *models.User, isAborted bool) {
	userInterface, exist := c.Get("USER")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": "login process is not completed, please login again!",
		})
		return nil, true
	}
	user, ok := userInterface.(*models.User)
	if !ok {
		log.Warnf("gin USER key is not a model.User instance, wtf? ", userInterface)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error": internalServerError,
		})
		return nil, true
	}

	return user, false
}



func NotImplementedYet(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "not implemented yet",
	})
}

func postgresErrorHandler(c *gin.Context, error error) {
	if error != nil {
		log.Errorf("cannot connect to postgres databases, error: %s\n", error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": internalServerError,
		})
	}
}