package mysql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

func Init(tableName string, connection *sql.DB) error {
	stmt, err := connection.Prepare(fmt.Sprintf(`
		CREATE TABLE %s (
			id INTEGER(16) UNIQUE AUTO_INCREMENT NOT NULL,
			name VARCHAR(512) UNIQUE NOT NULL,
			up TEXT NOT NULL,
			down TEXT NOT NULL,
			error TEXT,
			status VARCHAR(16) NOT NULL,
			applied_at DATETIME,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		) Engine=InnoDB;
	`, tableName))
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func stringEndsWith(test, postfix string) bool {
	startIndexOfPostfix := strings.Index(test, postfix)
	if startIndexOfPostfix != -1 && startIndexOfPostfix == len(test)-len(postfix) {
		return true
	}
	return false
}

func stringSliceContains(test []string, item string) bool {
	for i := 0; i < len(test); i++ {
		if test[i] == item {
			return true
		}
	}
	return false
}

func GetMigrationNamesFromFilenames(filenameList []string) ([]string, []error) {
	var acceptedFilenames []string
	var rejectedFilenames []error
	for i := 0; i < len(filenameList); i++ {
		var filenameWithoutExtension string
		filename := filenameList[i]
		accepted := true
		switch true {
		case stringEndsWith(filename, ".up.sql"):
			filenameWithoutExtension = filename[:len(filename)-len(".up.sql")]
		case stringEndsWith(filename, ".down.sql"):
			filenameWithoutExtension = filename[:len(filename)-len(".down.sql")]
		default:
			rejectedFilenames = append(rejectedFilenames, fmt.Errorf("filename %s does not end with .{up, down}.sql", filename))
			accepted = false
		}
		acceptedBefore := stringSliceContains(acceptedFilenames, filenameWithoutExtension)
		if accepted && !acceptedBefore {
			migrationPairIsFound := stringSliceContains(filenameList, filenameWithoutExtension+".up.sql") &&
				stringSliceContains(filenameList, filenameWithoutExtension+".down.sql")
			if migrationPairIsFound {
				acceptedFilenames = append(acceptedFilenames, filenameWithoutExtension)
			} else {
				rejectedFilenames = append(rejectedFilenames, fmt.Errorf("could not find reverse migration for %s", filename))
			}
		}
	}
	if len(rejectedFilenames) > 0 {
		return acceptedFilenames, rejectedFilenames
	}
	return acceptedFilenames, nil
}

func New(name, up, down string) *Migration {
	return &Migration{
		Name: name,
		Up:   up,
		Down: down,
	}
}

func NewFromDB(name, tableName string, connection *sql.DB) (*Migration, error) {
	migration := Migration{}
	stmt, err := connection.Prepare(fmt.Sprintf(
		`SELECT
			id,
			name,
			up,
			down,
			error,
			status,
			applied_at,
			created_at
			FROM %s
				WHERE name = ?`,
		tableName,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query for retrieving migration entry: '%s'", err)
	}
	if err := stmt.QueryRow(name).Scan(
		&migration.ID,
		&migration.Name,
		&migration.Up,
		&migration.Down,
		&migration.Error,
		&migration.Status,
		&migration.AppliedAt,
		&migration.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to retrieve row from database for migration '%s': '%s'", name, err)
	}
	return &migration, nil
}

func NewFromFile(name, upFilePath, downFilePath string) (*Migration, error) {
	upFileContents, err := ioutil.ReadFile(upFilePath)
	if err != nil {
		return nil, err
	}
	downFileContents, err := ioutil.ReadFile(downFilePath)
	if err != nil {
		return nil, err
	}
	return New(name, string(upFileContents), string(downFileContents)), nil
}

func NewFromDirectory(directoryPath string) (Migrations, error) {
	directoryListing, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}
	var filenames []string
	for i := 0; i < len(directoryListing); i++ {
		file := directoryListing[i]
		filenames = append(filenames, file.Name())
	}
	filenames, errs := GetMigrationNamesFromFilenames(filenames)
	var migrations Migrations
	for i := 0; i < len(filenames); i++ {
		filename := filenames[i]
		migration, err := NewFromFile(filename, path.Join(directoryPath, filename+".up.sql"), path.Join(directoryPath, filename+".down.sql"))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		migrations = append(migrations, migration)
	}
	if errs != nil {
		var errString strings.Builder
		for i := 0; i < len(errs); i++ {
			err = errs[i]
			errString.WriteString("\n")
			errString.WriteString(err.Error())
		}
		return migrations, fmt.Errorf("following errors/warnings happened: %s", errString.String())
	}
	return migrations, nil
}

func NormalizeQuery(query string) string {
	return strings.Trim(
		strings.ReplaceAll(
			strings.ReplaceAll(
				query,
				"\t", " ",
			),
			"\n", " ",
		), "\t\n\r ",
	)
}
