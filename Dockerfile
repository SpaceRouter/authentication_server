FROM golang

COPY src /web
WORKDIR /web

RUN apt update && apt install libpam0g-dev -y
RUN go get
RUN go build

CMD go run