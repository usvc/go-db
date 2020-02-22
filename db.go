package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// SupportedDrivers is a list of driver names which the `db` package supports
	SupportedDrivers = []string{"mysql"}
	// connection is the singleton instance for the database connection
	connection = map[string]*sql.DB{}
)

// Check verifies that a connection can be made using the configuration
// provided in the connection :options parameter
func Check(optionalConnectionName ...string) error {
	connectionName := DefaultConnectionName
	if len(optionalConnectionName) > 0 {
		connectionName = optionalConnectionName[0]
	}
	var (
		selectedConnection *sql.DB
		exists             bool
	)
	if selectedConnection, exists = connection[connectionName]; !exists {
		if selectedConnection, exists = connection[DefaultConnectionName]; !exists {
			return fmt.Errorf("connection with id '%s' does not exist", connectionName)
		}
	}
	if err := selectedConnection.Ping(); err != nil {
		return err
	}
	return nil
}

// CloseAll attempts to close all established connections, returning
// a list of errors in cases where the connection could not be closed.
func CloseAll() []error {
	errs := []error{}
	for key, value := range connection {
		if err := value.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error while closing connection '%s': '%s'", key, err))
		} else {
			delete(connection, key)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Get returns the instance of the database connection
func Get(optionalConnectionName ...string) *sql.DB {
	connectionName := DefaultConnectionName
	if len(optionalConnectionName) > 0 {
		connectionName = optionalConnectionName[0]
	}
	if val, ok := connection[connectionName]; ok {
		return val
	} else if val, ok := connection[DefaultConnectionName]; ok {
		return val
	}
	return nil
}

// Import imports an existing database connection if another connection
// with the same name does not exist (an error is returned if so)
func Import(existingConnection *sql.DB, optionalConnectionName ...string) error {
	connectionName := DefaultConnectionName
	if len(optionalConnectionName) > 0 {
		connectionName = optionalConnectionName[0]
	}
	if _, ok := connection[connectionName]; ok {
		return fmt.Errorf("unable to import connection with id '%s' - another connection with the same id already exists", connectionName)
	}
	connection[connectionName] = existingConnection
	return nil
}

// Init initialises the `db` module so that Get can be called
// anywhere in the consuming package
func Init(options Options) error {
	var err error
	options.AssignDefaults()
	if _, ok := connection[options.ConnectionName]; ok {
		defer connection[options.ConnectionName].Close()
	}
	connection[options.ConnectionName], err = sql.Open(
		options.Driver,
		generateDSN(options),
	)
	if err != nil {
		return err
	}
	return nil
}
