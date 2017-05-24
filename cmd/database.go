package cmd

import (
	"fmt"
	"log"

	"github.com/logvoyage/logvoyage/models/migrations"
	"github.com/logvoyage/logvoyage/shared/config"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates database",
	Long:  "Creates database",
	Run: func(cmd *cobra.Command, args []string) {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s sslmode=%s password=%s",
			config.Get("db.address"),
			config.Get("db.port"),
			config.Get("db.user"),
			config.Get("db.sslmode"),
			config.Get("db.password"),
		)
		db, err := gorm.Open("postgres", dsn)
		if err != nil {
			log.Fatal(err)
		}
		sql := fmt.Sprintf("CREATE DATABASE %s", config.Get("db.database"))
		db.Exec(sql)
		if db.Error != nil {
			fmt.Println("Error creating database", db.Error)
		}
		fmt.Printf("Database %s created\n", config.Get("db.database"))
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply database migrations",
	Long:  "Apply database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.Migrate()
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback last migration",
	Long:  "Rollback last migration",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.Rollback()
	},
}

func init() {
	var databaseCmd = &cobra.Command{Use: "database", Short: "Set of commands to work with database"}

	databaseCmd.AddCommand(
		createCmd,
		migrateCmd,
		rollbackCmd,
	)

	RootCmd.AddCommand(databaseCmd)
}
