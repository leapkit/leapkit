---
index: 1
title: "Creating new app"
---

To start a new Leapkit app, you need to run the following command:

```bash
$ kit new <app_name>
```

This will create a new directory with the name of your app into the current path.

## Folder structure

The Leapkit directory will contain the following structure:

```text
├── bin/
├── cmd/
│    ├── app/
│    │    └── main.go
│    ├── migrate/
│    │    └── migrate.go
│    └── setup/
│         └── main.go
├── internal/
│    ├── assets/
│    │    ├── application.css
│    │    └── application.js
│    ├── home/
│    │    ├── index.go
│    │    └── index.html
│    ├── migrations/
│    │    ├── 0_pragmas.sql
│    │    └── migrations.go
│    ├── app.go
│    └── layout.html
├── public/
│    └── public.go
├── .dockerignore
├── database.db
├── .Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── README.md
└── tailwind.config.js
```