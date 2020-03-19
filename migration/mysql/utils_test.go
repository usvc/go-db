package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/usvc/go-config"
	"github.com/usvc/go-db"
)

type UtilsTests struct {
	connection     *sql.DB
	migrationTable string
	suite.Suite
}

func TestUtils(t *testing.T) {
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
	suite.Run(t, &UtilsTests{
		connection:     db.Get(),
		migrationTable: "migrations",
	})
}

func (s *UtilsTests) SetupTest() {
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

func (s *UtilsTests) TestGetMigrationNamesFromFilenames() {
	migrationNames, err := GetMigrationNamesFromFilenames(
		[]string{
			"a.up.sql",
			"a.down.sql",
			"b.up.sql",
			"b.down.sql",
			"A.up.sql",
			"A.down.sql",
			"B.up.sql",
			"B.down.sql",
			"1.up.sql",
			"1.down.sql",
			"2.up.sql",
			"2.down.sql",
		},
	)
	s.Nil(err)
	s.EqualValues([]string{
		"a",
		"b",
		"A",
		"B",
		"1",
		"2",
	}, migrationNames)
}

func (s *UtilsTests) TestGetMigrationNamesFromFilenames_error() {
	migrationNames, err := GetMigrationNamesFromFilenames([]string{
		"a.up.sql",
		"b.up.sql",
		"B.up.sql",
		"b.down.sql",
		"randomFile",
	})
	s.NotNil(err)
	s.EqualValues([]string{"b"}, migrationNames)
	s.EqualValues([]error{
		errors.New("could not find reverse migration for a.up.sql"),
		errors.New("could not find reverse migration for B.up.sql"),
		errors.New("filename randomFile does not end with .{up, down}.sql"),
	}, err)
}

func (s *UtilsTests) TestNew() {
	migration := New("name", "up", "down")
	s.Equal(migration.Name, "name")
	s.Equal(migration.Up, "up")
	s.Equal(migration.Down, "down")
}

func (s *UtilsTests) TestNewFromDB() {
	expectedName := "test_new_from_db"
	expectedUpwardMigration := "CREATE TABLE test_new_from_db (id INTEGER)"
	expectedDownwardMigration := "DROP TABLE test_new_from_db"
	migration := New(
		expectedName,
		expectedUpwardMigration,
		expectedDownwardMigration,
	)
	err := migration.Apply(s.migrationTable, s.connection)
	s.Nil(err)
	migration, err = NewFromDB(expectedName, s.migrationTable, s.connection)
	s.Nil(err)
	s.Equal(migration.Name, expectedName)
	s.Equal(migration.Up, expectedUpwardMigration)
	s.Equal(migration.Down, expectedDownwardMigration)
}

func (s *UtilsTests) TestNewFromDirectory() {
	wd, _ := os.Getwd()
	migrations, err := NewFromDirectory(path.Join(wd, "./test/migrations/set01"))
	s.Nil(err)
	s.Len(migrations, 1)
	migrations, err = NewFromDirectory(path.Join(wd, "./test/migrations/set02"))
	s.Nil(err)
	s.Len(migrations, 2)
}

func (s *UtilsTests) TestNewFromDirectory_error() {
	wd, _ := os.Getwd()
	migrations, err := NewFromDirectory(path.Join(wd, "./test/migrations/set03"))
	s.NotNil(err)
	s.Contains(err.Error(), "b.up.sql")
	s.Len(migrations, 1)
	migrations, err = NewFromDirectory(path.Join(wd, "./test/migrations/set04"))
	s.NotNil(err)
	s.Contains(err.Error(), "b.down.sql")
	s.Len(migrations, 1)
	migrations, err = NewFromDirectory(path.Join(wd, "./test/migrations/set05"))
	s.NotNil(err)
	s.Contains(err.Error(), "not an sql file")
	s.Len(migrations, 1)
}

func (s *UtilsTests) Test_stringEndsWith() {
	s.True(stringEndsWith("a.up.sql", ".up.sql"))
	s.False(stringEndsWith("a.sql", ".up.sql"))
	s.False(stringEndsWith("a.up", ".up.sql"))
	s.False(stringEndsWith("a", ".up.sql"))
	s.False(stringEndsWith("Makefile", ".down.sql"))
}

func (s *UtilsTests) Test_stringSliceContains() {
	s.True(stringSliceContains([]string{"a", "b"}, "b"))
	s.True(stringSliceContains([]string{"a", "b"}, "a"))
	s.False(stringSliceContains([]string{"a", "b"}, "c"))
}
