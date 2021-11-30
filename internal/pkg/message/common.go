package message

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"wxrobot/internal/pkg/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	kCtxUserKey    = "user"
	kCtxUserId     = "userId"
	kCtxFromUserId = "fromUserId"
)

type MessageType string

const (
	TEXT_MESSAGE        MessageType = "text"
	IMAGE_MESSAGE       MessageType = "image"
	VOICE_MESSAGE       MessageType = "voice"
	VIDEO_MESSAGE       MessageType = "video"
	SHORT_VIDEO_MESSAGE MessageType = "shortvideo"
	LOCATION_MESSAGE    MessageType = "location"
	LINK_MESSAGE        MessageType = "link"
	EVENT_MESSAGE       MessageType = "event"
)

type MessageCommon struct {
	XMLName      xml.Name    `xml:"xml"` // 指定最外层的标签为config
	ToUserName   string      `xml:"ToUserName"`
	FromUserName string      `xml:"FromUserName"`
	CreateTime   int64       `xml:"CreateTime"`
	MsgType      MessageType `xml:"MsgType"`
}

var gMsgHandler = map[MessageType]func(c *gin.Context, body []byte){
	TEXT_MESSAGE:  handleTextMessage,
	IMAGE_MESSAGE: handleImageMessage,
	EVENT_MESSAGE: handleEventMessage,
}

func CommonMessageProxy(c *gin.Context) {
	log.Debug("request:", c.Request)

	body, _ := ioutil.ReadAll(c.Request.Body)
	msg := MessageCommon{}
	err := xml.Unmarshal(body, &msg)
	if err != nil {
		log.Error("HandleMessage err:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Set(kCtxUserId, msg.FromUserName)
	c.Set(kCtxFromUserId, msg.ToUserName)
	wxUser, err := model.GetWxUser(msg.FromUserName)
	if err == nil {
		c.Set(kCtxUserKey, wxUser)
	}
	handleCb := gMsgHandler[msg.MsgType]
	handleCb(c, body)
}

func GetWxUserFromCtx(c *gin.Context) *model.WxUser {
	if v, _ := c.Get(kCtxUserKey); v != nil {
		if user, ok := v.(*model.WxUser); ok {
			return user
		}
	}
	return nil
}

func getKeyFromCtx(c *gin.Context, key string) string {
	if v, _ := c.Get(key); v != nil {
		if vstr, ok := v.(string); ok {
			return vstr
		}
	}
	return ""
}

func GetUserIdFromCtx(c *gin.Context) string {
	return getKeyFromCtx(c, kCtxUserId)
}

func GetFromUserIdFromCtx(c *gin.Context) string {
	return getKeyFromCtx(c, kCtxFromUserId)
}
