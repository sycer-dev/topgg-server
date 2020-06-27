FROM golang:1.14-alpine

WORKDIR /usr/src/topgg-server

COPY go.mod go.sum main.go ./

RUN go mod download
RUN go build -o server .

ENTRYPOINT ["./server"]