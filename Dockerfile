FROM golang

ENV APP_DIR /go/src/github.com/mikkel-larsen/orchid
ADD . $APP_DIR

WORKDIR $APP_DIR

RUN go get -t ./...
RUN go install

EXPOSE 80
