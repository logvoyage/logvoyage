// Usage:
// # Initialize migations table
// > go run *.go init
//
// # Run migrations
// > go run *.go
//
// # Create new migration file
// > go run *.go create add_email_to_users
package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/firstrow/logvoyage/models"
	"github.com/go-pg/migrations"
)

const verbose = true

func main() {
	flag.Parse()

	db := models.GetConnection()

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		exitf(err.Error())
	}
	if verbose {
		if newVersion != oldVersion {
			fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
		} else {
			fmt.Printf("version is %d\n", oldVersion)
		}
	}
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
