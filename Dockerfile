FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
RUN mkdir -p /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY pkg/modifiedquickfix ./modifiedquickfix
RUN go mod edit -replace github.com/quickfixgo/quickfix=/app/modifiedquickfix
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ddexgo
RUN chmod +x ./ddexgo
CMD ./ddexgo
