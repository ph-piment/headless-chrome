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

ENV DATABASES="postgres mysql redshift cassandra spanner cockroachdb clickhouse mongodb sqlserver firebird sqlite3 neo4j"
ENV SOURCES="file go_bindata github github_ee aws_s3 google_cloud_storage godoc_vfs gitlab"

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.10.0/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin && mv /usr/local/bin/migrate.linux-amd64 /usr/local/bin/migrate
