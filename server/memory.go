package server

import (
	"RedisFox/dataprovider"
	"RedisFox/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s *Server) memory(context *gin.Context) {
	ctx := util.NewContextWithGinContext(context)
	serverId := context.Query("server")
	fromDate := context.DefaultQuery("from", "")
	toDate := context.DefaultQuery("to", "")

	var start string
	var end string
	now := time.Now()
	layout := "2006-01-02 15:04:05"
	if fromDate == "" || toDate == "" {
		end = now.Format(layout)
		endTmp, _ := time.ParseDuration("-60s")
		start = now.Add(endTmp).Format(layout)
	} else {
		start = fromDate
		end = toDate
	}

	sqlDb := dataprovider.NewProvider()

	memoryList, _ := sqlDb.GetMemoryInfo(ctx, serverId, start, end)

	var data [][]string
	for _, v := range memoryList {
		data = append(data, []string{v["datetime"].(string), v["peak"].(string), v["used"].(string)})
	}

	context.JSON(http.StatusOK, gin.H{
		"data":      data,
		"timestamp": now.Format(layout),
	})
}
