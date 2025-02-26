# Указываем базовый образ с поддержкой Go
FROM golang:alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Создаем исполняемую и служебные директории
RUN mkdir cmd internal

# Копируем в них файлы
COPY ./cmd ./cmd
COPY ./internal ./internal

# Устанавливаем зависимости
RUN go mod download

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/main.go

# Указываем точку входа
ENTRYPOINT ["./main"]