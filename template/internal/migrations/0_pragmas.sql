-- This folder holds the database migrations for the application
-- These migrations will run and then update the database schema
-- as the application is updated

-- These are some default settings for the SQLite database,
-- these can be changed or removed as needed.
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 2000;
PRAGMA temp_store = memory;
