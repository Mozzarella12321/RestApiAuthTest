CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login TEXT,
    hash BYTEA,
    unsuccessful_logins INT
);

CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    token INT,
    created_at TIMESTAMP
);