package mihoyo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Akegarasu/cocogoat-signin/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
	"sync"
	"time"
)

type Tasks struct {
	Signin       int
	LikePostsNum int
	ReadPostsNum int
	Share        int
}

type Post struct {
	PostID  string
	Subject string
}

type MihoyoBBS struct {
	Stuid       string
	Stoken      string
	LoginTicket string
	Tasks       *Tasks
	Posts       []Post
	Wg          sync.WaitGroup
}

func NewMihoyoBBS(LoginTicket string, stuid string, stoken string) *MihoyoBBS {
	return &MihoyoBBS{
		LoginTicket: LoginTicket,
		Stuid:       stuid,
		Stoken:      stoken,
		Tasks:       &Tasks{Signin: 0, LikePostsNum: 0, ReadPostsNum: 0, Share: 0},
		Posts:       make([]Post, 0),
		Wg:          sync.WaitGroup{},
	}
}

func (m *MihoyoBBS) Login() error {
	r, err := m.doLogin()
	if err != nil {
		return err
	}
	msg := r.Get("data.msg").String()
	if !strings.Contains(msg, "成功") {
		return errors.New(msg)
	}
	accountID := r.Get("data.cookie_info.account_id").String()
	m.Stuid = accountID
	r, err = m.getMultiTokenByLoginTicket()
	if err != nil {
		return err
	}
	m.Stoken = r.Get("data.list.0.token").String()
	return nil
}

func (m *MihoyoBBS) GetHeaders() map[string]string {
	headers := map[string]string{
		"x-rpc-client_type":  "2",
		"x-rpc-app_version":  "2.7.0",
		"x-rpc-sys_version":  "6.0.1",
		"x-rpc-channel":      "mihoyo",
		"x-rpc-device_id":    "f30a320c-aa0f-43a8-a3d5-971a9e2efcc0",
		"x-rpc-device_name":  "Mi 10",
		"x-rpc-device_model": "Mi 10",
		"Referer":            "https://app.mihoyo.com",
		"Host":               "bbs-api.mihoyo.com",
		"User-Agent":         "okhttp/4.8.0",
	}

	headers["DS"] = utils.DS(0)
	headers["cookie"] = fmt.Sprintf("stuid=%s; stoken=%s", m.Stuid, m.Stoken)
	return headers
}

func (m *MihoyoBBS) doLogin() (gjson.Result, error) {
	url := fmt.Sprintf("https://webapi.account.mihoyo.com/Api/cookie_accountinfo_by_loginticket?login_ticket=%s", m.LoginTicket)
	return utils.GetJson(url, nil)
}

func (m *MihoyoBBS) getMultiTokenByLoginTicket() (gjson.Result, error) {
	url := fmt.Sprintf("https://api-takumi.mihoyo.com/auth/api/getMultiTokenByLoginTicket?login_ticket=%s&token_types=3&uid=%s", m.LoginTicket, m.Stuid)
	return utils.GetJson(url, nil)
}

func (m *MihoyoBBS) GetTaskList() error {
	log.Info("正在获取任务列表")
	url := "https://bbs-api.mihoyo.com/apihub/sapi/getUserMissionsState"
	b, err := utils.GetBytes(url, m.GetHeaders())
	if err != nil {
		return err
	}
	t := new(TaskListResp)
	err = json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	if t.Retcode != 0 || strings.Contains(t.Message, "err") {
		return errors.New("获取任务列表失败 可能是cookie过期了请重新登录获取")
	}
	if t.Data.CanGetPoints == 0 {
		m.Tasks.Share = 1
		m.Tasks.Signin = 1
		m.Tasks.LikePostsNum = 5
		m.Tasks.ReadPostsNum = 3
		return nil
	}
	// 所有任务都没做
	if t.Data.States[0].MissionID >= 62 {
		return nil
	}

	for _, s := range t.Data.States {
		switch s.MissionID {
		case 58:
			if s.IsGetAward {
				m.Tasks.Signin = 1
			}
		case 59:
			if s.IsGetAward {
				m.Tasks.ReadPostsNum = 3
			} else {
				m.Tasks.ReadPostsNum += s.HappenedTimes
			}
		case 60:
			if s.IsGetAward {
				m.Tasks.LikePostsNum = 5
			} else {
				m.Tasks.LikePostsNum += s.HappenedTimes
			}
		case 61:
			if s.IsGetAward {
				m.Tasks.Share = 1
			}
		}
	}
	return nil
}

