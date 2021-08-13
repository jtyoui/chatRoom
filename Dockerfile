FROM golang:alpine
MAINTAINER Jtyoui <jtyoui@qq.com>

RUN mkdir /app
WORKDIR /app

COPY ./main.go ./main.go
COPY ./proxy.go ./proxy.go
COPY ./go.mod ./go.mod
COPY ./dist ./dist

ENV GOPROXY https://goproxy.cn,direct
ENV GOSUMDB off
ENV GO111MODULE on

RUN go mod tidy && go build

CMD ["./chatRoom"]