package util

import (
	"os"
	"fmonitor/flog"
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

func CheckError(err error) bool {
	if err != nil {
		flog.Fatalf(err.Error())
		return false
	}
	return true
}
