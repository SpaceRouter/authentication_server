FROM golang:alpine

COPY src /source
WORKDIR /source

RUN apk add gcc --no-cache --purge -uU linux-pam-dev musl-dev

RUN go get && \
    go build -o /usr/bin/auth_server && \
    rm $GOPATH -rf

RUN mkdir /config && cp config/*.yaml /config -r

WORKDIR /

ENV GIN_MODE=release

CMD auth_server