FROM golang:latest

RUN mkdir -p /go/src/work

WORKDIR /go/src/work

ADD . /go/src/work

RUN go get github.com/pkg/errors
RUN go get github.com/chromedp/chromedp
RUN go get github.com/orisano/pixelmatch
RUN go get github.com/ph-piment/headless-chrome/internal/browser
