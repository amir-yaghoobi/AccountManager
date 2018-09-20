package models

import (
	"github.com/jinzhu/gorm"
)


type InviteAccount struct {
	gorm.Model
	from User // pk
	to User // pk
	Account Account // pk
	status uint // 0 pending, 1 approved, 2 rejected
}