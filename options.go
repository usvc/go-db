package db

const (
	// DriverMySQL is the key for the MySQL driver
	DriverMySQL = "mysql"
	// DriverPostgreSQL is the key for the PostgreSQL driver
	DriverPostgreSQL = "postgres"
	// DriverMSSQL is the key for the Microsoft SQL Server driver
	DriverMSSQL = "sqlserver"

	// DefaultConnectionName is the assigned driver name when no .ConnectionName property is specified in Options
	DefaultConnectionName = "default"
	// DefaultDatabaseName is the default schema which the connection will connect to
	DefaultDatabaseName = "database"
	// DefaultDriver is the assigned driver name when no .Driver property is specified in Options
	DefaultDriver = DriverMySQL
	// DefaultHostname is the default hostname at which the database server is reachable at
	DefaultHostname = "127.0.0.1"
	// DefaultUser is the default username to use to login to the database service
	DefaultUser = "user"
	// DefaultPassword is the default password to use to login to the database service
	DefaultPassword = "password"
	// DefaultPortMySQL is the default MySQL port
	DefaultPortMySQL = uint16(3306)
	// DefaultPortPostgreSQL is the default PostgreSQL port
	DefaultPortPostgreSQL = uint16(5432)
)

// Options stores the database options that we use while establishing a connection
// to the database
type Options struct {
	// ConnectionName defines a local name of the connection, defaults to DefaultConnectionName
	ConnectionName string
	// Hostname defines the hostname where the database service can be reached
	Hostname string
	// Port defines the port which the database service is listening on
	Port uint16
	// Username defines the username of the user used to login to the database server
	Username string
	// Password defines the password of the user represented in the Username property
	Password string
	// Database defines the name of the database schema to use
	Database string
	// Driver defines the database driver to use, defaults to DefaultConnectionName (see SupportedDrivers for a list of supported drivers)
	Driver string
	// Params define connection parameters to use in the data source name (DSN)
	Params map[string]string
}

// AssignDefaults takes in a pointer to a connection :options parameter
// and updates optional fields with the defaults if they haven't been
// specified
func (o *Options) AssignDefaults() {
	if stringNotSet(o.ConnectionName) {
		o.ConnectionName = DefaultConnectionName
	}
	if stringNotSet(o.Driver) {
		o.Driver = DefaultDriver
	}
	if stringNotSet(o.Hostname) {
		o.Hostname = DefaultHostname
	}
	if uint16NotSet(o.Port) {
		switch o.Driver {
		case "postgresql":
			o.Port = DefaultPortPostgreSQL
		case "mysql":
			fallthrough
		default:
			o.Port = DefaultPortMySQL
		}
	}
	if stringNotSet(o.Username) {
		o.Username = DefaultUser
	}
	if stringNotSet(o.Password) {
		o.Password = DefaultPassword
	}
	if stringNotSet(o.Database) {
		o.Database = DefaultDatabaseName
	}
}
