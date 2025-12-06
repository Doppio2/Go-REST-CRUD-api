FROM golang:1.24.7 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cli/api_service/main.go

FROM alpine:3.19
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
