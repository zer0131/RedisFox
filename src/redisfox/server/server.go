package server

import (
	"redisfox/conf"
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"redisfox/flog"
	"time"
	"context"
)

type Server struct {
	config *conf.Config
	srv *http.Server
}

func NewServer(config *conf.Config) *Server {
	server := new(Server)
	server.config = config
	server.srv = new(http.Server)

	router := gin.Default()

	//静态文件处理
	router.Static("/static", server.config.Staticdir)

	//模板变量标识
	router.Delims("{[{", "}]}")

	//首页
	router.LoadHTMLFiles(server.config.Tpldir+"index.html")
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})

	//接口
	router.GET("/api/servers", server.serverList)
	router.GET("/api/info", server.info)
	router.GET("/api/memory", server.memory)
	router.GET("/api/commands", server.commands)
	router.GET("/api/topcommands", server.topcommands)
	router.GET("/api/topkeys", server.topkeys)

	//srv设置
	server.srv.Addr = server.config.Serverip+":"+strconv.Itoa(server.config.Serverport)
	server.srv.Handler = router

	go server.start()
	flog.Infof("web server start")

	return server
}

func (this *Server) start() {
	err := this.srv.ListenAndServe()
	if err != nil {
		flog.Fatalf("web server start error: "+err.Error())
	}
}

//graceful stop需Go1.8+
func (this *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := this.srv.Shutdown(ctx); err != nil {
		flog.Fatalf("web server shutdown error: "+err.Error())
	}
	flog.Infof("web server shutdown")
}
