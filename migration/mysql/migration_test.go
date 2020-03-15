package mysql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/usvc/go-config"
	"github.com/usvc/go-db"
)

type MigrationTests struct {
	connection     *sql.DB
	migrationTable string
	suite.Suite
}

func TestMigration(t *testing.T) {
	conf := config.Map{
		"mysql_hostname": &config.String{
			Default: "127.0.0.1",
		},
		"mysql_database": &config.String{
			Default: "database",
		},
		"mysql_root_password": &config.String{
			Default: "toor",
		},
	}
	conf.LoadFromEnvironment()
	db.Init(db.Options{
		Hostname: conf.GetString("mysql_hostname"),
		Port:     3306,
		Username: "root",
		Password: conf.GetString("mysql_root_password"),
		Database: conf.GetString("mysql_database"),
	})
	suite.Run(t, &MigrationTests{
		connection:     db.Get(),
		migrationTable: "migrations",
	})
}

func (s *MigrationTests) SetupTest() {
	stmt, err := s.connection.Prepare("SHOW TABLES")
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
		stmt, err := s.connection.Prepare(fmt.Sprintf("DROP TABLE %s", tableName))
		if err != nil {
			panic(err)
		}
		_, err = stmt.Exec()
		if err != nil {
			panic(err)
		}
	}
	Init(s.migrationTable, s.connection)
}

func (s *MigrationTests) TestApply() {
	migration := Migration{
		Name: "test_apply",
		Up:   "CREATE TABLE test_apply (id INTEGER) Engine=InnoDB;",
		Down: "DROP TABLE test_apply",
	}
	err := migration.Apply(s.migrationTable, s.connection)
	s.Nil(err)
	err = migration.Validate(s.migrationTable, s.connection)
	s.Nil(err)
}

func (s *MigrationTests) TestApply_error() {
	migration := Migration{
		Name: "test_apply_error",
		Up:   "CREATE TABLE test_apply_error (id INTEGER AUTO_INCREMENT) Engine=InnoDB;",
		Down: "DROP TABLE test_apply",
	}
	err := migration.Apply(s.migrationTable, s.connection)
	s.Contains(err.Error(), "there can be only one auto column")
	err = migration.Validate(s.migrationTable, s.connection)
	s.Contains(err.Error(), "there can be only one auto column")
}

func (s *MigrationTests) TestResolve() {
	migration := Migration{
		Name: "test_resolve",
		Up:   "CREATE TABLE test_resolve (id INTEGER AUTO_INCREMENT) Engine=InnoDB;",
		Down: "DROP TABLE test_resolve",
	}
	err := migration.Apply(s.migrationTable, s.connection)
	s.NotNil(err)
	err = migration.Resolve(s.migrationTable, s.connection)
	s.Nil(err)

}

func (s *MigrationTests) TestRollback() {
	migration := Migration{
		Name: "test_rollback",
		Up:   "CREATE TABLE test_rollback (id INTEGER) Engine=InnoDB;",
		Down: "DROP TABLE test_rollback",
	}
	err := migration.Apply(s.migrationTable, s.connection)
	s.Nil(err)
	err = migration.Rollback(s.migrationTable, s.connection)
	s.Nil(err)
}

func (s *MigrationTests) TestRollback_error_notApplied() {
	migration := Migration{
		Name: "test_rollback_error_not_applied",
		Up:   "CREATE TABLE test_rollback_ena (id INTEGER) Engine=InnoDB;",
		Down: "DROP TABLE test_rollback_ena",
	}
	err := migration.Rollback(s.migrationTable, s.connection)
	s.Contains(err.Error(), "Unknown table")
}

func (s *MigrationTests) TestValidate() {
	migration := Migration{
		Name: "test_validate",
		Up:   "CREATE TABLE test_validate (id INTEGER) Engine=InnoDB;",
		Down: "DROP TABLE test_validate",
	}
	err := migration.Validate(s.migrationTable, s.connection)
	s.Equal(err, NoErrDoesNotExist)
	err = migration.Apply(s.migrationTable, s.connection)
	s.Nil(err)
	err = migration.Validate(s.migrationTable, s.connection)
	s.Nil(err)
}

func (s *MigrationTests) TestValidate_error() {
	up := "CREATE TABLE test_validate_error (id INTEGER) Engine=InnoDB;"
	down := "DROP TABLE test_validate_error;"
	migration := Migration{
		Name: "test_validate_error",
		Up:   up,
		Down: down,
	}
	err := migration.Apply(s.migrationTable, s.connection)
	s.Nil(err)

	migration.Up = "CREATE TABLE test_validate_error (id INTEGER)"
	err = migration.Validate(s.migrationTable, s.connection)
	s.Contains(err.Error(), "failed to reconcile upward migration query")
	migration.Up = up
	migration.Down = "DROP TABLE test_validate;"
	err = migration.Validate(s.migrationTable, s.connection)
	s.Contains(err.Error(), "failed to reconcile downward migration query")
}
