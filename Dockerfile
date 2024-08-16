# Используем базовый образ для сборки Go приложений
FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию для сборки
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем все файлы проекта
COPY . .

# Собираем бинарник
RUN go build -o discord-bot main.go

# Используем минимальный образ для запуска приложения
FROM debian:bullseye-slim

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/discord-bot /discord-bot

# Запускаем бот
CMD ["/discord-bot"]
