package migration

import "database/sql"

type Step interface {
	Apply(string, *sql.DB) error
	Rollback(string, *sql.DB) error
	Resolve(string, *sql.DB) error
	Validate(string, *sql.DB) error
}
