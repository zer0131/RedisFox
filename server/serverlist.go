package server

import (
	"RedisFox/conf"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) serverList(context *gin.Context) {
	var ret []gin.H
	for _, v := range conf.ConfigVal.BaseVal.Servers {
		ret = append(ret, gin.H{
			"id":       v["server"] + ":" + v["port"],
			"server":   v["server"],
			"port":     v["port"],
			"password": v["password"],
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"servers": ret,
	})
}
