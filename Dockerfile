FROM golang

LABEL "space.opengate.vendor"="SpaceRouter"
LABEL org.opencontainers.image.source https://github.com/SpaceRouter/authentication_server
LABEL space.opengate.image.authors="theo.lefevre@edu.esiee.fr"

COPY src /source
WORKDIR /source

RUN apt update && apt install libpam0g-dev -y && apt clean

RUN go get && \
    go build -o /usr/bin/auth_server && \
    rm $GOPATH -rf

RUN mkdir /config && cp config/*.yaml /config -r

WORKDIR /

ENV GIN_MODE=release

CMD auth_server
