FROM golang:1.25.0-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /main .

EXPOSE 8089

CMD ["./main"]