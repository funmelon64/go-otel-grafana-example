# Используем официальный образ Golang как базовый
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum в рабочую директорию
COPY ./price-calcs/go.mod ./price-calcs/go.sum ./
COPY pkg /pkg
# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта в рабочую директорию
COPY ./price-calcs/ .

# Собираем приложение
RUN go build -o price-calcs ./cmd/start

# Команда для запуска приложения
CMD ["./price-calcs"]