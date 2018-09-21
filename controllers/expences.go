package controllers

type AddExpenseRequest  struct {
	Amount     uint    `json:"amount"      binding:"required"`
	AccountId  uint    `json:"accountId"   binding:"required"`
	CategoryId uint    `json:"categoryId"  binding:"required"`
}

