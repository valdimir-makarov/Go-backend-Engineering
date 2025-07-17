CREATE EXTENSION IF NOT EXISTS "uuid-ossp";          -- only once per DB
CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    uuid       UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);