package server

import "github.com/gin-gonic/gin"

func (this *Server) memory(context *gin.Context)  {
	//serverId := context.Query("server")
	fromDate := context.DefaultQuery("from", "")
	toDate := context.DefaultQuery("to", "")

	if fromDate == "" || toDate == "" {
		//
	} else {
		//
	}
}
