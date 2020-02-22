# DB

[![pipeline status](https://gitlab.com/usvc/modules/go/db/badges/master/pipeline.svg)](https://gitlab.com/usvc/modules/go/db/-/commits/master)

A Go package to handle database connections.

- [DB](#db)
  - [Usage](#usage)
    - [Importing](#importing)
    - [Creating a new database connection](#creating-a-new-database-connection)
    - [Creating a new, named database connection](#creating-a-new-named-database-connection)
    - [Retrieving a database connection](#retrieving-a-database-connection)
    - [Retrieving a named database connection](#retrieving-a-named-database-connection)
    - [Importing an existing connection](#importing-an-existing-connection)
    - [Verifying a connection works](#verifying-a-connection-works)
  - [Development Runbook](#development-runbook)
    - [Getting Started](#getting-started)
    - [Continuous Integration (CI) Pipeline](#continuous-integration-ci-pipeline)
  - [Licensing](#licensing)

## Usage

### Importing

```go
import "github.com/usvc/db"
```

### Creating a new database connection

```go
if err := db.Init(Options{
  Username: "user",
  Password: "password",
  Database: "schema",
  Hostname: "localhost",
  Port:     3306,
}); err != nil {
  log.Printf("an error occurred while creating the connection: %s", err)
}
```

### Creating a new, named database connection

```go
if err := db.Init(Options{
  ConnectionName: "non-default",
  Username: "user",
  Password: "password",
  Database: "schema",
  Hostname: "localhost",
  Port:     3306,
}); err != nil {
  log.Printf("an error occurred while creating the connection: %s", err)
}
```

### Retrieving a database connection

```go
// retrieve the 'default' connection
connection := db.Get()
if connection == nil {
  log.Println("connection 'default' does not exist")
}
```

### Retrieving a named database connection

```go
// retrieve the 'non-default' connection
connection := db.Get("non-default")
if connection == nil {
  log.Println("connection 'non-default' does not exist")
}
```

### Importing an existing connection

```go
var existingConnection *sql.DB
// ... initialise `existingConnection` by other means ...
if err := db.Import(existingConnection, "default"); err != nil {
  log.Printf("connection 'default' seems to already exist: %s\n", err)
}
```

### Verifying a connection works

```go
var existingConnection *sql.DB
// ... initialise `existingConnection` by other means ...
db.Import(existingConnection, "existing-connection")
if err := db.Check("existing-connection"); err != nil {
  log.Printf("connection 'existing-connection' could not connect: %s\n", err)
}
```

## Development Runbook

### Getting Started

1. Clone this repository
2. Run `make deps` to pull in external dependencies
3. Write some awesome stuff
4. Run `make test` to ensure unit tests are passing
5. Push

### Continuous Integration (CI) Pipeline

To set up the CI pipeline in Gitlab:

1. Run `make .ssh`
2. Copy the contents of the file generated at `./.ssh/id_rsa.base64` into an environment variable named **`DEPLOY_KEY`** in **Settings > CI/CD > Variables**
3. Navigate to the **Deploy Keys** section of the **Settings > Repository > Deploy Keys** and paste in the contents of the file generated at `./.ssh/id_rsa.pub` with the **Write access allowed** checkbox enabled

- **`DEPLOY_KEY`**: generate this by running `make .ssh` and copying the contents of the file generated at `./.ssh/id_rsa.base64`

## Licensing

Code in this package is licensed under the [MIT license (click to see full text))](./LICENSE)

