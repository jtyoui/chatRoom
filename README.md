# chatRoom

一个简单使用方便部署的聊天系统

## 前提条件

1. 安装git
2. 安装docker

## 直接启动

```shell
git clone https://github.com/jtyoui/chatRoom.git
LINUX点击运行： ./chatRoom
WINDOW点击运行： chatRoom.exe
```

## Docker部署

```shell
git clone https://github.com/jtyoui/chatRoom.git
docker build -t chatroom .
docker run -d -p 28181:28181 -p 11280:11280 chatroom
```

## 聊天界面

![聊天界面](./home.png)

## 前端代码

[点击跳转到界面地址](https://github.com/jtyoui/chatRoomFront)

## 功能

- [x] 支持群聊
- [x] 支持私聊
- [ ] 支持emoji表情
- [ ] 支持文件传输
