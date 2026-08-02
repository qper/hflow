CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(64) NOT NULL,
    email VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NULL,
    timezone VARCHAR(64) NOT NULL DEFAULT 'UTC',
    theme VARCHAR(16) NOT NULL DEFAULT 'system',
    parent_user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX ux_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX ux_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX ix_users_parent_user_id ON users(parent_user_id);

CREATE TABLE habits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    habit_type VARCHAR(32) NOT NULL,
    frequency VARCHAR(32) NOT NULL DEFAULT 'daily',
    target_value NUMERIC(10,2) NULL,
    unit VARCHAR(32) NULL,
    color VARCHAR(32) NULL,
    icon VARCHAR(64) NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    archived_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX ix_habits_user_id ON habits(user_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_habits_archived_at ON habits(archived_at);
CREATE INDEX ix_habits_sort_order ON habits(user_id, sort_order);
CREATE INDEX ix_habits_name_trgm ON habits USING gin (name gin_trgm_ops);

CREATE TABLE habit_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    entry_date DATE NOT NULL,
    value NUMERIC(10,2) NULL,
    note TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    UNIQUE(habit_id, entry_date)
);

CREATE INDEX ix_habit_entries_habit_date ON habit_entries(habit_id, entry_date DESC);
CREATE INDEX ix_habit_entries_user_date ON habit_entries(user_id, entry_date DESC);
CREATE INDEX ix_habit_entries_date ON habit_entries(entry_date);

CREATE TABLE streaks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, habit_id)
);

CREATE INDEX ix_streaks_habit_id ON streaks(habit_id);
CREATE INDEX ix_streaks_user_id ON streaks(user_id);
