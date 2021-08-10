FROM golang:alpine
MAINTAINER Jtyoui <jtyoui@qq.com>

RUN mkdir /app
WORKDIR /app

COPY ./main.go /app/main.go
COPY ./go.mod /app/go.mod

ENV GOPROXY https://goproxy.cn,direct
ENV GOPRIVATE *.gitlab.com,*.gitee.com
ENV GOSUMDB off
ENV GO111MODULE on

RUN go mod tidy