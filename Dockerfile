# Этап сборки
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Установка необходимых утилит
RUN apk add --no-cache git ca-certificates

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# Финальный этап
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Создание непривилегированного пользователя
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Копируем бинарный файл
COPY --from=builder /app/main .

# Устанавливаем владельца
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 4200

CMD ["./main"]