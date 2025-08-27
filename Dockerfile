FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/member-handler-go /root/member-handler-go
RUN chmod +x /root/member-handler-go
ENTRYPOINT ["/root/member-handler-go"]
