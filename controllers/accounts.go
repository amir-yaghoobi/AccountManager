package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin/binding"
	"github.com/amir-yaghoobi/accountManager/db"
	"github.com/amir-yaghoobi/accountManager/models"
)

type CreateAccountRequest struct {
	Name         string   `json:"name"        binding:"required"`
	Icon         string   `json:"icon"        binding:"required"`
	Description  string   `json:"description" binding:"required"`
}

/*
	@params
		* name: string
		* icon: string
 */
func CreateNewAccount(c *gin.Context) {
	createAccountForm := CreateAccountRequest{}

	err := c.MustBindWith(&createAccountForm, binding.JSON)
	if err != nil { // missing required fields
		log.Errorf("error on binding to createAccountForm, error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}

	user, aborted := getUserFromContext(c)
	if aborted {
		return
	}

	pConn, err := db.GetPostgres()
	if err != nil {
		postgresErrorHandler(err, c)
		return
	}

	account := &models.Account{
		Name:         createAccountForm.Name,
		Icon:         createAccountForm.Icon,
		Description:  createAccountForm.Description,
	}

	q := pConn.Save(account)
	if q.Error != nil {
		log.Errorf("error on saving new account %q\n", account)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error": internalServerError,
		})
		return
	}
	go pConn.Model(&account).Association("Users").Append([]models.User{*user})

	c.JSON(http.StatusOK, account)
}


func RemoveAccount(c *gin.Context) {
	// TODO remove all expense and incomes from account
	c.JSON(200, gin.H{
		"status": "not implemented yet",
	})
}


func InviteUserToAccount(c *gin.Context) {
	// TODO create a new table for this
	// invite from user to user
	// account
	c.JSON(200, gin.H{
		"status": "not implemented yet",
	})
}

