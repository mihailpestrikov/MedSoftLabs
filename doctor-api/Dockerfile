FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o doctor-api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/doctor-api .

EXPOSE 8081

CMD ["./doctor-api"]
