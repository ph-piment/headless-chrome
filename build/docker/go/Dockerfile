FROM golang:latest

RUN mkdir -p /go/src/work

WORKDIR /go/src/work

ADD . /go/src/work

RUN go get -u -t github.com/pkg/errors
RUN go get -u -t github.com/chromedp/chromedp
RUN go get -u -t github.com/orisano/pixelmatch
RUN go get -u -t github.com/BurntSushi/toml
RUN go get -u -t github.com/lib/pq
RUN go get -u -t github.com/jinzhu/gorm