func (m *MihoyoBBS) GetPostList(forumID string) error {
	url := fmt.Sprintf("https://bbs-api.mihoyo.com/post/api/getForumPostList?forum_id=%s&is_good=false&is_hot=false&page_size=20&sort_type=1", forumID)
	b, err := utils.GetBytes(url, m.GetHeaders())
	if err != nil {
		return err
	}
	p := new(PostListResp)
	err = json.Unmarshal(b, p)
	if err != nil {
		return err
	}
	for _, p := range p.Data.List[:5] {
		m.Posts = append(m.Posts, Post{PostID: p.Post.PostID, Subject: p.Post.Subject})
	}
	log.Infof("获取帖子成功, 共获取 %d 个帖子", len(p.Data.List))
	return nil
}

func (m *MihoyoBBS) ReadPosts() error {
	defer m.Wg.Done()
	if m.Tasks.ReadPostsNum >= 3 {
		return nil
	}
	for i := 1; i <= 3-m.Tasks.ReadPostsNum; i++ {
		url := fmt.Sprintf("https://bbs-api.mihoyo.com/post/api/getPostFull?post_id=%s", m.Posts[i].PostID)
		r, err := utils.GetJson(url, m.GetHeaders())
		if err != nil {
			log.Warnf("看第%d个帖子失败了: %v", i, err)
		}
		msg := r.Get("message").String()
		if msg == "OK" {
			log.Infof("看第 %d 个帖子成功~ 帖子主题: %s", i, m.Posts[i].Subject)
		} else {
			log.Warnf("看第 %d 个帖子失败了! 报错: %s", i, msg)
		}
		time.Sleep(time.Second * 3)
	}
	return nil
}

func (m *MihoyoBBS) LikePosts() error {
	defer m.Wg.Done()
	url := "https://bbs-api.mihoyo.com/apihub/sapi/upvotePost"
	type LikeData struct {
		PostID   string `json:"post_id"`
		IsCancel bool   `json:"is_cancel"`
	}
	for i := 4 - m.Tasks.LikePostsNum; i >= 0; i-- {
		data := &LikeData{
			PostID:   m.Posts[i].PostID,
			IsCancel: false,
		}
		bd, _ := json.Marshal(data)
		b, err := utils.PostBytes(url, bd, m.GetHeaders())
		if err != nil {
			return err
		}
		msg := gjson.ParseBytes(b).Get("message").String()
		if msg == "OK" {
			log.Infof("点赞成功 帖子主题: %s", m.Posts[i].Subject)
		} else {
			log.Warnf("点赞失败 报错: %s", msg)
		}
		time.Sleep(time.Second * 3)
	}
	return nil
}

func (m *MihoyoBBS) SharePosts() error {
	defer m.Wg.Done()
	url := fmt.Sprintf("https://bbs-api.mihoyo.com/apihub/api/getShareConf?entity_id=%s&entity_type=1", m.Posts[0].PostID)
	r, err := utils.GetJson(url, m.GetHeaders())
	if err != nil {
		log.Warnln("分享帖子失败了: ", err)
	}
	msg := r.Get("message").String()
	if msg == "OK" {
		log.Infof("分享帖子成功~ 帖子主题: %s", m.Posts[0].Subject)
	} else {
		log.Warnf("分享帖子失败了! 报错: %s", msg)
	}
	return nil
}

func (m *MihoyoBBS) Signin() error {
	defer m.Wg.Done()
	url := fmt.Sprintf("https://bbs-api.mihoyo.com/apihub/sapi/signIn?gids=%s", "26")
	b, err := utils.PostBytes(url, []byte(`{}`), m.GetHeaders())
	if err != nil {
		return err
	}
	msg := gjson.ParseBytes(b).Get("message").String()
	if msg == "OK" {
		log.Info("签到成功")
	} else {
		log.Warnf("签到失败~ 可能是 cookie 过期了请重新获取cookie! 报错: %s", msg)
	}
	return nil
}
