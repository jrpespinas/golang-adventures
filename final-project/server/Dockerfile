FROM golang:1.17-alpine

WORKDIR /server

COPY database ./database
COPY data ./data
COPY server.go .

RUN go mod init server
RUN go get github.com/dgrijalva/jwt-go

RUN go build -o server

EXPOSE 8080

CMD ["./server"]


