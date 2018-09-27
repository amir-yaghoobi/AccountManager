package models

import "github.com/jinzhu/gorm"

type ExpenseCategory struct {
	gorm.Model
	Name string
	Icon string
	Account Account                 `json:"-"`
	AccountID uint                  `json:"-"`
	ParentCategory uint
	SubCategories []ExpenseCategory `gorm:"foreignkey:ParentCategory"`
}

type IncomeCategory struct {
	gorm.Model
	Name string
	Icon string
	Account Account                 `json:"-"`
	AccountID uint                  `json:"-"`
	ParentCategory uint
	SubCategories []IncomeCategory  `gorm:"foreignkey:ParentCategory"`
}