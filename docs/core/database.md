---
index: 2
title: "Database"
---

Leapkit's `db` package streamlines database management operations. It offers you support for both [**PostgreSQL**](https://www.postgresql.org/) and [**SQLite3**](https://www.sqlite.org/) database engines.

These are the tools that you can find within the `db` package:

- **Management**: Create and delete your database.
- **Connection**: Establish and manage database connections.
- **Migrations**: Modify the database data and schema structures.

Let's explore each tool in detail.

## Management

The `db` package offers two essential functions: `Create()` and `Drop()`. These functions enable you to set up and manage your database for your project. Simply provide the database URL as a parameter to create it. These functions internally identify, based on the URL structure, which *database engine* will be used to process the action.

```go
import "github.com/leapkit/leapkit/core/db"

var databaseURL = "postgres://user:password@host:port/db_name"

// Creating the database...
if err := db.Create(databaseURL); err != nil {
    // handle the error
}

// Dropping the database...
if err := db.Drop(databaseURL); err != nil {
    // handle the error
}
```

## Connection

To establish a connection with your database in your app, use the `ConnectionFn` function, which returns the `*sqlx.DB` and an `error` if the connection cannot be established.

```go
var (
    databaseURL = "postgres://user:password@host:port/db_name"

    DB = db.ConnectionFn(databaseURL)
)

conn, err := DB()
if err != nil {
    // handle the error
}
```

### Options

By default, the `ConnectionFn` function is configured to connect to a **PostgreSQL** database engine. But if you want to connect to your **SQLite3** database, you should use the `db.WithDriver()` option. With this you can define the driver to stablish the connection with the database. The allowed driver names are:

- `postgres`
- `sqlite` or `sqlite3`

```go
databaseURL = "/path/to/your/sqlite3.db"

DB = db.ConnectionFn(databaseURL, db.WithDriver("sqlite3"))
```

### Drivers

So that the `db` package can connect to your database, ensure that you import the specific driver package for the database engine you want to use.

- PostgreSQL: `github.com/lib/pq`
- SQLite3: `github.com/mattn/go-sqlite3`

### Examples

#### Connecting to a PostgreSQL database engine

```go
package main

import (
    "github.com/leapkit/leapkit/core/db"
    _ "github.com/lib/pq" // importing the PostgreSQL driver.
)

var (
    databaseURL = "postgres://user:password@host:port/db_name"

    DB = db.ConnectionFn(databaseURL)
)

func myAwesomeFunc() {
    conn, err := DB()
    if err != nil {
        // handle the error
    }
    // ...
}
```

#### Connecting to a SQLite3 database engine

```go
package main

import (
    "github.com/leapkit/leapkit/core/db"
    _ "github.com/mattn/go-sqlite3" // importing the SQLite3 driver.
)

var (
    databaseURL = "/path/to/your/sqlite3.db"

    DB = db.ConnectionFn(databaseURL, db.WithDriver("sqlite3"))
)

// ...
func myAwesomeFunc() {
    conn, err := DB()
    if err != nil {
        // handle the error
    }
    // ...
}
```
## Migrations

Within the `db` package, Leapkit offers tools for database migration operations.

### Creating migrations

To create a migration, the `db` package provides you with the `GenerateMigration()` function, which creates a file with a name with the convention `YYYYMMDDHHmmss_migration_name.sql`, for instance `20260626100405_my_awesome_migration.sql`.

By default, migrations are created on the `./internal/app/database/migrations` path. However, you can set a custom migration path by applying the `UseMigrationFolder()` option.

```go
import (
    "github.com/leapkit/leapkit/core/db"
	"github.com/leapkit/leapkit/core/db/migrations"
)

// ...
err := db.GenerateMigration("adding_new_column_to_my_awesome_table",
    migrations.UseMigrationFolder("/my/awesome/migration/path"),
)

if err != nil {
    // handle the error
}
```

### Running migrations

To run your database migrations, you can make use of `RunMigrations` function, which receives an `embed.FS` file to locate your migrations files, and the `*sqlx.DB` connection.

```go
// /my/awesome/migration/path/migrations.go
package migrations

import "embed"

//go:embed *.sql */*.sql
var FS embed.FS
```


```go
package main

import (
    "github.com/leapkit/leapkit/core/db"
    "github.com/leapkit/leapkit/core/db/migrations"
    _ "github.com/lib/pq"

    "/my/awesome/migration/path"
)

var (
    databaseURL = "postgres://user:password@host:port/db_name"

    DB = db.ConnectionFn(databaseURL)
)

func applyDatabaseMigrations() {
    conn, err := DB()
    if err != nil {
        // handle the error
    }

    if err := db.RunMigrations(migrations.FS, conn); err != nil {
        // handle the error
    }
    // ...
}
```
