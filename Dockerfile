FROM golang:1.25 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --From=builder /app/main .
ENTRYPOINT ["./main"]