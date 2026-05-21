# Этап сборки
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Устанавливаем необходимые утилиты
RUN apk add --no-cache git

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main cmd/main.go

# Финальный этап
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Копируем бинарный файл
COPY --from=builder /app/main .
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/migrations ./migrations

# Делаем скрипт исполняемым
RUN chmod +x ./scripts/wait-for-it.sh

EXPOSE 4200

CMD ["./scripts/wait-for-it.sh", "postgres:5432", "--", "./main"]