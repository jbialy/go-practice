# Build stage
FROM golang:1.25 AS builder
WORKDIR /app

COPY *.go .
COPY go.mod .
RUN go build -o myapp .

# Runtime stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/myapp .

ENTRYPOINT [ "./myapp" ]
