#base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frontApp ./cmd/web

RUN chmod +x /app/frontApp

#build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/frontApp /app

CMD [ "/app/frontApp" ]