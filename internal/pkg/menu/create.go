package menu

import (
	"encoding/json"
	"fmt"
	"wxrobot/internal/app/utils"
	"wxrobot/internal/pkg/access"

	log "github.com/sirupsen/logrus"
)

const (
	kCreateMenuUrlFormat = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s"
)

type ButtonType string

const (
	DEFAULT_BUTTON           ButtonType = ""
	CLICK_BUTTON             ButtonType = "click"                //点击按钮
	VIEW_BUTTON              ButtonType = "view"                 //搜索
	MINI_PROGRAM_BUTTON      ButtonType = "miniprogram"          //网页
	SCANCODE_WAIT_MSG_BUTTON ButtonType = "scancode_waitmsg"     //扫码接收信息
	SCANCODE_PUSH_BUTTON     ButtonType = "scancode_push"        //扫码推送事件
	PIC_SYSTEM_PHOTO_BUTTON  ButtonType = "pic_sysphoto"         //调用系统相机拍照
	PIC_PHOTO_ALBUM_BUTTON   ButtonType = "pic_photo_or_album"   //调用相机或相册
	PIC_WECHAT_BUTTON        ButtonType = "pic_weixin"           //调用微信相册
	LOCATION_SELECT_BUTTON   ButtonType = "location_select"      //调用定位
	MEDIA_ID_BUTTON          ButtonType = "media_id"             //根据素材id下发素材
	ARTICLE_ID_BUTTON        ButtonType = "article_id"           //根据文章id返回文章
	ARTICLE_VIEW_LIMITED     ButtonType = "article_view_limited" //类似articleid
)

type ButtonItem struct {
	Type       *ButtonType   `json:"type,omitempty"`
	Name       *string       `json:"name"`
	Key        *string       `json:"key,omitempty"`
	Url        *string       `json:"url,omitempty"`
	MediaId    *string       `json:"media_id,omitempty"`
	ArticleId  *string       `json:"article_id,omitempty"`
	SubButtons []*ButtonItem `json:"sub_button,omitempty"`
}

type Menu struct {
	Buttons []*ButtonItem `json:"button"`
}

func NewButton(tp ButtonType, name string, key string, url string) *ButtonItem {
	var pkey *string
	var purl *string
	var ptp *ButtonType
	if len(key) > 0 {
		pkey = &key
	}
	if len(url) > 0 {
		purl = &url
	}
	if tp != DEFAULT_BUTTON {
		ptp = &tp
	}
	return &ButtonItem{
		Type: ptp,
		Name: &name,
		Key:  pkey,
		Url:  purl,
	}
}

func TestCreateMenu() error {

	token, err := access.GetAccessToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf(kCreateMenuUrlFormat, token)
	log.Info(url)

	b1 := NewButton(CLICK_BUTTON, "今日歌曲", "TODAY_MUSIC", "")
	b2 := NewButton(DEFAULT_BUTTON, "菜单", "", "")
	subb21 := NewButton(VIEW_BUTTON, "搜索", "", "https://www.baidu.com/")
	subb22 := NewButton(CLICK_BUTTON, "点赞助力", "CLICK_FIGHTING", "")
	b2.SubButtons = append(b2.SubButtons, subb21, subb22)

	menu := Menu{}
	menu.Buttons = append(menu.Buttons, b1, b2)

	body, err := json.Marshal(&menu)
	if err != nil {
		return err
	}

	resp, err := utils.HttpPost(url, body, "application/json")
	if err != nil {
		return err
	}

	log.Info("resp:", resp)
	return nil
}
