package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"time"
	"github.com/amir-yaghoobi/accountManager/db"
	"github.com/amir-yaghoobi/accountManager/models"
	"github.com/amir-yaghoobi/accountManager/controllers"
	"github.com/amir-yaghoobi/accountManager/middlewares"
)

func apiStatusController(c *gin.Context) {
	db, err := db.GetPostgres()
	if err != nil {
		log.Errorf("cannot connect to postgres!, error: %s\n", err.Error())
	}

	var usersCount uint
	db.Model(&models.User{}).Count(&usersCount)

	var accountsCount uint
	db.Model(&models.Account{}).Count(&accountsCount)

	uptime := time.Now()
	duration := uptime.Sub(serverStartedAt)
	response := map[string]interface{}{
		"uptime":   duration.Seconds(),
		"health":   "OK",
		"users":    usersCount,
		"accounts": accountsCount,
	}
	c.JSON(200, response)
}

func getApiRoutes() (router *gin.Engine) {
	router = gin.Default()

	router.GET("/", apiStatusController)

	authMiddleWare := middlewares.SetupJWT()
	router.POST("/auth",                                 authMiddleWare.LoginHandler)

	userGroup := router.Group("/user")
	userGroup.POST("/register",                           controllers.Register)
	{
		userGroup.Use(authMiddleWare.MiddlewareFunc())
		userGroup.GET("/:userId",                         controllers.GetUser)
	}

	accountGroup := router.Group("/account")
	{
		accountGroup.Use(authMiddleWare.MiddlewareFunc())
		accountGroup.GET("/",                             controllers.NotImplementedYet)
		accountGroup.POST("/",                            controllers.CreateNewAccount )
		accountGroup.POST("/inviteUser",                  controllers.NotImplementedYet)
		accountGroup.DELETE("/:accountId",                controllers.NotImplementedYet)
	}

	incomeGroup := router.Group("/income")
	{
		incomeGroup.Use(authMiddleWare.MiddlewareFunc())
		incomeGroup.POST("/",                             controllers.AddIncome        )
		incomeGroup.GET("/:accountId",                    controllers.NotImplementedYet)
		incomeGroup.GET("/:accountId/category",           controllers.GetAccountCategories)
		incomeGroup.GET("/:accountId/category/:cId",      controllers.GetAccountCategories)
		incomeGroup.POST("/:accountId/category",          controllers.AddCategory      )
	}

	expenseGroup := router.Group("/expense")
	{
		expenseGroup.Use(authMiddleWare.MiddlewareFunc())
		expenseGroup.POST("/",                            controllers.NotImplementedYet)
		expenseGroup.GET("/:accountId",                   controllers.NotImplementedYet)
		expenseGroup.GET("/:accountId/category",          controllers.NotImplementedYet)
		expenseGroup.GET("/:accountId/category/:cId",     controllers.NotImplementedYet)
		expenseGroup.POST("/:accountId/category",         controllers.NotImplementedYet)
	}

	budgetGroup := router.Group("/budget")
	{
		budgetGroup.Use(authMiddleWare.MiddlewareFunc())
		budgetGroup.POST("/",                             controllers.NotImplementedYet)
		budgetGroup.GET("/:accountId",                    controllers.NotImplementedYet)
		budgetGroup.PUT("/:accountId/:budgetId",          controllers.NotImplementedYet)
		budgetGroup.DELETE("/:accountId/:budgetId",       controllers.NotImplementedYet)
	}

	return router
}