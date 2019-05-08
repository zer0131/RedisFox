package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (this *Server) serverList(context *gin.Context)  {
	var ret []gin.H
	for _,v := range this.config.Servers {
		ret = append(ret, gin.H{
			"id":v["server"]+":"+v["port"],
			"server":v["server"],
			"port":v["port"],
			"password":v["password"],
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"servers": ret,
	})
}