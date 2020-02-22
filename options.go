package db

const (
	// DefaultConnectionName is the assigned driver name when no .ConnectionName property is specified in Options
	DefaultConnectionName = "default"
	// DefaultDriver is the assigned driver name when no .Driver property is specified in Options
	DefaultDriver = "mysql"
)

// Options stores the database options that we use while establishing a connection
// to the database
type Options struct {
	// ConnectionName defines a local name of the connection, defaults to DefaultConnectionName
	ConnectionName string
	// Hostname defines the hostname where the database service can be reached
	Hostname string
	// Port defines the port which the database service is listening on
	Port int
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
	if len(o.ConnectionName) == 0 {
		o.ConnectionName = DefaultConnectionName
	}
	if len(o.Driver) == 0 {
		o.Driver = DefaultDriver
	}
}
