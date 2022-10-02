FROM golang:1.19

WORKDIR /usr/src/fan-control

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go .
RUN go build -v -o /usr/local/bin/fan-control /usr/src/fan-control/

CMD ["fan-control"]
