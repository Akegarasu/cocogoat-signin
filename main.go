package main

import (
	"flag"
	"fmt"
	"github.com/Akegarasu/cocogoat-signin/mihoyo"
	"github.com/Akegarasu/cocogoat-signin/utils"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

var rushMysGood bool

func init() {
	flag.BoolVar(&rushMysGood, "g", false, "抢米游社商品")
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
	flag.Parse()
	err = configCheck()
	if err != nil {
		log.Error(err)
		Exit()
	}
	log.Info("欢迎使用椰羊签到~")
	if rushMysGood {
		RushMysGood(config.Accounts[0].Tickets.Cookie)
	} else {
		for pos, account := range config.Accounts {
			if account.BBSTaskConfig.Enable {
				log.Info("开始进行米游社任务")
				BBSTask(account, pos)
			}
			if account.SignTask.Genshin {
				log.Info("开始进行原神签到")
				GenshinTask(account, pos)
			}
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
			log.Errorf("账户 %d cookie 错误: 未包含 login_ticket, 请重新按照教程填写", pos)
			Exit()
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
		saveConfig()
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
	log.Info("今日任务已经完成")
}

func GenshinTask(account *Account, pos int) {
	var err error
	g := mihoyo.NewGenshin(account.Tickets.Cookie)
	err = g.GetAccountList()
	if err != nil {
		log.Error("获取原神账号列表失败", err)
	}
	log.Infof("共获取到 %d 个绑定的原神账号", len(g.Accounts))
	if len(g.Accounts) == 0 {
		log.Errorf("账户 %d 没有绑定原神账号", pos)
		return
	}
	g.SignIn()
}

func RushMysGood(cookie string) {
	var ch, uch int
	var uid string
	h := mihoyo.NewHomuShop(cookie)
	err := h.GetGoodsList()
	if err != nil {
		log.Error("获取兑换列表出错", err)
		Exit()
	}
	for i := 0; i < len(h.GoodList); i++ {
		g := h.GoodList[i]
		log.Infof("%d) %s", i+1, g.GoodsName)
	}
	log.Infof("请选择...")
	for {
		ch, err = utils.ReadNumber()
		if err == nil {
			ch = ch - 1
			break
		}
		log.Warnf("输入非数字 请重新选择")
	}
	log.Infof("选择了 %d 号商品", ch)
	if h.GoodList[ch].Type == 2 {
		log.Info("选择了虚拟商品 正在获取绑定账号")
		genshin := mihoyo.NewGenshin(cookie)
		err = genshin.GetAccountList()
		if err != nil {
			log.Error("获取绑定的原神账号失败")
		}
		if len(genshin.Accounts) > 1 {
			for i := 0; i < len(genshin.Accounts); i++ {
				log.Infof("%d) uid: %s 昵称: %s", i+1, genshin.Accounts[i].Uid, genshin.Accounts[i].NickName)
			}
			log.Info("请选择 uid")
			for {
				uch, err = utils.ReadNumber()
				if err == nil {
					uid = genshin.Accounts[uch-1].Uid
					break
				}
				log.Warnf("输入非数字 请重新选择")
			}

		} else {
			uid = genshin.Accounts[0].Uid
		}
		log.Infof("选择了uid: %s", uid)
	}
	h.Rush(ch, uid)
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
