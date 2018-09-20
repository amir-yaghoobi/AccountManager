package models

import "github.com/jinzhu/gorm"

type Expense struct {
	gorm.Model
	Amount uint
	Account Account
	AccountID uint
	User User
	UserID uint
	ExpenseCategory ExpenseCategory
	ExpenseCategoryID uint
}