FROM golang:1.17-alpine

WORKDIR /go/src/github.com/bazaglia/warehouse
COPY . .

RUN go mod tidy && \
    go build -ldflags="-s -w" -o bin/server cmd/server/main.go

CMD ["./bin/server"]