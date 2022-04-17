package mihoyo

import (
	"encoding/json"
	"github.com/Akegarasu/cocogoat-signin/utils"
	log "github.com/sirupsen/logrus"
	"net/url"
	"sync"
	"time"
)

const (
	success = iota
)

type HomuShop struct {
	cookie   string
	GoodList []*HomuShopGood
	Wg       sync.WaitGroup
}

func NewHomuShop(cookie string) *HomuShop {
	return &HomuShop{
		cookie: cookie,
	}
}

func (h *HomuShop) GetHeaders() map[string]string {
	return map[string]string{
		"x-rpc-client_type":  "5",
		"x-rpc-app_version":  "2.10.0",
		"X-Requested-With":   "com.mihoyo.hyperion",
		"x-rpc-device_id":    utils.UUID,
		"x-rpc-device_name":  "iPhone",
		"x-rpc-device_model": "iPhone13",
		"Referer":            "https://user.mihoyo.com",
		"Host":               "api-takumi.mihoyo.com",
		"User-Agent":         "Mozilla/5.0 (Linux; Android 9; Unspecified Device) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/39.0.0.0 Mobile Safari/537.36 miHoYoBBS/2.7.0",
		"Accept":             "application/json, text/plain, */*",
		"Accept-Encoding":    "deflate",
		"Connection":         "keep-alive",
		"Accept-Language":    "zh-CN,en-US;q=0.8",
		"Origin":             "https://user.mihoyo.com",
		"Cookie":             h.cookie,
	}
}

func (h *HomuShop) GetGoodsList() error {
	base, _ := url.Parse("https://api-takumi.mihoyo.com/common/homushop/v1/web/goods/list")
	params := url.Values{
		"app_id":    {"1"},
		"gids":      {"2"},
		"page":      {"1"},
		"page_size": {"20"},
		"point_sn":  {"myb"},
	}
	base.RawQuery = params.Encode()
	b, err := utils.GetBytes(base.String(), h.GetHeaders())
	if err != nil {
		return err
	}
	g := new(HomuShopGoodListResp)
	err = json.Unmarshal(b, g)
	h.GoodList = g.Data.List
	return nil
}

func (h *HomuShop) exchange(choice int) (int, error) {
	log.Info(choice)
	return success, nil
}

func (h *HomuShop) Rush(choice int) {
	goodChoice := h.GoodList[choice]
	for {
		now := int(time.Now().Unix())
		if goodChoice.NextTime-now > 3 {
			log.Infof("抢购尚未开始 剩余 %d 秒", goodChoice.NextTime-now)
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	log.Infof("距抢购开始小于3秒 开始")
	for i := 1; i < 50; i++ {
		log.Infof("正在抢第 %d 次", i)
		ret, err := h.exchange(choice)
		if ret == success {
			log.Infof("------------第 %d 次成功抢到了哦-------------", i)
			break
		} else {
			log.Infof("抢购失败 原因: %s", err)
		}
	}
}
