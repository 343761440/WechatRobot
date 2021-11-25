package access

import (
	"encoding/json"
	"fmt"
	"wxrobot/internal/app/common"
	"wxrobot/internal/app/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//GetAccessToken Success, token:51_KBIJELcopd5RHT8t5o6uOPDeg8ARFW4FBKyRvNmY4Xt3E0H6R4bnnKNCcE_T51gDDF8XpE5qierFKzMcJPD-x_xewH7l4-CgUa0XoDxIRsZ2IsKaNsj3VWRe0xlCAqx8LIPETiKkcmYW_QIaWRXfAAADYE
//expire:7200

func GetAccessToken() (token string, expireTime int, err error) {
	type Model struct {
		Token  string `json:"access_token"`
		Expire int    `json:"expires_in"`
	}
	appid := viper.GetString(common.CFG_APPID)
	appsecret := viper.GetString(common.CFG_APPSECRET)
	host := viper.GetString(common.CFG_ACS_TOKEN_URL)
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", host, appid, appsecret)
	log.Info("url:", url)
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

	//写到redis缓存里面
	return
}
