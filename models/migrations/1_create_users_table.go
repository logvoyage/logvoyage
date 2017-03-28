package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating table users...")
		_, err := db.Exec(`
		  CREATE TABLE users(
		    id serial PRIMARY KEY,
		    email VARCHAR (255) UNIQUE NOT NULL,
		    name VARCHAR (255),
		    password VARCHAR (255),
		    created_on timestamp default current_timestamp
		  );`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table users...")
		_, err := db.Exec(`DROP TABLE users`)
		return err
	})
}
