-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    discord_id TEXT NOT NULL UNIQUE,
    role_id INTEGER NOT NULL DEFAULT 0,
    points INTEGER DEFAULT 0,
    points_today INTEGER DEFAULT 0,
    respect INTEGER DEFAULT 0,
    respect_today INTEGER DEFAULT 0,
    daily_messages INTEGER DEFAULT 0,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    respect_required INTEGER NOT NULL,
    privileges TEXT
);

CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER,
    receiver_id INTEGER,
    type TEXT NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS lotteries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fund INTEGER DEFAULT 0,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    winner_id INTEGER,
    active BOOLEAN DEFAULT 1
);

CREATE TABLE IF NOT EXISTS lottery_tickets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    lottery_id INTEGER NOT NULL,
    UNIQUE(user_id, lottery_id)
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS lotteries;
DROP TABLE IF EXISTS lottery_tickets;

