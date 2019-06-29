FROM golang:1.11.11-alpine3.10


RUN apk add --no-cache git mercurial

RUN go get github.com/alabianca/rapi-api
WORKDIR /go/src/github.com/alabianca/rapi-api
# RUN go install ./...

COPY . .

# RUN go build -o bin/rapid
RUN go install ./...
CMD [ "rapi-api" ]
