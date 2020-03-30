# DB

[![latest release](https://badge.fury.io/gh/usvc%2Fgo-db.svg)](https://github.com/usvc/go-db/releases)
[![release status](https://travis-ci.org/usvc/go-db.svg?branch=master)](https://travis-ci.org/usvc/go-db)
[![pipeline status](https://gitlab.com/usvc/modules/go/db/badges/master/pipeline.svg)](https://gitlab.com/usvc/modules/go/db/-/commits/master)
[![test coverage](https://api.codeclimate.com/v1/badges/90b27fd4714b0ce75111/test_coverage)](https://codeclimate.com/github/usvc/go-db/test_coverage)
[![maintainability](https://api.codeclimate.com/v1/badges/90b27fd4714b0ce75111/maintainability)](https://codeclimate.com/github/usvc/go-db/maintainability)

A Go package to handle database connections.

This package is a convenience wrapper around other libraries and currently supports databases utilising the following protocols:

1. MySQL ([`github.com/go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql))
2. PostgreSQL ([`github.com/lib/pq`](https://github.com/lib/pq))
3. MSSQL ([`github.com/denisenkom/go-mssqldb`](https://github.com/denisenkom/go-mssqldb))

- **Github**: [https://github.com/usvc/go-db](https://github.com/usvc/go-db)
- **Gitlab**: [https://gitlab.com/usvc/modules/go/db](https://gitlab.com/usvc/modules/go/db)

> The main database tool by `usvc` can be found at [`github.com/usvc/db`](https://github.com/usvc/db) if that's what you were looking for

- [DB](#db)
- [Usage](#usage)
  - [Importing](#importing)
  - [Creating a new database connection](#creating-a-new-database-connection)
  - [Creating a new, named database connection](#creating-a-new-named-database-connection)
  - [Importing an existing connection](#importing-an-existing-connection)
  - [Verifying a connection works](#verifying-a-connection-works)
  - [Retrieving a database connection](#retrieving-a-database-connection)
  - [Closing a database connection](#closing-a-database-connection)
  - [Closing all connections](#closing-all-connections)
- [Configuration](#configuration)
  - [`db.Options`](#dboptions)
- [Development Runbook](#development-runbook)
  - [Getting Started](#getting-started)
  - [Continuous Integration (CI) Pipeline](#continuous-integration-ci-pipeline)
    - [On Github](#on-github)
      - [Releasing](#releasing)
    - [On Gitlab](#on-gitlab)
      - [Version Bumping](#version-bumping)
      - [DockerHub Publishing](#dockerhub-publishing)
- [Licensing](#licensing)

- - -

# Usage

## Importing

```go
import "github.com/usvc/go-db"
```

## Creating a new database connection

The following registers a new connection instance (but does not estabish a connection):

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

## Creating a new, named database connection

The following registers a new connection instance (but does not estabish a connection) named `"my-connection"`:

```go
if err := db.Init(Options{
  ConnectionName: "my-connection",
  Username: "user",
  Password: "password",
  Database: "schema",
  Hostname: "localhost",
  Port:     3306,
}); err != nil {
  log.Printf("an error occurred while creating the connection: %s", err)
}
```

## Importing an existing connection

The following imports an existing `*sql.DB` connection and names it `"connection-name"`

```go
var existingConnection *sql.DB
// ... initialise `existingConnection` by other means ...
if err := db.Import(existingConnection, "connection-name"); err != nil {
  log.Printf("connection 'connection-name' seems to already exist: %s\n", err)
}
```

## Verifying a connection works

The following checks if the connection instance named `"existing-connection"` can establish a connection to the database server:

```go
var existingConnection *sql.DB
// ... initialise `existingConnection` by other means ...
if err := db.Check("existing-connection"); err != nil {
  log.Printf("connection 'existing-connection' could not connect: %s\n", err)
}
```

## Retrieving a database connection

The following retrieves the default `*sql.DB` connection instance:

```go
// retrieve the 'default' connection
connection := db.Get()
if connection == nil {
  log.Println("connection 'default' does not exist")
}
```

The following retrieves the `*sql.DB` connection instance named `"my-connection"`:

```go
// retrieve the connection with name 'my-connection'
connection := db.Get("my-connection")
if connection == nil {
  log.Println("connection 'my-connection' does not exist")
}
```

## Closing a database connection

The following closes the default database connection:

```go
// ... run db.Init ...
if err := db.Close(); err != nil {
  log.Println("the default connection could not be closed")
}
```

The following closes the database connection named `"my-connection"`:

```go
// ... run db.Init for a connection named 'my-connection'...
if err := db.Close("my-connection"); err != nil {
  log.Println("the connection named 'my-connection' could not be closed")
}
```

## Closing all connections

The following closes all database connections:

```go
// ... run db.Init for a connection named 'my-connection'...
if err := db.Close("my-connection"); err != nil {
  for i := 0; i < len(err); i++ {
    log.Println(err[i])
  }
}
```

- - -

# Configuration

## `db.Options`

- **`ConnectionName`** `string`: Defines a local name of the connection. Defaults to `"default"`
- **`Hostname`** `string`: Defines the hostname where the database service can be reached. Defaults to `"127.0.0.1"`
- **`Port`** `string`: Defines the port which the database service is listening on. Defaults to `uint16(3306)`
- **`Username`** `string`: Defines the username of the user used to login to the database server. Defaults to `"user"`
- **`Password`** `string`: Defines the password of the user represented in the Username property. Defaults to `"password"`
- **`Database`** `string`: Defines the name of the database schema to use. Defaults to `"database"`
- **`Driver`** `string`: Defines the database driver to use. One of `db.DriverMySQL`, `db.DriverPostgreSQL`, or `db.DriverMSSQL`. Defaults to `db.DriverMySQL`
- **`Params`** `map[string]string`: Defines connection parameters to use in the data source name (DSN).

- - -

# Development Runbook

## Getting Started

1. Clone this repository
2. Run `make deps` to pull in external dependencies
3. Write some awesome stuff
4. Run `make test` to ensure unit tests are passing
5. Push

## Continuous Integration (CI) Pipeline

### On Github

Github is used to deploy binaries/libraries because of it's ease of access by other developers.

#### Releasing

Releasing of the binaries can be done via Travis CI.

1. On Github, navigate to the [tokens settings page](https://github.com/settings/tokens) (by clicking on your profile picture, selecting **Settings**, selecting **Developer settings** on the left navigation menu, then **Personal Access Tokens** again on the left navigation menu)
2. Click on **Generate new token**, give the token an appropriate name and check the checkbox on **`public_repo`** within the **repo** header
3. Copy the generated token
4. Navigate to [travis-ci.org](https://travis-ci.org) and access the cooresponding repository there. Click on the **More options** button on the top right of the repository page and select **Settings**
5. Scroll down to the section on **Environment Variables** and enter in a new **NAME** with `RELEASE_TOKEN` and the **VALUE** field cooresponding to the generated personal access token, and hit **Add**

### On Gitlab

Gitlab is used to run tests and ensure that builds run correctly.

#### Version Bumping

1. Run `make .ssh`
2. Copy the contents of the file generated at `./.ssh/id_rsa.base64` into an environment variable named **`DEPLOY_KEY`** in **Settings > CI/CD > Variables**
3. Navigate to the **Deploy Keys** section of the **Settings > Repository > Deploy Keys** and paste in the contents of the file generated at `./.ssh/id_rsa.pub` with the **Write access allowed** checkbox enabled

- **`DEPLOY_KEY`**: generate this by running `make .ssh` and copying the contents of the file generated at `./.ssh/id_rsa.base64`

#### DockerHub Publishing

1. Login to [https://hub.docker.com](https://hub.docker.com), or if you're using your own private one, log into yours
2. Navigate to [your security settings at the `/settings/security` endpoint](https://hub.docker.com/settings/security)
3. Click on **Create Access Token**, type in a name for the new token, and click on **Create**
4. Copy the generated token that will be displayed on the screen
5. Enter the following varialbes into the CI/CD Variables page at **Settings > CI/CD > Variables** in your Gitlab repository:

- **`DOCKER_REGISTRY_URL`**: The hostname of the Docker registry (defaults to `docker.io` if not specified)
- **`DOCKER_REGISTRY_USERNAME`**: The username you used to login to the Docker registry
- **`DOCKER_REGISTRY_PASSWORD`**: The generated access token

- - -

# Licensing

Code in this package is licensed under the [MIT license (click to see full text))](./LICENSE)

