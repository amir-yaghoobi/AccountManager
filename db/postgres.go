package db

import (
	"github.com/jinzhu/gorm"
	"github.com/amir-yaghoobi/accountManager/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/amir-yaghoobi/accountManager/config"
)


var pgDB *gorm.DB


func InitializePostgres(db *gorm.DB) {
	db.AutoMigrate(&models.Budget{})
	db.AutoMigrate(&models.Income{})
	db.AutoMigrate(&models.Expense{})
	db.AutoMigrate(&models.IncomeCategory{})
	db.AutoMigrate(&models.ExpenseCategory{})
	db.AutoMigrate(&models.User{}, &models.Account{})
}

func GetPostgres() (*gorm.DB, error) {
	if pgDB != nil {
		return pgDB, nil
	}

	cfg := config.GetConfig()

	var err error
	pgDB, err = gorm.Open("postgres", cfg.PostgresDB.GetConnectionString())
	if err != nil {
		return nil, err
	}

	return pgDB, nil
}