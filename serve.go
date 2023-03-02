// Package chatRoom
// @Time  : 2023/3/2 15:03
// @Email: jtyoui@qq.com
// @Author: 张伟
package chatRoom

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gounits/gohtml"
	"log"
	"net/http"
	"strings"
	"time"
)

// Data webSocket传播的用户基本数据
type Data struct {
	IP       string   `json:"ip"`        // ip地址
	User     string   `json:"user"`      // 用户名
	Type     string   `json:"type"`      // 类型：login 登录，handshake 打开网页的状态,system 系统信息,logout 登出,user 普通信息
	Content  string   `json:"content"`   // 数据内容
	UserList []string `json:"user_list"` // 用户列表
	Time     string   `json:"time"`      // 信息发布时间
}

// Connection 连接的基本信息
type Connection struct {
	WebSocket *websocket.Conn //  websocket的连接信息
	Channel   chan []byte     //  管道信息
	Data      *Data           // 基本信息
}

// Hub websocket的连接器
type Hub struct {
	Connection map[*Connection]bool // 已经注册成功的连接器
	Send       chan []byte          // 连接器发送的数据
	Register   chan *Connection     // 注册的请求
	Unregister chan *Connection     // 注销的请求
}

var (
	userList []string // 用户列表

	Upgrade = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	// 初始化一个全局变量
	hub = Hub{
		Connection: make(map[*Connection]bool),
		Send:       make(chan []byte),
		Register:   make(chan *Connection),
		Unregister: make(chan *Connection),
	}
)

func (h *Hub) run() {
	for {
		select {
		case r := <-h.Register: // 获取注册信息
			h.Connection[r] = true
			r.Data.IP = r.WebSocket.RemoteAddr().String()
			r.Data.Type = "handshake"
			r.Data.UserList = userList
			js, _ := json.Marshal(r.Data)
			r.Channel <- js
		case u := <-h.Unregister: // 获取注销的信息
			if _, ok := h.Connection[u]; ok {
				delete(h.Connection, u)
				close(u.Channel)
			}
		case data := <-h.Send: // 获取发送的信息
			for c := range h.Connection {
				if strings.Contains(string(data), "@") { //增加匿名聊天
					tag := &Data{}
					if err := json.Unmarshal(data, tag); err != nil {
						fmt.Println(err)
						continue
					}
					prefix := strings.Split(tag.Content, "@")
					if c.Data.User == prefix[0] || c.Data.User == tag.User {
						c.Channel <- data
					}
				} else {
					select {
					case c.Channel <- data:
					}
				}
			}
		}
	}
}

// Writer 向websocket里面写数据
func (c *Connection) Writer() {
	for message := range c.Channel {
		if err := c.WebSocket.WriteMessage(websocket.TextMessage, message); err != nil { // 写入客户端（网页）
			fmt.Println("写入客户端数据异常")
			break
		}
	}
	if err := c.WebSocket.Close(); err != nil {
		fmt.Println("关闭WebSocket异常！")
	}
}

// Reader 向websocket里面拿数据
func (c *Connection) Reader() {
	for {
		_, message, err := c.WebSocket.ReadMessage() // 从客户端（网页）读数据
		if err != nil {
			hub.Unregister <- c // 读取数据失败默认移除用户
			break
		}
		if err = json.Unmarshal(message, c.Data); err != nil {
			fmt.Println("无法解析网页数据:" + err.Error())
		}
		if c.Data.Type == "login" {
			if len(strings.Trim(c.Data.Content, " ")) > 0 {
				c.Data.User = c.Data.Content
				userList = append(userList, c.Data.User)
				c.Data.UserList = userList
			}
		}
		c.Data.Time = time.Now().Format("2006-01-02 15:04:05")
		js, _ := json.Marshal(c.Data)
		hub.Send <- js
	}
}

// deleteUserList 删除用户列表中的某一个用户
func deleteUserList(userList []string, user string) []string {
	var slice []string
	fmt.Println(userList, user)
	for _, username := range userList {
		if username != user {
			slice = append(slice, username)
		} else {
			break
		}
	}
	return slice
}

// handle 路由
func handle(c *gin.Context) {
	ws, _ := Upgrade.Upgrade(c.Writer, c.Request, nil)
	connection := &Connection{WebSocket: ws, Data: &Data{}, Channel: make(chan []byte, 1024)}
	hub.Register <- connection
	go connection.Writer()
	connection.Reader()
	defer func() {
		connection.Data.Type = "logout"
		fmt.Println("删除：" + connection.Data.User)
		userList = deleteUserList(userList, connection.Data.User)
		connection.Data.UserList = userList
		hub.Unregister <- connection
		js, _ := json.Marshal(c.Data)
		hub.Send <- js
	}()
}

//go:embed web/dist
var efs embed.FS

func Serve() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	go hub.run()
	r.Use(gohtml.NewFs(efs))
	r.GET("/ws", handle)
	fmt.Println("启动成功： http://127.0.0.1:11280")
	if err := r.Run("0.0.0.0:11280"); err != nil {
		log.Fatalln("端口被占用！")
	}
}
