package controllers

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin/binding"
	"github.com/amir-yaghoobi/accountManager/db"
	"github.com/amir-yaghoobi/accountManager/models"
)

type AddIncomeRequest  struct {
	Amount     uint    `json:"amount"      binding:"required"`
	AccountId  uint    `json:"accountId"   binding:"required"`
	CategoryId uint    `json:"categoryId"  binding:"required"`
}


func AddIncome(c *gin.Context) {
	newIncomeRequest := AddIncomeRequest{}

	err := c.MustBindWith(&newIncomeRequest, binding.JSON)
	if err != nil { // missing required fields
		log.Errorf("error on binding to incomeRequest, error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}

	user, aborted := abortOnInvalidAccount(c, newIncomeRequest.AccountId)
	if aborted {
		return
	}


	pConn, err := db.GetPostgres()
	if err != nil {
		postgresErrorHandler(c, err)
		return
	}

	var category models.IncomeCategory
	category.AccountID = newIncomeRequest.AccountId
	category.ID = newIncomeRequest.CategoryId

	q := pConn.Where(&category).First(&category)
	if q.Error != nil && q.Error != gorm.ErrRecordNotFound {
		log.Errorf("cannot get category:%d from postgres database, error: %s\n", q.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error": internalServerError,
		})
		return

	} else if q.Error == gorm.ErrRecordNotFound {
		log.Warnf("category:%d does not exist. request from userId: %d", newIncomeRequest.CategoryId, user.ID)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error": "Category does not exist",
		})
		return
	}

	income := models.Income{
		UserID:           user.ID,
		Amount:           newIncomeRequest.Amount,
		AccountID:        newIncomeRequest.AccountId,
		IncomeCategoryID: category.ID,
	}

	pConn.Save(&income)
	if q := pConn.Save(&category); q.Error != nil {
		log.Errorf("cannot save new income into postgres, error:%s\n", q.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error": internalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"income": income,
	})
}


