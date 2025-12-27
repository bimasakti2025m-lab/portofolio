-- Active: 1725517175358@@127.0.0.1@5432@catatan_keuangan_db@public
CREATE DATABASE catatan_keuangan_db;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE transaction_type AS ENUM ('CREDIT', 'DEBIT');
CREATE TYPE roles AS ENUM ('user', 'admin', 'superadmin');

CREATE TABLE expenses (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    date DATE NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    transaction_type transaction_type,
    balance DOUBLE PRECISION NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE budgets (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    amount DOUBLE PRECISION NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

SELECT * FROM expenses;

ALTER TABLE expenses ADD COLUMN user_id uuid;

ALTER TABLE users ALTER COLUMN role TYPE roles;

ALTER TABLE users ADD COLUMN role roles;

ALTER TABLE users DROP COLUMN role;

SELECT * FROM users;

-- ALTER TABLE expenses ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

SELECT id, username, password, role, created_at, updated_at FROM users WHERE username= ' or 1=1--