FROM golang:1.11.11-alpine3.10

WORKDIR /go/src/github.com/alabianca/rapi-api

COPY . .

RUN apk add --no-cache git mercurial
RUN go get -v ./...
RUN go install ./...

CMD [ "rapi-api" ]
