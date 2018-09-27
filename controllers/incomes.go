package controllers

import (
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin/binding"
	"github.com/amir-yaghoobi/accountManager/db"
	"github.com/amir-yaghoobi/accountManager/models"
	"strconv"
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


// param accountId
//
// options:
//   catId integer
//   limit integer default 25
//   offset integer default 0
//
// example:
//  ?catId=4&limit=25&offset=0
func GetIncomes(c *gin.Context) {
	accountIdString := c.Param("accountId")
	accountId, err := strconv.ParseUint(accountIdString, 10, 64)
	if err != nil {
		log.Warnf("invalid accountId, error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": invalidAccountId,
		})
		return
	}

	_, aborted := abortOnInvalidAccount(c, uint(accountId))
	if aborted {
		return
	}

	pConn, err := db.GetPostgres()
	if err != nil {
		postgresErrorHandler(c, err)
		return
	}

	where := models.Income{AccountID: uint(accountId)}

	catIdString := c.Query("catId")
	if len(catIdString) > 0 {
		if catId, err := strconv.ParseUint(catIdString, 10, 64); err == nil {
			where.IncomeCategoryID = uint(catId)
		} else {
			log.Warnf("invalid categoryId:%s query proceed without categoryId, error:%s", catId, err.Error())
		}
	}

	limit := parseUintWithDefault(c.Query("limit"), 25)
	offset := parseUintWithDefault(c.Query("offset"), 0)

	var expenses []models.Income
	q := pConn.Where(&where).Limit(limit).Offset(offset).Find(&expenses)
	if q.Error != nil && q.Error != gorm.ErrRecordNotFound {
		log.Errorf("cannot get expenses from postgres, error:%s", q.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error": internalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, expenses)
}