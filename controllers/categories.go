package controllers

import (
	"strconv"
	"strings"
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin/binding"
	"github.com/amir-yaghoobi/accountManager/db"
	"github.com/amir-yaghoobi/accountManager/models"
)

type AddCategoryRequest struct {
	Name     string    `json:"name"        binding:"required"`
	Icon     string    `json:"icon"        binding:"required"`
	ParentId uint      `json:"parentId"                      `
}


func isIncomeCategory(pathUrl string) bool {
	return strings.Index(pathUrl, "/income/") > -1
}

// Add a new category to account
//
// params:
//		* name: string
//		* icon: string
//		* parentId: uint request (optional)
func AddCategory(c *gin.Context) {
	accountIdString := c.Param("accountId")
	accountId, err := strconv.ParseUint(accountIdString, 10, 64)
	if err != nil {
		log.Warnf("invalid accountId, error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": "AccountId must be an integer!",
		})
		return
	}

	var newCategoryRequest AddCategoryRequest
	err = c.MustBindWith(&newCategoryRequest, binding.JSON)
	if err != nil { // missing required fields
		log.Errorf("error on binding to AddCategoryRequest, error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": err.Error(),
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


	if isIncomeCategory(c.Request.URL.Path) {
		category := models.IncomeCategory{
			AccountID: uint(accountId),
			Icon: newCategoryRequest.Icon,
			Name: newCategoryRequest.Name,
		}


		if newCategoryRequest.ParentId > 0 {
			parentCategory          := models.IncomeCategory{}
			parentCategory.ID        = newCategoryRequest.ParentId
			parentCategory.AccountID = uint(accountId)

			q := pConn.Where(&parentCategory).First(&parentCategory)
			if q.Error != nil && q.Error != gorm.ErrRecordNotFound {
				log.Errorf("cannot get category:%d from postgres, error:%s",
					newCategoryRequest.ParentId, q.Error.Error())

				c.JSON(http.StatusInternalServerError, gin.H{
					"status": http.StatusInternalServerError,
					"error": internalServerError,
				})
				return

			} else if q.Error == gorm.ErrRecordNotFound {
				log.Warnf("category:%d does not belongs to account:%d")
				c.JSON(http.StatusNotFound, gin.H{
					"status": http.StatusNotFound,
					"error": "Parent category does not exist",
				})
				return
			}
			category.ParentCategory = newCategoryRequest.ParentId
		}

		if q := pConn.Save(&category); q.Error != nil {
			log.Errorf("cannot save category into postgres, error:%s\n", q.Error.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error": internalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"category": category,
		})
	} else {
		category := models.ExpenseCategory{
			AccountID: uint(accountId),
			Icon: newCategoryRequest.Icon,
			Name: newCategoryRequest.Name,
		}


		if newCategoryRequest.ParentId > 0 {
			parentCategory          := models.ExpenseCategory{}
			parentCategory.ID        = newCategoryRequest.ParentId
			parentCategory.AccountID = uint(accountId)

			q := pConn.Where(&parentCategory).First(&parentCategory)
			if q.Error != nil && q.Error != gorm.ErrRecordNotFound {
				log.Errorf("cannot get category:%d from postgres, error:%s",
					newCategoryRequest.ParentId, q.Error.Error())

				c.JSON(http.StatusInternalServerError, gin.H{
					"status": http.StatusInternalServerError,
					"error": internalServerError,
				})
				return

			} else if q.Error == gorm.ErrRecordNotFound {
				log.Warnf("category:%d does not belongs to account:%d")
				c.JSON(http.StatusNotFound, gin.H{
					"status": http.StatusNotFound,
					"error": "Parent category does not exist",
				})
				return
			}
			category.ParentCategory = newCategoryRequest.ParentId
		}

		if q := pConn.Save(&category); q.Error != nil {
			log.Errorf("cannot save category into postgres, error:%s\n", q.Error.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error": internalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"category": category,
		})
	}

}


// fetch all parent categories of give account
//
// params:
//		* accountId: uint url param
//		* user: models.User  request
//		* parentId: uint request (optional)
func GetAccountCategories(c *gin.Context) {
	accountIdString := c.Param("accountId")
	accountId, err := strconv.ParseUint(accountIdString, 10, 64)
	if err != nil {
		log.Warnf("invalid accountId, error: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": "AccountId must be an integer!",
		})
		return
	}

	catIdString := c.Param("cId")
	var catId uint64 = 0
	if len(catIdString) > 0 {
		catId, err = strconv.ParseUint(catIdString, 10, 64)
		if err != nil {
			log.Warnf("invalid accountId, error: %s\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error": "CategoryId must be an integer!",
			})
			return
		}
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

	var query string
	if catId > 0 {
		query = "account_id = ? AND ID = ?"
	} else {
		query = "account_id = ? AND parent_category = ?"
	}

	isIncome := isIncomeCategory(c.Request.URL.Path)
	if isIncome {
		var categories []models.IncomeCategory
		pConn.Preload("SubCategories").Where(query, uint(accountId), catId).Find(&categories)
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"categories": categories,
		})
	} else {
		var categories []models.ExpenseCategory
		pConn.Preload("SubCategories").Where(query, uint(accountId), catId).Find(&categories)
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"categories": categories,
		})
	}
}