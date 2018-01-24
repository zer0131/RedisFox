package server

import (
	"redisfox/conf"
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
)

type Server struct {
	config *conf.Config
	router *gin.Engine
}

func NewServer(config *conf.Config)  {
	server := new(Server)
	server.config = config
	server.router = gin.Default()

	//静态文件处理
	server.router.Static("/static", server.config.Staticdir)

	//模板变量
	server.router.Delims("{[{", "}]}")

	//首页
	server.router.LoadHTMLFiles(server.config.Tpldir+"index.html")
	server.router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})

	//接口
	server.router.GET("/api/servers", server.serverList)
	server.router.GET("/api/info", server.info)
	server.router.GET("/api/memory", server.memory)
	server.router.GET("/api/commands", server.commands)
	server.router.GET("/api/topcommands", server.topcommands)
	server.router.GET("/api/topkeys", server.topkeys)

	server.router.Run(server.config.Serverip+":"+strconv.Itoa(server.config.Serverport))
}
