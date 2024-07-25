# Первый этап: сборка приложения
FROM golang:1.22-alpine AS build_stage

# Копируем весь исходный код в рабочую директорию контейнера
COPY . /app

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod .
COPY go.sum .

# Копируем файл окружения
COPY .env .env

# Скачиваем зависимости
RUN go mod download

# Сборка приложения с отключением CGO
RUN CGO_ENABLED=0 go build -o /app_binary/EMtask ./cmd

# Второй этап: создание минимального конечного образа
FROM alpine:latest AS run_stage

# Устанавливаем рабочую директорию
WORKDIR /app_binary

# Копируем скомпилированное приложение из первого этапа
COPY --from=build_stage /app_binary/EMtask /app_binary/

# Делаем бинарный файл исполняемым
RUN chmod +x ./EMtask

# Указываем команду по умолчанию для запуска контейнера
ENTRYPOINT ["./EMtask"]
CMD ["EMtask"]
