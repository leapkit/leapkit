---
index: 2
title: "Database"
---

The kit CLI provides several commands for managing the database.

```bash
$ kit [database|db] <commands>
```

## Creating a new database

To create your database, you need to run the `kit [database|db] create` command. This command creates a new database using the value from the `DATABASE_URL` environment variable. The type of database engine will depend on the type of database URL set. Currently, Leapkit supports the creation of `SQLite3` and `PostgreSQL` databases.
```toml
# .env
DATABASE_URL=postgres://user:password@host:5432/database
```

```bash
$ kit [database|db] create
✅ Database created successfully
```


## Running migrations

To migrate the database to the latest version, you need to run the `kit db migrate`. This command applies any pending migrations to the database, ensuring it is up-to-date with the latest changes

```bash
$ kit [database|db] migrate
✅ Migrations ran successfully
```

## Deleting database

To delete the existing database, you need to run the `kit db drop` command. This command permanently deletes the database, **so use with caution**.

```bash
$ kit [database|db] drop
✅ Database dropped successfully
```

## Resetting database

The `kit db reset` command drops the existing database, creates a new one, and runs pending migrations. This command is useful for quickly resetting the database to a clean state.

```bash
$ kit db reset
✅ Database reset successfully
```