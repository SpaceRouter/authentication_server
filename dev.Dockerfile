FROM golang

VOLUME /web
EXPOSE 8080

WORKDIR /web

RUN apt update && apt full-upgrade -y && apt install libpam0g-dev -y

CMD go get && go run main.go