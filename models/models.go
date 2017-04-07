package models

import (
	"fmt"
	"log"

	"bitbucket.org/firstrow/logvoyage/shared/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
	db = NewConnection()
}

// NewConnection creates new database connection
func NewConnection() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		config.Get("db.address"),
		config.Get("db.port"),
		config.Get("db.user"),
		config.Get("db.database"),
		config.Get("db.sslmode"),
		config.Get("db.password"),
	)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// GetConnection returns database connection instance
func GetConnection() *gorm.DB {
	return db
}
