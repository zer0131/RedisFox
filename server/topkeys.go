package server

import (
	"github.com/gin-gonic/gin"
	"time"
	"redisfox/dataprovider"
	"strconv"
	"net/http"
)

func (this *Server) topkeys(context *gin.Context)  {
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

	sqlDb,_ := dataprovider.NewProvider(this.config)
	defer sqlDb.Close()

	topKeyStats,_ := sqlDb.GetTopKeysStats(serverId, start, end)

	var data [][]interface{}
	for _,v := range topKeyStats {
		count,_ := strconv.Atoi(v["total"].(string))
		data = append(data, []interface{}{v["keyname"].(string),count})
	}

	context.JSON(http.StatusOK, gin.H{
		"data":data,
		"timestamp":now.Format(layout),
	})
}
