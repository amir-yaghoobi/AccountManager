package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/amir-yaghoobi/accountManager/db"
	"time"
	"github.com/amir-yaghoobi/accountManager/models"
	"github.com/jinzhu/gorm"
	"github.com/amir-yaghoobi/accountManager/config"
	"fmt"
)

var serverStartedAt time.Time

func initLogger() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func createSampleUserAndAccounts(conn *gorm.DB) {
	//admin := models.User{
	//	FirstName: "امیرحسین",
	//	LastName: "یعقوبی",
	//	Username: "amir",
	//	Password: "battle2021",
	//	Avatar: "awesome.jpg",
	//	LastIp: "127.0.0.1",
	//	LastAccess: time.Now(),
	//}
	//
	//conn.Save(&admin)
	//
	//amir := models.User{
	//	FirstName: "امیرحسین",
	//	LastName: "یعقوبی",
	//	Username: "amir",
	//	Password: "battle2021",
	//	Avatar: "awesome.jpg",
	//	LastIp: "127.0.0.1",
	//	LastAccess: time.Now(),
	//}
	//
	//conn.Save(&amir)
	//
	acc := models.Account{
		Name: "خانه",
		Icon: "family.png",
	}
	conn.Save(&acc)

	var user models.User
	conn.First(&user, 1)
	//conn.Model(&acc).Association("Users").Append([]models.User{amir, admin})
	conn.Model(&acc).Association("Users").Append([]models.User{user})
}

func loadAccountsWithUsers(conn *gorm.DB) {
	var account models.Account
	conn.Preload("Users").First(&account)

	log.Info(account)
	for _, user := range account.Users {
		log.Info(user)
	}
}

func createCategoriesAndSubCategories(conn *gorm.DB) {
	var account models.Account
	conn.First(&account)

	cat := models.IncomeCategory{
		Account: account,
		Name: "H@$@",
		Icon: "SAD.ico",
	}

	conn.Save(&cat)

	subCats := []models.IncomeCategory{
		{Name: "hd1", Icon: "hd1", Account: account},
		{Name: "hd2", Icon: "hd2", Account: account},
		{Name: "hd3", Icon: "hd3", Account: account},
	}

	conn.Model(&cat).Association("SubCategories").Append(subCats)
}

func main() {
	initLogger()
	if err := config.Initialize(); err != nil {
		log.Fatalf("cannot load configs, error=%s", err.Error())
		os.Exit(-1)
	}

	conn, err := db.GetPostgres()
	if err != nil {
		log.Fatalf("cannot connect to postgres database, error: %v\n", err.Error())
		os.Exit(-1)
	}
	defer conn.Close()
	log.Info("connected to postgres database")

	// TODO if MODE == debug
	db.InitializePostgres(conn)

	cfg := config.GetConfig()

	router := getApiRoutes()
	err = router.Run(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("cannot start API server. error: %s\n", err)
	}
	serverStartedAt = time.Now()
	log.Infof("API server started at %s\n", serverStartedAt.String())
}