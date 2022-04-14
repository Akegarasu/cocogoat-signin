package main

import (
	_ "embed"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Config struct {
	Accounts []*Account `yaml:"accounts"`
}

type Account struct {
	Tickets struct {
		Cookie      string `yaml:"cookie"`
		Stuid       string `yaml:"stuid"`
		Stoken      string `yaml:"stoken"`
		LoginTicket string `yaml:"loginTicket"`
	} `yaml:"tickets"`
	BBSTaskConfig struct {
		Enable    bool `yaml:"enable"`
		ReadPosts bool `yaml:"readPosts"`
		LikePosts bool `yaml:"likePosts"`
		Unlike    bool `yaml:"unlike"`
		Share     bool `yaml:"share"`
	} `yaml:"BBSTaskConfig"`
	SignTask struct {
		Genshin bool `yaml:"genshin"`
	} `yaml:"SignTask"`
}

var config = &Config{}

// defaultConfig 默认配置文件
//go:embed default_config.yml
var defaultConfig string

func configCheck() error {
	for pos, c := range config.Accounts {
		if c.Tickets.Cookie == "" {
			return errors.New(fmt.Sprintf("第 %d 个账户未配置cookie", pos))
		}
	}
	return nil
}

func generateDefaultConfig() {
	sb := strings.Builder{}
	sb.WriteString(defaultConfig)
	_ = os.WriteFile("config.yml", []byte(sb.String()), 0o644)
	log.Info("默认配置文件已生成，请修改 config.yml 后重新启动!")
}
