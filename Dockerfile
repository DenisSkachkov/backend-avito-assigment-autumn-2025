FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

ENV PORT=8080

EXPOSE 8080

CMD ["./app"]
