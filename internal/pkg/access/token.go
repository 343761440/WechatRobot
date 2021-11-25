package access

import (
	"encoding/json"
	"fmt"
	"time"
	"wxrobot/internal/app/common"
	"wxrobot/internal/app/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	kWxAkRediskey = "WxAccessToken"
)

func getAccessToken() (token string, expireTime int, err error) {
	type Model struct {
		Token  string `json:"access_token"`
		Expire int    `json:"expires_in"`
	}
	appid := viper.GetString(common.CFG_SECRET_APPID)
	appsecret := viper.GetString(common.CFG_SECRET_APPSECRET)
	host := viper.GetString(common.CFG_URL_ACS_TOKEN)
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", host, appid, appsecret)
	resp, err := utils.HttpGet(url)
	if err != nil {
		log.Warn("get failed")
		return
	}

	var m Model
	err = json.Unmarshal(resp, &m)
	if err != nil {
		return
	}
	token = m.Token
	expireTime = m.Expire
	log.Info("GetAccessToken Success, token:", token, " expire:", expireTime)
	return
}

//export
func GetAccessToken() (token string, err error) {
	rdb, err := common.GetRedisClient()
	if err != nil {
		log.Warn("init redis failed, err=", err)
		return
	}

	isExist, err := rdb.Exists(kWxAkRediskey).Result()
	if err != nil {
		return
	}
	if isExist {
		leavetime, err1 := rdb.TTL(kWxAkRediskey).Result()
		if err1 != nil {
			err1 = err
			return
		}
		//提前10分钟获取
		if leavetime >= time.Minute*10 {
			acsToken, err2 := rdb.Get(kWxAkRediskey).Result()
			if err2 != nil {
				err = err2
				return
			}
			token = acsToken
			return
		}
	}

	token, expire, err := getAccessToken()
	if err != nil {
		return
	}
	err = rdb.Set(kWxAkRediskey, token, time.Second*time.Duration(expire)).Err()
	if err != nil {
		return
	}
	return
}
