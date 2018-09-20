package config

import (
	"github.com/spf13/viper"
	"fmt"
)

type PostgresConfig struct {
	Host string
	Port string
	DB string
	User string
	Password string
	SslMode string
}

func (c *PostgresConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.DB, c.User, c.Password, c.SslMode)
}


type Config struct {
	Port int64
	DevMode bool
	SecretKey string
	PostgresDB *PostgresConfig
}

var config *Config

func Initialize() error {
	viper.SetEnvPrefix("acc")
	viper.AutomaticEnv()

	viper.SetDefault("env", "development")
	viper.SetDefault("secretKey", "t3sT_-@SecR3t")
	viper.SetDefault("port", 8081)
	viper.SetDefault("devMode", true)

	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.accountManager")

	viper.SetConfigName("development")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	config = new(Config)
	config.Port = viper.GetInt64("port")
	config.DevMode = viper.GetBool("devMode")
	config.SecretKey = viper.GetString("secretKey")

	pgConfig := viper.GetStringMapString("postgresDB")
	postgres := PostgresConfig{
		Host: pgConfig["host"],
		Port: pgConfig["port"],
		DB:   pgConfig["db"],
		User: pgConfig["user"],
		Password: pgConfig["password"],
		SslMode: pgConfig["sslmode"],
	}
	config.PostgresDB = &postgres
	return err
}

func GetConfig() *Config {
	return config
}