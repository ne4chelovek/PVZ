# Dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o pvz-app cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/pvz-app .

COPY --from=builder /app/configs /root/configs
EXPOSE 8080 9000
CMD ["./pvz-app"]