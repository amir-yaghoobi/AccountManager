package models

import "github.com/jinzhu/gorm"

type Income struct {
	gorm.Model
	Amount uint
	Account Account
	AccountID uint
	User User
	UserID uint
	IncomeCategory IncomeCategory
	IncomeCategoryID uint
}