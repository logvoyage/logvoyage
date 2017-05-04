package models

import (
	"fmt"
	"log"
	"time"

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
		log.Println("Database connection error:", err)
	}
	db.LogMode(true)
	return db
}

// GetConnection returns database connection instance
func GetConnection() *gorm.DB {
	return db
}

type BaseModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
