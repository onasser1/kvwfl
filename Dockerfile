FROM golang:1.26.3 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o validating-kontroller .

FROM alpine

WORKDIR /app

COPY --from=builder /app/validating-kontroller /usr/local/bin/validating-kontroller
COPY --from=builder /app /app
CMD ["validating-kontroller"]