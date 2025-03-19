package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func Fatal(what string, err error) {
	// 启动错误多次打印
	for index := 0; index < 10; index++ {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"action": what,
		}).Errorf("BOOT FAILED[%d]", index+1)
	}

	// 延迟退出，避免日志丢失
	time.Sleep(time.Second * 3)

	os.Exit(-1)
}
