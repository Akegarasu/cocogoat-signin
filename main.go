package main

import (
	"bytes"
	"fmt"
	"github.com/Akegarasu/cocogoat-signin/mihoyo"
	"github.com/Akegarasu/cocogoat-signin/utils"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[椰羊签到][%time%][%lvl%]: %msg% \n",
	})
	log.SetLevel(log.InfoLevel)
	file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Error("读取配置文件失败, 正在尝试为你生成默认配置文件")
		generateDefaultConfig()
		Exit()
	}
	err = yaml.NewDecoder(strings.NewReader(string(file))).Decode(config)
	if err != nil {
		log.Error("配置文件不合法!", err)
		Exit()
	}
	log.Infof("加载配置文件成功: 共 %d 个账户", len(config.Accounts))
}

func main() {
	var err error
	err = configCheck()
	if err != nil {
		log.Error(err)
		Exit()
	}
	log.Info("欢迎使用椰羊签到~")
	for pos, account := range config.Accounts {
		if account.BBSTaskConfig.Enable {
			BBSTask(account, pos)
		}
		if account.SignTask.Genshin {
			GenshinTask()
		}
	}
	log.Info("运行完毕~")
	Exit()
}

func BBSTask(account *Account, pos int) {
	var err error
	if account.Tickets.LoginTicket == "" {
		log.Infof("账户 %d loginTicket 未配置, 尝试从 cookie 中读取", pos)
		cookieMap := utils.ParseCookie(account.Tickets.Cookie)
		if loginTicket, ok := cookieMap["login_ticket"]; ok {
			account.Tickets.LoginTicket = loginTicket
		} else {
			log.Fatalf("账户 %d cookie 错误: 未包含 login_ticket, 请重新按照教程填写", pos)
		}
	}
	m := mihoyo.NewMihoyoBBS(account.Tickets.LoginTicket, account.Tickets.Stuid, account.Tickets.Stoken)
	if m.Stuid == "" || m.Stoken == "" {
		err = m.Login()
		if err != nil {
			log.Error("登录出错, 可能是 cookie 过期了请重新登录 err: ", err)
			return
		}
		log.Info("登录成功, 正在保存相关 ticket 至配置文件")
		account.Tickets.Stoken = m.Stoken
		account.Tickets.Stuid = m.Stuid
		buf := new(bytes.Buffer)
		err = yaml.NewEncoder(buf).Encode(config)
		if err != nil {
			log.Error("格式化配置文件出错", err)
			Exit()
		}
		err = ioutil.WriteFile("config.yml", buf.Bytes(), 0644)
		if err != nil {
			log.Error("保存配置文件出错", err)
		}
	}
	err = m.GetTaskList()
	if err != nil {
		log.Error("获取任务列表失败")
	}
	log.Info("正在获取帖子")
	err = m.GetPostList("26") // 26 原神分区
	if err != nil {
		log.Error("获取帖子失败")
	}
	log.Infof("今日米游社任务: 点赞 (%d/5) 看帖子 (%d/3)", m.Tasks.LikePostsNum, m.Tasks.ReadPostsNum)
	log.Infof("分享 (%d/1) 签到 (%d/1)", m.Tasks.Share, m.Tasks.Signin)
	if m.Tasks.LikePostsNum < 5 && account.BBSTaskConfig.LikePosts {
		log.Info("点赞任务开始")
		go warp(m.LikePosts)
		m.Wg.Add(1)
	}
	if m.Tasks.Share == 0 && account.BBSTaskConfig.Share {
		log.Info("分享任务开始")
		go warp(m.SharePosts)
		m.Wg.Add(1)
	}
	if m.Tasks.ReadPostsNum < 3 && account.BBSTaskConfig.ReadPosts {
		log.Info("阅读帖子任务开始")
		go warp(m.ReadPosts)
		m.Wg.Add(1)
	}
	if m.Tasks.Signin == 0 {
		log.Info("签到任务开始")
		go warp(m.Signin)
		m.Wg.Add(1)
	}
	m.Wg.Wait()
}

func GenshinTask() {
	// todo: add genshin
}

func warp(f func() error) {
	err := f()
	if err != nil {
		log.Error(err)
	}
}

func Exit() {
	var input string
	log.Infoln("按回车退出...")
	_, _ = fmt.Scanln(&input)
	os.Exit(0)
}
