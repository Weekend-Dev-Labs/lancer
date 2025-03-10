FROM node:20-alpine AS vite-builder
WORKDIR /dashboard

COPY dashboard/package.json dashboard/package-lock.json ./
RUN npm install --force

COPY dashboard/ .
RUN npm run build

FROM golang:1.23.3 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=vite-builder /dashboard/dist ./dashboard/dist

RUN CGO_ENABLED=0 GOOS=linux go build -o lancer .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/lancer .

EXPOSE 8080

# No ENTRYPOINT or CMD here (let docker-compose handle it)