package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/firstrow/logvoyage/models/migrations"
	"bitbucket.org/firstrow/logvoyage/shared/config"

	"github.com/go-pg/pg"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates database",
	Long:  "Creates database",
	Run: func(cmd *cobra.Command, args []string) {
		conn := pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("%s:%s", config.Get("db.address"), config.Get("db.port")),
			User:     config.Get("db.user"),
			Password: config.Get("db.password"),
		})
		sql := fmt.Sprintf("CREATE DATABASE %s", config.Get("db.database"))
		_, err := conn.Exec(sql)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("Created database %s\n", config.Get("db.database"))
		}
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate [up, down, version]",
	Short: "Apply database migrations",
	Long:  "Apply database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.Migrate("init") // Always check for migrations table
		migrations.Migrate(args...)
	},
}

func init() {
	var databaseCmd = &cobra.Command{Use: "database", Short: "Set of commands to work with database"}

	databaseCmd.AddCommand(
		createCmd,
		migrateCmd,
	)

	RootCmd.AddCommand(databaseCmd)
}
