FROM golang:1.25 AS builder
WORKDIR /app
RUN go mod init github.com/halushko/member-handler-go
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /app/member-handler-go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/member-handler-go .
CMD ["./member-handler-go"]
