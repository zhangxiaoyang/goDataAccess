FROM golang:latest
MAINTAINER zhangxiaoyang <zhangxiaoyang.hit@gmail.com>

# Initialize env
ENV DA_NAME github.com/zhangxiaoyang/goDataAccess
ENV DA_REPO https://$DA_NAME
ENV DA_PATH /go/src/$DA_NAME
ENV DA_VERSION 1.0
RUN mkdir -p $DA_PATH
RUN git clone $DA_REPO $DA_PATH

# Install dep
RUN go get github.com/gpmgo/gopm
RUN go get github.com/mattn/go-sqlite3
RUN gopm get -g golang.org/x/net/html/charset

# Enjoy
WORKDIR $DA_PATH
