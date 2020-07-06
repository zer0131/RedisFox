package util

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zer0131/logfox"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewContextWithGinContext(c *gin.Context) context.Context {
	logid, ok := c.Get(logfox.LogIDKey)
	if !ok {
		logid = logfox.GenLogId()
	}
	ctx := logfox.NewContextWithSpecifyLogID(c, logid.(string))
	return ctx
}
