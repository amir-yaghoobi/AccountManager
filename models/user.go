package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	FirstName string	 `json:"firstname"`
	LastName string		 `json:"lastname"`
	Username string      `gorm:"type:varchar(100);unique_index" json:"username"`
	Password string      `gorm:"not null;" json:"password"`
	Avatar string		 `json:"avatar"`
	Accounts []Account   `gorm:"many2many:user_accounts;" json:"accounts"`
	LastAccess time.Time `json:"last_access"`
	LastIp string        `json:"last_ip"`
}

func (user *User) HashPassword(password string) (error) {
	saltBytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(saltBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	hashBytes := []byte(user.Password)
	passwordBytes  := []byte(password)
	err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	return err == nil
}

func (user *User) GetAccount(accountId uint) *Account {
	for _, account := range user.Accounts {
		if account.ID == accountId {
			return &account
		}
	}
	return nil
}