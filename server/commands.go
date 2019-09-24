package server

import (
	"github.com/gin-gonic/gin"
	"time"
	"RedisFox/dataprovider"
	"net/http"
)

func (s *Server) commands(context *gin.Context)  {
	serverId := context.Query("server")
	fromDate := context.DefaultQuery("from", "")
	toDate := context.DefaultQuery("to", "")

	var start string
	var end string
	now := time.Now()
	layout := "2006-01-02 15:04:05"
	if fromDate == "" || toDate == "" {
		end = now.Format(layout)
		endTmp,_ := time.ParseDuration("-120s")
		start = now.Add(endTmp).Format(layout)
	} else {
		start = fromDate
		end = toDate
	}
	startTime,_ := time.Parse(layout, start)
	endTime,_ := time.Parse(layout, end)
	difference := endTime.Unix()-startTime.Unix()

	minutes := difference/60
	hours := minutes/60

	//group by设置
	var groupBy string

	if hours > 120 {
		groupBy = "day"
	} else if minutes > 120 {
		groupBy = "hour"
	} else if difference > 120 {
		groupBy = "minute"
	} else {
		groupBy = "second"
	}

	sqlDb,_ := dataprovider.NewProvider(s.config)
	defer sqlDb.Close()

	commandStats,_ := sqlDb.GetCommandStats(serverId, start, end, groupBy)

	var data [][]string
	for _,v := range commandStats {
		data = append(data, []string{v["datetime"].(string),v["total"].(string)})
	}

	context.JSON(http.StatusOK, gin.H{
		"data":data,
		"timestamp":now.Format(layout),
	})
}
