# Используем официальный образ Golang как базовый
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum в рабочую директорию
COPY ./bench/go.mod ./

# Копируем остальные файлы проекта в рабочую директорию
COPY ./bench/ .

# Собираем приложение
RUN go build -o bench ./main.go

# Команда для запуска приложения
CMD ["./bench"]