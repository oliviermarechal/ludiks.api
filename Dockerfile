FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .


FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]