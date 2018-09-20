package models

import "github.com/jinzhu/gorm"

type Budget struct {
	gorm.Model
	Name string
	Amount uint
	Period uint
	Account Account
	AccountID uint
	User User        `gorm:"foreignkey:CreatedBy"`
	CreatedBy uint
	ExpenseCategory ExpenseCategory
	ExpenseCategoryID uint
}