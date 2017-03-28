package models

import (
	"fmt"
	"log"
	"time"

	"bitbucket.org/firstrow/logvoyage/shared/config"
	"github.com/go-pg/pg"
)

var db *pg.DB

func init() {
	db = NewConnection()
}

// NewConnection creates new database connection
func NewConnection() *pg.DB {
	conn := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Get("db.address"), config.Get("db.port")),
		User:     config.Get("db.user"),
		Password: config.Get("db.password"),
		Database: config.Get("db.database"),
	})

	return conn
}

// GetConnection returns database connection instance
func GetConnection() *pg.DB {
	return db
}

// EnableSQLLogging enables detailed sql logging to stdout
func EnableSQLLogging() {
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			log.Println("Query error:", err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})
}
