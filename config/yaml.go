package config

import (
	"fmt"

	"github.com/toolkits/pkg/file"
)

type Config struct {
	Logger   loggerSection   `yaml:"logger"`
	DingTalk dingTalkSection `yaml:"dingtalk"`
	Consumer consumerSection `yaml:"consumer"`
	Redis    redisSection    `yaml:"redis"`
}

type loggerSection struct {
	Dir       string `yaml:"dir"`
	Level     string `yaml:"level"`
	KeepHours uint   `yaml:"keepHours"`
}

type redisSection struct {
	Addr    string         `yaml:"addr"`
	Pass    string         `yaml:"pass"`
	DB      int            `yaml:"db"`
	Idle    int            `yaml:"idle"`
	Timeout timeoutSection `yaml:"timeout"`
}

type timeoutSection struct {
	Conn  int `yaml:"conn"`
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
}

type consumerSection struct {
	Queue  string `yaml:"queue"`
	Worker int    `yaml:"worker"`
}

type dingTalkSection struct {
	Token   string   `yaml:"token"`
	Mobiles []string `yaml:"mobiles"`
	IsAtAll bool     `yaml:"is_at_all"`
}

var yaml Config

func Get() Config {
	return yaml
}

func ParseConfig(yf string) error {
	err := file.ReadYaml(yf, &yaml)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", yf, err)
	}
	return nil
}
