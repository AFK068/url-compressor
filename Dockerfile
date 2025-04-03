FROM golang:1.23.2 as builder

ENV GOPATH=/
WORKDIR /app
COPY ./ /app

RUN go mod download && CGO_ENABLED=0 go build -o /compressor ./cmd/run/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /compressor /app/compressor
COPY ./config/dev.yaml /app/config/dev.yaml

CMD ["/app/compressor"]