package main

import (
	"github.com/gin-gonic/gin"
	"github.com/amir-yaghoobi/accountManager/controllers"
	"github.com/amir-yaghoobi/accountManager/middlewares"
)


func getApiRoutes() (router *gin.Engine) {
	router = gin.Default()

	authMiddleWare := middlewares.SetupJWT()
	router.POST("/auth",                                  authMiddleWare.LoginHandler)

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
		incomeGroup.POST("/",                             controllers.AddIncome           )
		incomeGroup.GET("/:accountId",                    controllers.NotImplementedYet   )
		incomeGroup.GET("/:accountId/category",           controllers.GetAccountCategories)
		incomeGroup.GET("/:accountId/category/:cId",      controllers.GetAccountCategories)
		incomeGroup.POST("/:accountId/category",          controllers.AddCategory         )
	}

	expenseGroup := router.Group("/expense")
	{
		expenseGroup.Use(authMiddleWare.MiddlewareFunc())
		expenseGroup.POST("/",                            controllers.AddExpense          )
		expenseGroup.GET("/:accountId",                   controllers.NotImplementedYet   )
		expenseGroup.GET("/:accountId/category",          controllers.GetAccountCategories)
		expenseGroup.GET("/:accountId/category/:cId",     controllers.GetAccountCategories)
		expenseGroup.POST("/:accountId/category",         controllers.AddCategory         )
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