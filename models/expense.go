package models

import "github.com/jinzhu/gorm"

type Expense struct {
	gorm.Model
	Amount uint
	Account Account                     `json:"-"`
	AccountID uint                      `json:"-"`
	User User                           `json:"-"`
	UserID uint                         `json:"-"`
	ExpenseCategory ExpenseCategory     `json:"-"`
	ExpenseCategoryID uint
}