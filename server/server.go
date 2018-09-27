package main

import (
	"os"
	"fmt"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/amir-yaghoobi/accountManager/db"
	"github.com/amir-yaghoobi/accountManager/config"
)

var serverStartedAt time.Time

func initLogger() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	initLogger()
	if err := config.Initialize(); err != nil {
		log.Fatalf("cannot load configs, error=%s", err.Error())
		return
	}

	conn, err := db.GetPostgres()
	if err != nil {
		log.Fatalf("cannot connect to postgres database, error: %v\n", err.Error())
		return
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