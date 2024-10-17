FROM golang:1.23-alpine AS builder

COPY . /github.com/Oleg-Pro/chat-server
WORKDIR /github.com/Oleg-Pro/chat-server



RUN go mod download
RUN go build -o ./bin/chat_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/Oleg-Pro/chat-server/bin/chat_server .

ADD .env .

CMD ["./chat_server"]