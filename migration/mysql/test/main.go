package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
	"github.com/usvc/go-db"
	"github.com/usvc/go-db/migration/mysql"
)

const (
	migrationTableName = "migrations"
)

var conf = config.Map{
	"db_host": &config.String{
		Default: "127.0.0.1",
	},
	"db_port": &config.Uint{
		Default: 3306,
	},
	"db_username": &config.String{
		Default: "root",
	},
	"db_password": &config.String{
		Default: "toor",
	},
	"db_database": &config.String{
		Default: "database",
	},
}

var migrations = []*mysql.Migration{
	mysql.New(
		"create_table",
		`CREATE TABLE mesa (
			id INTEGER PRIMARY KEY AUTO_INCREMENT
		) Engine=InnoDB;`,
		`DROP TABLE mesa;`,
	),
}

var rootCommand = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func getDBOptions() db.Options {
	return db.Options{
		Hostname: conf.GetString("db_host"),
		Port:     uint16(conf.GetUint("db_port")),
		Username: conf.GetString("db_username"),
		Password: conf.GetString("db_password"),
		Database: conf.GetString("db_database"),
	}
}

func main() {
	conf.ApplyToCobra(rootCommand)
	rootCommand.AddCommand(&cobra.Command{
		Use: "check",
		Run: func(cmd *cobra.Command, args []string) {
			db.Init(getDBOptions())
			connection := db.Get()
			stmt, err := connection.Prepare("SHOW TABLES")
			if err != nil {
				panic(err)
			}
			rows, err := stmt.Query()
			if err != nil {
				panic(err)
			}
			for rows.Next() {
				var tableName string
				rows.Scan(&tableName)
				fmt.Println(tableName)
			}
		},
	})
	rootCommand.AddCommand(&cobra.Command{
		Use: "init",
		Run: func(cmd *cobra.Command, args []string) {
			db.Init(getDBOptions())
			connection := db.Get()
			err := mysql.Init(migrationTableName, connection)
			if err != nil {
				fmt.Println(err)
			}
		},
	})
	rootCommand.AddCommand(&cobra.Command{
		Use: "apply",
		Run: func(cmd *cobra.Command, args []string) {
			db.Init(getDBOptions())
			connection := db.Get()
			for i := 0; i < len(migrations); i++ {
				migration := migrations[i]
				err := migration.Validate(migrationTableName, connection)
				if err != nil && err != mysql.NoErrDoesNotExist {
					fmt.Println(err)
					err = migration.Resolve(migrationTableName, connection)
					if err != nil {
						fmt.Println(err)
						break
					}
				}
				err = migration.Apply(migrationTableName, connection)
				if err != nil {
					if err == mysql.NoErrAlreadyApplied {
						fmt.Printf("[apply:%s] migration already applied\n", migration.Name)
						continue
					}
					fmt.Println(err)
				}
			}
		},
	})
	rootCommand.AddCommand(&cobra.Command{
		Use: "rollback",
		Run: func(cmd *cobra.Command, args []string) {
			db.Init(getDBOptions())
			connection := db.Get()
			for i := 0; i < len(migrations); i++ {
				migration := migrations[i]
				err := migration.Validate(migrationTableName, connection)
				if err != nil {
					if err == mysql.NoErrDoesNotExist {
						fmt.Printf("[rollback:%s] %s", migration.Name, err)
						break
					}
				}
				err = migration.Rollback(migrationTableName, connection)
				if err != nil {
					fmt.Println(err)
				}
			}
		},
	})
	rootCommand.Execute()
}
