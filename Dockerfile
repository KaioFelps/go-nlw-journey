FROM golang:1.22.5-alpine

WORKDIR /nlw-journey

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR /nlw-journey/cmd/journey

RUN go build -o /nlw-journey/bin/journey .

EXPOSE 8080
ENTRYPOINT [ "/nlw-journey/bin/journey" ]