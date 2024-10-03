---
index: 3
title: "Generators"
---

The Leapkit CLI also provides generator that helps you to create different components such as migrations, handlers or actions.


```bash
$ kit [generate|g] commands
```

## Generating migrations

To generate a new migration, you need to run the `kit g migration <migration_name>`. It will generate a `.sql` file where you can place the change that wou will apply to the database.

```bash
$ kit g migration create_users_table
✅ Migration file `create_users_table` generated
```

```text
├── internal/
│    └── migrations/
│         └── 20060102030405_create_users_table.sql
```

The new migration file will follow the naming convention `yyyyMMddhhmmss_migration_name.sql`, and will be placed in the `internal/migrations` folder by default.