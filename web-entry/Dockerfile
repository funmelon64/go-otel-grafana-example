# Используем официальный образ Golang как базовый
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum в рабочую директорию
COPY ./web-entry/go.mod ./web-entry/go.sum ./
COPY pkg /pkg
# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта в рабочую директорию
COPY ./web-entry/ .

# Собираем приложение
RUN go build -o webentry ./cmd/start

# Команда для запуска приложения
CMD ["./webentry"]