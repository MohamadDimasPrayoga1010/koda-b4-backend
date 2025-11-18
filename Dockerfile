FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o coffeeder main.go

FROM alpine:latest 

WORKDIR /backend

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/coffeeder /backend/coffeeder

EXPOSE 8080

ENTRYPOINT ["/backend/coffeeder"]