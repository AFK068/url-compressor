FROM golang:1.23.2 as builder

ENV GOPATH=/
WORKDIR /app
COPY ./ /app

RUN go mod download && CGO_ENABLED=0 go build -o /compressor ./cmd/run/main.go

FROM alpine:latest

RUN apk add --no-cache postgresql-client

WORKDIR /app
COPY --from=builder /compressor /app/compressor
COPY ./config/dev.yaml /app/config/dev.yaml
COPY ./migrations /app/migrations

COPY wait-for-postgres.sh /app/wait-for-postgres.sh
RUN chmod +x /app/wait-for-postgres.sh

CMD ["/app/compressor"]