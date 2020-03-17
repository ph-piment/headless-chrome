FROM golang:latest

RUN mkdir -p /go/src/work

WORKDIR /go/src/work

ADD . /go/src/work

RUN go get github.com/chromedp/chromedp
