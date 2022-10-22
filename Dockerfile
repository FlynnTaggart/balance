FROM golang:1.18-alpine

WORKDIR /go/src/balance

COPY go.mod go.sum /go/src/balance/
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/server

CMD ["app"]