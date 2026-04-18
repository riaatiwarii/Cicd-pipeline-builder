FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o processor main.go

# Final stage
FROM docker:latest

RUN apk add --no-cache ca-certificates git bash

WORKDIR /root/

COPY --from=builder /app/processor .

CMD ["./processor"]
