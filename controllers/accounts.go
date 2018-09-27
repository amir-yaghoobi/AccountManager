package controllers

import (
	"time"
	"strconv"
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


func abortOnInvalidAccount(c *gin.Context, accountId uint) (user *models.User, isAborted bool) {
	user, aborted := getUserFromContext(c)
	if aborted {
		return nil, true
	}

	account := user.GetAccount(uint(accountId))
	if account == nil {
		log.Warnf("userID:%d attempts to add income to accountId:%d," +
			" but account does not belongs to him",
			user.ID, accountId)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error": "Account not found",
		})
		return nil, true
	}
	return user, false
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
		postgresErrorHandler(c, err)
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


 // return list of accounts that related to this user
func GetUserAccounts(c *gin.Context) {
	user, aborted := getUserFromContext(c)
	if aborted {
		return
	}

	accounts := user.Accounts
	c.JSON(http.StatusOK, accounts)
}


func DashboardStats(c *gin.Context) {
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

	now := time.Now()
	lastDay := now.AddDate(0, 0, -1)
	lastWeek := now.AddDate(0, 0, -7)
	lastMonth := now.AddDate(0, -1, 0)

	var incomes []models.Income
	pConn.Where("created_at > ? AND account_id = ?", lastMonth, accountId).Find(&incomes)

	todayIncome := models.SumIncomes(incomes, lastDay)
	weekIncome  := models.SumIncomes(incomes, lastWeek)
	monthIncome := models.SumIncomes(incomes, lastMonth)
	mvcIncome   := models.MustValuableIncome(incomes)

	var expenses []models.Expense
	pConn.Where("created_at > ? AND account_id = ?", lastMonth, accountId).Find(&expenses)

	todayExpense := models.SumExpenses(expenses, lastDay)
	weekExpense  := models.SumExpenses(expenses, lastWeek)
	monthExpense := models.SumExpenses(expenses, lastMonth)
	mvcExpense   := models.MustValuableExpense(expenses)

	c.JSON(http.StatusOK, gin.H{
		"incomes": gin.H{
			"count": len(incomes),
			"today": todayIncome,
			"week":  weekIncome,
			"month": monthIncome,
			"mvcId": mvcIncome,
		},
		"expenses": gin.H{
			"count": len(expenses),
			"today": todayExpense,
			"week":  weekExpense,
			"month": monthExpense,
			"mvcId": mvcExpense,
		},
	})
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

