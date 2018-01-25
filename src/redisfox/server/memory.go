package server

import (
	"github.com/gin-gonic/gin"
	"time"
	"redisfox/dataprovider"
	"net/http"
)

func (this *Server) memory(context *gin.Context)  {
	serverId := context.Query("server")
	fromDate := context.DefaultQuery("from", "")
	toDate := context.DefaultQuery("to", "")

	var start string
	var end string
	now := time.Now()
	if fromDate == "" || toDate == "" {
		end = now.Format("2006-01-02 15:04:05")
		endTmp,_ := time.ParseDuration("-60s")
		start = now.Add(endTmp).Format("2006-01-02 15:04:05")
	} else {
		start = fromDate
		end = toDate
	}

	sqlDb,_ := dataprovider.NewProvider(this.config)
	defer sqlDb.Close()

	memoryList,_ := sqlDb.GetMemoryInfo(serverId, start, end)

	var data [][]string
	for _,v := range memoryList {
		data = append(data, []string{v["datetime"].(string),v["peak"].(string),v["used"].(string)})
	}

	context.JSON(http.StatusOK, gin.H{
		"data":data,
		"timestamp":now.Format("2006-01-02 15:04:05"),
	})
}
