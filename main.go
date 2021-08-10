package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

// Data webSocket传播的用户基本数据
type Data struct {
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
	Upgrade  = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true }}
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
					}
					prefix := strings.Split(tag.Content, "@")
					if c.Data.User == prefix[0] || c.Data.User == tag.User {
						c.Channel <- data
					}
				} else {
					select {
					case c.Channel <- data:
					default: // 防止长时间获取不到数据
						delete(h.Connection, c)
						close(c.Channel)
					}
				}
			}
		}
	}
}

// Writer 向websocket里面写数据
func (c *Connection) Writer() {
	for message := range c.Channel {
		fmt.Println(string(message))
		if err := c.WebSocket.WriteMessage(websocket.TextMessage, message); err != nil { // 写入客户端（网页）
			log.Fatalln("写入客户端数据异常")
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
			fmt.Println("获取网页信息识别，删除用户")
			hub.Unregister <- c // 读取数据失败默认移除用户
			break
		}
		if err = json.Unmarshal(message, c.Data); err != nil {
			log.Fatalln("无法解析网页数据")
		}
		if c.Data.Type == "login" {
			if c.Data.User != "" && len(c.Data.User) > 0 {
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
	count := len(userList)
	var slice []string
	if count > 1 { // 当用户列表为0和1的时候，删除一个，表示为空
		for index := range userList {
			if userList[index] == user {
				if index == count { // 最后一个是当前用户
					slice = userList[:index]
				} else {
					slice = append(userList[:index], userList[index+1:]...)
				}
				break
			}
		}
	}
	return slice
}

// handle 路由
func handle(w http.ResponseWriter, r *http.Request) {
	ws, _ := Upgrade.Upgrade(w, r, nil)
	c := &Connection{WebSocket: ws, Data: &Data{}, Channel: make(chan []byte, 1024)}
	hub.Register <- c
	go c.Writer()
	c.Reader()
	defer func() {
		c.Data.Type = "logout"
		userList = deleteUserList(userList, c.Data.User)
		c.Data.UserList = userList
		hub.Unregister <- c
		js, _ := json.Marshal(c.Data)
		hub.Send <- js
	}()
}

func main() {
	route := mux.NewRouter()
	go hub.run()
	route.HandleFunc("/ws", handle)
	if err := http.ListenAndServe(":8080", route); err != nil {
		log.Fatalln("链接异常！")
	}
}
