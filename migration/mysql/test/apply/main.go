package main

import (
	"fmt"

	"github.com/usvc/go-db"
	"github.com/usvc/go-db/migration/mysql"
)

func main() {
	opts := db.Options{
		Hostname: "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "toor",
		Database: "default",
	}
	db.Init(opts)
	connection := db.Get()
	migrations := []*mysql.Migration{
		mysql.New(
			"create_table",
			`CREATE TABLE mesa (
				id INTEGER PRIMARY KEY AUTO_INCREMENT
			) Engine=InnoDB;`,
			`DROP TABLE mesa;`,
		),
	}
	for i := 0; i < len(migrations); i++ {
		migration := migrations[i]
		err := migration.Validate("migrations", connection)
		if err != nil {
			fmt.Println(err)
			err = migration.Resolve("migrations", connection)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		err = migration.Apply("migrations", connection)
		if err != nil {
			if err == mysql.NoErrAlreadyApplied {
				fmt.Printf("[migration:%s] already applied\n", migration.Name)
				continue
			}
			fmt.Println(err)
		}
	}
}
