package migrations

import (
	"fmt"

	"github.com/firstrow/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating table projects...")
		_, err := db.Exec(`
		  CREATE TABLE projects(
		    id serial PRIMARY KEY,
		    name VARCHAR (255),
		    uuid VARCHAR (36),
		    owner_id INTEGER,
		    created_at timestamp default current_timestamp
		  );`)
		if err != nil {
			return err
		}
		_, err = db.Exec(`
      CREATE INDEX uuid_idx ON projects (uuid);
      CREATE INDEX owner_id_idx ON projects (owner_id);
	  `)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table projects...")
		_, err := db.Exec(`DROP TABLE projects`)
		return err
	})
}
