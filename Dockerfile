# Dockerfile
FROM golang:1.24-alpine


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/subscription-aggregator

CMD ["./app"]