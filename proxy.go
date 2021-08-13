package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const DIR = "./dist"

func Proxy() {
	r := gin.Default()
	r.LoadHTMLGlob(DIR + "/*.html") // 添加入口index.html
	r.Static("/favicon.ico", DIR+"/favicon.ico")
	r.Static("/assets", DIR+"/assets")        // 添加资源路径
	r.StaticFile("/", DIR+"/index.html")      //前端接口
	r.StaticFile("/login", DIR+"/index.html") //前端接口
	fmt.Println("http://localhost:28181")
	err := r.Run(":28181")
	if err != nil {
		panic("端口被占用")
	}
}
