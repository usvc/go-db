package mysql

import (
	"database/sql"
	"fmt"
	"time"
)

type Migration struct {
	// ID will contain the database-assigned auto-incremented ID of the migration
	ID int64 `json:"id" yaml:"id"`
	// Name contains the name of the migration derived from the migration file name
	Name string `json:"name" yaml:"name"`
	// Up contains a SQL command that should result in a up-version shift of the database schema
	Up string `json:"up" yaml:"up"`
	// Down contains a SQL command that should result in a down-version shift of the database schema
	Down string `json:"down" yaml:"down"`
	// Error contains any error that happened
	Error *string `json:"error" yaml:"error"`
	// Status contains the status of the migration
	Status string `json:"status" yaml:"status`
	// AppliedAt holds the timestamp when the migration was successfully applied to the database
	AppliedAt *time.Time `json:"applied_at" yaml:"applied_at"`
	// CreatedAt holds the timestamp when the migration was initialised in the database
	CreatedAt *time.Time `json:"created_at" yaml:"created_at"`
}

func (m *Migration) Apply(tableName string, connection *sql.DB) error {
	exists, err := m.exists(tableName, connection)
	if err != nil {
		return fmt.Errorf("[apply:%s] failed to check if migration has already been applied: '%s'", m.Name, err)
	} else if exists {
		return NoErrAlreadyApplied
	}
	err = m.insert(tableName, connection)
	if err != nil {
		return fmt.Errorf("[apply:%s] failed to insert migration entry into migration table '%s': '%s'", m.Name, tableName, err)
	}
	stmt, err := connection.Prepare(m.Up)
	if err != nil {
		return fmt.Errorf("[apply:%s] failed to prepare migration: '%s'", m.Name, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		if setErrorErr := m.setError(err, tableName, connection); setErrorErr != nil {
			return fmt.Errorf("[apply:%s] failed to apply migration: '%s'", m.Name, setErrorErr)
		}
		return fmt.Errorf("[apply:%s] failed to apply migration: '%s'", m.Name, err)
	}
	stmt, err = connection.Prepare(fmt.Sprintf(
		"UPDATE %s SET applied_at = ?, status = ? WHERE id = ?", tableName,
	))
	if err != nil {
		return fmt.Errorf("[apply:%s] failed to prepare query to indicate success: '%s'", m.Name, err)
	}
	_, err = stmt.Exec(time.Now(), StatusApplied, m.ID)
	if err != nil {
		return fmt.Errorf("[apply:%s] failed to indicate success for migration: '%s'", m.Name, err)
	}
	return nil
}

func (m *Migration) Rollback(tableName string, connection *sql.DB) error {
	stmt, err := connection.Prepare(fmt.Sprintf(
		"UPDATE %s SET status = ?",
		tableName,
	))
	if err != nil {
		return fmt.Errorf("[rollback:%s] failed to prepare query for rollback migration status update: '%s'", m.Name, err)
	}
	if _, err = stmt.Exec("rolling back"); err != nil {
		return fmt.Errorf("[rollback:%s] failed to update status for rollback migration: '%s'", m.Name, err)
	}
	stmt, err = connection.Prepare(m.Down)
	if err != nil {
		return fmt.Errorf("[rollback:%s] failed to prepare query for rollback migration: '%s'", m.Name, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		stmt, err2 := connection.Prepare(fmt.Sprintf(
			`UPDATE %s SET error = ? WHERE name = ?`, tableName,
		))
		if err2 != nil {
			return fmt.Errorf("[rollback:%s] failed to prepare query to set error column for failed migration: '%s' > '%s'", m.Name, err, err2)
		}
		_, err2 = stmt.Exec(err.Error(), m.Name)
		if err2 != nil {
			return fmt.Errorf("[rollback:%s] failed to set error column for failed rollback migration: '%s' > '%s'", m.Name, err, err2)
		}
		return fmt.Errorf("[rollback:%s] failed to rollback migration: '%s'", m.Name, err)
	}
	if err := m.delete(tableName, connection); err != nil {
		return fmt.Errorf("[rollback:%s] failed to remove migration table entry: '%s'", m.Name, err)
	}
	return nil
}

func (m *Migration) Resolve(tableName string, connection *sql.DB) error {
	if err := m.delete(tableName, connection); err != nil {
		return fmt.Errorf("[resolve:%s] failed to resolve migration error: '%s'", m.Name, err)
	}
	return nil
}

func (m *Migration) Validate(tableName string, connection *sql.DB) error {
	remoteMigration, err := NewFromDB(m.Name, tableName, connection)
	if err != nil {
		if err == sql.ErrNoRows {
			return NoErrDoesNotExist
		}
		return fmt.Errorf("[validate:%s] failed to retrieve migration entry: %s", m.Name, err)
	}
	if remoteMigration.Error != nil && len(*remoteMigration.Error) > 0 {
		return fmt.Errorf("[validate:%s] migration exists but has been recorded as failed: '%s'", m.Name, *remoteMigration.Error)
	} else if NormalizeQuery(m.Up) != NormalizeQuery(remoteMigration.Up) {
		return fmt.Errorf("[validate:%s] failed to reconcile upward migration query local and remote versions:\n%s\n--\n%s", m.Name, m.Up, remoteMigration.Up)
	} else if NormalizeQuery(m.Down) != NormalizeQuery(remoteMigration.Down) {
		return fmt.Errorf("[validate:%s] failed to reconcile downward migration query local and remote versions:\n%s\n--\n%s", m.Name, m.Down, remoteMigration.Down)
	}
	return nil
}

func (m *Migration) delete(tableName string, connection *sql.DB) error {
	stmt, err := connection.Prepare(fmt.Sprintf(
		`DELETE FROM %s WHERE name = ?`, tableName,
	))
	if err != nil {
		return fmt.Errorf("failed to prepare query to delete migration entry: '%s'", err)
	}
	_, err = stmt.Exec(m.Name)
	if err != nil {
		return fmt.Errorf("failed to delete migration entry: '%s'", err)
	}
	return nil
}

func (m *Migration) exists(tableName string, connection *sql.DB) (bool, error) {
	stmt, err := connection.Prepare(fmt.Sprintf(
		"SELECT 1 FROM %s WHERE name = ?", tableName,
	))
	if err != nil {
		return false, fmt.Errorf("failed to prepare query: '%s'", err)
	}
	var exists int64
	if err = stmt.QueryRow(m.Name).Scan(&exists); err == nil {
		return exists == 1, nil
	} else if err == sql.ErrNoRows {
		return false, nil
	}
	return false, fmt.Errorf("failed to query row: '%s'", err)
}

func (m *Migration) insert(tableName string, connection *sql.DB) error {
	stmt, err := connection.Prepare(fmt.Sprintf(
		"INSERT INTO %s (name, up, down, status) VALUES (?, ?, ?, ?)", tableName,
	))
	if err != nil {
		return fmt.Errorf("[apply:%s] failed to prepare query for initial row insert into migrations table '%s': '%s'", m.Name, tableName, err)
	}
	res, err := stmt.Exec(m.Name, m.Up, m.Down, StatusApplying)
	if err != nil {
		return fmt.Errorf("[apply:%s] failed to insert initial row into migrations table '%s': '%s'", m.Name, tableName, err)
	}
	m.ID, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("[apply:%s] failed to get id of newly inserted row in the migrations table '%s': '%s'", m.Name, tableName, err)
	}
	return nil
}

func (m *Migration) setError(originalError error, tableName string, connection *sql.DB) error {
	stmt, updateError := connection.Prepare(fmt.Sprintf(
		"UPDATE %s SET error = ? WHERE id = ?", tableName,
	))
	if updateError != nil {
		return fmt.Errorf("failed to prepare query to set error column to '%s': '%s'", originalError, updateError)
	}
	_, updateError = stmt.Exec(originalError.Error(), m.ID)
	if updateError != nil {
		return fmt.Errorf("failed to set error column to '%s': '%s'", originalError, updateError)
	}
	return nil
}
