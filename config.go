package main

import (
	"bytes"
	_ "embed"
	"github.com/Akegarasu/cocogoat-signin/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
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
			log.Errorf("第 %d 个账户未配置cookie", pos)
			inputCookie(c)
			saveConfig()
		}
	}
	return nil
}

func inputCookie(a *Account) {
	log.Infof("请粘贴 cookie 后按回车: ")
	cookie := utils.ReadLine()
	pc := utils.ParseCookie(cookie)
	if _, ok := pc["login_ticket"]; !ok {
		log.Error("该 cookie 缺少 login_ticket 请确认按照教程登录了两个网站")
	}
	a.Tickets.Cookie = cookie
}

func generateDefaultConfig() {
	sb := strings.Builder{}
	sb.WriteString(defaultConfig)
	_ = os.WriteFile("config.yml", []byte(sb.String()), 0o644)
	log.Info("默认配置文件已生成, 请重新启动")
}

func saveConfig() {
	buf := new(bytes.Buffer)
	err := yaml.NewEncoder(buf).Encode(config)
	if err != nil {
		log.Error("格式化配置文件出错", err)
		Exit()
	}
	err = ioutil.WriteFile("config.yml", buf.Bytes(), 0644)
	if err != nil {
		log.Error("保存配置文件出错", err)
	}
}
