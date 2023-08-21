FROM golang:1.17-alpine

WORKDIR /app

COPY receipt-points-service/go.mod receipt-points-service/go.sum ./
RUN go mod download

COPY receipt-points-service/ ./

RUN go build -o app

EXPOSE 9000

CMD ["./app"]