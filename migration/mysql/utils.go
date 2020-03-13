package mysql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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

func New(name, up, down string) *Migration {
	return &Migration{
		Name: name,
		Up:   up,
		Down: down,
	}
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
