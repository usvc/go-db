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
	err := mysql.Init("migrations", connection)
	if err != nil {
		fmt.Println(err)
	}
}
