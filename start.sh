#!/bin/bash

DIR=${PWD}

docker-compose down --v

# 创建文件夹
if [ ! -d "$DIR"/app ]; then
  mkdir "$DIR"/app
fi

# 判断有没有编译好的ws二进制文件
if [ ! -f "$DIR"/app/ws ]; then
  docker build -t ws .
  docker run --rm -d -v "$DIR"/app:/home ws go build -o /home/ws -ldflags '-w -s' -gcflags '-N -l' main.go
fi

# 拉取前端代码
if [ ! -d "$DIR"/app/chatRoomFront ];then
   git clone https://github.com/jtyoui/chatRoomFront.git  "$DIR"/app/chatRoomFront
   docker build -f Dockerfile.front -t front .
fi

# 编译前端代码
if [ ! -d "$DIR"/app/dist ]; then
  docker run --rm -d -v "$DIR"/app/dist:/home/dist front npm run build
fi

docker-compose up -d
