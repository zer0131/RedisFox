package server

import (
	"RedisFox/conf"
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"time"
	"context"
	"github.com/zer0131/logfox"
)

type Server struct {
	config *conf.Config
	srv *http.Server
}

func NewServer(config *conf.Config) *Server {
	server := new(Server)
	server.config = config
	server.srv = new(http.Server)

	//关闭调试模式
	if server.config.Debugmode == 0 {
		gin.SetMode(gin.ReleaseMode)
	}

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
	logfox.Info("web server start")

	return server
}

func (s *Server) start() {
	if err := s.srv.ListenAndServe(); err != nil {
		logfox.Errorf("listen %s errmsg: %s", s.srv.Addr, err.Error())
	}
}

//graceful stop需Go1.8+
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		logfox.Errorf("web server shutdown errmsg: %s", err.Error())
	}
	logfox.Info("web server shutdown")
}
