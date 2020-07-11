package server

import (
	"RedisFox/dataprovider"
	"RedisFox/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) topkeys(context *gin.Context) {
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
		endTmp, _ := time.ParseDuration("-120s")
		start = now.Add(endTmp).Format(layout)
	} else {
		start = fromDate
		end = toDate
	}

	sqlDb := dataprovider.NewProvider()

	topKeyStats, _ := sqlDb.GetTopKeysStats(ctx, serverId, start, end)

	var data [][]interface{}
	for _, v := range topKeyStats {
		count, _ := strconv.Atoi(v["total"].(string))
		data = append(data, []interface{}{v["keyname"].(string), count})
	}

	context.JSON(http.StatusOK, gin.H{
		"data":      data,
		"timestamp": now.Format(layout),
	})
}
