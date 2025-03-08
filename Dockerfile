FROM golang:1.23.3 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o lancer .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/lancer .

EXPOSE 8080

ENTRYPOINT ["./lancer", "--config", "/root/config.yaml"]