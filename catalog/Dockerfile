FROM golang

USER root
RUN apt-get update

WORKDIR /go/src/app
COPY ../.. .

RUN go get github.com/pilu/fresh
RUN go get ./...

CMD [ "fresh" ]
