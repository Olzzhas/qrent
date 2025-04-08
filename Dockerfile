FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o qrent ./cmd/api

EXPOSE 4000

# Запускаем приложение
CMD ["./qrent"]
