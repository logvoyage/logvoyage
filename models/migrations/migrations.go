package migrations

import (
	"fmt"
	"os"

	"bitbucket.org/firstrow/logvoyage/models"
	"github.com/firstrow/migrations"
)

// Migrate executes migration command.
// cmd values can be: up, down, init, version, create
func Migrate(cmd ...string) (int64, int64) {
	path := fmt.Sprintf("%s/src/bitbucket.org/firstrow/logvoyage/models/migrations", os.Getenv("GOPATH"))
	migrations.SetMigratonsPath(path)

	db := models.GetConnection()

	oldVersion, newVersion, err := migrations.Run(db, cmd...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	}

	if len(cmd) > 0 {
		if cmd[0] == "version" {
			fmt.Printf("version: %v\n", oldVersion)
		}
	}
	return oldVersion, newVersion
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
