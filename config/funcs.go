package config

import (
	"fmt"
	"os"
	"time"

	"github.com/n9e/dingtalk-sender/corp"

	"github.com/toolkits/pkg/logger"
)

// InitLogger init logger toolkits
func InitLogger() {
	c := Get().Logger

	lb, err := logger.NewFileBackend(c.Dir)
	if err != nil {
		fmt.Println("cannot init logger:", err)
		os.Exit(1)
	}

	lb.SetRotateByHour(true)
	lb.SetKeepHours(c.KeepHours)

	logger.SetLogging(c.Level, lb)
}

func Test(args []string) {
	c := Get()
	var mobile []string

	dingClient := corp.New(c.DingTalk.Token, c.DingTalk.Mobiles,c.DingTalk.IsAtAll)

	if len(args) == 0 {
		fmt.Println("token not given")
		os.Exit(1)
	}

	for i := 0; i < len(args); i++ {
		token := args[i]
		err := dingClient.Send(token, mobile, fmt.Sprintf("test message from n9e at %v", time.Now()))
		if err != nil {
			logger.Error("test send to %s fail: %v\n", args[i], err)
		} else {
			logger.Error("test send to %s success!!!\n", args[i])
		}
		time.Sleep(time.Millisecond*10)
	}
}

