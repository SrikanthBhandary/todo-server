-- Down Migration: Add the email column back to todos
ALTER TABLE todos ADD COLUMN email VARCHAR(300) UNIQUE NOT NULL;
