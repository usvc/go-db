# DB

[![latest release](https://badge.fury.io/gh/usvc%2Fgo-db.svg)](https://github.com/usvc/go-db/releases)
[![pipeline status](https://gitlab.com/usvc/modules/go/db/badges/master/pipeline.svg)](https://gitlab.com/usvc/modules/go/db/-/commits/master)
[![release status](https://travis-ci.org/usvc/go-db.svg?branch=master)](https://travis-ci.org/usvc/go-db)
[![test coverage](https://api.codeclimate.com/v1/badges/90b27fd4714b0ce75111/test_coverage)](https://codeclimate.com/github/usvc/go-db/test_coverage)
[![maintainability](https://api.codeclimate.com/v1/badges/90b27fd4714b0ce75111/maintainability)](https://codeclimate.com/github/usvc/go-db/maintainability)

A Go package to handle database connections.

Currently supports databases utilising the following protocols:

1. MySQL
2. PostgreSQL
3. MSSQL

| | |
| --- | --- |
| Github | [https://github.com/usvc/go-db](https://github.com/usvc/go-db) |
| Gitlab | [https://gitlab.com/usvc/modules/go/db](https://gitlab.com/usvc/modules/go/db) |

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
      - [On Github](#on-github)
        - [Releasing](#releasing)
      - [On Gitlab](#on-gitlab)
        - [Version Bumping](#version-bumping)
        - [DockerHub Publishing](#dockerhub-publishing)
  - [Licensing](#licensing)

## Usage

### Importing

```go
import "github.com/usvc/go-db"
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

#### On Github

Github is used to deploy binaries/libraries because of it's ease of access by other developers.

##### Releasing

Releasing of the binaries can be done via Travis CI.

1. On Github, navigate to the [tokens settings page](https://github.com/settings/tokens) (by clicking on your profile picture, selecting **Settings**, selecting **Developer settings** on the left navigation menu, then **Personal Access Tokens** again on the left navigation menu)
2. Click on **Generate new token**, give the token an appropriate name and check the checkbox on **`public_repo`** within the **repo** header
3. Copy the generated token
4. Navigate to [travis-ci.org](https://travis-ci.org) and access the cooresponding repository there. Click on the **More options** button on the top right of the repository page and select **Settings**
5. Scroll down to the section on **Environment Variables** and enter in a new **NAME** with `RELEASE_TOKEN` and the **VALUE** field cooresponding to the generated personal access token, and hit **Add**

#### On Gitlab

Gitlab is used to run tests and ensure that builds run correctly.

##### Version Bumping

1. Run `make .ssh`
2. Copy the contents of the file generated at `./.ssh/id_rsa.base64` into an environment variable named **`DEPLOY_KEY`** in **Settings > CI/CD > Variables**
3. Navigate to the **Deploy Keys** section of the **Settings > Repository > Deploy Keys** and paste in the contents of the file generated at `./.ssh/id_rsa.pub` with the **Write access allowed** checkbox enabled

- **`DEPLOY_KEY`**: generate this by running `make .ssh` and copying the contents of the file generated at `./.ssh/id_rsa.base64`

##### DockerHub Publishing

1. Login to [https://hub.docker.com](https://hub.docker.com), or if you're using your own private one, log into yours
2. Navigate to [your security settings at the `/settings/security` endpoint](https://hub.docker.com/settings/security)
3. Click on **Create Access Token**, type in a name for the new token, and click on **Create**
4. Copy the generated token that will be displayed on the screen
5. Enter the following varialbes into the CI/CD Variables page at **Settings > CI/CD > Variables** in your Gitlab repository:

- **`DOCKER_REGISTRY_URL`**: The hostname of the Docker registry (defaults to `docker.io` if not specified)
- **`DOCKER_REGISTRY_USERNAME`**: The username you used to login to the Docker registry
- **`DOCKER_REGISTRY_PASSWORD`**: The generated access token

## Licensing

Code in this package is licensed under the [MIT license (click to see full text))](./LICENSE)

