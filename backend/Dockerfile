FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod tidy

COPY . .
RUN go mod tidy && go build -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 3000
CMD ["./server"]
