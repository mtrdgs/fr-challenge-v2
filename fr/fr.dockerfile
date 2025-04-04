FROM golang:1.23 AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frApp ./cmd/api

RUN chmod +x /app/frApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/frApp /app

CMD [ "/app/frApp" ]