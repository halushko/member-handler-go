FROM golang:1.25 AS builder
WORKDIR /app
RUN go mod init github.com/halushko/halushko-ist-chat-bot
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /app/halushko-ist-chat-bot

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/member-handler-go /root/member-handler-go
RUN chmod +x /root/member-handler-go
ENTRYPOINT ["/root/member-handler-go"]
