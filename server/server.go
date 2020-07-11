package server

import (
	"RedisFox/conf"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zer0131/logfox"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer(ctx context.Context) *Server {
	server := new(Server)
	server.srv = new(http.Server)

	//关闭调试模式
	if conf.ConfigVal.BaseVal.Debugmode == 0 {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(logger())
	router.Use(gin.Recovery())

	//静态文件处理
	router.Static("/static", conf.ConfigVal.BaseVal.Staticdir)

	//模板变量标识
	router.Delims("{[{", "}]}")

	//首页
	router.LoadHTMLFiles(conf.ConfigVal.BaseVal.Tpldir + "index.html")
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
	server.srv.Addr = conf.ConfigVal.BaseVal.Serverip + ":" + strconv.Itoa(conf.ConfigVal.BaseVal.Serverport)
	server.srv.Handler = router

	go server.start(ctx)
	logfox.InfoWithContext(ctx, "web server start")

	return server
}

//自定义日志
func logger() gin.HandlerFunc {

	return func(c *gin.Context) {
		//针对每个请求生成新的logid
		logid := logfox.GenLogId()
		c.Set(logfox.LogIDKey, logid)

		start := time.Now()
		uri := c.Request.RequestURI
		c.Next()
		end := time.Now()
		latencyTime := end.Sub(start)
		method := c.Request.Method
		statusCode := c.Writer.Status()
		clientIp := c.ClientIP()
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()
		userAgent := c.Request.UserAgent()

		logfox.Infof("[GIN] %v logid[%s] runtime[%v] uri[%s]  method[%s] status[%d] client_ip[%s] user_agent[%s] error_message[%s] \n", start.Format("2006/01/02 - 15:04:05"), logid, latencyTime, uri, method, statusCode, clientIp, userAgent, errMsg)
	}
}

func (s *Server) start(ctx context.Context) {
	if err := s.srv.ListenAndServe(); err != nil {
		logfox.ErrorfWithContext(ctx, "listen %s errmsg: %s", s.srv.Addr, err.Error())
	}
}

//graceful stop需Go1.8+
func (s *Server) Stop(ctx context.Context) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := s.srv.Shutdown(newCtx); err != nil {
		logfox.ErrorfWithContext(ctx, "web server shutdown errmsg: %s", err.Error())
	}
	logfox.InfoWithContext(ctx, "web server shutdown")
}
