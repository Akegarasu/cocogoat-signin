package mihoyo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Akegarasu/cocogoat-signin/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"time"
)

type Account struct {
	NickName string
	Uid      string
	Region   string
	SignInfo GenshinSignInfo
}

type Genshin struct {
	cookie   string
	Accounts []*Account
}

func NewGenshin(cookie string) *Genshin {
	return &Genshin{
		cookie: cookie,
	}
}

func (g *Genshin) GetHeaders() map[string]string {
	headers := map[string]string{
		"x-rpc-client_type": "5", // pc web=4 / mobile web=5
		"x-rpc-app_version": "2.3.0",
		"X-Requested-With":  "com.mihoyo.hyperion",
		"x-rpc-device_id":   utils.UUID,
		"Referer":           "https://webstatic.mihoyo.com/bbs/event/signin-ys/index.html?bbs_auth_required=true&act_id=e202009291139501&utm_source=bbs&utm_medium=mys&utm_campaign=icon",
		"Host":              "bbs-api.mihoyo.com",
		"User-Agent":        "Mozilla/5.0 (Linux; Android 9; Unspecified Device) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/39.0.0.0 Mobile Safari/537.36 miHoYoBBS/2.3.0",
		"Accept":            "application/json, text/plain, */*",
		"Accept-Encoding":   "gzip, deflate",
		"Accept-Language":   "zh-CN,en-US;q=0.8",
		"Origin":            "https://webstatic.mihoyo.com",
	}
	headers["DS"] = utils.DS(2)
	headers["cookie"] = g.cookie
	return headers
}

func (g *Genshin) GetAccountList() error {
	url := "https://api-takumi.mihoyo.com/binding/api/getUserGameRolesByCookie?game_biz=hk4e_cn"
	j, err := utils.GetBytes(url, g.GetHeaders())
	if err != nil {
		return err
	}
	ar := new(GenshinAccountsResp)
	err = json.Unmarshal(j, ar)
	if err != nil {
		return err
	}
	if ar.Retcode != 0 {
		return errors.New("米游社 cookie 错误")
	}
	for _, a := range ar.Data.List {
		g.Accounts = append(g.Accounts, &Account{
			NickName: a.Nickname,
			Uid:      a.GameUID,
			Region:   a.Region,
		})
	}
	return nil
}

func (g *Genshin) signInfo(a *Account) error {
	url := fmt.Sprintf("https://api-takumi.mihoyo.com/event/bbs_sign_reward/info?act_id=%s&region=%s&uid=%s", "e202009291139501", a.Region, a.Uid)
	b, err := utils.GetBytes(url, g.GetHeaders())
	if err != nil {
		return err
	}
	i := new(GenshinSignInfoResp)
	err = json.Unmarshal(b, i)
	if err != nil {
		return err
	}
	a.SignInfo = i.Data
	return nil
}

func (g *Genshin) doSignIn(a *Account) error {
	url := "https://api-takumi.mihoyo.com/event/bbs_sign_reward/sign"
	data := &GenshinSignPostData{
		ActID:  "e202009291139501",
		UID:    a.Uid,
		Region: a.Region,
	}
	d, _ := json.Marshal(data)
	r, err := utils.PostBytes(url, d, g.GetHeaders())
	if err != nil {
		return err
	}
	retcode := gjson.ParseBytes(r).Get("retcode").Int()
	switch retcode {
	case 0:
		log.Infof("UID: %s, 昵称: %s 签到成功", a.Uid, a.NickName)
	case -5003:
		log.Infof("UID: %s, 昵称: %s 今天已经签到过了", a.Uid, a.NickName)
	default:
		return errors.New("签到失败")
	}
	return nil
}

func (g *Genshin) SignIn() {
	for pos, a := range g.Accounts {
		err := g.signInfo(a)
		if err != nil {
			log.Errorf("第 %d 个原神账号检查签到状况失败: %v", pos, err)
			continue
		}
		if a.SignInfo.FirstBind {
			log.Infof("UID: %s, 昵称: %s 首次绑定米游社, 请先手动签到一次", a.Uid, a.NickName)
			continue
		}
		if a.SignInfo.IsSign {
			log.Infof("UID: %s, 昵称: %s 今天已经签到过了", a.Uid, a.NickName)
			continue
		}
		err = g.doSignIn(a)
		if err != nil {
			log.Errorf("UID: %s, 昵称: %s 签到失败", a.Uid, a.NickName)
		}
		if pos != len(g.Accounts)-1 {
			time.Sleep(time.Second * 5)
		}
	}
}
