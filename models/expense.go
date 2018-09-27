package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

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

func SumExpenses(expenses []Expense, from time.Time) uint {
	sum := uint(0)
	for _, expense := range expenses {
		if expense.CreatedAt.Unix() > from.Unix() {
			sum += expense.Amount
		}
	}
	return sum
}

func MustValuableExpense(expenses []Expense) uint {
	categories := make(map[uint]uint)
	for _, income := range expenses {
		categories[income.ExpenseCategoryID] += income.Amount
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