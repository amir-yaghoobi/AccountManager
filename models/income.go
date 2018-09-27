package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

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

func SumIncomes(incomes []Income, from time.Time) uint {
	sum := uint(0)
	for _, income := range incomes {
		if income.CreatedAt.Unix() > from.Unix() {
			sum += income.Amount
		}
	}
	return sum
}

func MustValuableIncome(incomes []Income) uint {
	categories := make(map[uint]uint)
	for _, income := range incomes {
		categories[income.IncomeCategoryID] += income.Amount
	}

	var maxValue uint = 0
	var maxKey uint = 0
	for key, value := range categories {
		if value > maxValue {
			maxValue = value
			maxKey   = key
		}
	}
	return maxKey
}