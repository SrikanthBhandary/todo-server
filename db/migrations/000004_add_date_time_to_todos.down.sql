-- File: migrations/xxxxxx_add_datetime_to_todos.down.sql
ALTER TABLE todos
DROP COLUMN IF EXISTS datetime;