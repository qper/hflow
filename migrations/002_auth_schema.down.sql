DROP TABLE IF EXISTS refresh_tokens;
ALTER TABLE users
    DROP COLUMN IF EXISTS failed_login_count,
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS password_hash;
