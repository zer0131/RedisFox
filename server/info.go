package server

import (
	"github.com/gin-gonic/gin"
	"RedisFox/dataprovider"
	"net/http"
	"strings"
	"strconv"
)

func (s *Server) info(context *gin.Context)  {
	serverId := context.Query("server")

	sqlDb, _ := dataprovider.NewProvider(s.config)
	defer sqlDb.Close()

	redisInfo, _ := sqlDb.GetInfo(serverId)

	var dataBases []map[string]string
	for k,v := range redisInfo {
		if strings.HasPrefix(k, "db") == true {
			vArr := strings.Split(v.(string), ",")
			keyArr := strings.Split(vArr[0], "=")
			keys := keyArr[1]
			expiresArr := strings.Split(vArr[1], "=")
			expires := expiresArr[1]
			database := map[string]string{
				"name":k,
				"keys":keys,
				"expires":expires,
			}
			dataBases = append(dataBases, database)
		}
	}

	totalKeys := 0
	for _,v := range dataBases {
		keyNum, _ := strconv.Atoi(v["keys"])
		totalKeys += keyNum
	}

	if totalKeys == 0 {
		dataBases = []map[string]string{{"name":"db0", "keys" : "0", "expires" : "0",}}
	}

	redisInfo["databases"] = dataBases
	redisInfo["total_keys"] = s.shortenNumber(totalKeys)

	uptimeSeconds, _ := strconv.Atoi(redisInfo["uptime_in_seconds"].(string))
	redisInfo["uptime"] = s.uptimeInSeconds(uptimeSeconds)

	commandsProcessed, _ := strconv.Atoi(redisInfo["total_commands_processed"].(string))
	redisInfo["total_commands_processed_human"] = s.shortenNumber(commandsProcessed)

	context.JSON(http.StatusOK, redisInfo)
}

func (s *Server) shortenNumber(number int) string {
	var val string
	if number < 1000 {
		val = strconv.Itoa(number)
	} else if number < 1000000 {
		if num := number/1000;num == 1000 {
			val = "1M"
		} else {
			val = strconv.Itoa(num) + "K"
		}
	} else if number < 1000000000000 {
		if num := number/1000000000;num == 1000 {
			val = "1T"
		} else {
			val = strconv.Itoa(num) + "G"
		}
	} else {
		num := number/1000000000000
		val = strconv.Itoa(num) + "T"
	}
	return val
}

func (s *Server) uptimeInSeconds(seconds int) string {
	var val string
	if seconds < 60 {
		val = strconv.Itoa(seconds) + "sec"
	} else if seconds <= 3600 {
		if num := seconds/60;num == 60 {
			val = "1hour"
		} else {
			val = strconv.Itoa(num) + "min"
		}
	} else if seconds <= 60*60*24 {
		if num := seconds/3600;num == 24 {
			val = "1day"
		} else {
			val = strconv.Itoa(num) + "hour"
		}
	} else {
		num := seconds/(60*60*24)
		val = strconv.Itoa(num) + "day"
	}
	return val
}

