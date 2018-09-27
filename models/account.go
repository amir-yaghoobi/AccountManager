package models

import (
	"github.com/jinzhu/gorm"
)


type Account struct {
	gorm.Model
	Name        string
	Description string
	Icon        string
	Users       []*User `gorm:"many2many:user_accounts;" json:"-"`
}