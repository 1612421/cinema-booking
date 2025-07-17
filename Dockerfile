# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Cài các tool cần thiết
RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./config/dev.yml ./config/local.yml
RUN go build -tags musl --ldflags "-extldflags -static" -o ./bin/cinema-booking ./cmd/cinema-booking

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/cinema-booking .
COPY wait-for-mysql.sh .

# Expose port (đổi nếu bạn dùng cổng khác)
EXPOSE 8080

CMD ["./wait-for-mysql.sh", "./cinema-booking"]
