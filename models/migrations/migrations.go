package migrations

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/logvoyage/logvoyage/models"
	gormigrate "gopkg.in/gormigrate.v1"
)

// Migrate migrate database schema up
func Migrate() {
	db := models.GetConnection()
	db.LogMode(true)

	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")
}

// Rollback rollbacks last migration
func Rollback() {
	db := models.GetConnection()
	db.LogMode(true)

	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())

	if err := m.RollbackLast(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")
}

func getMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "201608301400",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					gorm.Model
					Email    string
					Name     string
					Password string
				}
				return tx.AutoMigrate(&User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("users").Error
			},
		},
		{
			ID: "201608301430",
			Migrate: func(tx *gorm.DB) error {
				type Project struct {
					gorm.Model
					Name    string
					UUID    string
					OwnerID uint
				}
				return tx.AutoMigrate(&Project{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("projects").Error
			},
		},
	}
}
