FROM golang:1.11.11-alpine3.10 as builder

RUN apk add --no-cache git mercurial
# add gcc and g++ ro run tests
RUN apk add --update gcc
RUN apk add --update g++

RUN go get github.com/alabianca/rapi-api
WORKDIR /go/src/github.com/alabianca/rapi-api

COPY . .

# RUN go build -o bin/rapid
# RUN go install ./...
# CMD [ "rapi-api" ]
RUN go build -o bin/rapid

FROM alpine
WORKDIR /app

COPY --from=builder /go/src/github.com/alabianca/rapi-api/bin/rapid /app/
CMD ["./rapid"]

