package models

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	conn *gorm.DB
}

var DBClient *Client

// InitModels migrates users, and initialises database connection
func InitModels() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		viper.GetString("db.host"), viper.GetString("db.user"),
		viper.GetString("db.pass"), viper.GetString("db.name"),
		viper.GetString("db.port"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Guild{}, &Config{}, &ConfigItem{})
	DBClient = &Client{conn: db}
}

type DiscordResp struct {
	ID string `json:"id"`
}

type JWT struct {
	ID string
	jwt.StandardClaims
}
