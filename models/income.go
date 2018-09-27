package models

import "github.com/jinzhu/gorm"

type Income struct {
	gorm.Model
	Amount uint
	Account Account                     `json:"-"`
	AccountID uint                      `json:"-"`
	User User                           `json:"-"`
	UserID uint                         `json:"-"`
	IncomeCategory IncomeCategory       `json:"-"`
	IncomeCategoryID uint
}