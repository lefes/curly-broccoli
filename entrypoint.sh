#!/bin/bash

DB_FILE="./data/bot.db"

# Проверяем, существует ли файл базы данных
until [ -f "$DB_FILE" ]; do
  echo "Waiting for the SQLite database file to appear..."
  sleep 2
done

echo "Migration UP..."
goose -dir ./db/migrations sqlite3 "$DB_FILE" up

echo "Starting the bot..."
go run .

