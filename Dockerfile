FROM golang:1.17.4-alpine3.15

ENV CGO_ENABLED=0

WORKDIR /app
ADD . .
RUN go build -o metric-api ./cmd/main.go

ARG PORT

EXPOSE $PORT

CMD ["/app/metric-api"]